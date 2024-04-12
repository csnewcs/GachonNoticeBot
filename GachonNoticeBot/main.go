package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	// "github.com/clinet/discordgo-embed"
)

type SendConfig struct {
	All              []string `json:"all"`
	CloudEnginerring []string `json:"cloudEnginerring"`
}
type LastNotice struct {
	All              int `json:"all"`
	CloudEnginerring int `json:"cloudEnginerring"`
}
type Config struct {
	Token               string     `json:"token"`
	SendMessageChannels SendConfig `json:"sendMessageChannels"`
	LastNotice          LastNotice `json:"lastNotice"`
}

var conf Config
var session *discordgo.Session

func main() {
	var err error
	fmt.Println("Starting bot...")
	conf, err = getConfigFile()
	if err != nil {
		fmt.Println("Error getting config file: ", err)
		return
	}

	session, err = discordgo.New("Bot " + conf.Token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}
	// session.AddHandler(guildCreate)
	session.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		fmt.Println("인터렉션 받음: ", interaction.ApplicationCommandData().Name)

		if function, ok := slashCommandsExecuted[interaction.ApplicationCommandData().Name]; ok {
			function(session, interaction)
		}
	})
	err = session.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}
	makeSlashCommands(session)
	lastNumbers[NoticePageAll] = conf.LastNotice.All
	lastNumbers[NoticePageCloudEnginerring] = conf.LastNotice.CloudEnginerring
	go loopCheckingNewNotices(60)
	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)

	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	fmt.Println("\nbye")
	session.Close()
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

func loopCheckingNewNotices(delay int) {
	for {
		for noticePage, lastNumber := range lastNumbers {
			notices := GetNoticeList(noticePage)
			for _, notice := range notices {
				if notice.Number > lastNumber {
					fmt.Println("새로운 공지: ", notice.Number)
					sendNotice(notice, noticePage)
					lastNumbers[noticePage] = notice.Number
					break
				}
			}
		}
		time.Sleep(time.Duration(delay) * time.Second)
	}
}
func sendNotice(notice Notice, noticePage NoticePage) {
	var channels []string
	switch noticePage {
	case NoticePageAll:
		channels = conf.SendMessageChannels.All
		conf.LastNotice.All = notice.Number
	case NoticePageCloudEnginerring:
		channels = conf.SendMessageChannels.CloudEnginerring
		conf.LastNotice.CloudEnginerring = notice.Number
	}
	saveConfig()
	for _, channel := range channels {
		fileExist := ""
		if notice.File != "0" {
			fileExist = "📎"
		}
		embed := discordgo.MessageEmbed{
			Title: notice.Title,
			URL:   notice.Link,
			Color: 0x3a4480,
			Footer: &discordgo.MessageEmbedFooter{
				Text: notice.Date + " | " + notice.Auther + " " + fileExist,
			},
			// Description: getDescription(notice.ContentLink),
		}
		session.ChannelMessageSendEmbed(channel, &embed)
	}
}

func saveConfig() {
	file, err := os.Create("config.json")
	if err != nil {
		fmt.Println("Error opening config file: ", err)
		return
	}
	jsonData, err := json.Marshal(conf)
	// fmt.Println(string(jsonData))
	if err != nil {
		fmt.Println("Error marshalling json: ", err)
		return
	}
	file.Write(jsonData)
	defer file.Close()
}
