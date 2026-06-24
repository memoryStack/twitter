package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"

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

	if runEnv == initializers.EnvProduction {
		if strings.TrimSpace(os.Getenv("CORS_ALLOW_ORIGINS")) == "" {
			log.Fatal("CORS_ALLOW_ORIGINS is required in production")
		}
	}

	initializers.ConnectDB(runEnv)
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
