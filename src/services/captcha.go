// Package services 验证码服务
package services

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/Company-Automation-1/video-backend-go/src/infrastructure"
	"github.com/Company-Automation-1/video-backend-go/src/tools"
	"github.com/redis/go-redis/v9"
)

const (
	// CaptchaExpire 验证码有效期5分钟
	CaptchaExpire = 5 * time.Minute
	// CaptchaPrefix Redis key 前缀
	CaptchaPrefix = "captcha:"
)

// CaptchaType 验证码类型
type CaptchaType string

const (
	// CaptchaTypeRegister 注册验证码
	CaptchaTypeRegister CaptchaType = "register"
	// CaptchaTypeReset 重置密码验证码
	CaptchaTypeReset CaptchaType = "reset"
)

// CaptchaService 验证码服务
type CaptchaService struct {
	redisClient *infrastructure.Redis
	email       *infrastructure.Email
}

// NewCaptchaService 创建验证码服务
func NewCaptchaService(redisClient *infrastructure.Redis, email *infrastructure.Email) *CaptchaService {
	return &CaptchaService{
		redisClient: redisClient,
		email:       email,
	}
}

// GenerateCode 生成6位数字验证码
func (s *CaptchaService) GenerateCode() string {
	var b [4]byte
	if _, err := rand.Read(b[:]); err != nil {
		// 如果 crypto/rand 失败，使用时间戳作为后备
		return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	}
	return fmt.Sprintf("%06d", binary.BigEndian.Uint32(b[:])%1000000)
}

// checkCaptchaExists 检查验证码是否已存在
func (s *CaptchaService) checkCaptchaExists(ctx context.Context, email string, captchaType CaptchaType) (bool, error) {
	key := CaptchaPrefix + string(captchaType) + ":" + email
	exists, err := s.redisClient.Exists(ctx, key)
	if err != nil {
		return false, fmt.Errorf("检查验证码失败: %w", err)
	}
	return exists, nil
}

// SetCaptcha 设置验证码到 Redis
func (s *CaptchaService) SetCaptcha(ctx context.Context, email, code string, captchaType CaptchaType) error {
	key := CaptchaPrefix + string(captchaType) + ":" + email
	return s.redisClient.Set(ctx, key, code, CaptchaExpire)
}

// VerifyCaptcha 验证验证码
func (s *CaptchaService) VerifyCaptcha(ctx context.Context, email, code string, captchaType CaptchaType) (bool, error) {
	key := CaptchaPrefix + string(captchaType) + ":" + email

	storedCode, err := s.redisClient.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			return false, errors.New("验证码已过期或不存在")
		}
		return false, err
	}

	if storedCode != code {
		return false, errors.New("验证码错误")
	}

	// 验证成功后删除验证码（一次性使用）
	_ = s.redisClient.Del(ctx, key) //nolint:errcheck // 验证成功后删除验证码，失败不影响结果
	return true, nil
}

// SendCode 发送验证码
func (s *CaptchaService) SendCode(ctx context.Context, email string, captchaType CaptchaType) (string, error) {
	// 检查是否已存在验证码
	exists, err := s.checkCaptchaExists(ctx, email, captchaType)
	if err != nil {
		return "", err
	}
	if exists {
		return "", tools.ErrBadRequest("验证码已发送，请稍后再试") // 返回AppError，状态码400
	}

	code := s.GenerateCode() // 生成验证码

	if err := s.SetCaptcha(ctx, email, code, captchaType); err != nil {
		return "", err
	}

	if err := s.email.SendCaptchaEmail(email, code); err != nil {
		// 邮件发送失败，清理已存储的验证码，避免数据不一致
		key := CaptchaPrefix + string(captchaType) + ":" + email
		_ = s.redisClient.Del(ctx, key) //nolint:errcheck // 清理失败不影响错误返回
		return "", tools.ErrInternalServer("验证码发送失败")
	}

	return code, nil
}
