package crolling

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// 콘텐츠 위치: HTML > body > div.(sub _responsiveObj sub) > div.wrap-contents > div.container > div.contents > div.scroll-table > table.(board-table horizon), tbody
// tr > td.td-num: 번호 / td.td-subject > a > strong: 제목 / td.td-write: 작성자 / td.td-date: 작성일 / td.td-access: 조회수 / td.td-file: 첨부파일
// 링크: javascript:jf_viewArtcl('kor', '96372') 이런 형식 => https://www.gachon.ac.kr/commonNotice/kor/96372/artclView.do 여기로 연결
// fnct1|@@|%2FcommonNotice%2Fkor%2F96372%2FartclView.do%3Fpage%3D1%26srchColumn%3D%26srchWord%3D%26
// 위 링크에서 '96372'를 {숫자}라 할 때 fnct1|@@|%2FcommonNotice%2Fkor%2F{숫자}%2FartclView.do%3Fpage%3D1%26srchColumn%3D%26srchWord%3D%26 이걸 base64로 인코딩 한 다음 https://www.gachon.ac.kr/kor/7986/subview.do?enc={인코딩} 하면 해당 페이지로 이동

type Notice struct {
	Number      int
	Title       string
	Link        string
	ContentLink string
	Auther      string
	Date        string
	Views       string
	File        string
}

func GetNoticeList(noticePage NoticePage) []Notice {
	var notices []Notice

	// Request
	resp, err := http.Get(NoticeURLList[noticePage])
	if err != nil {
		// panic(err)
		now := time.Now()
		fmt.Printf("[%d-%d-%d %d:%d:%d] ERR: HTTP GET FAILED\n", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
		resp.Body.Close()
		return make([]Notice, 0)
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
		link, contentLink := getNoticeLinks[noticePage](sel)
		num, _ := strconv.Atoi(removeVoidText(sel.Find("td.td-num").Text()))
		notice := Notice{
			Number:      num,
			Link:        link,
			ContentLink: contentLink,
			Title:       sel.Find("strong").Text(),
			Auther:      removeVoidText(sel.Find("td.td-write").Text()),
			Date:        sel.Find("td.td-date").Text(),
			Views:       removeVoidText(sel.Find("td.td-access").Text()),
			File:        removeVoidText(sel.Find("td.td-file").Text()),
		}
		notices = append(notices, notice)
	})
	return notices
}

func AddToSendedNotices(title string, noticePage NoticePage) {
	for i := 49; i > 0; i-- {
		SendedNotices[noticePage][i] = SendedNotices[noticePage][i-1]
	}
	SendedNotices[noticePage][0] = title
}

func removeVoidText(text string) string {
	voidTexts := []string{string(rune(10)), string(rune(9)), string(rune(32))}
	for _, voidText := range voidTexts {
		text = strings.ReplaceAll(text, voidText, "")
	}
	return text
}

var getNoticeLinks = map[NoticePage]func(selection *goquery.Selection) (string, string){ // 학교공지와 학과공지 페이지가 살짝 달라 구분
	NoticePageAll: func(selection *goquery.Selection) (string, string) {
		aInfo := selection.Find("a").AttrOr("href", "")
		number, _ := strconv.Atoi(strings.Split(aInfo, "'")[3])
		textToEncode := fmt.Sprintf("fnct1|@@|%%2FcommonNotice%%2Fkor%%2F%d%%2FartclView.do%%3Fpage%%3D1%%26srchColumn%%3D%%26srchWord%%3D%%26", number)
		encoded := base64.StdEncoding.EncodeToString([]byte(textToEncode))
		return NoticeURLList[NoticePageAll] + "?enc=" + encoded, fmt.Sprintf("https://www.gachon.ac.kr/commonNotice/kor/%d/artclView.do", number)
	},
	NoticePageCloudEngineering: func(selection *goquery.Selection) (string, string) {
		aInfo := selection.Find("a").AttrOr("href", "")          // /bbs/ce/1496/96327/artclView.do
		onclickInfo := selection.Find("a").AttrOr("onclick", "") // jf_viewArtcl('ce', '1496', '96327')
		numbers := strings.Split(onclickInfo, "'")
		departmentNumber := numbers[3]
		articleNumber := numbers[5]

		textToEncode := fmt.Sprintf("fnct1|@@|/bbs/ce/%s/%s/artclView.do?page=1&srchColumn=&srchWrd=&bbsClSeq=&bbsOpenWrdSeq=&rgsBgndeStr=&rgsEnddeStr=&isViewMine=false&password=&", departmentNumber, articleNumber)
		encoded := base64.StdEncoding.EncodeToString([]byte(textToEncode))
		return NoticeURLList[NoticePageCloudEngineering] + "?enc=" + encoded, "https://www.gachon.ac.kr/" + aInfo
	},
}