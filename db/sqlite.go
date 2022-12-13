package db

import (
	"TelegramGroupBot/common"
	"log"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // 初始化gorm使用sqlite
)

var db *gorm.DB

type setting struct {
	gorm.Model
	Key   string `gorm:"unique;not null"`
	Value string
}

type rule struct {
	gorm.Model
	GroupId  int64 `gorm:"unique;not null"`
	RuleJson string
}

type peopleck struct {
	gorm.Model
	GroupId int64 //`gorm:"unique;not null"`
	Uid     int
	Answer  string
}

type BanAd struct {
	gorm.Model
	GroupId int64
	UserID  int
	Num     int
}

/*
//检查队列
type cking struct {
	id     int
	Gid    int64
	Uid    int
	Answer string
}
*/

// Init 数据库初始化，包括新建数据库（如果还没有建立），基本数据的读写
func Init(newToken string) (token string) {
	dbtmp, err := gorm.Open("sqlite3", "data.db")
	if err != nil {
		panic("failed to connect database")
	}
	db = dbtmp
	db.AutoMigrate(&setting{}, &rule{}, &peopleck{}, &BanAd{}) //自动初始化表
	var tokenSetting setting
	db.Find(&tokenSetting, "Key=?", "token")
	token = tokenSetting.Value
	if newToken != "" {
		token = newToken
		if tokenSetting.ID > 0 {
			tokenSetting.Value = newToken
			db.Model(&tokenSetting).Update(tokenSetting)
		} else {
			db.Create(&setting{
				Key:   "token",
				Value: newToken,
			})
		}
	}
	readAllGroupRules()
	return
}

// AddNewGroup 数据库中添加一条记录来记录新群组的规则
func AddNewGroup(groupId int64) {
	db.Create(&rule{
		GroupId:  groupId,
		RuleJson: "",
	})
}

// UpdateGroupRule 更新群组的规则
func UpdateGroupRule(groupId int64, ruleJson string) {
	db.Model(&rule{}).Where("group_id=?", groupId).Update("rule_json", ruleJson)
}

func readAllGroupRules() {
	var allGroupRules []rule
	db.Find(&allGroupRules)
	for _, rule := range allGroupRules {
		ruleStruct := common.Json2kvs(rule.RuleJson)
		common.AllGroupRules[rule.GroupId] = ruleStruct
		common.AllGroupId = append(common.AllGroupId, rule.GroupId)
	}
}

func AddCKpeople(gid int64, uid int, Answer string) bool {

	//db.Model(&peopleck{}).Create("")
	if err := db.Create(&peopleck{GroupId: gid, Uid: uid, Answer: Answer}).Error; err != nil {
		//error
		return false
	}
	return true
}

func CKpeopleProgress(gid int64, uid int, Answer string) bool {
	//db.Model(&peopleck{}).Delete("group_id=?", gid)
	rows, _ := db.Model(&peopleck{}).Where("group_id = ?", gid).Select("id, group_id, uid, answer").Rows() // (*sql.Rows, error)
	//defer rows.Close()
	var cking peopleck
	i := 0
	for rows.Next() {
		//在群内查找她
		db.ScanRows(rows, &cking)
		log.Println(cking)
		if cking.Uid == uid {
			rows.Close()
			break
		}
		i++
	}
	rows.Close()

	if cking.Answer == Answer {
		e := db.Where("id = ?", cking.ID).Unscoped().Delete(&peopleck{}).Error
		if e != nil {
			log.Println(strconv.Itoa(cking.Uid) + ": 验证Error!")
		}
		log.Println(strconv.Itoa(cking.Uid) + ": 验证成功!")
		return true
	} else {
		return false
	}
}

func IfPeopleck(gid int64, uid int) bool {
	log.Println("IfPeopleck")
	ifhave := false
	rows, _ := db.Model(&peopleck{}).Where("group_id = ?", gid).Select("id, group_id, uid, answer").Rows() // (*sql.Rows, error)
	defer rows.Close()
	for rows.Next() {
		var cking peopleck
		db.ScanRows(rows, &cking)
		log.Println(cking)
		if cking.Uid == uid {
			ifhave = true
		}
		// do something
	}
	log.Println("IfPeopleck END")
	return ifhave
}

func AddADBan(gid int64, uid int) int {
	rows, _ := db.Model(&BanAd{}).Where("group_id = ?", gid).Rows() //.Select("id, group_id, uid, answer").Rows() // (*sql.Rows, error)
	defer rows.Close()
	var bid BanAd
	have := false
	for rows.Next() {
		//在群内查找她
		db.ScanRows(rows, &bid)
		if bid.UserID == uid {
			have = true
			rows.Close()
			break
		}
	}
	if !have {
		if err := db.Create(&BanAd{GroupId: gid, UserID: uid, Num: 1}).Error; err != nil {
			return -1
		}
		return 1
	}
	if err := db.Model(&BanAd{}).Update(&BanAd{Num: bid.Num + 1}).Error; err != nil {
		return -1
	}
	return bid.Num + 1
}
