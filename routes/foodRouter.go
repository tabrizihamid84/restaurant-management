package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tabrizihamid84/restaurant-management/controllers"
)

func FoodRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/foods", controllers.GetFoods())
	incomingRoutes.GET("/foods/:food_id", controllers.GetFood())
	incomingRoutes.POST("/foods", controllers.CreateFood())
	incomingRoutes.PATCH("/foods/:food_id", controllers.UpdateFood())
}
