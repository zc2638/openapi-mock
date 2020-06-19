/**
 * @Author: fan
 * @Description:
 * @Date: 2020/6/19
 */
package network

import (
	"fmt"
	"testing"
)

func TestHostname(t *testing.T) {
	ip, err := ExternalIP()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ip.String())
}
