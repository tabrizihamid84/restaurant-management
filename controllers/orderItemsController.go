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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderItemPack struct {
	Table_id    *string
	Order_items []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

func GetOrderItems() gin.HandlerFunc {
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

		var orderItems []models.OrderItem
		if err := cur.All(ctx, &orderItems); err != nil {
			log.Fatal(err)
			return
		}

		c.JSON(http.StatusOK, orderItems[0])
	}
}

func GetOrderItemsByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("order_id")

		allOrderItems, err := ItemsByOrder(orderId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing order item by Order ID"})
			return
		}

		c.JSON(http.StatusOK, allOrderItems)
	}
}

func ItemsByOrder(id string) (OrdersItems []primitive.M, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	matchStage := bson.M{"$match": bson.M{"_order_id": id}}
	lookupStage := bson.M{"$lookup": bson.M{"from": "food", "localField": "food_id", "foreignField": "food_id", "as": "food"}}
	unwindStage := bson.M{"$unwind": bson.M{"path": "$food", "preserveNullAndEmptyArrays": true}}
	lookupOrderStage := bson.M{"$lookup": bson.M{"from": "order", "localField": "order_id", "foreignField": "order_id", "as": "order"}}
	unwindOrderStage := bson.M{"$unwind": bson.M{"path": "$order", "preserveNullAndEmptyArrays": true}}
	lookupTableStage := bson.M{"$lookup": bson.M{"from": "table", "localField": "order.table_id", "foreignField": "table_id", "as": "table"}}
	unwindTableStage := bson.M{"$unwind": bson.M{"path": "$table", "preserveNullAndEmptyArrays": true}}
	projectStage := bson.M{"#project": bson.M{"_id": 0, "amount": "$food.price", "total_count": 1, "food_name": "$food.name", "food_image": "$food.food_image", "table_number": "$table.table_number", "table_id": "$table.table_id", "order_id": "$order.order_id", "price": "$food.price", "quantity": 1}}
	groupStage := bson.M{"$group": bson.M{"_id": bson.M{"order_id": "$order_id", "table_id": "$table_id", "table_number": "$table_number", "payment_due": bson.M{"$sum": "$amount"}, "total_count": bson.M{"$sum": 1}, "order_items": bson.M{"$push": "$$ROOT"}}}}
	projectStage2 := bson.M{"$project": bson.M{"id": 0, "payment_due": 1, "total_count": 1, "table_number": "$_id.table_number", "order_items": 1}}

	pipeline := bson.A{
		matchStage, lookupStage, unwindStage, lookupOrderStage, unwindOrderStage, lookupTableStage, unwindTableStage, unwindTableStage, projectStage, groupStage, projectStage2,
	}

	cur, err := orderItemCollection.Aggregate(ctx, pipeline)
	defer cancel()
	if err != nil {
		panic(err)
	}

	if err := cur.All(ctx, &OrdersItems); err != nil {
		panic(err)
	}

	return OrdersItems, err
}

func GetOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderItemId := c.Param("orcer_item_id")
		var orderItem models.OrderItem

		err := orderItemCollection.FindOne(ctx, bson.M{"orderItem_id": orderItemId}).Decode(&orderItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing ordered item"})
		}
		c.JSON(http.StatusOK, orderItem)
	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var OrderItemPack OrderItemPack
		var order models.Order

		if err := c.BindJSON(&OrderItemPack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		orderItemsToBeInserted := []interface{}{}
		order.Table_id = OrderItemPack.Table_id
		order_id := OrderItemOrderCreator(order)

		for _, orderItem := range OrderItemPack.Order_items {
			orderItem.Order_id = order_id

			err := validate.Struct(orderItem)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			orderItem.ID = primitive.NewObjectID()
			orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Order_item_id = orderItem.ID.Hex()

			var num = toFixed(*orderItem.Unit_price, 2)
			orderItem.Unit_price = &num
			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
		}

		res, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "order item was not created"})
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orderItem models.OrderItem

		orderItemId := c.Param("order_item_id")

		var updateObj primitive.D

		if orderItem.Unit_price != nil {
			updateObj = append(updateObj, bson.E{Key: "unit_price", Value: orderItem.Unit_price})

		}

		if orderItem.Quantity != nil {
			updateObj = append(updateObj, bson.E{Key: "quantity", Value: orderItem.Quantity})

		}

		if orderItem.Food_id != nil {
			updateObj = append(updateObj, bson.E{Key: "food_id", Value: orderItem.Food_id})
		}

		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: orderItem.Updated_at})

		upsert := true
		filter := bson.M{"order_item_id": orderItemId}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		res, err := orderItemCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "order item update failed"})
			return
		}

		c.JSON(http.StatusOK, res)

	}
}
