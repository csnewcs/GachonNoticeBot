package main

import (
	"strconv"
)

var commands = map[string]func([]string) string{
	"setLastNoticeNumber": func(s []string) string {
		//changeLastNoticeNumber <NoticePage> <Number>
		if len(s) != 3 {
			return "Invalid arguments"
		}
		noticePage := NoticePage(s[1])
		number, err := strconv.Atoi(s[2])
		if err != nil {
			return "Invalid number"
		}
		lastNumbers[noticePage] = number
		return "Changed last notice number to " + s[2]
	},
	"getLastNoticeNumber": func(s []string) string {
		//getLastNoticeNumber <NoticePage>
		if len(s) != 2 {
			return "Invalid arguments"
		}
		noticePage := NoticePage(s[1])
		return "Last notice number is " + strconv.Itoa(lastNumbers[noticePage])
	},
	"resetSendedNotices": func(s []string) string {
		//resetSendNotices <NoticePage>
		if len(s) != 2 {
			return "Invalid arguments"
		}
		noticePage := NoticePage(s[1])
		for i := range sendedNotices[noticePage] {
			sendedNotices[noticePage][i] = ""
		}
		return "Reset sended notices of " + s[1]
	},
}
