package telegram_bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// TelegramBot ..
type TelegramBot struct {
	*tgbotapi.BotAPI
	commands map[string]struct {
		i     int
		cmd   string
		intro string
		fn    func(*tgbotapi.Message) string
	}
	keyboard           tgbotapi.ReplyKeyboardMarkup
	inlineLinkKeyboard tgbotapi.InlineKeyboardMarkup

}

// NewTelegramBot ..
func NewTelegramBot(token, proxy string, debug bool) (*TelegramBot, error) {
	client := &http.Client{}
	if proxy != "" {
		url, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(url),
		}
	}

	bot, err := tgbotapi.NewBotAPIWithClient(token, client)
	if err != nil {
		return nil, err
	}
	bot.Debug = debug

	imp := &TelegramBot{
		bot,
		make(map[string]struct {
			i     int
			cmd   string
			intro string
			fn    func(*tgbotapi.Message) string
		}),
		tgbotapi.ReplyKeyboardMarkup{},
		tgbotapi.InlineKeyboardMarkup{},
	}

	//imp.RegisterCommand("links", "快捷跳转链接", func(msg *tgbotapi.Message) string {
	//	return fmt.Sprintf("%v", "快捷跳转链接")
	//})

	//imp.RegisterCommand("current_chat_id", "返回当前会话ID", func(msg *tgbotapi.Message) string {
	//	return fmt.Sprintf("%v", msg.Chat.ID)
	//})
	return imp, nil
}

// RegisterCommand ..
func (bot *TelegramBot) RegisterCommand(command string, intro string, fn func(*tgbotapi.Message) string) {

	bot.commands[command] = struct {
		i     int
		cmd   string
		intro string
		fn    func(*tgbotapi.Message) string
	}{i: len(bot.commands) + 1, cmd: command, intro: intro, fn: fn}
}

// RegisterKeyboard ..
func (bot *TelegramBot) RegisterKeyboard(lines [][]string) {
	var rows [][]tgbotapi.KeyboardButton

	for _, line := range lines {
		var buttons []tgbotapi.KeyboardButton
		for i := 0; i < len(line); i++ {
			buttons = append(buttons, tgbotapi.NewKeyboardButton(line[i]))
		}
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(buttons...))
	}
	rows = append(rows, []tgbotapi.KeyboardButton{tgbotapi.NewKeyboardButton("keyboard-close")})

	bot.keyboard = tgbotapi.NewReplyKeyboard(rows...)
}

// RegisterInlineLinkKeyboard ..
func (bot *TelegramBot) RegisterInlineLinkKeyboard(lines [][]struct{ Text, Link string }) {
	var rows [][]tgbotapi.InlineKeyboardButton
	for _, line := range lines {
		var row []tgbotapi.InlineKeyboardButton
		for i := 0; i < len(line); i++ {
			row = append(row, tgbotapi.NewInlineKeyboardButtonURL(line[i].Text, line[i].Link))
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(row...))
	}

	bot.inlineLinkKeyboard = tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// Broadcast ..
func (bot *TelegramBot) Broadcast(cid int64, context string) {
	msg := tgbotapi.NewMessage(cid, context)
	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}

func (bot *TelegramBot)IsAdmin(msg *tgbotapi.Message) bool {
	members, err := bot.GetChatAdministrators(tgbotapi.ChatConfig{ChatID: msg.Chat.ID})
	if err != nil {return false}
	for _, member := range members {
		if member.User.ID == msg.From.ID {
			return true
		}
	}
	return false
}

// Run ..
func (bot *TelegramBot) Run() error {
	defer func() {
		if ec := recover(); ec != nil {
			log.Fatalf("[TelegramBot]RUN EC:%v", ec)
		}
	}()

	bot.commands["commands"] = struct {
		i     int
		cmd   string
		intro string
		fn    func(*tgbotapi.Message) string
	}{i: 0, cmd: "commands", intro: "命令集",
		fn: func(*tgbotapi.Message) string {
			var menu string
			for i := 0; i < len(bot.commands); i++ {
				command, intro := func() (string, string) {
					for k, v := range bot.commands {
						if v.i == i {
							return k, v.intro
						}
					}
					return "", ""
				}()
				menu += fmt.Sprintf("/%v  %v\n", command, intro)
			}
			return menu
		}}

	cfg := tgbotapi.NewUpdate(0)
	cfg.Timeout = 60

	updates, err := bot.GetUpdatesChan(cfg)
	if err != nil {
		log.Println(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if command := update.Message.Command(); command != "" {

			if _, ok := bot.commands[command]; !ok {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown Command")
				if _, err := bot.Send(msg); err != nil {
					log.Println(err)
				}

			} else {
				test := bot.commands[command].fn(update.Message)
				if test != "" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, test)
					if command == "links" {
						msg.ReplyMarkup = bot.inlineLinkKeyboard
					}
					if _, err := bot.Send(msg); err != nil {
						log.Println(err)
					}
				}
			}
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		switch update.Message.Text {
		case "kb":
			msg.ReplyMarkup = bot.keyboard
		case "keyboard-close":
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		}

		if update.Message.Document != nil {
			fn := update.Message.Document.FileName
			if fn == "settings.json" {
				f, err := bot.GetFile(tgbotapi.FileConfig{FileID: update.Message.Document.FileID})
				if err != nil {
					bot.Broadcast(update.Message.Chat.ID, err.Error())
				}

				link := f.Link(bot.Token)
				resp, err := http.Get(link)
				if err != nil {
					bot.Broadcast(update.Message.Chat.ID, err.Error())
				}
				defer resp.Body.Close()

				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					bot.Broadcast(update.Message.Chat.ID, err.Error())
				}

				if err := ioutil.WriteFile("./settings/settings.json", data, 0644); err != nil {
					bot.Broadcast(update.Message.Chat.ID, err.Error())
				}

				bot.Broadcast(update.Message.Chat.ID, "点击{/LoadSettings}载入当前配置开始抽奖进程")
			}
		}
	}
	return nil
}
