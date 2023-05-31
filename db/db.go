package main

import (
	"fmt"
	"github.com/rereandy/db"
)

func main() {
	options := &db.Options{
		UserName: "root",
		Password: "123456",
		Host:     "localhost",
		Port:     3306,
		DBName:   "test",
	}
	conn := db.Open(options)
	typeImpl := TAdUgYybTypeInfoImpl{
		Connection: conn,
	}
	data, _ := typeImpl.Select()
	fmt.Println(data[0].Material)
	ent := &TAdUgYybTypeInfoEntity{
		Material: "王者荣耀",
		Type:     "游戏",
	}
	_, err := typeImpl.Insert(ent)
	if err != nil {
		fmt.Println(err)
	}
}
