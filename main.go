package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/layers"
	"code.google.com/p/gopacket/pcap"
	"github.com/zionist/gossip/base"
	"github.com/zionist/gossip/parser"
	"github.com/zionist/sipparse/tpackage"

	//"code.google.com/p/gopacket/pcapgo"
)

type statistic map[string]uint64

func (s statistic) String() string {
	r := ""
	for k, v := range s {
		r = fmt.Sprintf("%s%s -> %d\n", r, k, v)
	}
	return r
}

//PrintCount stat each
func PrintCount(count *uint64) {
	for {
		time.Sleep(1 * time.Second)
		fmt.Fprintf(os.Stderr, fmt.Sprintf("done %d \n", *count))
	}
}

func doParse(inPackages chan gopacket.Packet, outPackages chan *tpackage.Tpackage, done chan string) {
	//uuid := uuid.NewRandom()
	for packet := range inPackages {
		if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
			udp, _ := udpLayer.(*layers.UDP)
			msg, err := parser.ParseMessage(udp.Payload)
			tmsg := tpackage.NewTpackage(msg, packet.Metadata().Timestamp)
			if err != nil {
				fmt.Println(err)
				//panic("Err in thread")
			} else {
				outPackages <- tmsg
			}
		}
	}
	done <- "s"
}

func getParse(packages chan *tpackage.Tpackage, done chan string) {
	var stat statistic = make(map[string]uint64)

	var c uint64
	count := &c

	go PrintCount(count)
	//Count only requests
	for msg := range packages {
		switch msg.SipPackage.(type) {
		case *base.Request:
			_, ok := stat[msg.GetID()]
			if !ok {
				stat[msg.GetID()] = 1
			} else {
				stat[msg.GetID()]++
			}
			c++
		//Do not count responces
		case *base.Response:
			continue
		}
	}

	fmt.Println(stat)

	done <- "s"
}

func main() {
	const HELP = `usage: sipparse -f <pcap_filename>`

	filename := flag.String("f", "", "pcap filename")
	flag.Parse()
	if len(*filename) == 0 {
		fmt.Println(HELP)
		os.Exit(1)
	}

	// here we go
	runtime.GOMAXPROCS(runtime.NumCPU())

	if handle, err := pcap.OpenOffline(*filename); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		packages := make(chan *tpackage.Tpackage)

		doneParse := make(chan string, runtime.NumCPU())
		doneGet := make(chan string, runtime.NumCPU())
		go getParse(packages, doneGet)
		for i := 0; i < runtime.NumCPU(); i++ {
			go doParse(packetSource.Packets(), packages, doneParse)
		}

		for i := 0; i < runtime.NumCPU(); i++ {
			<-doneParse
		}
		close(packages)

		<-doneGet

	}

}
