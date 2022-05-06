package main

import (
	"TelegramGroupBot/common"
	"TelegramGroupBot/db"
	"flag"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/robfig/cron"
)

var bot *api.BotAPI
var gcron *cron.Cron

var (
	debug       bool
	superUserId int
)

func main() {
	botToken := flag.String("t", "", "your bot Token")
	flag.IntVar(&superUserId, "s", 0, "super manager Id")
	flag.BoolVar(&debug, "d", false, "debug mode")
	flag.Parse()
	token := db.Init(*botToken)
	gcron = cron.New()
	gcron.Start()
	//开始工作
	start(token)

}
func st() {
	//Just Test
	i := 0
	for {
		db.AddCKpeople(int64(i), i, "Test")
		i++
	}
}
func start(botToken string) {
	var err error
	bot, err = api.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = debug
	log.Printf("Authorized on account: %s  ID: %d", bot.Self.UserName, bot.Self.ID)

	u := api.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic("Can't get Updates")
	}
	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			//non-Message
			log.Println("non-Message")
			continue
		}
		UpMessage := update.Message
		log.Printf("用户gid:%d用户uid:%dtext:%s", UpMessage.Chat.ID, UpMessage.From.ID, UpMessage.Text)

		//Update.....
		go processUpdate(&update)
	}
}

/**
 * 对于每一个update的单独处理
 */
func processUpdate(update *api.Update) {
	upmsg := update.Message
	gid := upmsg.Chat.ID
	uid := upmsg.From.ID
	//检查是不是新加的群或者新来的人
	in := checkInGroup(gid)
	if !in {

		//不在就需要加入, 内存中加一份, 数据库中添加一条空规则记录 进入新群
		common.AddNewGroup(gid)
		db.AddNewGroup(gid)
	}

	if upmsg.IsCommand() {
		go processCommond(update)
	} else {
		go processReplyCommond(update)
		go processReply(update)

		/*
			var msg api.MessageConfig
			msg = api.NewMessage(gid, "")
			msg.Text = " Join...."
			msg.ParseMode = "Markdown"
			msg.DisableWebPagePreview = true
			sendMessage(msg)
		*/
		//sendMessagedel(api.NewMessage(gid, upmsg.Text))

		if db.IfPeopleck(gid, uid) {
			if db.CKpeopleProgress(gid, uid, upmsg.Text) {
				sendMessagedel(api.NewMessage(gid, "欢迎欢迎! @"+upmsg.From.UserName))
			} else {
				sendMessagedel(api.NewMessage(gid, "回答错误! @"+upmsg.From.UserName))
			}
			_, _ = bot.DeleteMessage(api.NewDeleteMessage(gid, upmsg.MessageID))

		}

		// 新人入群 新用户通过用户名检查是否是清真
		if upmsg.NewChatMembers != nil {
			for _, auser := range *(upmsg.NewChatMembers) {
				if checkQingzhen(auser.UserName) ||
					checkQingzhen(auser.FirstName) ||
					checkQingzhen(auser.LastName) {
					banMember(gid, uid, -1)
				} else {
					if !db.IfPeopleck(gid, uid) {
						//Don't have so add it
						ChuTi(upmsg)
						//New process
						go PeopleCKdel(gid, uid, upmsg)
					}
					//db.AddCKpeople(gid,uid)
					/*
						log.Println("NewPeople:" + auser.UserName)
						sendMessagedel(api.NewMessage(gid,"欢迎新人: @" + auser.UserName))
					*/
					//sendPhoto(gid,"https://t.me/c/1472018167/53095")
				}
			}

		}

		//检查清真并剔除
		if checkQingzhen(upmsg.Text) {
			_, _ = bot.DeleteMessage(api.NewDeleteMessage(gid, upmsg.MessageID))
			banMember(gid, uid, -1)
		}
	}
}

func processReply(update *api.Update) {
	var msg api.MessageConfig
	upmsg := update.Message
	gid := upmsg.Chat.ID
	uid := upmsg.From.ID
	replyText := findKey(gid, upmsg.Text)
	if replyText == "delete" {
		_, _ = bot.DeleteMessage(api.NewDeleteMessage(gid, upmsg.MessageID))
		num := db.AddADBan(gid, uid)
		if num != -1 {
			if num >= 3 { //再一再二不能再三
				msg = api.NewMessage(gid, "用户:"+upmsg.From.UserName+"\r\n多次发广告被踢除群聊!")
				msg.DisableWebPagePreview = true
				sendMessagenodel(msg)
				kickMember(gid, uid)
			}
			return
		}
		//error

	} else if strings.HasPrefix(replyText, "ban") {
		_, _ = bot.DeleteMessage(api.NewDeleteMessage(gid, upmsg.MessageID))
		banMember(gid, uid, -1)
	} else if strings.HasPrefix(replyText, "kick") {
		_, _ = bot.DeleteMessage(api.NewDeleteMessage(gid, upmsg.MessageID))
		kickMember(gid, uid)
	} else if strings.HasPrefix(replyText, "photo:") {
		sendPhoto(gid, replyText[6:])
	} else if strings.HasPrefix(replyText, "gif:") {
		sendGif(gid, replyText[4:])
	} else if strings.HasPrefix(replyText, "video:") {
		sendVideo(gid, replyText[6:])
	} else if strings.HasPrefix(replyText, "file:") {
		sendFile(gid, replyText[5:])
	} else if replyText != "" {
		msg = api.NewMessage(gid, replyText)
		msg.DisableWebPagePreview = true
		msg.ReplyToMessageID = upmsg.MessageID
		sendMessagedel(msg)
	}
}

func processCommond(update *api.Update) {
	var msg api.MessageConfig
	upmsg := update.Message
	gid := upmsg.Chat.ID
	uid := upmsg.From.ID
	msg = api.NewMessage(update.Message.Chat.ID, "")
	_, _ = bot.DeleteMessage(api.NewDeleteMessage(update.Message.Chat.ID, upmsg.MessageID))
	switch upmsg.Command() {
	case "start", "help", "about":
		msg.Text = "TG群组机器人" +
			"\r\n/me 查看个人信息" +
			"\r\n/banme 禁言某人" +
			"\r\n/add 添加规则" +
			"\r\n/del 删除规则" +
			"\r\n/list 列出规则" +
			"\r\n机器人作者: @JiCode"
		//sendMessagedel(msg)
		sendMessagenodel(msg)
	case "add":
		if checkAdmin(gid, *upmsg.From) {
			order := upmsg.CommandArguments()
			if order != "" {
				addRule(gid, order)
				msg.Text = "规则添加成功: " + order
			} else {
				msg.Text = addText
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
			}
			sendMessagedel(msg)
		}
	case "del":
		if checkAdmin(gid, *upmsg.From) {
			order := upmsg.CommandArguments()
			if order != "" {
				delRule(gid, order)
				msg.Text = "规则删除成功: " + order
			} else {
				msg.Text = delText
				msg.ParseMode = "Markdown"
			}
			sendMessagedel(msg)
		}
	case "list":
		if checkAdmin(gid, *upmsg.From) {
			rulelists := getRuleList(gid)
			msg.Text = "ID: " + strconv.FormatInt(gid, 10)
			msg.ParseMode = "Markdown"
			msg.DisableWebPagePreview = true
			sendMessagedel(msg)
			for _, rlist := range rulelists {
				msg.Text = rlist
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
				sendMessagedel(msg)
			}
		}
	case "admin":
		msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.Itoa(uid) + ") 请求管理员出来打屁股\r\n\r\n" + getAdmins(gid)
		msg.ParseMode = "Markdown"
		sendMessagedel(msg)
		banMember(gid, uid, 30)
	case "banme":
		botme, _ := bot.GetChatMember(api.ChatConfigWithUser{ChatID: gid, UserID: bot.Self.ID})
		if botme.CanRestrictMembers {
			rand.Seed(time.Now().UnixNano())
			sec := rand.Intn(540) + 60
			banMember(gid, uid, int64(sec))
			msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.Itoa(upmsg.From.ID) + ")被禁言" + strconv.Itoa(sec) + "秒"
			msg.ParseMode = "Markdown"
		} else {
			msg.Text = "请给bot禁言权限"
		}
		sendMessagedel(msg)
	case "me":
		myuser := upmsg.From
		msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.Itoa(upmsg.From.ID) + ") 的账号信息" +
			"\r\nID: " + strconv.Itoa(uid) +
			"\r\nUseName: [" + upmsg.From.String() + "](tg://user?id=" + strconv.Itoa(upmsg.From.ID) + ")" +
			"\r\nLastName: " + myuser.LastName +
			"\r\nFirstName: " + myuser.FirstName +
			"\r\nIsBot: " + strconv.FormatBool(myuser.IsBot)
		msg.ParseMode = "Markdown"
		sendMessagedel(msg)
	default:
	}
}

func processReplyCommond(update *api.Update) {
	var msg api.MessageConfig
	upmsg := update.Message
	gid := upmsg.Chat.ID
	//回复类型的管理命令
	if upmsg.ReplyToMessage != nil {
		reolyToUserId := upmsg.ReplyToMessage.From.ID
		switch upmsg.Text {
		case "ban":
			if checkAdmin(gid, *upmsg.From) {
				banMember(gid, reolyToUserId, -1)
				mem, _ := bot.GetChatMember(api.ChatConfigWithUser{ChatID: gid, SuperGroupUsername: "", UserID: reolyToUserId})
				if !mem.CanSendMessages {
					msg = api.NewMessage(gid, "")
					msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.Itoa(upmsg.From.ID) + ") 禁言了 " +
						"[" + upmsg.ReplyToMessage.From.String() + "](tg://user?id=" + strconv.Itoa(reolyToUserId) + ") "
					msg.ParseMode = "Markdown"
					sendMessagedel(msg)
				}
			}
		case "unban":
			if checkAdmin(gid, *upmsg.From) {
				unbanMember(gid, reolyToUserId)
				//mem,_ := bot.GetChatMember(api.ChatConfigWithUser{gid, "", reolyToUserId})
				//
				msg = api.NewMessage(gid, "")
				msg.Text = "[" + upmsg.From.String() + "](tg://user?id=" + strconv.Itoa(upmsg.From.ID) + ") 解禁了 " +
					"[" + upmsg.ReplyToMessage.From.String() + "](tg://user?id=" + strconv.Itoa(reolyToUserId) + ") "
				msg.ParseMode = "Markdown"
				sendMessagedel(msg)
			}
		case "kick":
			if checkAdmin(gid, *upmsg.From) {
				kickMember(gid, reolyToUserId)
			}
		case "unkick":
			if checkAdmin(gid, *upmsg.From) {
				unkickMember(gid, reolyToUserId)
			}
		default:
		}
	}
}
