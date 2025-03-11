package main

import (
	"radar240/sender"
	"radar240/utils"

)

func main() {

	go sender.Sender() // websockets server
	go utils.ReadData()	// read data from pcap file
	go utils.ParseData() // parse data and send to websockets server
	select {}
}

