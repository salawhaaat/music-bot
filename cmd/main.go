package main

import (
	"log"
	"telegram-bot/internal/bot"
	"telegram-bot/internal/config"
	"telegram-bot/internal/database/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg := config.New()
	db, err := gorm.Open(mysql.Open(cfg.Dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к БД %v", err)
	}

	bot := bot.MusicBot{
		Bot:    bot.InitBot(cfg.Token),
		User:   &models.UserModel{Db: db},
		Config: cfg,
	}
	bot.Bot.Handle("/start", bot.StartHandler)
	bot.Bot.Handle("/Getspotifylink", bot.GetSpotifyLinkHandler)
	bot.Bot.Handle("/getyt", bot.GetYoutubeLinkHandler)
	bot.Bot.Start()

}
