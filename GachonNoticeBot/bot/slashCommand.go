package bot

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	. "github.com/csnewcs/gachonnoticebot/config"
	. "github.com/csnewcs/gachonnoticebot/crolling"
)

func makeSlashCommands(client *discordgo.Session) {
	var commands []*discordgo.ApplicationCommand = []*discordgo.ApplicationCommand{
		{
			Name:        "등록",
			Description: "이 채널을 공지를 가져올 채널으로 등록합니다.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "공지선택",
					Description: "어느 공지를 받을지 선택합니다.(서버 관리 역할 필요)",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "전체",
							Value: NoticePageAll,
						},
						{
							Name:  "클라우드공학과",
							Value: NoticePageCloudEngineering,
						},
					},
				},
			},
		},
		{
			Name:        "해제",
			Description: "이 채널을 공지를 가져오지 않도록 설정합니다.(서버 관리 역할 필요)",
		},
	}
	var oldCommands []*discordgo.ApplicationCommand;
	if Conf.IsTesting {
		for _, guildID := range Conf.TestingGuilds {
			oldCommands, _ = client.ApplicationCommands(client.State.User.ID, guildID)
			removeCommands(oldCommands, client)
		}
	}
	
	oldCommands, _ = client.ApplicationCommands(client.State.User.ID, "")
	removeCommands(oldCommands, client)

	for _, command := range commands {
		if Conf.IsTesting {
			for _, guildID := range Conf.TestingGuilds {
				client.ApplicationCommandCreate(client.State.User.ID, guildID, command)
			}
			testLog("Created Command: " + command.Name)
		} else {
			client.ApplicationCommandCreate(client.State.User.ID, "", command)
		}
	}
}

func removeCommands(commands []*discordgo.ApplicationCommand, client *discordgo.Session) {
	for _, command := range commands {
		err := client.ApplicationCommandDelete(client.State.User.ID, command.GuildID, command.ID)
		testLog("Remove Command: " + command.Name + " | " + command.ID)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

var slashCommandsExecuted = map[string]func(client *discordgo.Session, interactionCreated *discordgo.InteractionCreate){
	"등록": func(client *discordgo.Session, interactionCreated *discordgo.InteractionCreate) {
		if interactionCreated.Member.Permissions&discordgo.PermissionManageServer == 0 && interactionCreated.GuildID == interactionCreated.User.ID {
			client.InteractionRespond(interactionCreated.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "관리자 권한이 필요합니다.",
				},
			})
			return
		}
		content := ""
		noticePage := NoticePage(interactionCreated.ApplicationCommandData().Options[0].StringValue())
		fmt.Println(noticePage)
		channelID := interactionCreated.ChannelID
		if noticePage == NoticePageAll {
			if contains(&Conf.SendMessageChannels.All, channelID) {
				content = "이미 등록되어 있습니다"
			} else {
				Conf.SendMessageChannels.All = append(Conf.SendMessageChannels.All, channelID)
				content = fmt.Sprintf("해당 채널을 `%s` 공지를 가져올 채널로 등록했습니다.", noticePage)
			}
		} else if noticePage == NoticePageCloudEngineering {
			if contains(&Conf.SendMessageChannels.CloudEngineering, channelID) {
				content = "이미 등록되어 있습니다"
			} else {
				Conf.SendMessageChannels.CloudEngineering = append(Conf.SendMessageChannels.CloudEngineering, channelID)
				content = fmt.Sprintf("해당 채널을 `%s` 공지를 가져올 채널로 등록했습니다.", noticePage)
			}
		}
		SaveConfig()
		client.InteractionRespond(interactionCreated.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
			},
		})
	},
	"해제": func(client *discordgo.Session, interactionCreated *discordgo.InteractionCreate) {
		testLog("해제 | " + interactionCreated.ChannelID)
		if interactionCreated.Member.Permissions&discordgo.PermissionManageServer == 0 && interactionCreated.GuildID == interactionCreated.User.ID {
			client.InteractionRespond(interactionCreated.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "관리자 권한이 필요합니다.",
				},
			})
			return
		}
		content := ""
		channelID := interactionCreated.ChannelID

		index := indexOf(Conf.SendMessageChannels.All, channelID)
		testLog("해제 | indexOfAll: " + strconv.Itoa(index))
		if index != -1 {
			Conf.SendMessageChannels.All = append(Conf.SendMessageChannels.All[:index], Conf.SendMessageChannels.All[index+1:]...)
			content += fmt.Sprintf("해당 채널에 `%s` 공지가 오지 않도록 설정했습니다\n", NoticePageAll)
		}
		index = indexOf(Conf.SendMessageChannels.CloudEngineering, channelID)
		testLog("해제 | indexOfCloudEngineering: " + strconv.Itoa(index))
		if index != -1 {
			Conf.SendMessageChannels.CloudEngineering = append(Conf.SendMessageChannels.CloudEngineering[:index], Conf.SendMessageChannels.CloudEngineering[index+1:]...)
			content += fmt.Sprintf("해당 채널에 `%s` 공지가 오지 않도록 설정했습니다\n", NoticePageCloudEngineering)
		}
		SaveConfig()
		client.InteractionRespond(interactionCreated.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
			},
		})
	},
}

func indexOf(slice []string, element string) int {
	for i, sliceElement := range slice {
		if sliceElement == element {
			return i
		}
	}
	return -1
}

func contains(slices *[]string, element string) bool {
	for _, slice := range *slices {
		if slice == element {
			return true
		}
	}
	return false
}

