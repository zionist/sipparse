package main

import (
	"fmt"

	"code.google.com/p/gopacket/layers"
	"code.google.com/p/gopacket/pcap"
	//"code.google.com/p/gopacket/pcapgo"
)
import "code.google.com/p/gopacket"

import "github.com/zionist/gossip/parser"
import "github.com/zionist/gossip/base"

func main() {
	if handle, err := pcap.OpenOffline("/tmp/test.pcap"); err != nil {
		panic(err)
	} else {

		output := make(chan base.SipMessage)
		errs := make(chan error)
		pars := parser.NewParser(output, errs, false)
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

		runs := 0
		for packet := range packetSource.Packets() {
			if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
				fmt.Println(packet.Metadata().Timestamp)
				udp, _ := udpLayer.(*layers.UDP)
				//fmt.Println(string(udp.Payload))
				pars.Write(udp.Payload)
				//fmt.Printf("# write end %d %s", n, err)
				//if err := <-errs; err != nil {
				//	fmt.Println(err)
				//} else {

				if msg := <-output; msg != nil {
					runs++
					fmt.Println(msg.Headers("Call-Id"))
				}
				//err := <-errs
				//fmt.Println(err)

			}
			//fmt.Println(string(packet.Data()))

			//fmt.Printf(string(packet.Data()))
			// Do something with a packet here.
		}
		fmt.Println(runs)

	}

}
