package main

import (
	"context"
	"log"
	"os"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct{
	UserName string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type Post struct{
	Id string `json:"id" bson:"id"`
	Title string `json:"title" bson:"title"`
	Summary string `json:"summary" bson:"summary"`
	Image string `json:"image" bson:"image"`
	Content string `json:"content" bson:"content"`
	Author string `json:"_id" bson:"_id"`
	
}

type Data struct{
	Users []User `json:"users" bson:"users"`
	Posts []Post `json:"posts" bson:"posts"`
}

func main() {
	err := godotenv.Load(".env")
  if err != nil {
    log.Fatalf("Error loading .env file")
  }
  uri := os.Getenv("DATABASE")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
  opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
  client, err := mongo.Connect(context.TODO(), opts)
  if err != nil {
    panic(err)
  }
  defer func() {
    if err = client.Disconnect(context.TODO()); err != nil {
      panic(err)
    }
  }()

  userCollection :=client.Database("test").Collection("users")
  postCollection :=client.Database("test").Collection("posts")

  app := fiber.New(fiber.Config{
    Prefork:       true,
    CaseSensitive: true,
    StrictRouting: true,
    ServerHeader:  "Fiber",
    AppName: "Test App v1.0.1",
})
	app.Get("/", func(c *fiber.Ctx) error {
		cursoruser, err := userCollection.Find(context.TODO(), bson.D{})
		if err != nil {
			panic(err)
		}
		cursorpost, err := postCollection.Find(context.TODO(), bson.D{})
		if err != nil {
			panic(err)
		}
		var users []User
		var posts []Post
		cPost := make(chan []Post)
		cUser := make(chan []User)
		go func() {
			var tempPosts []Post
			if err = cursorpost.All(context.TODO(), &tempPosts); err != nil {
				panic(err)
			}
			cPost <- tempPosts
			close(cPost)
		}()
		go func ()  {
			var tempUsers []User
			if err = cursoruser.All(context.TODO(), &tempUsers); err != nil {
				panic(err)
			}
			cUser <- tempUsers
			close(cUser)
		}()
		posts = <-cPost
		users = <-cUser
		data := Data{Users: users , Posts: posts}
		jsonData, _ := json.MarshalIndent(data,"", "  ")
        return c.SendString(string(jsonData))
    })
	port:=os.Getenv("PORT")
	if port == "" {
		port="3000"
	}
	app.Listen(":" + port)
}