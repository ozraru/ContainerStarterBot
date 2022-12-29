package discord

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/ozraru/ContainerStarterBot/config"
)

var Session *discordgo.Session

func ConnectDiscord() {
	var err error
	Session, err = discordgo.New("Bot " + config.Config.Token)
	if err != nil {
		log.Panic("Failed to create discord session: ", err)
	}
	Session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = Session.Open()
	if err != nil {
		log.Panic("Failed to open discord session: ", err)
	}
	RegisterCommand()
	Session.AddHandler(SlashCommand)
}

func CloseDiscord() {
	Session.Close()
}
