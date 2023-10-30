package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func backgroundTask() {
	for {
		time.Sleep(15 * time.Minute)
		getRecentsCommits()
	}
}

func main() {
	// Load dotenv
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found!")
	}
	discordToken = os.Getenv("DISCORD_TOKEN")
	discordStatusChannel = os.Getenv("DISCORD_STATUS_CHANNEL")
	discordUpdateChannel = os.Getenv("DISCORD_UPDATES_CHANNEL")

	// get recent commits
	getRecentsCommits()

	// start Discord bot
	go startDiscordBot()

	// start background task
	go backgroundTask()

	// initialize fiber app
	app := fiber.New()
	app.Use(cors.New())

	// root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("online")
	})

	// status endpoint
	app.Get("/status", func(c *fiber.Ctx) error {
		return c.JSON(currentStatus)
	})

	// latest updates endpoint
	app.Get("/updates", func(c *fiber.Ctx) error {
		return c.JSON(currentUpdate)
	})

	// recent commits endpoint
	app.Get("/commits", func(c *fiber.Ctx) error {
		return c.JSON(recentCommits)
	})

	// start fiber app
	app.Listen(":3000")
}
