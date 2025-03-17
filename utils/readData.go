package utils

import (
	"fmt"
	"net"
	"radar240/global"
	"time"
)

// ReadData reads data from a pcap file
var IP_ADDRESS = "84.247.170.21:7055"
func ReadData() {
	
	for {
		connection, err  := net.Dial("tcp", IP_ADDRESS)
		if err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}
		fmt.Println("Connected to the server ", IP_ADDRESS)
		n , err := connection.Write([]byte("JTh0453YsksaCYo\n"))
		fmt.Println(n, err)
		// defer connection.Close()
		// read data from the connection
		for {
			var buffer = make([]byte, 1000000)
			// fmt.Println("before reading")
			n , err :=  connection.Read(buffer)
			if err != nil {
				fmt.Println(err)
				break
			}
			global.FilteredData <- buffer[:n]
			// time.Sleep(1 * time.Second)
		}
		connection.Close()
	}
}
