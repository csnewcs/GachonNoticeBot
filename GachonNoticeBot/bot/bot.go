package bot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	. "github.com/csnewcs/gachonnoticebot/config"
	. "github.com/csnewcs/gachonnoticebot/crolling"
)

var session *discordgo.Session

func Init() error {
	session, _ = discordgo.New("Bot " + Conf.Token)
	session.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if interaction.Type != discordgo.InteractionApplicationCommand {
			return
		}
		testLog("인터렉션 받음: " + interaction.ApplicationCommandData().Name + "(" + interaction.ApplicationCommandData().ID + ") | " + interaction.GuildID)
		if function, ok := slashCommandsExecuted[interaction.ApplicationCommandData().Name]; ok {
			function(session, interaction)
		} // 해당 인터렉션이 있을 때 그에 맞는 함수 실행(./slashCommand.go)
	})
	var err error = nil
	err = session.Open()
	makeSlashCommands(session)
	LastNumbers[NoticePageAll] = Conf.LastNotice.All
	LastNumbers[NoticePageCloudEngineering] = Conf.LastNotice.CloudEngineering
	return err
}
func Close() {
	session.Close()
}

func SendNotice(notice Notice, noticePage NoticePage) {
	var channels []string
	switch noticePage {
	case NoticePageAll:
		channels = Conf.SendMessageChannels.All
	case NoticePageCloudEngineering:
		channels = Conf.SendMessageChannels.CloudEngineering
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
		now := time.Now()
		fmt.Printf("[%d-%d-%d %d:%d:%d] INF: New notice {%s} detected. send to '{%s}' channel\n", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), notice.Title, channel)
		go session.ChannelMessageSendEmbed(channel, &embed)
	}
}
func testLog(message string) {
	if Conf.IsTesting {
		fmt.Println(message)
	}
}
