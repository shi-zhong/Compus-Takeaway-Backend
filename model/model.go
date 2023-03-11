package model

import (
	"time"
)

type User struct {
	ID              uint   `gorm:"<-:create;not null;uniqueIndex" json:"id"`
	NickName        string `gorm:"<-;not null;type:char(64)" json:"nick_name"`
	Avatar          string `gorm:"<-;type:varchar(1024)" json:"avatar"`
	Phone           string `gorm:"<-;not null;type:char(16)" json:"phone"`
	OpenID          string `gorm:"<-;not null;type:char(32)" json:"open_id"`
	Identity        uint   `gorm:"<-:create;not null" json:"identity"`
	LastUsedAddress uint   `gorm:"<-;not null" json:"last_used_address"`
}

type CustomerAddress struct {
	ID       uint   `gorm:"<-:create;not null;uniqueIndex" json:"id"`
	UserID   uint   `gorm:"<-:create;not null"`
	Receiver string `gorm:"<-;not null;type:char(32)"`
	Phone    string `gorm:"<-;not null;type:char(16)"`
	Address  string `gorm:"<-;not null;type:varchar(1024)"`
}

type PhysicalAddress struct {
	ID             uint `gorm:"<-:create;not null;uniqueIndex" json:"id"`
	Building       Building
	BuildingID     uint   `gorm:"<-:create;not null"`
	BuildingFloor  uint   `gorm:"<-:create;not null"`
	BuildingNumber string `gorm:"<-;create;not null;type:varchar(64)"`
}

type Shop struct {
	ID            uint `gorm:"<-:create;not null;uniqueIndex" json:"id"`
	ShopKeeper    User
	ShopKeeperID  uint   `gorm:"<-:create;not null" json:"shop_keeper_id"`
	ShopName      string `gorm:"<-;not null;unique;type:char(32)" json:"name"`
	ShopIntro     string `gorm:"<-;type:varchar(1024)" json:"intro"`
	ShopAvatar    string `gorm:"<-;type:varchar(1024)" json:"avatar"`
	Address       PhysicalAddress
	AddressID     uint    `gorm:"<-;not null" json:"address_id"`
	Star          float64 `gorm:"<-;not null" json:"star"`
	Monthly       uint    `gorm:"<-;not null" json:"monthly"`
	StartDeliver  float64 `gorm:"<-;not null" json:"start_deliver"`
	CanBeSearched uint    `gorm:"<-;not null" json:"can_be_searched"`
}

type Tag struct {
	ID     uint   `gorm:"<-:create;not null;uniqueIndex" json:"id"`
	Belong uint   `gorm:"<-:create;not null" json:"belong"`
	Tag    string `gorm:"<-;type:varchar(32)" json:"tag"`
}

type CommodityInfo struct {
	ID      uint `gorm:"<-:create;not null;uniqueIndex" json:"id"`
	Shop    Shop
	ShopID  uint    `gorm:"<-:create;not null" json:"shop_id"`
	Name    string  `gorm:"<-;not null;type:char(32)" json:"name"`
	Price   float64 `gorm:"<-;not null" json:"price"`
	Intro   string  `gorm:"<-;not null;type:varchar(256)" json:"intro"`
	Status  uint    `gorm:"<-;not null" json:"status"`
	Picture string  `gorm:"<-;not null;type:varchar(256)" json:"picture"`
	Tags    string  `gorm:"<-;type:text" json:"tags"`
}

type Order struct {
	ID          string `gorm:"<-:create;type:char(24);not null;uniqueIndex;primaryKey" json:"id"`
	CustomerID  uint   `gorm:"<-:create;not null" json:"customer_id"`
	Customer    User
	ShopID      uint `gorm:"<-:create;not null" json:"shop_id"`
	Shop        Shop
	RiderID     uint `gorm:"<-:create" json:"rider_id"`
	Rider       User
	CommodityID string `gorm:"<-:create;not null;type:varchar(1024)" json:"commodity_id"`
	Receiver string `gorm:"<-;not null;type:char(32)"`
	Phone    string `gorm:"<-;not null;type:char(16)"`
    Address  string `gorm:"<-;not null;type:varchar(1024)"`

	//	PayAt           time.Time `gorm:"<-:create;not null"`
	CreateAt        time.Time `gorm:"<-:create;not null" json:"create_at"`
	AcceptAt        time.Time `gorm:"<-;default:null" json:"accept_at"`
    RiderAccceptAt  time.Time `gorm:"<-;default:null" json:"rider_acccept_at"`
    CookFinishAt    time.Time `gorm:"<-;default:null" json:"cook_finish_at"`
    DeliverBeginAt  time.Time `gorm:"<-;default:null" json:"deliver_begin_at"`
    DeliverFinishAt time.Time `gorm:"<-;default:null" json:"deliver_finish_at"`
    FinishAt        time.Time `gorm:"<-;default:null" json:"finish_at"`

    Status uint `gorm:"<-;not null" json:"status"`
}

type Building struct {
	ID            uint   `gorm:"<-:create;not null;uniqueIndex" json:"id"`
	Name          string `gorm:"<-;type:varchar(128)" json:"name"`
	AcceptFloor   string `gorm:"<-;type:varchar(1024)" json:"accept_floor"`
	DetailAddress string `gorm:"<-;type:varchar(1024)" json:"detail_address"`
}
