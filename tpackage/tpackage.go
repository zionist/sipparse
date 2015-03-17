package tpackage

import (
	"time"

	"github.com/zionist/gossip/base"
)

//Tpackage is Struct wich contains SipMessage and timestamp from pcap file
type Tpackage struct {
	SipPackage base.SipMessage
	Timestamp  time.Time
}

//NewTpackage Constructor
func NewTpackage(sipPackage base.SipMessage, timestamp time.Time) *Tpackage {
	t := new(Tpackage)
	t.SipPackage = sipPackage
	t.Timestamp = timestamp
	return t
}

func (t Tpackage) String() string {
	return t.SipPackage.String()
}
