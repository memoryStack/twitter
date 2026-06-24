package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"

	"twitter/backend/auth"
	"twitter/backend/controllers"
	"twitter/backend/initializers"
	"twitter/backend/middlewares"
)

var runEnv string

func init() {
	env := flag.String("env", "", "environment to run (development or production)")
	flag.Parse()

	runEnv = *env
	if runEnv != initializers.EnvDevelopment && runEnv != initializers.EnvProduction {
		log.Fatal("flag -env is required and must be development or production")
	}

	initializers.LoadEnv(runEnv)

	if err := auth.Init(context.Background()); err != nil {
		log.Fatalf("auth: %v", err)
	}

	if runEnv == initializers.EnvProduction {
		if strings.TrimSpace(os.Getenv("CORS_ALLOW_ORIGINS")) == "" {
			log.Fatal("CORS_ALLOW_ORIGINS is required in production")
		}
	}

	initializers.ConnectDB(runEnv)
	initializers.SyncDB()
	initializers.InitRepositories()
}

func main() {
	app := fiber.New(fiber.Config{
		AppName:      "Twitter API",
		ServerHeader: "Fiber",
	})

	for _, handler := range middlewares.Stack(runEnv) {
		app.Use(handler)
	}

	app.Get("/health", controllers.Health)

	app.Get("/api/auth/login", controllers.AuthLogin)
	app.Get("/api/auth/callback", controllers.AuthCallback)
	app.Post("/api/auth/refresh", controllers.AuthRefresh)
	app.Post("/api/auth/logout", controllers.AuthLogout)
	app.Get("/api/auth/me", middlewares.RequireAuth, controllers.AuthMe)

	tweets := app.Group("/api/tweets", middlewares.RequireAuth)
	tweets.Post("/", controllers.CreateTweet)
	tweets.Get("/", controllers.GetMyTweets)
	tweets.Get("/:id", controllers.GetTweetByID)
	tweets.Patch("/:id", controllers.UpdateMyTweet)
	tweets.Delete("/:id", controllers.DeleteMyTweet)
	tweets.Post("/:id/like", controllers.LikeTweet)

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		if runEnv == initializers.EnvDevelopment {
			addr = ":3000"
		} else {
			addr = ":8080"
		}
	}

	log.Printf("starting server in %s mode on %s", runEnv, addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
