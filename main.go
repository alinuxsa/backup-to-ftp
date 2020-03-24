package main

import (
	"encoding/json"
	"fmt"
	"github.com/jlaffaye/ftp"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	CONFIG_FILE = "config.json"
	LOG_FILE    = "run.log"
)

type Config struct {
	Server   string `json:"server"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Pasv     bool   `json:"pasv"`
}

func init() {
	_, err := os.Stat(LOG_FILE)
	if err != nil {
		_, err := os.Create(LOG_FILE)
		if err != nil {
			log.Fatal("创建日志文件 %s 失败!\n", LOG_FILE)
		}
	}
	_, err = os.Stat(CONFIG_FILE)
	if err != nil {
		log.Fatal("不能读取配置文件 %s \n", CONFIG_FILE)
	}
}

func writeLog(msg string) {
	logFile, err := os.OpenFile(LOG_FILE, os.O_WRONLY|os.O_APPEND, 0666)
	defer logFile.Close()
	if err != nil {
		log.Fatal(err.Error())
	}
	LOG := log.New(logFile, "", log.LstdFlags)
	LOG.Println(msg)
}

func upload(dest string) {
	config := new(Config)
	data, err := ioutil.ReadFile(CONFIG_FILE)
	if err != nil {
		writeLog(err.Error())
	}
	err = json.Unmarshal([]byte(data), config)
	if err != nil {
		writeLog(err.Error())
	}
	c, err := ftp.Connect(config.Server + ":" + config.Port)
	if err != nil {
		writeLog(err.Error())
	}
	err = c.Login(config.User, config.Password)
	if err != nil {
		writeLog(err.Error())
	}
	/*
		fileList, err := c.List("/")
		if err != nil{
			writeLog(err.Error())
		}
		fmt.Println(fileList)
		for _, f := range fileList{
			fmt.Println(f.Name)
		}
	*/

	file, err := os.Open(dest)
	if err != nil {
		writeLog(err.Error())
		_ = c.Quit()
		return
	}
	defer file.Close()
	saveName := filepath.Base(dest)
	err = c.Stor(saveName, file)
	if err != nil {
		writeLog(err.Error())
	}
	_ = c.Quit()

}

func main() {
	dest := os.Args[1]
	_, err := os.Stat(dest)
	if err != nil {
		writeLog(fmt.Sprintf("%s 不是合法的文件路径!", dest))
		return
	}
	writeLog(fmt.Sprintf("上传 %s ", dest))
	upload(dest)

}
