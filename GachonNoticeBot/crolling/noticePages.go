package crolling

type NoticePage string

const (
	NoticePageAll              NoticePage = "all"
	NoticePageCloudEngineering NoticePage = "cloudEngineering"
)

var LastNumbers = map[NoticePage]int{
	NoticePageAll:              0,
	NoticePageCloudEngineering: 0,
}
var NoticeURLList = map[NoticePage]string{
	NoticePageAll:              "https://www.gachon.ac.kr/kor/7986/subview.do",
	NoticePageCloudEngineering: "https://www.gachon.ac.kr/ce/9514/subview.do",
}
var SendedNotices = map[NoticePage][]string{
	NoticePageAll:              make([]string, 50),
	NoticePageCloudEngineering: make([]string, 50),
}