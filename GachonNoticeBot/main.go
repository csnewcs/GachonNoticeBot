package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/csnewcs/gachonnoticebot/bot"
	"github.com/csnewcs/gachonnoticebot/config"
	"github.com/csnewcs/gachonnoticebot/crolling"
)

func main() {
	var err error
	fmt.Println("Starting bot...")
	err = config.GetConfig()
	if err != nil {
		fmt.Println("Error getting config file: ", err)
		return
	}
	err = bot.Init()
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}
	testLog("Starting with testing mode...")

	if config.Conf.IsTesting {
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
			bot.Close()
			fmt.Println("\nbye")
			return
		}
		commandArr := strings.Split(command, " ")
		if commandFunc, ok := bot.Commands[commandArr[0]]; ok {
			result := commandFunc(commandArr)
			fmt.Println(result)
		} else {
			fmt.Println("Invalid command")
		}
	}
}

func loopCheckingNewNotices(delay int) { //주기적으로 새로운 공지 확인
	for {
		for noticePage, lastNumber := range crolling.LastNumbers {
			notices := crolling.GetNoticeList(noticePage)
			checkNewNotice(notices, noticePage, lastNumber, false)
		}
		config.SaveConfig()
		time.Sleep(time.Duration(delay) * time.Second)
	}
}
func checkNewNotice(notices []crolling.Notice, noticePage crolling.NoticePage, lastNumber int, checkedLastNumber bool) {
	for i := range notices {
		notice := notices[len(notices)-i-1]
		if notice.Number > lastNumber {
			testLog("새로운 공지: " + strconv.Itoa(notice.Number) + " | " + notice.Title + " | " + notice.Auther + " | " + notice.Date + " | " + notice.Views + " | " + notice.File)
			if !checkedLastNumber {
				crolling.LastNumbers[noticePage] = notice.Number
			}
			if isNewNotice(notice.Title, noticePage) {
				bot.SendNotice(notice, noticePage)
				testLog("새로운 공지, 전송")
				crolling.AddToSendedNotices(notice.Title, noticePage)
				if checkedLastNumber {
					return
				}
			} else {
				testLog("새로운 공지 아님 이전 번호 탐색")
				checkNewNotice(notices, noticePage, lastNumber-1, true)
				break
			}
		}
	}
}
func isNewNotice(title string, noticePage crolling.NoticePage) bool {
	for i := 0; i < 50; i++ {
		if crolling.SendedNotices[noticePage][i] == title {
			return false
		}
	}
	return true
}

func testLog(message string) {
	if config.Conf.IsTesting {
		fmt.Println(message)
	}
}
