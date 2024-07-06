## 들어가기게 앞서
이 문서는 독자가 프로그래밍 지식이 있음을 전제하고 있습니다. Go언어, HTML과 JSON 등의 지식이 필요합니다.

## 변수와 함수 설명
### NoticePage
공지들의 구분에 있어 가장 중요한 변수입니다. 다른 언어들의 열거형(enum)저럼 사용합니다. main.go 18줄에서 확인할 수 있습니다.

---
- main.go
  - Config: 봇의 설정. `./config.json`을 읽어 그 데이터를 저장.
    |변수|설명|타입|
    |---|---|---|
    |Token|디스코드 봇의 토큰|string|
    |SendMessageChannels|공지를 보낼 채널들|공지 페이지, []string|
    |LastNotice|마지막으로 보낸 공지 번호|공지 페이지, int|
    |IsTesting|테스트 중인지 여부|bool|
    |TestingGuilds|테스트 모드에서 메세지를 보낼 서버들|[]string|
  - getConfig(), setConfig(): `./config.json`파일을 읽거나 씁니다.
  - loopCheckingNewNotice(): 주기적으로 `checkNewNotice()`함수를 호출해 새로운 공지를 확인합니다.
- crolling.go
  - lastNumbers: 마지막 전송한 공지 번호
  - noticeURLList: 공지 페이지, URL
  - sendedNotices: 페이지별 전송된 최근 50개의 공지들
  - Notice: 공지 정보
    |변수|설명|타입|
    |---|---|---|
    |Number|공지 번호|int|
    |Title|공지 제목|string|
    |Link|공지 링크|string|
    |ContentLink|공지 내용만 들어있는 페이지 링크|string|
    |Auther|공지 작성자|string|
    |Date|공지 작성 날짜|string|
    |Views|공지 조회수|string|
    |File|공지 페이지에 첨부된 파일 수|string|
  - GetNoticeList(): 공지 페이지들의 공지들을 불러옵니다.반환합니다.
  - getNoticeLinks[NoticePage, func]: 공지 링크를 뽑아냅니다. **공지 페이지들마다 조금씩 다르게 만들어야 하기** 때문에 NoticePage별로 나뉘어 있습니다.
- slashCommand.go
  - slashCommandExecuted[string, func]: 디스코드 슬래시 명령어들이 실행되는 곳
    - 등록: 명령어 사용 채널을 지정한 공지들의 공지가 오도록 설정
    - 해제: 해당 채널이 어떤 공지도 오지 않도록 설정
---
## 공지 페이지 추가 가이드
main.go
1. main.go에서 NoticePage 추가 (`추가할 NoticePage` NoticePage = "문자열")
2. `config.json`, `config.default.json`에 추가한 NoticePage의 **문자열**로 SendMessageChannels와 lastNotice아래에 추가
```json
"sendMessageChannels": {
    ...
    "추가할 NoticePage의 문자열": []
},
"lastNotice": {
    ...
    "추가할 NoticePage의 문자열": 0
}
```
3. SendMessageChannel과 LastNotice에 NoticePage에 맞게 각각 []string과 int로 추가
---
crolling.go
1. lastNumbers에 `추가한 NoticePage`: 0 추가
2. noticeURLList에 `추가한 NoticePage`: `공지 페이지` 추가
3. sendedNotices에 `추가한 NoticePage`: make([]string, 50) 추가
4. getNoticeLinks에 `추가한 NoticePage`: `해당 페이지에서 공지 링크를 (Link, ContentLink) 형식 반환하는 func` 추가
---
slashCommand.go
1. makeSlashCommands().commands에 Options > Choices아래 다음과 같은 내용 추가
```go
{
    Name: "추가할 페이지 (한글) 이름",
    Value: 추가한 NoticePage
},
``` 
2. slashCommandsExecuted["등록"]에서 `channalID := interactionCreated.ChannelID` 구문 아래 조건문에 다음 내용 추가
```go
else if noticePage == <추가한 NoticePage> {
    if contains(&conf.SendMessageChannels.<추가한 NoticePage>, channelID) {
        content = "이미 등록되어 있습니다"
    } else {
        conf.SendMessageChannels.<추가한 NoticePage> = append(conf.SendMessageChannels.<추가한 NoticePage>, channelID)
        content = fmt.Sprintf("해당 채널을 `%s` 공지를 가져올 채널로 등록했습니다.", noticePage)
    }
}
```
3. slashCommandExecuted["해제"]에서 `saveConfig()` 위에 다음과 내용 추가
```go
index = indexOf(conf.SendMessageChannels.<추가한 NoticePage>, channelID)
    testLog("해제 | indexOf<추가한 NoticePage>: " + strconv.Itoa(index))
    if index != -1 {
        conf.SendMessageChannels.<추가한 NoticePage> = append(conf.SendMessageChannels.<추가한 NoticePage>[:index], conf.SendMessageChannels.<추가한 NoticePage>[index+1:]...)
        content += fmt.Sprintf("해당 채널에 `%s` 공지가 오지 않도록 설정했습니다\n", <추가한 NoticePage>)
}
```
