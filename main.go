package main

import (
	"fmt"

	"code.google.com/p/gopacket/layers"
	"code.google.com/p/gopacket/pcap"
	//"code.google.com/p/gopacket/pcapgo"
)
import "code.google.com/p/gopacket"

import (
	"github.com/zionist/gossip/base"
	"github.com/zionist/gossip/parser"
)

func doParse(packages chan base.SipMessage, udp *layers.UDP) {
	fmt.Println("# start")
	msg, err := parser.ParseMessage(udp.Payload)
	if err != nil {
		fmt.Println(err)
		panic("Err in thread")
	}
	packages <- msg
	fmt.Println("# end")
}

func main() {
	if handle, err := pcap.OpenOffline("/home/slaviann/work/teligent/multifone/samples/big.pcap"); err != nil {
		panic(err)
	} else {

		//output := make(chan base.SipMessage)
		//errs := make(chan error)
		//pars := parser.NewParser(output, errs, false)
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

		runs := 0
		packages := make(chan base.SipMessage)
		for packet := range packetSource.Packets() {
			if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
				//fmt.Println(packet.Metadata().Timestamp)
				//parser.
				udp, _ := udpLayer.(*layers.UDP)
				go doParse(packages, udp)

				//fmt.Println(string(udp.Payload))
				//go pars.Write(udp.Payload)
				runs++
			}

			//fmt.Println(string(packet.Data()))

			//fmt.Printf(string(packet.Data()))
			// Do something with a packet here.
		}

		for msg := range packages {
			runs--
			// last run
			if runs == 0 {
				close(packages)
			}
			switch msg.(type) {
			case *base.Request:
				fmt.Printf("request %s \n", msg.Short())
			case *base.Response:
				fmt.Printf("responce %s \n", msg.Short())
			}

		}
		//}

	}

}
