package main

import (
	"fmt"
	"runtime"

	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/gopacket/layers"
	"code.google.com/p/gopacket/pcap"
	//"code.google.com/p/gopacket/pcapgo"
)
import "code.google.com/p/gopacket"

import (
	"github.com/zionist/gossip/base"
	"github.com/zionist/gossip/parser"
)

func doParse(inPackages chan gopacket.Packet, outPackages chan base.SipMessage, done chan string) {
	uuid := uuid.NewRandom()
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
				//fmt.Printf("msg %s, in goroutine %s \n", msg.Short(), uuid.String())
			}
		}
	}
	done <- uuid.String()
	//fmt.Printf("Goroutine %s finished\n", uuid.String())
}

func getParse(packages chan base.SipMessage, done chan string) {
	uuid := uuid.NewRandom()
	for msg := range packages {
		fmt.Println(msg.Short())
	}
	done <- uuid.String()
	//for range packages {
	//
	//	}
}

func main() {
	// here we go
	runtime.GOMAXPROCS(runtime.NumCPU())

	if handle, err := pcap.OpenOffline("/home/slaviann/work/teligent/multifone/samples/small.pcap"); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		packages := make(chan base.SipMessage)
		//udpPackages := make(chan *layers.UDP)

		doneParse := make(chan string, runtime.NumCPU())
		doneGet := make(chan string, runtime.NumCPU())
		for i := 0; i < runtime.NumCPU(); i++ {
			go doParse(packetSource.Packets(), packages, doneParse)
			go getParse(packages, doneGet)
		}

		for i := 0; i < runtime.NumCPU(); i++ {
			<-doneParse
			//fmt.Printf("goroutine %s finished \n", <-done)
		}
		close(packages)
		for i := 0; i < runtime.NumCPU(); i++ {
			<-doneGet
			//fmt.Printf("goroutine %s finished \n", <-done2)
		}
	}

}
