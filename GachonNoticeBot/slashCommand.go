package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
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
	for _, command := range commands {
		if conf.IsTesting {
			fmt.Println("Created Command: ", command.Name)
			for _, guildID := range conf.TestingGuilds {
				client.ApplicationCommandCreate(client.State.User.ID, guildID, command)
			}
		} else {
			client.ApplicationCommandCreate(client.State.User.ID, "", command)
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
		channelID := interactionCreated.ChannelID
		if noticePage == NoticePageAll {
			if contains(&conf.SendMessageChannels.All, channelID) {
				content = "이미 등록되어 있습니다"
			} else {
				conf.SendMessageChannels.All = append(conf.SendMessageChannels.All, channelID)
				content = fmt.Sprintf("해당 채널을 `%s` 공지를 가져올 채널로 등록했습니다.", noticePage)
			}
		} else if noticePage == NoticePageCloudEngineering {
			if contains(&conf.SendMessageChannels.CloudEngineering, channelID) {
				content = "이미 등록되어 있습니다"
			} else {
				conf.SendMessageChannels.CloudEngineering = append(conf.SendMessageChannels.CloudEngineering, channelID)
				content = fmt.Sprintf("해당 채널을 `%s` 공지를 가져올 채널로 등록했습니다.", noticePage)
			}
		}
		saveConfig()
		client.InteractionRespond(interactionCreated.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: content,
			},
		})
	},
	"해제": func(client *discordgo.Session, interactionCreated *discordgo.InteractionCreate) {
		if(conf.IsTesting) {
			fmt.Println("해제 | ", interactionCreated.ChannelID)
		}
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

		index := indexOf(conf.SendMessageChannels.All, channelID)
		if(conf.IsTesting) {
			fmt.Println("해제 | indexOfAll: ", index)
		}
		if index != -1 {
			conf.SendMessageChannels.All = append(conf.SendMessageChannels.All[:index], conf.SendMessageChannels.All[index+1:]...)
			content += fmt.Sprintf("해당 채널에 `%s` 공지가 오지 않도록 설정했습니다\n", NoticePageAll)
		}
		index = indexOf(conf.SendMessageChannels.CloudEngineering, channelID)
		if(conf.IsTesting) {
			fmt.Println("해제 | indexOfCloudEngineering: ", index)
		}
		if index != -1 {
			conf.SendMessageChannels.CloudEngineering = append(conf.SendMessageChannels.CloudEngineering[:index], conf.SendMessageChannels.CloudEngineering[index+1:]...)
			content += fmt.Sprintf("해당 채널에 `%s` 공지가 오지 않도록 설정했습니다\n", NoticePageCloudEngineering)
		}
		saveConfig()
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
