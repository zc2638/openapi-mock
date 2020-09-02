package main

import (
	"fmt"
	_ "mock/config"
	_ "mock/listener"
	"mock/route"
	"time"
	_ "time/tzdata"
)

/**
 * Created by zc on 2019-10-24.
 */
func main() {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err == nil {
		time.Local = location
	} else {
		fmt.Println("load location Asia/Shanghai error: ", err)
	}
	route.Start()
}
