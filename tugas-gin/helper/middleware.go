package helper

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil data dari header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format Authorization header diperlukan"})
			c.Abort() // Hentikan request agar tidak lanjut ke Controller
			return
		}

		// Format header biasanya: "Bearer <token_jwt>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format token harus berupa Bearer <token>"})
			c.Abort()
			return
		}

		// Validasi tokennya menggunakan fungsi yang kita buat di jwt.go
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid atau sudah kedaluwarsa"})
			c.Abort()
			return
		}

		// Simpan data username ke dalam context Gin agar bisa dibaca di handler/controller lain jika butuh
		c.Set("username", claims.Username)

		c.Next() // Token aman, lanjutkan request ke API tujuan
	}
}