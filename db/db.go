package db

import (
	"exam/tree"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

var DB *gorm.DB

//用户  黑名单
type BlackList struct {
	Name string `json:"name" gorm:"type:varchar(255)"`
	gorm.Model
}

//词汇
type Vocabularies struct {
	Id   int    `json:"id" gorm:"primary_key"`
	Word string `json:"word" gorm:"type:char(255)"`
}

//用户
type User struct {
	Name     string `json:"name" gorm:"type:varchar(255);index:na_p"`
	PassWord string `json:"pass_word" gorm:"index:na_p"`
	gorm.Model
}

const (
	USER_NAME = "root"
	PASS_WORD = "123"
	HOST      = "localhost"
	PORT      = "3306"
	DATABASE  = "bullet"
	CHARSET   = "utf8"
)

func InitDb() {
	//初始化
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true", USER_NAME, PASS_WORD, HOST, PORT, DATABASE, CHARSET)
	open, err := gorm.Open("mysql", dbDSN)
	if err != nil {
		log.Panicln("open mysql error")
	}

	if !open.HasTable("black_lists") {
		err := open.CreateTable(&BlackList{}).Error
		if err != nil {
			log.Panicln("Creat black_lists Table error")
		}
	}
	if !open.HasTable("users") {
		err := open.CreateTable(&User{}).Error
		if err != nil {
			log.Panicln("Creat users Table error")
		}
	}

	DB = open
	GetBlackList()

}

func AddUser(u User) error {
	if err := DB.Create(&u).Error; err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func FindUser(name string) *User {
	var u User
	DB.Model(&User{}).Where("name=?", name).Scan(&u)
	return &u
}

func AddWords(str string) {
	c := Vocabularies{Word: str}
	if err := DB.Model(&Vocabularies{}).Create(&c).Error; err != nil {
		fmt.Println(str)
		fmt.Println(err)
	}
}

//获取词汇
func GetWords() []string {

	words := make([]string, 4465)
	rows, _ := DB.Table("vocabularies").Rows()
	for rows.Next() {
		var id int
		var str string
		if err := rows.Scan(&id, &str); err != nil {
			log.Println(err)
			return words
		}
		if len(str) != 0 {
			if int(str[len(str)-1]) == 13 {
				words = append(words, str[:len(str)-1]) //每个取出来后面都有个 “/n"
			} else {
				words = append(words, str)
			}
		}
	}

	return words
}

//增加黑名单
func AddBlack(user BlackList) error {
	if err := DB.Model(&BlackList{}).Create(&user).Error; err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

//初始化黑名单
func GetBlackList() {
	rows, _ := DB.Table("black_lists").Rows()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Println(err)
			return
		}
		if len(name) != 0 {
			tree.NameTree.AddNote(name)
		}
	}
}
