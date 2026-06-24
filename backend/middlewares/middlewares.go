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

func devCORSOrigins() []string {
	raw := strings.TrimSpace(os.Getenv("AUTH0_CORS_ORIGINS"))
	if raw == "" {
		return []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
		}
	}
	var out []string
	for _, p := range strings.Split(raw, ",") {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func Stack(environment string) []fiber.Handler {
	corsHeaders := "Origin,Content-Type,Accept,Authorization,X-Request-ID,Cookie"
	var corsCfg cors.Config

	if environment == initializers.EnvDevelopment {
		corsCfg = cors.Config{
			AllowOrigins:     strings.Join(devCORSOrigins(), ","),
			AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
			AllowHeaders:     corsHeaders,
			ExposeHeaders:    "X-Request-ID",
			AllowCredentials: true,
		}
	} else {
		origins := strings.TrimSpace(os.Getenv("CORS_ALLOW_ORIGINS"))
		corsCfg = cors.Config{
			AllowOrigins:     origins,
			AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
			AllowHeaders:     corsHeaders,
			ExposeHeaders:    "X-Request-ID",
			AllowCredentials: origins != "" && origins != "*",
		}
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
