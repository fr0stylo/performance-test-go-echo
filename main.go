package main

import (
	"context"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	uri := os.Getenv("MONGO_URI")
	PORT := "8080"
	if os.Getenv("PORT") != "" {
		PORT = os.Getenv("PORT")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	e := echo.New()

	e.GET("/hello", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"hello": "world"})
	})

	e.GET("/one-item", func(c echo.Context) error {
		var result map[string]interface{}

		if err := client.Database("sample_airbnb").Collection("listingsAndReviews").FindOne(c.Request().Context(), bson.M{}).Decode(&result); err != nil {
			return c.String(http.StatusInternalServerError, "")
		}

		return c.JSON(http.StatusOK, result)
	})

	e.GET("/fifty-items", func(c echo.Context) error {
		var result []map[string]interface{}

		cursor, err := client.Database("sample_airbnb").Collection("listingsAndReviews").Find(c.Request().Context(), bson.M{}, options.Find().SetLimit(50))
		if err != nil {
			return c.String(http.StatusInternalServerError, "")
		}

		if err = cursor.All(c.Request().Context(), &result); err != nil {
			return c.String(http.StatusInternalServerError, "")
		}

		return c.JSON(http.StatusOK, result)
	})

	e.GET("/fibonacci", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]int{"fib": fibonacci(10)})
	})

	e.Logger.Fatal(e.Start(":" + PORT))
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}
