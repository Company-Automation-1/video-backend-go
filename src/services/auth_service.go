// Package services 认证服务
package services

import (
	"context"
	"errors"
	"time"

	"github.com/Company-Automation-1/video-backend-go/src/config"
	"github.com/Company-Automation-1/video-backend-go/src/infrastructure"
	"github.com/Company-Automation-1/video-backend-go/src/query"
	"github.com/Company-Automation-1/video-backend-go/src/tools"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService 认证服务
type AuthService struct {
	jwtConfig *config.JWTConfig
	redis     *infrastructure.Redis
}

// NewAuthService 创建认证服务
func NewAuthService(jwtConfig *config.JWTConfig, redis *infrastructure.Redis) *AuthService {
	return &AuthService{
		jwtConfig: jwtConfig,
		redis:     redis,
	}
}

// Claims JWT Claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Login 用户登录
func (s *AuthService) Login(
	ctx context.Context,
	username, password string,
) (accessToken string, expiresIn int64, err error) {
	// 查询用户
	user, err := query.User.Where(query.User.Username.Eq(username)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", 0, tools.ErrBadRequest("用户名或密码错误")
		}
		return "", 0, tools.ErrInternalServer("登录失败")
	}

	// 验证密码
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", 0, tools.ErrBadRequest("用户名或密码错误")
	}

	// 生成Token
	accessToken, expiresIn, err = s.generateToken(user.ID, user.Username)
	if err != nil {
		return "", 0, tools.ErrInternalServer("Token生成失败")
	}

	return accessToken, expiresIn, nil
}

// generateToken 生成Access Token
func (s *AuthService) generateToken(userID uint, username string) (tokenString string, expiresIn int64, err error) {
	now := time.Now()
	expiresIn = now.Add(time.Duration(s.jwtConfig.ExpireTime) * time.Hour).Unix()

	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expiresIn, 0)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresIn, nil
}

// VerifyToken 验证Token并返回Claims
func (s *AuthService) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return []byte(s.jwtConfig.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的Token")
}

// Logout 用户登出（将token加入黑名单）
func (s *AuthService) Logout(ctx context.Context, tokenString string) error {
	claims, err := s.VerifyToken(tokenString)
	if err != nil {
		return nil // token无效，无需处理
	}

	// 将token加入黑名单（存储到Redis，有效期与token过期时间一致）
	blacklistKey := "blacklist_token:" + tokenString
	expiresAt := claims.ExpiresAt.Time
	ttl := time.Until(expiresAt)
	if ttl > 0 {
		//nolint:errcheck // Redis设置失败不影响登出流程，token会在过期时间后自然失效
		_ = s.redis.Set(ctx, blacklistKey, "1", ttl)
	}

	return nil
}

// IsTokenBlacklisted 检查token是否在黑名单中
func (s *AuthService) IsTokenBlacklisted(ctx context.Context, tokenString string) bool {
	blacklistKey := "blacklist_token:" + tokenString
	_, err := s.redis.Get(ctx, blacklistKey)
	return err == nil
}
