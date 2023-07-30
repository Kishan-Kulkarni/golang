package main

import (
	"log"
	"os"

	"github.com/Kamva/mgm/v2"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"

	// "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	mgm.DefaultModel `bson:",inline"`
	UserName         string `json:"username" bson:"username"`
	Password         string `json:"password" bson:"password"`
}

type Post struct {
	mgm.DefaultModel `bson:",inline"`
	Id               string `json:"id" bson:"id"`
	Title            string `json:"title" bson:"title"`
	Summary          string `json:"summary" bson:"summary"`
	Image            string `json:"image" bson:"image"`
	Content          string `json:"content" bson:"content"`
	Author           string `json:"_id" bson:"_id"`
}

type Data struct {
	Users []User `json:"users" bson:"users"`
	Posts []Post `json:"posts" bson:"posts"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	uri := os.Getenv("DATABASE")
	err = mgm.SetDefaultConfig(nil, "test", options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	userCollection := mgm.Coll(&User{})
	postCollection := mgm.Coll(&Post{})

	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Test App v1.0.1",
	})
	app.Get("/", func(c *fiber.Ctx) error {
		users := []User{}
		posts := []Post{}
		err := userCollection.SimpleFind(&users, bson.D{})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"ok":    false,
				"error": err.Error(),
			})
		}
		err = postCollection.SimpleFind(&posts, bson.D{})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"ok":    false,
				"error": err.Error(),
			})
		}
		data := Data{Users: users, Posts: posts}
		jsonData, _ := json.Marshal(data)
		return c.SendString(string(jsonData))
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	app.Listen(":" + port)
}
