package main

import (
	"flag"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
)

type Bot struct {
	Config         *Configuration
	Permissions    *PermissionsManager
	Commands       Commands
	Player         *Player
	DiscordSession *discordgo.Session
}

var (
	Tanuki Bot
)

func init() {
	Tanuki.Config = &Configuration{}

	//configPath := flag.String("c", "config.yml", "Config file path") // allow using -c flag, flags will override config file
	Tanuki.Config.Load("config.yml")

	flag.StringVar(&Tanuki.Config.Token, "a", Tanuki.Config.Token, "Auth token")
	flag.StringVar(&Tanuki.Config.Guild, "g", Tanuki.Config.Guild, "Guild ID")
	flag.StringVar(&Tanuki.Config.TextChannel, "t", Tanuki.Config.TextChannel, "Text channel ID")
	flag.StringVar(&Tanuki.Config.Owner, "o", Tanuki.Config.Owner, "Owner ID")
	flag.StringVar(&Tanuki.Config.YoutubeAPIKey, "y", Tanuki.Config.YoutubeAPIKey, "Youtube API key")
	flag.Parse()
}

func (bot *Bot) Init() {
	if !bot.Config.Validate() {
		log.Fatal("Invalid configuration")
		return
	}

	bot.Commands.ByPermission = make(PermissionCommand)
	bot.Commands.ByName = make(NameCommand)

	bot.Permissions = bot.Commands.InitPermissions("permissions.json")
	bot.Commands.InitPlayer()

	var err error
	bot.DiscordSession, err = discordgo.New(bot.Config.Token)
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.DiscordSession.AddHandler(bot.ProcessCommand)
	err = bot.DiscordSession.Open()
	if err != nil {
		log.Fatal(err)
		return
	}
}

func main() {
	Tanuki.Init()

	log.Println("Up and running!")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	Tanuki.DiscordSession.Close()
}
