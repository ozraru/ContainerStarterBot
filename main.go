package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ozraru/ContainerStarterBot/config"
	"github.com/ozraru/ContainerStarterBot/discord"
)

func main() {
	config.LoadConfig()
	discord.ConnectDiscord()
	defer discord.CloseDiscord()
	log.Print("discord connected")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-stop
	log.Print("recieved Interrupt signal")
}
