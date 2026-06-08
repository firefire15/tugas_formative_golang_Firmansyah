package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	DB *sql.DB
}


type CategoryInput struct {
	Name string `json:"name" binding:"required"`
}

func (cc *CategoryController) GetCategories(c *gin.Context) {
	query := "SELECT id, name, created_at, created_by FROM categories"
	rows, err := cc.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data kategori"})
		return
	}
	defer rows.Close()

	var categories []gin.H
	for rows.Next() {
		var id int
		var name string
		var createdAt, createdBy sql.NullString // Menggunakan NullString jika kolom bisa kosong

		if err := rows.Scan(&id, &name, &createdAt, &createdBy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memindai data"})
			return
		}

		categories = append(categories, gin.H{
			"id":         id,
			"name":       name,
			"created_at": createdAt.String,
			"created_by": createdBy.String,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func (cc *CategoryController) CreateCategory(c *gin.Context) {
	var input CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	currentUser, _ := c.Get("username")
	usernameStr, _ := currentUser.(string)

	query := "INSERT INTO categories (name, created_by) VALUES ($1, $2) RETURNING id"
	var newID int
	err := cc.DB.QueryRow(query, input.Name, usernameStr).Scan(&newID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Kategori berhasil ditambahkan!",
		"data": gin.H{
			"id":   newID,
			"name": input.Name,
		},
	})
}

func (cc *CategoryController) GetCategoryByID(c *gin.Context) {
	id := c.Param("id")

	var catID int
	var name string
	query := "SELECT id, name FROM categories WHERE id = $1"
	err := cc.DB.QueryRow(query, id).Scan(&catID, &name)
	
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data kategori tidak ditemukan atau tidak tersedia"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil detail kategori"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":   catID,
			"name": name,
		},
	})
}

func (cc *CategoryController) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)"
	err := cc.DB.QueryRow(checkQuery, id).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal melakukan validasi data"})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gagal menghapus! Data kategori tidak tersedia"})
		return
	}

	deleteQuery := "DELETE FROM categories WHERE id = $1"
	_, err = cc.DB.Exec(deleteQuery, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus kategori karena masih terikat dengan data buku"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil dihapus!"})
}

func (cc *CategoryController) GetBooksByCategory(c *gin.Context) {
	id := c.Param("id")

	var exists bool
	_ = cc.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)", id).Scan(&exists)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori yang Anda maksud tidak tersedia"})
		return
	}

	query := "SELECT id, title, description, price, total_page FROM books WHERE category_id = $1"
	rows, err := cc.DB.Query(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat data buku"})
		return
	}
	defer rows.Close()

	var books []gin.H
	for rows.Next() {
		var bID, price, totalPage int
		var title, description string

		if err := rows.Scan(&bID, &title, &description, &price, &totalPage); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memindai data buku"})
			return
		}

		books = append(books, gin.H{
			"id":         bID,
			"title":      title,
			"description": description,
			"price":      price,
			"total_page": totalPage,
		})
	}

	if books == nil {
		books = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{"data": books})
}