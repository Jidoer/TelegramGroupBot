package main

import (
	"TelegramGroupBot/db"
	"log"
	"math/rand"
	"strconv"
	"time"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
)

/**
 * 发送文字消息
 */
func sendMessagedel(msg api.MessageConfig) api.Message {
	if msg.Text == "" {
		return api.Message{}
	}
	mmsg, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
	go deleteMessage(msg.ChatID, mmsg.MessageID)
	return mmsg
}
func sendMessagenodel(msg api.MessageConfig) api.Message {
	if msg.Text == "" {
		return api.Message{}
	}
	mmsg, err := bot.Send(msg)
	if err != nil {
		log.Println(err)
	}
	//go deleteMessage(msg.ChatID, mmsg.MessageID)
	return mmsg
}

/**
 * 发送图片消息, 需要是已经存在的图片链接
 */
func sendPhoto(chatid int64, photoid string) api.Message {
	file := api.NewPhotoShare(chatid, photoid)
	mmsg, err := bot.Send(file)
	if err != nil {
		log.Println(err)
	}
	go deleteMessage(chatid, mmsg.MessageID)
	return mmsg
}

/**
 * 发送动图, 需要是已经存在的链接
 */
func sendGif(chatid int64, gifid string) api.Message {
	file := api.NewAnimationShare(chatid, gifid)
	mmsg, err := bot.Send(file)
	if err != nil {
		log.Println(err)
	}
	go deleteMessage(chatid, mmsg.MessageID)
	return mmsg
}

/**
 * 发送视频, 需要是已经存在的视频连接
 */
func sendVideo(chatid int64, videoid string) api.Message {
	file := api.NewVideoShare(chatid, videoid)
	mmsg, err := bot.Send(file)
	if err != nil {
		log.Println(err)
	}
	go deleteMessage(chatid, mmsg.MessageID)
	return mmsg
}

/**
 * 发送文件, 必须是已经存在的文件链接
 */
func sendFile(chatid int64, fileid string) api.Message {
	file := api.NewDocumentShare(chatid, fileid)
	mmsg, err := bot.Send(file)
	if err != nil {
		log.Println(err)
	}
	go deleteMessage(chatid, mmsg.MessageID)
	return mmsg
}

func deleteMessage(gid int64, mid int) {
	//新线程内执行
	time.Sleep(time.Second * 240)
	_, _ = bot.DeleteMessage(api.NewDeleteMessage(gid, mid))
}

func ChuTi(Messg *api.Message) {
	rand.Seed(time.Now().UnixNano())
	i1 := rand.Intn(30)
	rand.Seed(time.Now().UnixNano())
	i2 := rand.Intn(20) + 10
	if !db.AddCKpeople(Messg.Chat.ID, Messg.From.ID, strconv.Itoa(i1+i2)) {
		log.Println("ERROR:添加验证操作失败！")
		return
	}
	var msg api.MessageConfig
	msg = api.NewMessage(Messg.Chat.ID, "")
	msg.Text = "<i>请回答题目用于验证</i> "+"[" + Messg.From.String() + "](tg://user?id=" + strconv.Itoa(Messg.From.ID) + ")"  +
		"\r\n请180秒内完成，否则会删除拉黑" +
		"\r\n <b>" + strconv.Itoa(i1) + "+" + strconv.Itoa(i2) + "= ?</b>"
	msg.ParseMode = "Markdown" 
	msg.DisableWebPagePreview = false
	sendMessagedel(msg)

}

func PeopleCKdel(gid int64, uid int, Messg *api.Message) {
	time.Sleep(time.Second * 180)
	if db.IfPeopleck(gid, uid) {
		//banMember(gid, uid, -1)
		var msg api.MessageConfig
		msg = api.NewMessage(Messg.Chat.ID, "")
		msg.Text = "<i>Ban!</i> @" +  "[" + Messg.From.String() + "](tg://user?id=" + strconv.Itoa(Messg.From.ID) + ")" +
			"\n验证失败被移除群聊!"
		msg.ParseMode = "Markdown" //HTML
		msg.DisableWebPagePreview = false
		sendMessagenodel(msg)
		kickMember(gid, uid)
		log.Println(strconv.Itoa(uid) + " 被ban!!!")
	} else {
		log.Println(strconv.Itoa(uid) + " 挺过180s!!!")
	}

}
