package model

import (
	"backend/utils"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var dsn = formatDsn()

type Database struct {
	Self *gorm.DB
}

var Db *Database

func formatDsn() string {

	dbConfig := utils.GlobalConfig.DbConfig

	fDsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.Host,
		dbConfig.Password,
		dbConfig.Port,
		dbConfig.DatabaseName)

	return fDsn
}

func getDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
		return nil, err

	}
	return db, err
}

func (db *Database) Init() {
	newDb, err := getDatabase()
	if err != nil {
		log.Print("A error occurred when trying to init a Database!")
		log.Println(err)
	}
	Db = &Database{Self: newDb}
	Db.AutoMigrate()
}

func (db *Database) AutoMigrate() {
//	 author
	db.Self.AutoMigrate(&User{})

//	 basic_info
	db.Self.AutoMigrate(&CustomerAddress{})
	db.Self.AutoMigrate(&Shop{})
//	 chat
	db.Self.AutoMigrate(&PhysicalAddress{})
	db.Self.AutoMigrate(&Tag{})
//	 commodity
	db.Self.AutoMigrate(&CommodityInfo{})
	db.Self.AutoMigrate(&Building{})

//	 order
    db.Self.AutoMigrate(&Order{})
}

func (db *Database) Close() {

	//	if err := Db.Self.Close(); err != nil {
	//		log.Print("A error occurred when trying to close the Database!")
	//		log.Println(err)
	//	}
	fmt.Println("当前数据库驱动无关闭函数")
	return
}
