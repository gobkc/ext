package gext

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

func Encrypt(codeData string, saltKey string) string {
	dataArr := []rune(codeData)
	keyArr := []byte(saltKey)
	keyLen := len(keyArr)

	var tmpList []int

	for index, value := range dataArr {
		base := int(value)
		dataString := base + int(0xFF&keyArr[index%keyLen])
		tmpList = append(tmpList, dataString)
	}

	var str string

	for _, value := range tmpList {
		str += "_" + fmt.Sprintf("%d", value)
	}
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Decrypt(ntData string, saltKey string) string {
	decodeStr, err := base64.StdEncoding.DecodeString(ntData)
	if err != nil {
		return ""
	}
	ntData = string(decodeStr)
	strLen := len(ntData)
	newData := []rune(ntData)
	resultData := string(newData[1:strLen])
	dataArr := strings.Split(resultData, "_")
	keyArr := []byte(saltKey)
	keyLen := len(keyArr)

	var tmpList []int

	for index, value := range dataArr {
		base, _ := strconv.Atoi(value)
		dataString := base - int(0xFF&keyArr[index%keyLen])
		tmpList = append(tmpList, dataString)
	}

	var str string

	for _, val := range tmpList {
		str += string(rune(val))
	}
	return str
}

func Random(length int) (str string) {
	if length == 0 {
		return
	}
	var (
		randByte  = make([]byte, length)
		formatStr []string
		outPut    []interface{}
		byteHalf  uint8 = 127
	)
	rand.Read(randByte)
	for _, b := range randByte {
		if b > byteHalf {
			formatStr = append(formatStr, "%X")
		} else {
			formatStr = append(formatStr, "%x")
		}
		outPut = append(outPut, b)
	}
	if str = fmt.Sprintf(strings.Join(formatStr, ""), outPut...); len(str) > length {
		str = str[:length]
	}
	return
}
