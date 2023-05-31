package main

import (
	"fmt"
	"github.com/rereandy/db"
)

type TypeInfo struct {
	Material string `db:"material"`
	Type     string `db:"type"`
}

func main() {
	options := &db.Options{
		UserName: "root",
		Password: "123456",
		Host:     "localhost",
		Port:     3306,
		DBName:   "test",
	}
	conn := db.Open(options)
	ss := conn.NewSession()
	var values []*TypeInfo
	_, err := ss.Select("material,type").From("t_ad_ug_yyb_type_info").Load(&values)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(values[0].Type)
}
