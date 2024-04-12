package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func makeSlashCommands(client *discordgo.Session) {
	var commands []*discordgo.ApplicationCommand = []*discordgo.ApplicationCommand{
		{
			Name:        "도움말",
			Description: "도움말을 표시합니다.",
		},
		{
			Name:        "등록",
			Description: "이 채널을 공지를 가져올 채널으로 등록합니다.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "공지선택",
					Description: "어느 공지를 받을지 선택합니다.",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "전체",
							Value: NoticePageAll,
						},
						{
							Name:  "클라우드공학과",
							Value: NoticePageCloudEnginerring,
						},
					},
				},
			},
			// DefaultMemberPermissions: discordgo.PermissionAdministrator,
		},
	}
	for _, command := range commands {
		client.ApplicationCommandCreate(client.State.User.ID, "", command)
	}
}

var slashCommandsExecuted = map[string]func(client *discordgo.Session, interactionCreated *discordgo.InteractionCreate){
	"도움말": func(client *discordgo.Session, interactionCreated *discordgo.InteractionCreate) {
	},
	"등록": func(client *discordgo.Session, interactionCreated *discordgo.InteractionCreate) {
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
		} else if noticePage == NoticePageCloudEnginerring {
			if contains(&conf.SendMessageChannels.CloudEnginerring, channelID) {
				content = "이미 등록되어 있습니다"
			} else {
				conf.SendMessageChannels.CloudEnginerring = append(conf.SendMessageChannels.CloudEnginerring, channelID)
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
}

func contains(slices *[]string, element string) bool {
	for _, slice := range *slices {
		if slice == element {
			return true
		}
	}
	return false
}
