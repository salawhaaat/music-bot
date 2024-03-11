package bot

import (
	"log"
	"os"
	"path/filepath"
	"telegram-bot/internal/database/models"
	"telegram-bot/internal/utils"
	"time"

	"telegram-bot/internal/config"

	"gopkg.in/telebot.v3"
)

type MusicBot struct {
	Bot    *telebot.Bot
	User   *models.UserModel
	Config *config.Config
}

func (bot *MusicBot) StartHandler(ctx telebot.Context) error {
	newUser := models.User{
		Name:       ctx.Sender().Username,
		TelegramId: ctx.Chat().ID,
		FirstName:  ctx.Sender().FirstName,
		LastName:   ctx.Sender().LastName,
		ChatId:     ctx.Chat().ID,
	}
	msg := "\nSend me spotify link to the album or playlist\n/Getspotifylink link"

	existUser, err := bot.User.FindOne(ctx.Chat().ID)

	if err != nil {
		log.Printf("Ошибка получения пользователя %v", err)
	}

	if existUser == nil {
		err := bot.User.Create(newUser)

		if err != nil {
			log.Printf("Ошибка создания пользователя %v", err)
		}
	}
	return ctx.Send("Привет, " + ctx.Sender().FirstName + msg)

}

func InitBot(token string) *telebot.Bot {
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)

	if err != nil {
		log.Fatalf("Ошибка при инициализации бота %v", err)
	}

	return b
}
func (bot *MusicBot) GetYoutubeLinkHandler(ctx telebot.Context) error {
	var musicpath string
	youtubeLink := utils.ExtractYouTubeLink(ctx.Text())
	if youtubeLink == "" {
		log.Printf("Not Valid Link")
		return nil
		// msg = "Link is not valid. Example: https://www.youtube.com/playlist?list=example123idg456"
	} else {
		log.Printf("Received Youtube link: %s", youtubeLink)
		downloadComplete := make(chan string)
		go func() {
			err := utils.DownloadMusic(youtubeLink, downloadComplete)
			if err != nil {
				log.Printf("Error downloading music: %v", err)
				downloadComplete <- "" // Signal completion with an empty string
			}
		}()
		musicpath = <-downloadComplete
	}
	if musicpath == "" {
		_, err := bot.Bot.Send(ctx.Chat(), "Failed to download music. Please try again.")
		if err != nil {
			log.Printf("Error sending reply: %v", err)
		}
		return nil
	}

	musicFile, err := os.Open(musicpath)
	if err != nil {
		log.Printf("Error opening music file: %v", err)
		return nil
	}
	defer musicFile.Close()

	audio := &telebot.Audio{
		File:     telebot.FromReader(musicFile),
		FileName: filepath.Base(musicpath),
	}

	_, err = bot.Bot.Send(ctx.Chat(), audio)
	if err != nil {
		log.Printf("Error sending audio message: %v", err)
	}

	return nil
}

func (bot *MusicBot) GetSpotifyLinkHandler(ctx telebot.Context) error {
	msg := "Thanks! I received your Spotify link: "
	spotifyLink := utils.ExtractPlaylistID(ctx.Text())
	if spotifyLink == "" {
		log.Printf("Not Valid Link")
		msg = "Link is not valid. Example: https://open.spotify.com/playlist/example123idg456"
	} else {
		log.Printf("Received Spotify link: %s", spotifyLink)
	}
	_, err := bot.Bot.Send(ctx.Chat(), msg+spotifyLink)
	if err != nil {
		log.Printf("Error sending reply: %v", err)

	}
	return nil
}
