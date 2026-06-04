package controllers

import(
	"net/http"
	"github.com/gin-gonic/gin"
	"tugas-gin/database"
	_ "github.com/lib/pq"
)

type Bioskop struct{
	ID int `json:"id"`
	Nama string `json:"nama" binding:"required"`
	Lokasi string `json:"lokasi" binding:"required"`
	Rating float32 `json:"rating"`
}

var BioskopDatas = []Bioskop{}

func CreateBioskop(ctx *gin.Context){
	var newBioskop Bioskop

	if err := ctx.ShouldBindJSON(&newBioskop); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal memproses data. 'nama' dan 'lokasi' wajib diisi dan tidak boleh kosong.",
		})
		return
	}

	query := `INSERT INTO bioskop (nama, lokasi, rating) VALUES ($1, $2, $3) RETURNING id`
	
	var lastInsertId int
	err := database.DB.QueryRow(query, newBioskop.Nama, newBioskop.Lokasi, newBioskop.Rating).Scan(&lastInsertId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data ke database: " + err.Error()})
		return
	}
	newBioskop.ID = lastInsertId

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Bioskop berhasil ditambahkan!",
		"data":    newBioskop,
	})
}
