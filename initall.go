package main

import (
	"bufio"
	"exam/db"
	"exam/server"
	"exam/timer"
	"exam/tree"
	"fmt"
	"io"
	"log"
	"os"
)


func InitAll() {
	//数据库
	db.InitDb()
	if !db.DB.HasTable("vocabularies") {
		err := db.DB.CreateTable(&db.Vocabularies{}).Error
		if err != nil {
			log.Println("Creat vocabularies Table error")
		}
		setVocabularies()
	}
	s := db.GetWords()

	for i := 0; i < len(s); i++ {
		if len(s[i]) > 0 && int(s[i][len(s[i])-1]) == 13 {
			fmt.Println(s[i])
		}
		tree.T.AddNote(s[i])
	}
	//服务
	server.InitManger()
	//计时器
	go timer.StartTimer()
}
//文件内容的写入
func setVocabularies() {

	//打开文件
	f, err := os.Open("other.txt")
	if err != nil {
		panic(err)
	}

	//关闭文件
	defer f.Close()

	//新建一个缓冲区，把内容先放在缓冲区
	r := bufio.NewReader(f)

	for {
		//遇到'\n'结束读取, 但是'\n'也读取进入
		buf, err := r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF { //文件已经结束
				break
			}
			fmt.Println("err = ", err)
		}
		s := string(buf[:len(buf)-1])
		db.AddWords(s)
	}
}