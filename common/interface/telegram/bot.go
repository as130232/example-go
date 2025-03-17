package telegram

import (
	"log"
	"strings"
	"time"

	telegramConst "linebot-go/common/infrastructure/consts/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot (Token 請找 @BotFather 申請一個, 有藍色勾勾才是官方的)
//
// 限制 :
//
// 1. 發送訊息有傳送速度的限制。（同個群組內限每分鐘最多傳送二十則訊息，而全域限制為每秒上限三十則訊息。）
//
// 2. 機器人只能以管理員的身分被加入到頻道。
//
// 3. 超級群組中機器人的上限數量為 20 個，一般群組則不限制。
//
// 4. 一組 token 只能有一個 bot instance 在監聽，但可以多個 bot instance 同時送訊息
type Bot struct {
	botAPI      *tgbotapi.BotAPI
	tokenChatId *telegramConst.TokenChatId
}

func NewBot(appEnv string, replyCallback func(chatId int64, messageId int, message string) *tgbotapi.MessageConfig) *Bot {
	if appEnv == "local" || appEnv == "dev" {
		return nil
	}

	tokenChatId := telegramConst.EnvTokenChatIdMap[appEnv]

	// 1. 建立 Telegram Bot
	var bot *tgbotapi.BotAPI
	var err error

	for retry := 1; retry <= 3; retry++ {
		bot, err = tgbotapi.NewBotAPI(tokenChatId.Token)
		if err == nil {
			break
		}
		log.Printf("NewBotAPI token:%s,error:%+v,retry:%d", tokenChatId.Token, err, retry)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		panic(err)
	}

	// 2. 設定 Bot 的 Debug Mode
	//bot.Debug = true

	// 3. 獲取 Bot 的資訊
	var botInfo tgbotapi.User

	for retry := 1; retry <= 3; retry++ {
		botInfo, err = bot.GetMe()
		if err == nil {
			break
		}
		log.Printf("GetMe token:%s,error:%+v,retry:%d", tokenChatId.Token, err, retry)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		panic(err)
	}

	log.Printf("Bot 姓名：%s / Bot 使用者名稱：%s", botInfo.FirstName, botInfo.UserName)

	// 4. 有設定 callback 才監聽訊息
	if nil != replyCallback {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates := bot.GetUpdatesChan(u)

		go func() {
			defer func() {
				r := recover()
				if r != nil {
					log.Printf("%+v", r)
					return
				}
			}()

			// 5. 處理收到的訊息
			for update := range updates {
				if update.Message == nil { // 確認訊息不為空
					continue
				}

				log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

				// 6. 回應訊息
				var replyConfig = replyCallback(update.Message.Chat.ID, update.Message.MessageID, update.Message.Text)
				_, sendErr := bot.Send(replyConfig)
				if sendErr != nil {
					log.Println(sendErr)
				}
			}
		}()
	}

	return &Bot{botAPI: bot, tokenChatId: &tokenChatId}
}

func (b *Bot) SendMessage(msg string) {
	_, err := b.botAPI.Send(tgbotapi.NewMessage(b.tokenChatId.ChatId, msg))
	if err != nil {
		log.Println(err)
	}
}

// ReplayCallbackDemo : 自訂機器人回應的範例
//
// - chatId 發送頻道 ID
//
// - messageId 對機器人問訊息 ID
//
// - message 對機器人問的訊息
func ReplayCallbackDemo(chatId int64, messageId int, message string) *tgbotapi.MessageConfig {
	var (
		FearunSwaggerBtn  = tgbotapi.NewInlineKeyboardButtonURL("開啟Swagger 1", "https://www.google.com/")
		SolastaSwaggerBtn = tgbotapi.NewInlineKeyboardButtonURL("開啟Swagger 2", "https://www.google.com/")
		FearunAPIBtn      = tgbotapi.NewInlineKeyboardButtonURL("開啟API文件 1", "https://www.google.com/")
		SolastaAPIBtn     = tgbotapi.NewInlineKeyboardButtonURL("開啟API文件 2", "https://www.google.com/")
	)

	var replyConfig = tgbotapi.NewMessage(chatId, "你好！感謝你的訊息。")
	if strings.Contains(strings.ToLower(message), "help") {
		replyConfig = tgbotapi.NewMessage(chatId, "可使用的指令如下\n/api API文件鏈結\n/swagger Swagger文件鏈結")
	}
	if strings.Contains(strings.ToLower(message), "swagger") {
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(FearunSwaggerBtn, SolastaSwaggerBtn),
			//tgbotapi.NewInlineKeyboardRow(SolastaSwaggerBtn),
		)
		replyConfig = tgbotapi.NewMessage(chatId, "以下是 Swagger 文件鏈結")
		replyConfig.ReplyMarkup = keyboard
	} else if strings.Contains(strings.ToLower(message), "api") {
		keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(FearunAPIBtn, SolastaAPIBtn))
		replyConfig = tgbotapi.NewMessage(chatId, "以下是 api 文件鏈結")
		replyConfig.ReplyMarkup = keyboard
	}
	//replyConfig.ReplyToMessageID = messageId

	return &replyConfig
}
