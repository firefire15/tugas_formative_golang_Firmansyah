package controllers

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"tugas-gin/helper"
)

type UserController struct {
	DB *sql.DB
}

type AuthInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (uc *UserController) Register(c *gin.Context) {
	var input AuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses password"})
		return
	}

	query := "INSERT INTO users (username, password) VALUES ($1, $2)"
	_, err = uc.DB.Exec(query, input.Username, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username sudah terdaftar"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registrasi berhasil! Silakan login."})
}

func (uc *UserController) Login(c *gin.Context) {
	var input AuthInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dbPassword string
	query := "SELECT password FROM users WHERE username = $1"
	err := uc.DB.QueryRow(query, input.Username).Scan(&dbPassword)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses data"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(input.Password))
	if err != nil {
		// Jika password tidak cocok
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	token, err := helper.GenerateJWT(input.Username) 
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token akses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil!",
		"token":   token,
	})
}