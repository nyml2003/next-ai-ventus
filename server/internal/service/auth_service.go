package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthService 认证服务
type AuthService struct {
	jwtSecret string
}

// NewAuthService 创建认证服务
func NewAuthService(secret string) *AuthService {
	return &AuthService{jwtSecret: secret}
}

// Claims JWT 声明
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func (s *AuthService) GenerateToken(username string) (string, error) {
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ValidateToken 验证 JWT Token
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// ValidateCredentials 验证用户名密码（MVP 硬编码）
func (s *AuthService) ValidateCredentials(username, password string) bool {
	// MVP 版本使用硬编码账号
	return username == "admin" && password == "admin"
}

// GetSecret 获取 JWT Secret
func (s *AuthService) GetSecret() string {
	return s.jwtSecret
}
