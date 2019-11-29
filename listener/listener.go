package listener

import (
	"log"
	"mock/lib/logger"
)

/**
 * Created by zc on 2019-11-28.
 */

// 注册监听
func init() {
	go storeApply()
	go storeForce()
	go cubaRole()
}

func printf(format string, args ...interface{}) {
	log.Printf(format, args...)
	logger.Info.Printf(format, args...)
}

func errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
	logger.Error.Errorf(format, args...)
}

func errorln(args ...interface{}) {
	log.Println(args...)
	logger.Error.Errorln(args...)
}