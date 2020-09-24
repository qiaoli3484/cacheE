package cacheE

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Ecache struct {
	path string
	name string
}

//创建库
func NewCache(path, name string) (*Ecache, error) {
	if len(name) < 1 {
		return nil, errors.New("0001|数据库名必须1字节以上")
	}

	err := os.MkdirAll(path+"/"+name, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return &Ecache{
		path,
		name,
	}, nil
}

//打开库
func OpenCache(path, name string) (*Ecache, error) {
	err := PathExist(path + "/" + name)
	if err != nil {
		return nil, err
	}
	return &Ecache{
		path,
		name,
	}, nil
}

//删除库
func DelCache(path, name string) error {
	err := os.Remove(path + "/" + name)
	return err
}

//创建表
func (e *Ecache) CreateTable(name string) error {
	path := e.path + "/" + e.name + "/" + name
	err := PathExist(path)
	if err == nil {
		return nil
	}
	err = os.Mkdir(path, os.ModePerm)
	if err == nil {
		return nil
	}
	return err
}

//删除表
func (e *Ecache) DelTable(name string) error {
	path := e.path + "/" + e.name + "/" + name
	err := PathExist(path)
	if err != nil {
		return nil
	}

	err = os.Remove(path)
	if err == nil {
		return nil
	}
	return err
}

//读文本
func (e *Ecache) GetString(name, key, field string) (string, error) {

	path := e.keyToPath([]byte(key), name)

	err := PathExist(path)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	var next int
	var res []byte
	res, _ = ioutil.ReadFile(path + field + ".数据")
	for len(res) == 0 {
		time.Sleep(time.Second)
		if next >= 1000 {
			return "", errors.New("")
		}
		next++
		res, _ = ioutil.ReadFile(path + field + ".数据")
	}
	str, err := gbkToUtf8(res)

	return string(str), err
}

//读逻辑
func (e *Ecache) GetBool(name, key, field string) (bool, error) {
	path := e.keyToPath([]byte(key), name)
	err := PathExist(path)
	if err != nil {
		return false, err
	}

	var next int
	var res []byte
	res, _ = ioutil.ReadFile(path + field + ".数据")
	for len(res) == 0 {
		time.Sleep(time.Second)
		if next >= 1000 {
			return false, errors.New("")
		}
		next++
		res, _ = ioutil.ReadFile(path + field + ".数据")
	}
	var num int32
	bytetoint(res, &num)

	return intTobool(num), nil
}

//写文本
func (e *Ecache) SetString(name, key, field, value string) error {

	path := e.keyToPath([]byte(key), name)

	err := PathExist(path)
	if err != nil {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	str, err := utf8ToGbk([]byte(value))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path+field+".数据", str, os.ModePerm)
	var next int
	for err != nil {
		time.Sleep(time.Second)
		if next >= 1000 {
			return err
		}
		next++
		err = ioutil.WriteFile(path+field+".数据", str, os.ModePerm)
	}
	return nil
}

//写逻辑
func (e *Ecache) SetBool(name, key, field string, value bool) error {
	path := e.keyToPath([]byte(key), name)
	err := PathExist(path)
	if err != nil {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	str := booltoint(value)
	by := inttobyte(&str)
	err = ioutil.WriteFile(path+field+".数据", by, os.ModePerm)
	var next int
	for err != nil {
		time.Sleep(time.Second)
		if next >= 1000 {
			return err
		}
		next++
		err = ioutil.WriteFile(path+field+".数据", by, os.ModePerm)
	}
	return nil
}

func (e *Ecache) keyToPath(key []byte, name string) string {
	aa, _ := utf8ToGbk(key)
	cc := len(aa)
	fmt.Println(aa)
	//{178,226,202,212}
	//[230 181 139 232 175 149]
	var path string
	flag := false
	for i := 0; i < cc; i++ {
		if flag == false {
			if i+1 <= cc-1 {
				path = path + bytostr(aa[i]) + bytostr(aa[i+1]) + "/"
				flag = true
				continue
			} else {
				path = path + bytostr(aa[i]) + "/"
			}
		} else if flag {
			flag = false
		}
	}
	path = strings.Replace(path, "//", "/", -1)
	path = e.path + "/" + e.name + "/" + name + "/" + path
	return path
}
