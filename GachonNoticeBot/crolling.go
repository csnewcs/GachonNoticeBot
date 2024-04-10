package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// var noticeURL string = "https://www.gachon.ac.kr/kor/7986/subview.do"
// var cloudEnginerringNoticeURL string = "https://www.gachon.ac.kr/ce/9514/subview.do"
var lastNumbers int = 0
var noticeURLList = map[NoticePage]string{
	NoticePageAll:   "https://www.gachon.ac.kr/kor/7986/subview.do",
	NoticePageCloudEnginerring: "https://www.gachon.ac.kr/ce/9514/subview.do",
}

type NoticePage int

const (
	NoticePageAll NoticePage = iota
	NoticePageCloudEnginerring NoticePage = iota
)

// 콘텐츠 위치: HTML > body > div.(sub _responsiveObj sub) > div.wrap-contents > div.container > div.contents > div.scroll-table > table.(board-table horizon), tbody
// tr > td.td-num: 번호 / td.td-subject > a > strong: 제목 / td.td-write: 작성자 / td.td-date: 작성일 / td.td-access: 조회수 / td.td-file: 첨부파일
// 링크: javascript:jf_viewArtcl('kor', '96372') 이런 형식 => https://www.gachon.ac.kr/commonNotice/kor/96372/artclView.do 여기로 연결
// fnct1|@@|%2FcommonNotice%2Fkor%2F96372%2FartclView.do%3Fpage%3D1%26srchColumn%3D%26srchWord%3D%26
// 위 링크에서 '96372'를 {숫자}라 할 때 fnct1|@@|%2FcommonNotice%2Fkor%2F{숫자}%2FartclView.do%3Fpage%3D1%26srchColumn%3D%26srchWord%3D%26 이걸 base64로 인코딩 한 다음 https://www.gachon.ac.kr/kor/7986/subview.do?enc={인코딩} 하면 해당 페이지로 이동

type Notice struct {
	Number      string
	Title       string
	Link        string
	ContentLink string
	Auther      string
	Date        string
	Views       string
	FileLinks   []string
}

func GetNoticeList(noticePage NoticePage) []Notice {
	var notices []Notice

	// Request
	resp, err := http.Get(noticeURLList[noticePage])
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Parsing
	notices = parsingNoticeList(resp.Body, noticePage)
	return notices
}

func parsingNoticeList(page io.Reader, noticePage NoticePage) []Notice {
	notices := make([]Notice, 0)
	html, _ := goquery.NewDocumentFromReader(page)
	table := html.Find("table.board-table")
	tbody := table.Find("tbody")
	tbody.Find("tr").Each(func(i int, sel *goquery.Selection) {
		// voidText := string(sel.Find("td.td-write").Text()[0])
		link, contentLink := GetNoticeLinks[NoticePageAll](sel)
		notice := Notice{
			Number:      sel.Find("td.td-num").Text(),
			Link:        link,
			ContentLink: contentLink,
			Title:       sel.Find("strong").Text(),
			Auther:      removeVoidText(sel.Find("td.td-write").Text()),
			Date:        sel.Find("td.td-date").Text(),
			Views:       removeVoidText(sel.Find("td.td-access").Text()),
		}
		notices = append(notices, notice)
		fmt.Println(notice)
	})
	return notices
}

func removeVoidText(text string) string {
	voidTexts := []string{string(10), string(9), string(32)}
	for _, voidText := range voidTexts {
		text = strings.ReplaceAll(text, voidText, "")
	}
	return text
}

var GetNoticeLinks = map[NoticePage]func(selection *goquery.Selection) (string, string){
	NoticePageAll: func(selection *goquery.Selection) (string, string) {
		number, _ := strconv.Atoi(strings.Split(selection.AttrOr("href", ""), "'")[3])
		textToEncode := fmt.Sprintf("fnct1|@@|%%2FcommonNotice%%2Fkor%%2F%d%%2FartclView.do%%3Fpage%%3D1%%26srchColumn%%3D%%26srchWord%%3D%%26", number)
		encoded := base64.StdEncoding.EncodeToString([]byte(textToEncode))
		return "https://www.gachon.ac.kr/kor/7986/subview.do?enc=" + encoded, fmt.Sprintf("https://www.gachon.ac.kr/commonNotice/kor/%d/artclView.do", number)
	},
	NoticePageCloudEnginerring: func(selection *goquery.Selection) (string, string) {
		
		return "", ""	
	},
}
