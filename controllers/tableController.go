package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tabrizihamid84/restaurant-management/database"
	"github.com/tabrizihamid84/restaurant-management/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil || limit < 1 {
			limit = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * limit
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.M{"$match": bson.M{}}
		groupStage := bson.M{"$group ": bson.M{"_id": "null", "total_count": bson.M{"$sum": 1}, "data": bson.M{"$push": "$$ROOT"}}}
		projectStage := bson.M{"$project": bson.M{"_id": 0, "total_count": 1, "order_items": bson.M{"$slice": []interface{}{"$data", startIndex, limit}}}}

		pipeline := bson.A{
			matchStage, groupStage, projectStage,
		}

		cur, err := orderCollection.Aggregate(ctx, pipeline)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing order items"})
			return
		}

		var orders []models.Order
		if err := cur.All(ctx, &orders); err != nil {
			log.Fatal(err)
			return
		}

		c.JSON(http.StatusOK, orders[0])
	}
}
func GetTable() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
func UpdateTable() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
