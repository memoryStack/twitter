package middlewares

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/timeout"

	"twitter/backend/initializers"
)

func Stack(environment string) []fiber.Handler {
	corsCfg := cors.Config{
		AllowMethods:  "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:  "Origin,Content-Type,Accept,Authorization,X-Request-ID",
		ExposeHeaders: "X-Request-ID",
	}

	if environment == initializers.EnvDevelopment {
		corsCfg.AllowOrigins = "*"
		corsCfg.AllowCredentials = false
	} else {
		corsCfg.AllowOrigins = os.Getenv("CORS_ALLOW_ORIGINS")
		origins := strings.TrimSpace(corsCfg.AllowOrigins)
		corsCfg.AllowCredentials = origins != "" && origins != "*"
	}

	return []fiber.Handler{
		recover.New(),
		requestid.New(requestid.Config{
			Header:     "X-Request-ID",
			ContextKey: "requestid",
		}),
		cors.New(corsCfg),
		logger.New(logger.Config{
			Format: "${time} | ${status} | ${latency} | ${ip} | ${method} ${path} ${error}\n",
		}),
		timeout.NewWithContext(func(c *fiber.Ctx) error {
			return c.Next()
		}, 30*time.Second),
	}
}
