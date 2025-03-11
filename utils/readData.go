package utils

import (
	"fmt"
	"radar240/global"
	"time"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// ReadData reads data from a pcap file
func ReadData() {
	// open pcap file
	for {
		handle, err := pcap.OpenOffline("data/ASTERIX_CAT240_4_20230517190733 (1).pcap")
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer handle.Close()
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			global.FilteredData <- packet.ApplicationLayer().Payload()
			time.Sleep(1 * time.Second)
		}
	}
}