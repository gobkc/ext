package gext

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
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

func GzipEncode(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	gzipWriter, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	_, err := gzipWriter.Write(input)
	if err != nil {
		_ = gzipWriter.Close()
		return nil, err
	}
	if err = gzipWriter.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GzipDecode(input []byte) ([]byte, error) {
	bytesReader := bytes.NewReader(input)
	gzipReader, err := gzip.NewReader(bytesReader)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = gzipReader.Close()
	}()
	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(gzipReader); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func MarshalGzipJson(data interface{}) ([]byte, error) {
	marshalData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	gzipData, err := GzipEncode(marshalData)
	if err != nil {
		return nil, err
	}
	return gzipData, err
}

func UnmarshalGzipJson(input []byte, output interface{}) error {
	decodeData, err := GzipDecode(input)
	if err != nil {
		return err
	}

	err = json.Unmarshal(decodeData, output)
	if err != nil {
		return err
	}
	return nil
}
