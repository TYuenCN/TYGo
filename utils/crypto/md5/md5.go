// 用于获取MD5值的字符串的，相关工具函数
//
// Created by Yuen Cheng on 13-11-4.
// Copyright (c) 2013年 YuenSoft.com. All rights reserved.
package YSMD5Util

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"
)

// 获取某字符串的对应的MD5值
// 		import (
//			"yuensoft.com/utils/crypto/YSMD5Util"
//		)
//		func main() {
// 			fmt.Println(YSMD5Util.MD5FromString("Hello"))
// 		}
func MD5FromString(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// 获取基于当前Unix风格时间(Since 1,1,1970)的MD5值
// 		import (
//			"yuensoft.com/utils/crypto/YSMD5Util"
//		)
//		func main() {
// 			fmt.Println(YSMD5Util.MD5FromUnixTimeNow())
// 		}
func MD5FromUnixTimeNow() string {
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	return fmt.Sprintf("%x", h.Sum(nil))
}
