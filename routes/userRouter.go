package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tabrizihamid84/restaurant-management/controllers"
)

func UserRoutes(incommingRoutes *gin.Engine) {
	incommingRoutes.GET("/users", controllers.GetUsers())
	incommingRoutes.GET("/users/:user_id", controllers.GetUser())
	incommingRoutes.POST("/users/signup", controllers.SignUp())
	incommingRoutes.POST("/users/signin", controllers.SignIn())
}
