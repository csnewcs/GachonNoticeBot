# GachonNoticeBot

![](https://cdn.discordapp.com/attachments/1226816980531679272/1229358660086661172/2024-04-15_17-50-43.png?ex=662f6478&is=661cef78&hm=20ba47515014b5fdd9d580088a22bbeb320e07c253412b3c3aa06046af351e7c&)
학교 공지가 올라오면 자동으로 이를 알려주는 디스코드 봇

---

## 사용 패키지
- go 1.21.8
- discordgo 0.28.1
- discordgo-embed 0.0.0-20220113222025-bafe0c917646
- goquery 1.9.1
---
## config.json
config.default.json의 이름을 config.json으로 변경해 사용하면 된다.
```
token: 디스코드 봇 토큰                                        string
isTesting: 이 봇이 테스트 모드인지 여부                         bool
testingGuilds: 테스트 모드이면 커맨드를 생성할 서버들 목록       string[]
sendMessageChannels: 메세지를 보낼 디스코드 채널들              string[]
lastNotice: 마지막 공지 번호들                                 int
```
## 실행
```
go mod download
go run .
```
## 구조
- main.go
  - main: 디스코드 봇 생성
  - getConfig/saveConfig: 설정 파일 관리
  - loopCheckingNewNotices: 새로운 공지 확인
  - sendNotice: 공지 전송
- crolling.go
  - GetNoticeList: 학교 사이트에서 공지 가져오기
    - parsingNoticeList: 학교 페이지에서 내용 뽑아내기
      - removeVoidText: 공백 문자들(9, 10, 32) 삭제
      - getNoticeLinks: 공지 페이지 링크 만들기
- slashCommand.go
  - makeSlashCommands: 슬래시 커맨드 생성
  - slashCommandExecuted: 슬래시 커맨드 처리