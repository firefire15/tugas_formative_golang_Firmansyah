package routers

import(
	"tugas-gin/controllers"
	"tugas-gin/database"
	"github.com/gin-gonic/gin"
)

func StartBioskopServer() *gin.Engine{
	database.ConnectDB()

	router := gin.Default()

	router.POST("/bioskop", controllers.CreateBioskop)
	router.GET("/bioskop", controllers.GetBioskop) 
	router.GET("/bioskop/:id", controllers.GetBioskopByID)
	router.PUT("/bioskop/:id", controllers.UpdateBioskop)
	router.DELETE("/bioskop/:id", controllers.DeleteBioskop)
	
	return router
}