// Package services 用户服务层（仅包含业务逻辑）
package services

import (
	"context"
	"strings"

	"github.com/Company-Automation-1/video-backend-go/src/api/dto"
	"github.com/Company-Automation-1/video-backend-go/src/models"
	"github.com/Company-Automation-1/video-backend-go/src/query"
	"github.com/Company-Automation-1/video-backend-go/src/tools"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// UserService 用户业务服务
type UserService struct {
	captcha *CaptchaService
}

// NewUserService 创建用户业务服务
func NewUserService(captcha *CaptchaService) *UserService {
	return &UserService{
		captcha: captcha,
	}
}

// GetListWithQuery 获取用户列表（条件查询、模糊查询、范围查询）
func (s *UserService) GetListWithQuery(queryReq *dto.UserListQueryRequest) ([]*models.User, int64, error) {
	// 1. 构建基础查询条件
	conditions := tools.NewConditionBuilder().
		EqUint(&query.User.ID, queryReq.ID).
		EqString(&query.User.Username, queryReq.Username).
		EqString(&query.User.Email, queryReq.Email).
		EqBool(&query.User.EmailVerified, queryReq.EmailVerified).
		Like(&query.User.Username, queryReq.UsernameLike).
		Like(&query.User.Email, queryReq.EmailLike).
		GteInt64(&query.User.CreatedAt, queryReq.CreatedAtMin).
		LteInt64(&query.User.CreatedAt, queryReq.CreatedAtMax).
		Build()

	q := query.User.Where(conditions...)

	// 2. 积分范围查询（特殊处理：NULL 在业务逻辑上等于 0）
	// points_min: 当值为 0 时匹配 NULL 或 0，当值 > 0 时只匹配 points >= min（不包含 NULL）
	if queryReq.PointsMin != nil {
		if *queryReq.PointsMin == 0 {
			q = q.Where(query.User.Points.IsNull()).Or(query.User.Points.Eq(0))
		} else {
			q = q.Where(query.User.Points.Gte(*queryReq.PointsMin))
		}
	}
	// points_max: 当值为 0 时匹配 NULL 或 0，当值 > 0 时匹配 NULL 或 points <= max（包含 NULL）
	if queryReq.PointsMax != nil {
		if *queryReq.PointsMax == 0 {
			q = q.Where(query.User.Points.IsNull()).Or(query.User.Points.Eq(0))
		} else {
			q = q.Where(query.User.Points.IsNull()).Or(query.User.Points.Lte(*queryReq.PointsMax))
		}
	}

	// 3. 排序（字段存在时应用排序）
	if queryReq.OrderBy != "" {
		if orderField, ok := query.User.GetFieldByName(queryReq.OrderBy); ok {
			if strings.EqualFold(queryReq.Order, "asc") {
				q = q.Order(orderField.Asc())
			} else {
				q = q.Order(orderField.Desc())
			}
		}
	}

	// 4. 执行分页查询
	users, count, err := q.FindByPage(queryReq.GetOffset(), queryReq.GetLimit())
	return users, count, err
}

// GetOne 获取单个用户
func (s *UserService) GetOne(conditions ...gen.Condition) (*models.User, error) {
	user, err := query.User.Where(conditions...).First()
	if err == gorm.ErrRecordNotFound {
		return nil, tools.ErrNotFound("用户不存在")
	}
	return user, err
}

// Update 更新用户（用户自己更新，不允许修改积分）
func (s *UserService) Update(ctx context.Context, id uint, user *models.User, emailCode *string) (*models.User, error) {
	excludeID := &id

	// 组合调用：按需校验和处理
	if err := s.validateUsername(user.Username, excludeID); err != nil {
		return nil, err
	}

	if err := s.validateEmail(user.Email, excludeID); err != nil {
		return nil, err
	}

	// 如果更新邮箱，需要验证码（verifyEmailCaptcha 内部会检查 email 是否为空）
	if user.Email != "" {
		code := ""
		if emailCode != nil {
			code = *emailCode
		}
		if err := s.verifyEmailCaptcha(ctx, user.Email, code); err != nil {
			return nil, err
		}
		user.EmailVerified = true
	}

	hashedPassword, err := s.encryptPassword(user.Password)
	if err != nil {
		return nil, err
	}
	if hashedPassword != "" {
		user.Password = hashedPassword
	}

	// 使用 Updates 只更新有值的字段（忽略零值字段，自动更新 UpdatedAt）
	_, err = query.User.Where(query.User.ID.Eq(id)).Updates(user)
	if err != nil {
		return nil, tools.ErrInternalServer("用户更新失败")
	}

	// 返回更新后的用户
	return s.GetOne(query.User.ID.Eq(id))
}

// UpdatePoints 更新用户积分（管理员权限）
func (s *UserService) UpdatePoints(id uint, points *int) (*models.User, error) {
	// 如果提供了积分（包括置空），单独更新
	if points != nil {
		_, err := query.User.Where(query.User.ID.Eq(id)).Update(query.User.Points, points)
		if err != nil {
			return nil, tools.ErrInternalServer("积分更新失败")
		}
	}

	// 返回更新后的用户
	return s.GetOne(query.User.ID.Eq(id))
}

// Delete 删除用户
func (s *UserService) Delete(id uint) error {
	_, err := query.User.Where(query.User.ID.Eq(id)).Delete()
	return err
}

// SendVerificationCode 发送验证码
func (s *UserService) SendVerificationCode(ctx context.Context, email string) error {
	if err := s.validateEmail(email, nil); err != nil {
		return err
	}

	_, err := s.captcha.SendCode(ctx, email, CaptchaTypeRegister)
	if err != nil {
		// 保留原始错误的 code 和 message（SendCode 已使用 AppError 处理内部错误）
		return tools.WrapError(err)
	}
	return nil
}

// Register 用户注册（业务逻辑：验证码验证、密码加密、唯一性校验）
func (s *UserService) Register(ctx context.Context, username, email, password, captcha string) error {
	if err := s.validateUsername(username, nil); err != nil {
		return err
	}

	if err := s.validateEmail(email, nil); err != nil {
		return err
	}

	if err := s.verifyEmailCaptcha(ctx, email, captcha); err != nil {
		return err
	}

	hashedPassword, err := s.encryptPassword(password)
	if err != nil {
		return err
	}

	// 创建用户（验证码已验证，邮箱已验证）
	user := &models.User{
		Username:      username,
		Email:         email,
		Password:      hashedPassword,
		EmailVerified: true,
	}

	if err := query.User.Create(user); err != nil {
		return tools.ErrInternalServer("用户创建失败")
	}
	return nil
}

// validateUsername 校验用户名唯一性
func (s *UserService) validateUsername(username string, excludeID *uint) error {
	if username == "" {
		return nil
	}
	conditions := []gen.Condition{query.User.Username.Eq(username)}
	if excludeID != nil {
		conditions = append(conditions, query.User.ID.Neq(*excludeID))
	}
	existingUser, err := query.User.Where(conditions...).First()
	if err == nil && existingUser != nil {
		return tools.ErrBadRequest("用户名已存在")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return tools.ErrInternalServer("用户名检查失败")
	}
	return nil
}

// validateEmail 校验邮箱唯一性
func (s *UserService) validateEmail(email string, excludeID *uint) error {
	if email == "" {
		return nil
	}
	conditions := []gen.Condition{query.User.Email.Eq(email)}
	if excludeID != nil {
		conditions = append(conditions, query.User.ID.Neq(*excludeID))
	}
	existingUser, err := query.User.Where(conditions...).First()
	if err == nil && existingUser != nil {
		return tools.ErrBadRequest("邮箱已被使用")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return tools.ErrInternalServer("邮箱检查失败")
	}
	return nil
}

// verifyEmailCaptcha 验证邮箱验证码
func (s *UserService) verifyEmailCaptcha(ctx context.Context, email, code string) error {
	if email == "" {
		return nil
	}
	if code == "" {
		return tools.ErrBadRequest("更新邮箱需要验证码")
	}
	valid, err := s.captcha.VerifyCaptcha(ctx, email, code, CaptchaTypeRegister)
	if err != nil {
		return tools.ErrBadRequest(err.Error())
	}
	if !valid {
		return tools.ErrBadRequest("验证码错误")
	}
	return nil
}

// encryptPassword 加密密码
func (s *UserService) encryptPassword(password string) (string, error) {
	if password == "" {
		return "", nil
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", tools.ErrInternalServer("密码设置失败")
	}
	return string(hashedPassword), nil
}
