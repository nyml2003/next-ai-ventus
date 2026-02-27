package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/next-ai-ventus/server/internal/interfaces/http/response"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key" // 默认值，生产环境应使用环境变量
	}

	return func(c *gin.Context) {
		// 从 Cookie 获取 token
		tokenString, err := c.Cookie("token")
		if err != nil {
			// 尝试从 Authorization header 获取
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				response.Error(c, response.CodeTokenMissing)
				c.Abort()
				return
			}

			// Bearer token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Error(c, response.CodeTokenInvalid)
				c.Abort()
				return
			}
			tokenString = parts[1]
		}

		// 解析 token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil {
			if err == jwt.ErrTokenExpired {
				response.Error(c, response.CodeTokenExpired)
			} else {
				response.Error(c, response.CodeTokenInvalid)
			}
			c.Abort()
			return
		}

		if !token.Valid {
			response.Error(c, response.CodeTokenInvalid)
			c.Abort()
			return
		}

		// 提取 claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if username, ok := claims["username"].(string); ok {
				c.Set("username", username)
			}
		}

		c.Next()
	}
}
