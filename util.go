package cacheE

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
	"strconv"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func booltoint(bo bool) int32 {
	if bo {
		return 1
	} else {
		return 0
	}
}

func intTobool(b int32) bool {
	if b == 0 {
		return false
	} else {
		return true
	}
}

// byte 转化 int
func bytetoint(by []byte, num *int32) {
	b_buf := bytes.NewBuffer(by)
	binary.Read(b_buf, binary.LittleEndian, num)
}

// 数字 转化 byte
func inttobyte(num *int32) []byte {
	b_buf := new(bytes.Buffer)
	binary.Write(b_buf, binary.LittleEndian, num)
	return b_buf.Bytes()
}

func bytostr(by byte) string {
	return strconv.Itoa(int(by))
}

//判断目录或文件存在
func PathExist(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		return nil
	}
	return err
}

//utf8->gbk
func utf8ToGbk(str []byte) (b []byte, err error) {
	r := transform.NewReader(bytes.NewReader(str), simplifiedchinese.GBK.NewEncoder())
	b, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}
	return
}

//gbk->utf8
func gbkToUtf8(str []byte) (b []byte, err error) {
	r := transform.NewReader(bytes.NewReader(str), simplifiedchinese.GBK.NewDecoder())
	b, err = ioutil.ReadAll(r)
	if err != nil {
		return
	}
	return
}
