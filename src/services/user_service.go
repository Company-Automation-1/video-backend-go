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

// GetListWithPagination 获取用户列表（分页）
func (s *UserService) GetListWithPagination(
	offset, limit int,
	conditions ...gen.Condition,
) ([]*models.User, int64, error) {
	users, count, err := query.User.Where(conditions...).FindByPage(offset, limit)
	return users, count, err
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

// Update 更新用户
func (s *UserService) Update(id uint, user *models.User) (*models.User, error) {
	// 检查用户名是否已存在
	if user.Username != "" {
		existingUser, err := query.User.Where(query.User.Username.Eq(user.Username), query.User.ID.Neq(id)).First()
		if err == nil && existingUser != nil {
			return nil, tools.ErrBadRequest("用户名已存在")
		}
	}

	// 密码加密
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, tools.ErrInternalServer("密码设置失败")
		}
		user.Password = string(hashedPassword)
	}

	// 如果提供了积分（包括置空），单独更新
	if user.Points != nil {
		_, err := query.User.Where(query.User.ID.Eq(id)).Update(query.User.Points, user.Points)
		if err != nil {
			return nil, tools.ErrInternalServer("积分更新失败")
		}
	}

	// 移除Points字段，避免在Updates中处理（已单独处理）
	userForUpdate := *user
	userForUpdate.Points = nil

	// 使用 Updates 只更新有值的字段（忽略零值字段，自动更新 UpdatedAt）
	_, err := query.User.Where(query.User.ID.Eq(id)).Updates(&userForUpdate)
	if err != nil {
		return nil, tools.ErrInternalServer("用户更新失败")
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
	// 检查邮箱是否已存在
	existingUser, err := query.User.Where(query.User.Email.Eq(email)).First()
	if err == nil && existingUser != nil {
		return tools.ErrBadRequest("邮箱已被使用")
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return tools.ErrInternalServer("邮箱检查失败")
	}

	_, err = s.captcha.SendCode(ctx, email, CaptchaTypeRegister)
	if err != nil {
		// 保留原始错误的 code 和 message（SendCode 已使用 AppError 处理内部错误）
		return tools.WrapError(err)
	}
	return nil
}

// Register 用户注册（业务逻辑：验证码验证、密码加密、唯一性校验）
func (s *UserService) Register(ctx context.Context, username, email, password, captcha string) error {
	// 检查用户名或邮箱是否已存在（使用 OR 条件，只需查询一次）
	existingUser, err := query.User.Where(query.User.Username.Eq(username)).
		Or(query.User.Email.Eq(email)).
		First()
	if err == nil && existingUser != nil {
		// 判断是用户名冲突还是邮箱冲突
		if existingUser.Username == username {
			return tools.ErrBadRequest("用户名已存在")
		}
		if existingUser.Email == email {
			return tools.ErrBadRequest("邮箱已被使用")
		}
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return tools.ErrInternalServer("用户检查失败")
	}

	// 验证验证码
	valid, err := s.captcha.VerifyCaptcha(ctx, email, captcha, CaptchaTypeRegister)
	if err != nil {
		return tools.ErrBadRequest(err.Error())
	}
	if !valid {
		return tools.ErrBadRequest("验证码错误")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return tools.ErrInternalServer("密码设置失败")
	}

	// 创建用户（验证码已验证，邮箱已验证）
	user := &models.User{
		Username:      username,
		Email:         email,
		Password:      string(hashedPassword),
		EmailVerified: true,
	}

	if err := query.User.Create(user); err != nil {
		return tools.ErrInternalServer("用户创建失败")
	}

	return nil
}

// UpdateEmail 更新邮箱
func (s *UserService) UpdateEmail(ctx context.Context, userID uint, newEmail, code string) error {
	// 验证验证码（针对新邮箱，使用注册验证码类型）
	valid, err := s.captcha.VerifyCaptcha(ctx, newEmail, code, CaptchaTypeRegister)
	if err != nil {
		return tools.ErrBadRequest(err.Error())
	}
	if !valid {
		return tools.ErrBadRequest("验证码错误")
	}

	// 检查新邮箱是否已被使用
	existingUser, err := query.User.Where(query.User.Email.Eq(newEmail)).First()
	if err == nil && existingUser != nil {
		return tools.ErrBadRequest("邮箱已被使用")
	}

	// 更新用户邮箱
	_, err = query.User.Where(query.User.ID.Eq(userID)).Update(query.User.Email, newEmail)
	if err != nil {
		return tools.ErrInternalServer("邮箱更新失败")
	}

	return nil
}
