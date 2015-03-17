package main

import (
	"fmt"
	"runtime"

	"code.google.com/p/gopacket/layers"
	"code.google.com/p/gopacket/pcap"
	//"code.google.com/p/gopacket/pcapgo"
)
import "code.google.com/p/gopacket"

import (
	"github.com/zionist/gossip/base"
	"github.com/zionist/gossip/parser"
)

func doParse(inPackages chan gopacket.Packet, outPackages chan base.SipMessage, done chan int) {
	for packet := range inPackages {
		if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
			//fmt.Println(packet.Metadata().Timestamp)
			udp, _ := udpLayer.(*layers.UDP)
			msg, err := parser.ParseMessage(udp.Payload)
			if err != nil {
				fmt.Println(err)
				//panic("Err in thread")
			} else {
				outPackages <- msg
				fmt.Printf("msg %s, in goroutine %d \n", msg.Short(), runtime.NumGoroutine())
			}
		}
	}
	done <- runtime.NumGoroutine()
}

func getParse(packages chan base.SipMessage) {
	for msg := range packages {
		fmt.Println(msg.Short())
	}
	//for range packages {
	//
	//	}
}

func main() {
	// here we go
	runtime.GOMAXPROCS(runtime.NumCPU())

	if handle, err := pcap.OpenOffline("/home/slaviann/work/teligent/multifone/samples/big.pcap"); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		packages := make(chan base.SipMessage)
		//udpPackages := make(chan *layers.UDP)

		done := make(chan int, runtime.NumCPU())
		for i := 0; i < runtime.NumCPU(); i++ {
			go doParse(packetSource.Packets(), packages, done)
			go getParse(packages)
		}
		//wait all doParse gouroutines finished
		for i := 0; i < runtime.NumCPU(); i++ {
			<-done
		}
		close(packages)
	}

}
