package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tabrizihamid84/restaurant-management/routes"
)

// var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()

	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins: []string{"*"},
	// 	// AllowMethods:     []string{"PUT", "PATCH"},
	// 	// AllowHeaders:     []string{"Origin"},
	// 	// ExposeHeaders:    []string{"Content-Length"},
	// 	// AllowCredentials: true,
	// 	// AllowOriginFunc: func(origin string) bool {
	// 	// return origin == "https://github.com"
	// 	// },
	// 	// MaxAge: 12 * time.Hour,
	// }))

	router.Use(gin.Logger())
	routes.UserRoutes(router)
	// router.Use(middleware.Authentication())

	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	// routes.TableRoutes(router)
	// routes.OrderRoutes(router)
	// routes.OrderItemRoutes(router)
	// routes.InvoiceRoutes(router)

	router.Run(":" + port)

	fmt.Println("running")

}
