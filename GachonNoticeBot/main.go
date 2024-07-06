package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type NoticePage string

const (
	NoticePageAll              NoticePage = "all"
	NoticePageCloudEngineering NoticePage = "cloudEngineering"
)
type Config struct {
	Token               string             `json:"token"`
	SendMessageChannels SendMessageChannel `json:"sendMessageChannels"`
	LastNotice          LastNotice         `json:"lastNotice"`
	IsTesting           bool               `json:"isTesting"`
	TestingGuilds       []string           `json:"testingGuilds"`
}

type SendMessageChannel struct {
	All              []string `json:"all"`
	CloudEngineering []string `json:"cloudEngineering"`
}
type LastNotice struct {
	All              int `json:"all"`
	CloudEngineering int `json:"cloudEngineering"`
}

var conf Config
var session *discordgo.Session

func main() {
	var err error
	fmt.Println("Starting bot...")
	conf, err = getConfig()
	if err != nil {
		fmt.Println("Error getting config file: ", err)
		return
	}

	session, err = discordgo.New("Bot " + conf.Token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	testLog("Starting with testing mode...")

	session.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		testLog("인터렉션 받음: " + interaction.ApplicationCommandData().Name + "(" + interaction.ApplicationCommandData().ID + ") | " + interaction.GuildID)
		if function, ok := slashCommandsExecuted[interaction.ApplicationCommandData().Name]; ok {
			function(session, interaction)
		} // 해당 인터렉션이 있을 때 그에 맞는 함수 실행(./slashCommand.go)
	})
	err = session.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}
	makeSlashCommands(session)
	lastNumbers[NoticePageAll] = conf.LastNotice.All
	lastNumbers[NoticePageCloudEngineering] = conf.LastNotice.CloudEngineering
	if conf.IsTesting {
		go loopCheckingNewNotices(10)
	} else {
		go loopCheckingNewNotices(60)
	}
	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	// 프로그램 명령어
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		var command string = scanner.Text()
		if command == "exit" {
			session.Close()
			fmt.Println("\nbye")
			return
		}
		commandArr := strings.Split(command, " ")
		if commandFunc, ok := commands[commandArr[0]]; ok {
			result := commandFunc(commandArr)
			fmt.Println(result)
		} else {
			fmt.Println("Invalid command")
		}
	}
}

func getConfig() (Config, error) {
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

func loopCheckingNewNotices(delay int) { //주기적으로 새로운 공지 확인
	for {
		for noticePage, lastNumber := range lastNumbers {
			notices := GetNoticeList(noticePage)
			checkNewNotice(notices, noticePage, lastNumber, false)
		}
		saveConfig()
		time.Sleep(time.Duration(delay) * time.Second)
	}
}
func checkNewNotice(notices []Notice, noticePage NoticePage, lastNumber int, checkedLastNumber bool) {
	for i := range notices {
		notice := notices[len(notices)-i-1]
		if notice.Number > lastNumber {
			testLog("새로운 공지: " + strconv.Itoa(notice.Number) + " | " + notice.Title + " | " + notice.Auther + " | " + notice.Date + " | " + notice.Views + " | " + notice.File)
			if !checkedLastNumber {
				lastNumbers[noticePage] = notice.Number
			}
			if isNewNotice(notice.Title, noticePage) {
				sendNotice(notice, noticePage)
				testLog("새로운 공지, 전송")
				addToSendedNotices(notice.Title, noticePage)
				if checkedLastNumber {
					return
				}
			} else {
				testLog("새로운 공지 아님 이전 번호 탐색")
				checkNewNotice(notices, noticePage, lastNumber - 1, true)
				break
			}
		}
	}
}
func isNewNotice(title string, noticePage NoticePage) bool {
	for i := 0; i < 50; i++ {
		if sendedNotices[noticePage][i] == title {
			return false
		}
	}
	return true
}
func sendNotice(notice Notice, noticePage NoticePage) {
	var channels []string
	switch noticePage {
	case NoticePageAll:
		channels = conf.SendMessageChannels.All
	case NoticePageCloudEngineering:
		channels = conf.SendMessageChannels.CloudEngineering
	}
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
		}
		go session.ChannelMessageSendEmbed(channel, &embed)
	}
}

func saveConfig() {
	conf.LastNotice.All = lastNumbers[NoticePageAll]
	conf.LastNotice.CloudEngineering = lastNumbers[NoticePageCloudEngineering]

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

func testLog(message string) {
	if conf.IsTesting {
		fmt.Println(message)
	}
}
