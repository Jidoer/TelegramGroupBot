package common

import (
	"encoding/json"
)

// RuleMap 群组规则的字典，匹配规则=>回复内容
type RuleMap map[string]string

var (
	// AllGroupRules 所有群组的规则字典
	AllGroupRules = make(map[int64]RuleMap)
	// AllGroupId 目前服务的所有群组的id
	AllGroupId []int64
)

func (rm RuleMap) String() string {
	s, err := json.Marshal(rm)
	if err != nil {
		return ""
	}
	return string(s)
}

// Json2kvs 将json字符串转化为规则字典
func Json2kvs(rulesJson string) RuleMap {
	tkvs := make(RuleMap)
	_ = json.Unmarshal([]byte(rulesJson), &tkvs)
	return tkvs
}

// AddNewGroup 在内存中添加新群组的条目
func AddNewGroup(gid int64) {
	AllGroupId = append(AllGroupId, gid)
	AllGroupRules[gid] = make(RuleMap)
}


/*
	var msg api.MessageConfig
	msg = api.NewMessage(gid, "")
	msg.Text = " Join...."
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true
	sendMessage(msg)
*/
