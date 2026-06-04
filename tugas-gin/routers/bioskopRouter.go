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
	
	return router
}