package main

import (
	_ "mock/config"
	_ "mock/listener"
	"mock/route"
)

/**
 * Created by zc on 2019-10-24.
 */
func main() {
	route.Start()
}
