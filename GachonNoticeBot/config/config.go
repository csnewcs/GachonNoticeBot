package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	. "github.com/csnewcs/gachonnoticebot/crolling"
)

type SendMessageChannel struct {
	All              []string `json:"all"`
	CloudEngineering []string `json:"cloudEngineering"`
}
type LastNotice struct {
	All              int `json:"all"`
	CloudEngineering int `json:"cloudEngineering"`
}

type Config struct {
	Token               string             `json:"token"`
	SendMessageChannels SendMessageChannel `json:"sendMessageChannels"`
	LastNotice          LastNotice         `json:"lastNotice"`
	IsTesting           bool               `json:"isTesting"`
	TestingGuilds       []string           `json:"testingGuilds"`
}

var Conf Config
func GetConfig() (error) {
	//config file path: ./config.json
	file, err := os.Open("config.json")
	if err != nil {
		return err
	}
	readFile, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	json.Unmarshal(readFile, &Conf)
	return nil
}

func SaveConfig() {
	Conf.LastNotice.All = LastNumbers[NoticePageAll]
	Conf.LastNotice.CloudEngineering = LastNumbers[NoticePageCloudEngineering]

	file, err := os.Create("config.json")
	if err != nil {
		fmt.Println("Error opening config file: ", err)
		return
	}
	jsonData, err := json.Marshal(Conf)
	// fmt.Println(string(jsonData))
	if err != nil {
		fmt.Println("Error marshalling json: ", err)
		return
	}
	file.Write(jsonData)
	defer file.Close()
}
