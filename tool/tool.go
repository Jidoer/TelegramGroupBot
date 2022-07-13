package tool

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"unicode"
	//"golang.org/x/text/encoding/simplifiedchinese"
	//"golang.org/x/text/transform"
)


func interface2String(inter interface{}) string {

	switch inter.(type) {

	case string:
		// rt.Println("string", inter.(string))
		return inter.(string)
	case int:
		fmt.Println("int", inter.(int))
		return ""
	case float64:
		fmt.Println("float64", inter.(float64))
		return ""
	}
	return ""

}

func Isnumber(str string) bool {
	for _, x := range []rune(str) {
		if !unicode.IsDigit(x) {
			return false
		}
	}
	return true
}

func String2Int(str5 string) int {
	int5, err := strconv.Atoi(str5)
	if err != nil {
		fmt.Println(err)
		return 100000000 //error
	} else {

		return int5
	}
}

func MapToJson(param map[string]map[string]string/*interface{}*/) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

func InterfaceToJson(param interface{}) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}


func JsonToMap(str string) map[string]interface{} {
	var tempMap map[string]interface{}
	err := json.Unmarshal([]byte(str), &tempMap)
	if err != nil {
		panic(err)
	}
	return tempMap
}


/*
func Utf8ToGBK(utf8str string) string {
    result, _, _ := transform.String(simplifiedchinese.GBK.NewEncoder(), utf8str)
    return result
}*/


func URLCode(yoururl string) string{
	return url.QueryEscape(yoururl)
}
func UnURLCode(yoururl string) string{
	decodeurl,err := url.QueryUnescape(yoururl)
	if err != nil {
		fmt.Println(err)
	}
	return decodeurl
}


func StringRepalceALL(str string,old string, new string) string{
//替换两次
//fmt.Println(strings.Replace("oink oink oink", "k", "ky", 2)) 
//全部替换
return strings.Replace(str, old, new, -1)
}

func Md5(s string) string {
    h := md5.New()
    h.Write([]byte(s))
    return hex.EncodeToString(h.Sum(nil))
}