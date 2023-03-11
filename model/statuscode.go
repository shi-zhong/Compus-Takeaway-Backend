package model

var (
    IdentityCustomer   uint = 1
    IdentityShopKeeper uint = 2
    IdentityRider      uint = 3
    IdentityManager    uint = 4
    IdentityAdmin      uint = 5
)


var (
	AddressNotExsit   uint8 = 0
	AddressDefault    uint8 = 1
	AddressNotDefault uint8 = 2
	AddressDelete     uint8 = 3
)

var (
	CommodityStatusNotCreate    uint = 0
	CommodityStatusOnShelf      uint = 1
	CommodityStatusOffShelf     uint = 2
	CommodityStatusForceOnShelf uint = 3
	CommodityStatusNotEnought   uint = 4
	CommodityStatusDeleted      uint = 5
)

var (
	OrderCreate        uint = 1
	OrderPay           uint = 2
	OrderAccept        uint = 3
	OrderRiderAcccept  uint = 4
	OrderCookFinish    uint = 5
	OrderDeliverBegin  uint = 6
	OrderDeliverFinish uint = 7
	OrderFinish        uint = 8
	OrderCancel        uint = 9
)

var (
    ShopIsOpen uint = 2
    ShopIsClose uint = 3
)
