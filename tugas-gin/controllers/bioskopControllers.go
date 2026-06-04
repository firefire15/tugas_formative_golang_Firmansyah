package controllers

import(
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	"tugas-gin/database"
	"database/sql"
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
		"data": newBioskop,
	})
}

func GetBioskopByID(ctx *gin.Context){
	bioskopID := ctx.Param("id")
	query := "SELECT id, nama, lokasi, rating FROM bioskop where id=$1"

	var b Bioskop
	err := database.DB.QueryRow(query, bioskopID).Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating)

	if err != nil{
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Bioskop dengan ID " + bioskopID + " tidak ditemukan",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data: " + err.Error()})
		return	
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message":"Berhasil mengambil detail bioskop",
		"data":b,
	})
}

func GetBioskop(ctx *gin.Context){
	query := "SELECT id, nama, lokasi, rating FROM bioskop"

	rows, err := database.DB.Query(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data: " + err.Error()})
		return
	}
	defer rows.Close()

	var daftarBioskop []Bioskop
	for rows.Next() {
		var b Bioskop
		err := rows.Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca baris data: " + err.Error()})
			return
		}
		daftarBioskop = append(daftarBioskop, b)
	}
	if daftarBioskop == nil {
		daftarBioskop = []Bioskop{}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil daftar bioskop",
		"data":    daftarBioskop,
	})
}

func UpdateBioskop(ctx *gin.Context){
	idBioskop := ctx.Param("id")

	id, err := strconv.Atoi(idBioskop)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid. Id merupakan integer (bilangan bulat positif)",
		})
		return
	}

	var updatedBioskop Bioskop
	if err := ctx.ShouldBindJSON(&updatedBioskop); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Gagal memproses data. 'nama' dan 'lokasi' wajib diisi dan tidak boleh kosong.",
		})
		return
	}

	query := "UPDATE bioskop SET nama = $1, lokasi = $2, rating = $3 WHERE id = $4 RETURNING id"	
	var checkId int
	err = database.DB.QueryRow(query, updatedBioskop.Nama, updatedBioskop.Lokasi, updatedBioskop.Rating, id).Scan(&checkId)
	
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Bioskop dengan ID " + idBioskop + " tidak ditemukan.",
			})
			return
		}
		
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal memperbarui data ke database: " + err.Error(),
		})
		return
	}

	updatedBioskop.ID = id

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Bioskop berhasil diperbarui!",
		"data":    updatedBioskop,
	})
}

func DeleteBioskop(ctx *gin.Context) {
	idBioskop := ctx.Param("id")

	id, err := strconv.Atoi(idBioskop)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID tidak valid. Parameter ID harus berupa angka bulat.",
		})
		return
	}

	query := "DELETE FROM bioskop WHERE id = $1 RETURNING id"
	
	var deletedId int
	err = database.DB.QueryRow(query, id).Scan(&deletedId)
	
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Bioskop dengan ID " + idBioskop + " tidak ditemukan.",
			})
			return
		}
		
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus data dari database: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Bioskop dengan ID " + idBioskop + " berhasil dihapus!",
	})
}





