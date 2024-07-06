# GachonNoticeBot

![](https://cdn.discordapp.com/attachments/1226816980531679272/1229358660086661172/2024-04-15_17-50-43.png?ex=662f6478&is=661cef78&hm=20ba47515014b5fdd9d580088a22bbeb320e07c253412b3c3aa06046af351e7c&)
학교 공지가 올라오면 자동으로 이를 알려주는 디스코드 봇

---

## 사용 패키지
- go 1.21.8
- discordgo 0.28.1
- goquery 1.9.1
---
## config.json
config.default.json의 이름을 config.json으로 변경해 사용하면 된다.
```
token: 디스코드 봇 토큰                                        string
isTesting: 이 봇이 테스트 모드인지 여부                         bool
testingGuilds: 테스트 모드이면 커맨드를 생성할 서버들 목록       string[P]
sendMessageChannels: 메세지를 보낼 디스코드 채널들              string[]
lastNotice: 마지막 공지 번호들                                 int
```
## 실행
```
go mod download
go run .
```
## 새 공지 페이지 만들기 가이드라인
[여기](https://github.com/csnewcs/GachonNoticeBot/blob/main/CreateNewNoticeGuideLine.md)를 참조하세요
