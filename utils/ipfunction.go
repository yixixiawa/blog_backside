package utils

import (
	"fmt"
	"net"
	"os"
)

func GetUserIp() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Print(addrs)
}
