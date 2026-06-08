package routers

import(
	"tugas-gin/controllers"
	"tugas-gin/db"
	"tugas-gin/helper"
	"github.com/gin-gonic/gin"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func StartAPIServer() *gin.Engine {
	// 1. Jalankan koneksi DB (Jika Anda menggunakan Opsi A yang menggabungkan fungsi migrasi,
	// database akan otomatis termigrasi saat fungsi ini dipanggil).
	database := db.ConnectDB()
	
	// Opsional: Jika Anda memisahkan fungsinya (Opsi B), jalankan baris ini:
	// db.RunMigrations(database) 
	
	_ = database // Menghindari error "unused variable" jika objek DB diatur secara global di package db

	router := gin.Default()

	userCtrl := &controllers.UserController{DB: database}
	categoryCtrl := &controllers.CategoryController{DB: database}
	bookCtrl := &controllers.BookController{DB: database}

	
	router.POST("api/users/register", userCtrl.Register) 
	router.POST("api/users/login", userCtrl.Login)


	// Public API
	router.POST("/bioskop", controllers.CreateBioskop)
	router.GET("/bioskop", controllers.GetBioskop) 
	router.GET("/bioskop/:id", controllers.GetBioskopByID)
	router.PUT("/bioskop/:id", controllers.UpdateBioskop)
	router.DELETE("/bioskop/:id", controllers.DeleteBioskop)
	

	// Protected API
	protected := router.Group("/api")
	protected.Use(helper.JWTMiddleware())
	{

		protected.GET("/categories", categoryCtrl.GetCategories)
		protected.POST("/categories", categoryCtrl.CreateCategory)
		protected.GET("/categories/:id", categoryCtrl.GetCategoryByID)
		protected.DELETE("/categories/:id", categoryCtrl.DeleteCategory)
		protected.GET("/categories/:id/books", categoryCtrl.GetBooksByCategory)

		protected.GET("/books", bookCtrl.GetBooks)
		protected.POST("/books", bookCtrl.CreateBook)
		protected.GET("/books/:id", bookCtrl.GetBookByID)
		protected.PUT("/books/:id", bookCtrl.UpdateBook)
		protected.DELETE("/books/:id", bookCtrl.DeleteBook)

	}

	return router
}