package bot

import (
	"strconv"

	"github.com/csnewcs/gachonnoticebot/crolling"
)

var Commands = map[string]func([]string) string{
	"setLastNoticeNumber": func(s []string) string {
		//changeLastNoticeNumber <NoticePage> <Number>
		if len(s) != 3 {
			return "Invalid arguments"
		}
		noticePage := crolling.NoticePage(s[1])
		number, err := strconv.Atoi(s[2])
		if err != nil {
			return "Invalid number"
		}
		crolling.LastNumbers[noticePage] = number
		return "Changed last notice number to " + s[2]
	},
	"getLastNoticeNumber": func(s []string) string {
		//getLastNoticeNumber <NoticePage>
		if len(s) != 2 {
			return "Invalid arguments"
		}
		noticePage := crolling.NoticePage(s[1])
		return "Last notice number is " + strconv.Itoa(crolling.LastNumbers[noticePage])
	},
	"resetSendedNotices": func(s []string) string {
		//resetSendNotices <NoticePage>
		if len(s) != 2 {
			return "Invalid arguments"
		}
		noticePage := crolling.NoticePage(s[1])
		for i := range crolling.SendedNotices[noticePage] {
			crolling.SendedNotices[noticePage][i] = ""
		}
		return "Reset sended notices of " + s[1]
	},
}
