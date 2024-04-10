package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	// "github.com/clinet/discordgo-embed"
)

type Config struct {
	Token               string   `json:"token"`
	SendMessageChannels []string `json:"sendMessageChannels"`
}

var conf Config
var notices []Notice

func main() {
	var err error
	notices = GetNoticeList(NoticePageAll)
	fmt.Println("Starting bot...")
	conf, err = getConfigFile()
	if err != nil {
		fmt.Println("Error getting config file: ", err)
		return
	}

	discordBot, err := discordgo.New("Bot " + conf.Token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}
	discordBot.AddHandler(guildCreate)

	err = discordBot.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)

	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	fmt.Println("\nbye")
	discordBot.Close()
}

func getConfigFile() (Config, error) {
	//config file path: ./config.json
	file, err := os.Open("config.json")
	if err != nil {
		return Config{}, err
	}
	readFile, err := io.ReadAll(file)
	if err != nil {
		return Config{}, err
	}
	var jsonConfig Config
	json.Unmarshal(readFile, &jsonConfig)

	defer file.Close()
	return jsonConfig, nil
}

func guildCreate(discordBot *discordgo.Session, guild *discordgo.GuildCreate) {
	fmt.Println("Joined guild: ", guild.Name)
	for _, channel := range conf.SendMessageChannels {
		discordBot.ChannelMessageSend(channel, fmt.Sprintf("공지가 있습니다!```md\n# 제목: %s\n작성자: %s\n업로드일: %s```\n링크: %s", notices[0].Title, notices[0].Auther, notices[0].Date, notices[0].Link))
	}
}
