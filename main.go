package main

import (
	"log/slog"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	slogecho "github.com/samber/slog-echo"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		logger.Error("TELEGRAM_BOT_TOKEN is not set")
		os.Exit(1)
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		logger.Error("Failed to create Telegram bot", "error", err)
		os.Exit(1)
	}

	e := echo.New()
	e.Use(slogecho.New(logger))
	e.Use(middleware.Recover())

	e.GET("/v1/telegram/file/:file_id", func(c echo.Context) error {
		return getFile(c, bot, logger)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	e.Logger.Fatal(e.Start(":" + port))
}

func getFile(c echo.Context, bot *tgbotapi.BotAPI, logger *slog.Logger) error {
	fileID := c.Param("file_id")
	if fileID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "The file ID is missing."})
	}

	fileURLWithBotToken, err := bot.GetFileDirectURL(fileID)
	if err != nil {
		logger.Error("Failed to get direct file URL from telegram", "error", err, "file_id", fileID)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	resp, err := http.Get(fileURLWithBotToken)
	if err != nil {
		logger.Error("Failed to make a request to the direct file URL from telegram")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	defer resp.Body.Close()

	c.Response().Header().Set(echo.HeaderContentType, resp.Header.Get(echo.HeaderContentType))
	return c.Stream(resp.StatusCode, resp.Header.Get(echo.HeaderContentType), resp.Body)
}
