package discord

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ozraru/ContainerStarterBot/config"
	"github.com/ozraru/ContainerStarterBot/docker"
)

var zero float64 = 0

var optionTimestamps = "timestamps"
var optionTail = "tail"
var optionShutdownTimerOverride = "timeroverride"

func RegisterCommand() {
	addCommand(&discordgo.ApplicationCommand{
		GuildID:     config.Config.Guild,
		Name:        "start",
		Description: "Start container related to this channel",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        optionShutdownTimerOverride,
				Description: "Override shutdown timer. Suffixes: h, m, s. Put \"off\" to disable timer for this activation.",
				Required:    false,
			},
		},
	})
	addCommand(&discordgo.ApplicationCommand{
		GuildID:     config.Config.Guild,
		Name:        "log",
		Description: "Get log of container related to this channel",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        optionTail,
				Description: "Only return this number of log lines from the end of the logs. Defaults to max",
				Required:    false,
				MinValue:    &zero,
				MaxValue:    float64(config.Config.MaxTail),
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        optionTimestamps,
				Description: "Add timestamps to every log line. Defaults to false",
				Required:    false,
			},
		},
	})
	addCommand(&discordgo.ApplicationCommand{
		GuildID:     config.Config.Guild,
		Name:        "status",
		Description: "Check status of container related to this channel",
	})
	if config.Config.EnableStop {
		addCommand(&discordgo.ApplicationCommand{
			GuildID:     config.Config.Guild,
			Name:        "stop",
			Description: "Stop container related to this channel",
		})
	}
}

func addCommand(command *discordgo.ApplicationCommand) {
	_, err := Session.ApplicationCommandCreate(config.Config.AppId, config.Config.Guild, command)
	if err != nil {
		log.Print("Failed to make ", command.Name, " command: ", err)
	}
}

func SlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.GuildID != config.Config.Guild {
		return
	}
	if i.ChannelID != config.Config.Channel {
		return
	}
	if i.Type == discordgo.InteractionApplicationCommand {
		switch i.ApplicationCommandData().Name {
		case "start":
			start(s, i)
		case "log":
			getLog(s, i)
		case "status":
			getStatus(s, i)
		case "stop":
			stop(s, i)
		default:
			log.Print("Unexpected command: ", i.ApplicationCommandData().Name)
		}
	} else {
		log.Print("Unexpected intaraction type: ", i.Type)
	}
}

func start(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ack(s, i.Interaction)
	data := i.ApplicationCommandData()
	err := docker.StartContainer()
	shutdownTimer := ""
	if err != nil {
		respond(s, i.Interaction, "Failed to start container: "+err.Error())
		return
	}
	if data.Options != nil || config.Config.EnableAutoShutdown {
		for _, v := range data.Options {
			switch v.Name {
			case optionShutdownTimerOverride:
				shutdownTimer = v.StringValue()
			default:
				respond(s, i.Interaction, "Unknown options specified")
				return
			}
		}
		go setShutdownTimer(s, shutdownTimer)
	}
	respond(s, i.Interaction, "Container start successful")
}

func setShutdownTimer(s *discordgo.Session, shutdownTimer string) {
	if shutdownTimer == "" {
		shutdownTimer = config.Config.AutoShutdownDuration
	}
	if shutdownTimer == "off" {
		sendMessage(s, "Timer disabled by override")
		return
	}
	timerDuration, err := time.ParseDuration(shutdownTimer)
	if err != nil {
		sendMessage(s, "**Timer is not set; timer duration was invalid, check the config file**")
		return
	}
	go autoShutdown(s, timerDuration)
	sendMessage(s, "Shutdown timer set: "+shutdownTimer)

}

func autoShutdown(s *discordgo.Session, duration time.Duration) {
	time.Sleep(duration)
	err := docker.StopContainer()

	if err != nil {
		sendMessage(s, "Failed to stop container: "+err.Error())
		return
	}

	sendMessage(s, "Container stop successful")
}

func getLog(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ack(s, i.Interaction)

	data := i.ApplicationCommandData()

	timestamps := false
	tail := config.Config.MaxTail
	for _, v := range data.Options {
		switch v.Name {
		case optionTimestamps:
			timestamps = v.BoolValue()
		case optionTail:
			tail = v.IntValue()
		default:
			respond(s, i.Interaction, "Unknown options specifyed")
			return
		}
	}

	reader, err := docker.GetLog(timestamps, tail)
	if err != nil {
		respond(s, i.Interaction, "Failed to start container: "+err.Error())
		return
	}
	defer reader.Close()

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Files: []*discordgo.File{
			{
				Name:        "log.txt",
				ContentType: "text/plain",
				Reader:      reader,
			},
		},
	})
	if err != nil {
		log.Print("Failed to interaction respond log: ", err)
	}
}

func getStatus(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ack(s, i.Interaction)

	data, err := docker.ContainerStatus()
	if err != nil {
		respond(s, i.Interaction, "Failed to check status of container: "+err.Error())
		return
	}

	respond(s, i.Interaction, data)
}

func stop(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ack(s, i.Interaction)

	if !config.Config.EnableStop {
		respond(s, i.Interaction, "Stop command is disabled for this container")
		return
	}

	err := docker.StopContainer()

	if err != nil {
		respond(s, i.Interaction, "Failed to stop container: "+err.Error())
		return
	}

	respond(s, i.Interaction, "Container stop successful")
}

func ack(s *discordgo.Session, i *discordgo.Interaction) {
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		log.Print("Failed to interaction ack: ", err)
	}
}

func respond(s *discordgo.Session, i *discordgo.Interaction, msg string) {
	_, err := s.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Content: &msg,
	})
	if err != nil {
		log.Print("Failed to interaction respond: ", err)
	}
}

func sendMessage(s *discordgo.Session, msg string) {
	_, err := s.ChannelMessageSend(config.Config.Channel, msg)
	if err != nil {
		log.Print("Failed to send message: ", err)
	}
}
