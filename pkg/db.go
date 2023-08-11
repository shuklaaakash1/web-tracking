package pkg

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(mysql.Open("root:8955507476@Mn@tcp(127.0.0.1:3306)/myappdb?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	DB = db
}
