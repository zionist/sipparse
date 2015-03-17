package main

import (
	"fmt"
	"runtime"

	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/layers"
	"code.google.com/p/gopacket/pcap"
	"github.com/zionist/gossip/parser"
	"github.com/zionist/sipparse/call"
	"github.com/zionist/sipparse/tpackage"

	//"code.google.com/p/gopacket/pcapgo"
)

func doParse(inPackages chan gopacket.Packet, outPackages chan *tpackage.Tpackage, done chan string) {
	uuid := uuid.NewRandom()
	for packet := range inPackages {
		if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
			//fmt.Println(packet.Metadata().Timestamp)
			udp, _ := udpLayer.(*layers.UDP)
			msg, err := parser.ParseMessage(udp.Payload)
			tmsg := tpackage.NewTpackage(msg, packet.Metadata().Timestamp)
			if err != nil {
				fmt.Println(err)
				//panic("Err in thread")
			} else {
				outPackages <- tmsg
				//fmt.Printf("msg %s, in goroutine %s \n", msg.Short(), uuid.String())
			}
		}
	}
	done <- uuid.String()
	//fmt.Printf("Goroutine %s finished\n", uuid.String())
}

func getParse(packages chan *tpackage.Tpackage, done chan string) {
	uuid := uuid.NewRandom()
	calls := make(map[string]*call.Call)

	for msg := range packages {
		// Create Call
		l := len(msg.SipPackage.Headers("Call-Id"))
		// Workaround for empty Call-Id header
		if l != 0 {
			_, ok := calls[msg.SipPackage.Headers("Call-Id")[0].String()]
			//fmt.Println(ok)
			//fmt.Println(mm)
			if !ok {
				//fmt.Println(mm)
				c := call.NewCall(msg)
				calls[msg.SipPackage.Headers("Call-Id")[0].String()] = c
			}
			//fmt.Println(msg.SipPackage.Short())
		}
		for _, c := range calls {
			if c.CheckPackageInCall(msg) {
				fmt.Println("yes")
				c.AddPackage(msg)
			} else {
				fmt.Println("no")
			}
		}

		// Check each package does it belong to Call
	}
	//fmt.Println(calls)

	for _, call := range calls {
		fmt.Println(call)
		//fmt.Println(call.GetPackages())
	}

	//}
	//for range packages {
	//
	//	}
	done <- uuid.String()
}

func main() {
	// set max
	runtime.GOMAXPROCS(runtime.NumCPU())

	if handle, err := pcap.OpenOffline("/home/slaviann/work/teligent/multifone/samples/2calls.pcap"); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		packages := make(chan *tpackage.Tpackage)
		//udpPackages := make(chan *layers.UDP)

		doneParse := make(chan string, runtime.NumCPU())
		doneGet := make(chan string, runtime.NumCPU())
		go getParse(packages, doneGet)
		for i := 0; i < runtime.NumCPU(); i++ {
			go doParse(packetSource.Packets(), packages, doneParse)
		}

		for i := 0; i < runtime.NumCPU(); i++ {
			<-doneParse
			//fmt.Printf("goroutine %s finished \n", <-done)
		}
		close(packages)
		for i := 0; i < 1; i++ {
			<-doneGet
			//fmt.Printf("goroutine %s finished \n", <-done2)
		}
	}

}
