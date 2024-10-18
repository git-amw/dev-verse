package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(ctx *gin.Context) {
	atuhHeader := ctx.GetHeader("Authorization")
	if atuhHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Auth Header Missing"})
		ctx.Abort()
		return
	}
	tokenString := strings.Split(atuhHeader, " ")
	if len(tokenString) != 2 || tokenString[0] != "Bearer" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
		ctx.Abort()
		return
	}
	token, err := jwt.Parse(tokenString[1], func(token *jwt.Token) (interface{}, error) {
		return []byte("SECRET_KEY"), nil
	})
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		ctx.Abort()
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId, ok := claims["userid"]
		if !ok {
			log.Println("user id not found - auth middelware")
			ctx.Abort()
			return
		}
		ctx.Set("userid", userId)
	}
	ctx.Next()
}
