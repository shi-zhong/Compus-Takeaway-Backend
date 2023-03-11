package dbop

import (
	"backend/model"
	//	"backend/model"
		"fmt"
	"strconv"
)

type SearchModel struct {
	SearchKeys   []string `json:"search_keys"`
	Building     uint     `json:"building"`
	BuildngFloor uint     `json:"buildng_floor"`
}

func SearchCommoditiesLimitPage(limit, page string, search *SearchModel) []*model.Shop {

	var shops []*model.Shop

	//	limitInt, _ := strconv.Atoi(limit)
	//	pageInt, _ := strconv.Atoi(page)
	//
	//	if limitInt == 0 {
	//		limitInt = -1
	//		pageInt = 1
	//	}

	//	fmt.Println(pageInt)

	buildingWhere := ""
	if search.Building != 999 {
		buildingWhere += "Where `physical_addresses`.`building_id` = " + strconv.Itoa(int(search.Building))
		if search.BuildngFloor != 999 {
			buildingWhere += " and `physical_addresses`.`building_floor` = " + strconv.Itoa(int(search.BuildngFloor))
		}
	}

	db := model.Db.Self

	// 1. 地址内的所有店家
	var ids []int
	shopidstring := ""
	shopIds := "select id from `shops` where " + "`can_be_searched` = 2 and "+
		"`address_id` in (" +
		"SELECT `physical_addresses`.`id` FROM `physical_addresses` " +
		"left join `buildings` on `buildings`.`id` = `physical_addresses`.`building_id` " +
		buildingWhere +
		")"

	db.Raw(shopIds).Scan(&ids)

    //  有地点条件 但是 商铺已经为 0 直接返回空
	if search.Building != 999 && len(ids) == 0 {
		return shops
	}

	for index, id := range ids {
		shopidstring += strconv.Itoa(int(id))
		if index != len(ids)-1 {
			shopidstring += ","
		}
	}

	shopidstring = "(" + shopidstring + ")"

	searcKey := ""

	if len(search.SearchKeys) != 0 {
		for _, key := range search.SearchKeys {
			searcKey += key + "|"
		}

		searcKey = " '" + searcKey[:len(searcKey)-1] + "' "
	}

	// 有楼栋搜索

	shopsql := "SELECT * from `shops` where `can_be_searched` = 2 "
	commoditysql := "SELECT shop_id as id from `commodity_infos` "
	sql := ""

	// 有条件
	if search.Building != 999 || len(search.SearchKeys) != 0 {
		commoditysql += "where "
	}

	if search.Building != 999 {
		shopsql += "and id in " + shopidstring
		commoditysql += "shop_id in " + shopidstring
		if len(search.SearchKeys) != 0 {
			shopsql += " and "
			commoditysql += " and "
		}
	}

	if len(search.SearchKeys) != 0 {
		shopsql += "and `shop_name` REGEXP" + searcKey
		commoditysql += "`name` REGEXP" + searcKey
	}

	sql = shopsql

	if len(search.SearchKeys) != 0 {
		sql += "union select * from `shops` where `can_be_searched` = 2 and id in (" + commoditysql + ")"
	}

	fmt.Println(sql)

	db.Raw(sql).Scan(&shops)

	return shops
}
