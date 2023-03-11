package utils

import (
    "crypto/md5"
    "encoding/hex"
)

func MD5(str string) string {
	//创建MD5算法
	h := md5.New()
	//写入待加密数据
	h.Write([]byte(str))
	//获取哈希值字符切片
	bts := h.Sum(nil)
	//转化为16进制字符串
	return hex.EncodeToString(bts)
}
