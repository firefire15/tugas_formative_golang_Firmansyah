package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookController struct {
	DB *sql.DB
}

type BookInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	ImageURL    string `json:"image_url" binding:"required"`
	ReleaseYear int    `json:"release_year" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	TotalPage   int    `json:"total_page" binding:"required"`
	CategoryID  int    `json:"category_id" binding:"required"`
}

func (bc *BookController) GetBooks(c *gin.Context) {
	query := `SELECT id, title, description, image_url, release_year, price, total_page, thickness, category_id FROM books`
	rows, err := bc.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data buku"})
		return
	}
	defer rows.Close()

	var books []gin.H
	for rows.Next() {
		var id, releaseYear, price, totalPage, categoryID int
		var title, description, imageURL, thickness string

		err := rows.Scan(&id, &title, &description, &imageURL, &releaseYear, &price, &totalPage, &thickness, &categoryID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memindai data buku"})
			return
		}

		books = append(books, gin.H{
			"id":           id,
			"title":        title,
			"description":  description,
			"image_url":    imageURL,
			"release_year": releaseYear,
			"price":        price,
			"total_page":   totalPage,
			"thickness":    thickness,
			"category_id":  categoryID,
		})
	}

	if books == nil {
		books = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{"data": books})
}

func (bc *BookController) CreateBook(c *gin.Context) {
	var input BookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.ReleaseYear < 1980 || input.ReleaseYear > 2024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tahun rilis (release_year) harus berada di antara 1980 dan 2024"})
		return
	}

	var thickness string
	if input.TotalPage > 100 {
		thickness = "tebal"
	} else {
		thickness = "tipis"
	}

	// Ambil data creator dari JWT token
	currentUser, _ := c.Get("username")
	usernameStr, _ := currentUser.(string)

	query := `INSERT INTO books (title, description, image_url, release_year, price, total_page, thickness, category_id, created_by) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	
	var newID int
	err := bc.DB.QueryRow(query, input.Title, input.Description, input.ImageURL, input.ReleaseYear, input.Price, input.TotalPage, thickness, input.CategoryID, usernameStr).Scan(&newID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan buku baru, pastikan category_id valid"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Buku berhasil ditambahkan!",
		"data": gin.H{
			"id":           newID,
			"title":        input.Title,
			"thickness":    thickness, // Hasil konversi otomatis
			"release_year": input.ReleaseYear,
		},
	})
}

func (bc *BookController) GetBookByID(c *gin.Context) {
	id := c.Param("id")

	var bID, releaseYear, price, totalPage, categoryID int
	var title, description, imageURL, thickness string

	query := `SELECT id, title, description, image_url, release_year, price, total_page, thickness, category_id FROM books WHERE id = $1`
	err := bc.DB.QueryRow(query, id).Scan(&bID, &title, &description, &imageURL, &releaseYear, &price, &totalPage, &thickness, &categoryID)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data buku tidak ditemukan atau tidak tersedia"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil detail buku"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":           bID,
			"title":        title,
			"description":  description,
			"image_url":    imageURL,
			"release_year": releaseYear,
			"price":        price,
			"total_page":   totalPage,
			"thickness":    thickness,
			"category_id":  categoryID,
		},
	})
}

func (bc *BookController) UpdateBook(c *gin.Context) {
	id := c.Param("id")

	var input BookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exists bool
	_ = bc.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM books WHERE id = $1)", id).Scan(&exists)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gagal mengupdate! Data buku tidak tersedia"})
		return
	}

	if input.ReleaseYear < 1980 || input.ReleaseYear > 2024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tahun rilis (release_year) harus berada di antara 1980 dan 2024"})
		return
	}

	var thickness string
	if input.TotalPage > 100 {
		thickness = "tebal"
	} else {
		thickness = "tipis"
	}

	currentUser, _ := c.Get("username")
	usernameStr, _ := currentUser.(string)

	query := `UPDATE books SET title=$1, description=$2, image_url=$3, release_year=$4, price=$5, total_page=$6, thickness=$7, category_id=$8, modified_by=$9, modified_at=NOW() WHERE id=$10`
	_, err := bc.DB.Exec(query, input.Title, input.Description, input.ImageURL, input.ReleaseYear, input.Price, input.TotalPage, thickness, input.CategoryID, usernameStr, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data buku"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data buku berhasil diperbarui!"})
}

func (bc *BookController) DeleteBook(c *gin.Context) {
	id := c.Param("id")

	var exists bool
	_ = bc.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM books WHERE id = $1)", id).Scan(&exists)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gagal menghapus! Data buku tidak tersedia"})
		return
	}

	_, err := bc.DB.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus buku"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Buku berhasil dihapus!"})
}