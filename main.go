package main

import (
	"NewYearLuckyDraw/business"
	"NewYearLuckyDraw/libray/logger"
	"NewYearLuckyDraw/libray/storage"
	"NewYearLuckyDraw/libray/telegram_bot"
	"NewYearLuckyDraw/model"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 5095242734:AAHf8kwOBKSmWXnsW--YZxu7jLsInVSnycQ


var Bot *LuckyDrawBot

// LuckyDrawBot ..
type LuckyDrawBot struct {
	*telegram_bot.TelegramBot
	groups map[int64]business.LuckyDraw
}


func init(){
	tg, err := telegram_bot.NewTelegramBot("5095242734:AAHf8kwOBKSmWXnsW--YZxu7jLsInVSnycQ", "", false)
	if err != nil {
		logger.Fatal(err)
		return
	}

	Bot = &LuckyDrawBot{
		tg,make(map[int64]business.LuckyDraw),
	}
	Bot.RegisterCommands()

	logger.Info("bot init successful")
}

func main()  {
	if err := Bot.Run(); err != nil {
		logger.Fatal(err)
	}
}

func (bot *LuckyDrawBot)RegisterCommands() {

	bot.RegisterCommand("init", "初始化", func(msg *tgbotapi.Message) string {
		var result = "[提示]\n"

		if admin := bot.IsAdmin(msg); !admin {
			return "403"
		}
		cid := msg.Chat.ID

		count, err := bot.GetChatMembersCount(tgbotapi.ChatConfig{ChatID: cid})
		if err != nil {
			return err.Error()
		}

		tbNameUser := fmt.Sprintf("%v_%v", model.UserTable, cid)
		if exist := storage.Database.TableExists(tbNameUser); exist {
			storage.Database.DropTable(tbNameUser)
		}

		tbNameRecord := fmt.Sprintf("%v_%v", model.RecordTable, cid)
		if exist := storage.Database.TableExists(tbNameRecord); exist {
			storage.Database.DropTable(tbNameRecord)
		}

		handler := business.NewLuckyDraw(model.Settings{
			ID: fmt.Sprintf("%v", cid),
			LotteryUserNumber: count-1,
		})
		bot.groups[cid] = handler

		result += fmt.Sprintf("[SUCC]GroupChatID[%v] AdminID[%v] ThisGroupMembersCount[%v]\n", msg.Chat.ID, msg.From.ID, count-1)
		return result
	})

	bot.RegisterCommand("set_lottery_user_number", "设置抽奖人数 {人数}", func(msg *tgbotapi.Message) string {
		var result = "[提示]\n"

		if admin := bot.IsAdmin(msg); !admin {
			return "403"
		}
		cid := msg.Chat.ID

		args := strings.Fields(msg.CommandArguments())
		if len(args) < 1 {
			return "400"
		}

		n, err := strconv.Atoi(args[0])
		if err != nil {
			return err.Error()
		}

		handler := bot.groups[cid]
		handler.SetLotteryUserNumber(n)
		bot.groups[cid] = handler

		result += fmt.Sprintf("[SUCC]%+v", handler.Settings())
		return result
	})

	bot.RegisterCommand("prize_list", "奖品清单", func(msg *tgbotapi.Message) string {
		var result = "[提示]\n"
		if admin := bot.IsAdmin(msg); !admin {
			return "403"
		}
		cid := msg.Chat.ID
		if _, ok := bot.groups[cid]; !ok {
			return "404"
		}

		handler := bot.groups[cid]
		prizes, err := handler.QueryPrizes(true)
		if err != nil {
			return err.Error()
		}
		sort.Sort(prizes)

		var stats = make(map[string]int)
		for _, prize := range prizes {
			stats[prize.Title]++
		}
		var named = make(map[string]int)
		for _, prize := range prizes {
			if _, ok := named[prize.Title]; !ok {
				result += fmt.Sprintf("%v\t*%v\n", prize.Title, stats[prize.Title])
				named[prize.Title]++
			}
		}
		return result
	})

	bot.RegisterCommand("append_prize", "追加奖品 {奖品名称} {奖品封面URL} {奖品数量}", func(msg *tgbotapi.Message) string {
		var result = "[提示]\n"
		if admin := bot.IsAdmin(msg); !admin {
			return "403"
		}
		cid := msg.Chat.ID
		if _, ok := bot.groups[cid]; !ok {
			return "404"
		}

		var (
			title, cover, count string
		)
		args := strings.Fields(msg.CommandArguments())
		if len(args) < 1 {
			return "400"
		}
		title = args[0]
		if len(args) > 1 {
			cover = args[1]
		}
		if len(args) > 2 {
			count = args[2]
		} else {
			count = "1"
		}

		n, err := strconv.Atoi(count)
		if err != nil {
			return err.Error()
		}

		handler := bot.groups[cid]
		if err := handler.AppendPrize(title, cover, n); err != nil {
			return err.Error()
		}

		result += fmt.Sprintf("[SUCC]Title[%v]Cover[%v]Count[%v]", title, cover, n)
		return result
	})

	bot.RegisterCommand("delete_prize", "删除奖品 {奖品ID}", func(msg *tgbotapi.Message) string {
		var result = "[提示]\n"
		if admin := bot.IsAdmin(msg); !admin {
			return "403"
		}
		cid := msg.Chat.ID
		if _, ok := bot.groups[cid]; !ok {
			return "404"
		}

		args := strings.Fields(msg.CommandArguments())
		if len(args) < 1 {
			return "400"
		}
		n, err := strconv.Atoi(args[0])
		if err != nil {
			return err.Error()
		}

		handler := bot.groups[cid]
		if err := handler.DeletePrize(n); err != nil {
			return err.Error()
		}

		result += fmt.Sprintf("[SUCC]PrizeID[%v]", n)
		return result
	})

	bot.RegisterCommand("remove_prize", "清空奖池", func(msg *tgbotapi.Message) string {
		var result = "[提示]\n"
		if admin := bot.IsAdmin(msg); !admin {
			return "403"
		}
		cid := msg.Chat.ID
		if _, ok := bot.groups[cid]; !ok {
			return "404"
		}
		tbNamePrize := fmt.Sprintf("%v_%v", model.PrizeTable, cid)
		if exist := storage.Database.TableExists(tbNamePrize); exist {
			storage.Database.DropTable(tbNamePrize)
		}
		result += "[SUCC]清空奖池"
		return result
	})

	bot.RegisterCommand("user_list", "用户列表", func(msg *tgbotapi.Message) string {
		var result = "[用户列表]\n"
		if admin := bot.IsAdmin(msg); !admin {
			return "403"
		}
		cid := msg.Chat.ID
		if _, ok := bot.groups[cid]; !ok {
			return "404"
		}

		handler := bot.groups[cid]

		users, err := handler.QueryUsers(true)
		if err != nil {
			return err.Error()
		}
		for _, user := range users {
			exist, err := handler.ExistUserRecord(user.ID)
			if err != nil {
				return err.Error()
			}
			if exist {
				result += fmt.Sprintf("[%v]Name{%v} [已抽奖]\n", user.ID, user.Name)
			} else {
				result += fmt.Sprintf("[%v]Name{%v}\n", user.ID, user.Name)
			}
		}

		r, err  := handler.QueryUsers(false)
		if err != nil {
			return err.Error()
		}

		result += fmt.Sprintf("[SUCC]当前参与数量[%v] 预设参与数量[%v] 剩余未抽奖用户数量[%v]", len(users), handler.Settings().LotteryUserNumber, len(r))
		return result
	})

	bot.RegisterCommand("register", "报名抽奖", func(msg *tgbotapi.Message) string {
		cid := msg.Chat.ID
		if _, ok := bot.groups[cid]; !ok {
			return "404"
		}
		handler := bot.groups[cid]

		if handler.Started() {
			return "[提示]抽奖已开始暂停报名"
		}

		userName := msg.From.UserName
		if err := handler.AppendUser(userName); err != nil {
			return err.Error()
		}
		count, err := handler.CountUser()
		if err != nil {
			return err.Error()
		}

		if count == handler.Settings().LotteryUserNumber {
			time.Sleep(time.Second * 1)
			handler.Start()
			bot.groups[cid] = handler
			bot.Broadcast(cid, "[提示]报名结束")
		}

		return fmt.Sprintf("[提示][%v]已绑定成功, 当前参与数量[%v / %v]", userName, count, handler.Settings().LotteryUserNumber)
	})

	bot.RegisterCommand("start", "开始进程", func(msg *tgbotapi.Message) string {
		var result = "[提示]\n"
		if admin := bot.IsAdmin(msg); !admin {
			return "403"
		}
		cid := msg.Chat.ID
		if _, ok := bot.groups[cid]; !ok {
			return "404"
		}

		handler := bot.groups[cid]

		go func() {
			for {

				if handler.Started() {
					return
				}

				count, err := handler.CountUser()
				if err != nil {
					logger.Error(err)
					return
				}

				if handler.Settings().LotteryUserNumber-count > 0 {
					bot.Broadcast(cid,
						fmt.Sprintf("[提示]预计[%v]位参与抽奖 当前已报名[%v]位用户 剩余[%v]位用户请点击{/register}进行报名",
						handler.Settings().LotteryUserNumber, count, handler.Settings().LotteryUserNumber-count))

				} else {
					bot.Broadcast(cid, "[提示]报名已完成, 请管理员手动设置当前抽奖用户 \\/set_current_lottery_user @用户名")
					return
				}

				time.Sleep(time.Second * 30)
			}
		}()

		result += "[SUCC]开始报名"
		return result
	})


	bot.RegisterCommand("set_current_lottery_user", "设置当前抽奖用户 {@UserName}", func(msg *tgbotapi.Message) string {

		if admin := bot.IsAdmin(msg); !admin {
			return "403"
		}
		cid := msg.Chat.ID
		if _, ok := bot.groups[cid]; !ok {
			return "404"
		}

		args := strings.Fields(msg.CommandArguments())
		if len(args) < 1 {
			return "400"
		}

		userName := args[0][1:]

		handler := bot.groups[cid]
		user, err := handler.QueryByUserName(userName)
		if err != nil {
			return err.Error()
		}

		uid := user.ID
		if err := handler.SetRound(uid); err != nil {
			return err.Error()
		}
		handler.SetRound(uid)
		bot.groups[cid] = handler

		go func() {
			time.Sleep(time.Second * 1)
			bot.Broadcast(cid, fmt.Sprintf("@%v 请点击➡️{/lottery}⬅️开始抽奖", user.Name))
		}()

		result  := fmt.Sprintf("[提示][SUCC]当前抽奖用户[%v]", user.Name)
		return result
	})

	bot.RegisterCommand("lottery", "点击抽奖", func(msg *tgbotapi.Message) string {
		cid := msg.Chat.ID
		if _, ok := bot.groups[cid]; !ok {
			return "404"
		}
		handler := bot.groups[cid]

		userName := msg.From.UserName

		user, err := handler.QueryByUserName(userName)
		if err != nil {
			return err.Error()
		}

		record, err := handler.Lottery(user.ID)
		if err != nil {
			return err.Error()
		}

		prize, err := handler.QueryPrizeByID(record.PID)
		if err != nil {
			return err.Error()
		}


		for i:=0;i<3;i++{
			bot.Broadcast(cid, fmt.Sprintf("开奖中....%v", 3-i))
			time.Sleep(time.Second * 1)
		}

		handler.SetRound(0)
		bot.groups[cid] = handler


		r, err  := handler.QueryUsers(false)
		if err != nil {
			return err.Error()
		}
		if len(r) == 0 {
			go func() {
				time.Sleep(time.Second * 2)
				bot.Broadcast(cid, "所有人已完成抽奖!")
				report, _ := handler.Report()
				bot.Broadcast(cid, fmt.Sprintf("[本次抽奖活动获奖名单]\n%v", report))
			}()

		} else {
			go func() {
				time.Sleep(time.Second * 2)
				var result = "[提示]剩余奖品清单:\n"
				prizes, err := handler.QueryPrizes(true)
				if err != nil {
					bot.Broadcast(cid, err.Error())
					return
				}
				sort.Sort(prizes)

				records, err := handler.QueryRecords()
				if err != nil {
					bot.Broadcast(cid, err.Error())
					return
				}

				var stats = make(map[string]int)
				for _, prize := range prizes {
					var used bool
					for _, record := range records {
						if record.PID == prize.ID {
							used = true
							break
						}
					}
					if !used {
						stats[prize.Title]++
					}
				}

				var named = make(map[string]int)
				for _, prize := range prizes {
					if _, ok := named[prize.Title]; !ok {
						count, ok := stats[prize.Title]
						if ok {
							result += fmt.Sprintf("%v\t*%v\n", prize.Title, count)
						}
						named[prize.Title]++
					}
				}
				bot.Broadcast(cid, result)
			}()
		}

		if prize.Cover != "" {
			photo := tgbotapi.NewInputMediaPhoto(prize.Cover)
			photo.Caption = prize.Title
			_, err := bot.Send(tgbotapi.NewMediaGroup(cid, []interface{}{photo}))
			if err != nil {
				logger.ErrorF("PrizeCover[%v] error: %v", prize.Cover, err)
			}
		}
		return fmt.Sprintf("[抽奖结果]恭喜[@%v]获得奖品[%v]", user.Name, prize.Title)
	})
}