// Package vo 认证相关VO
package vo

// TokenVO Token响应值对象
type TokenVO struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"` // 过期时间戳
}

