package tpackage

import (
	"fmt"
	"time"

	"github.com/zionist/gossip/base"
)

//Tpackage is Struct wich contains SipMessage and timestamp from pcap file
type Tpackage struct {
	SipPackage base.SipMessage
	Timestamp  time.Time
	Method     base.Method
	XclientIP  string
}

//GetID Get package identificator
func (t Tpackage) GetID() string {
	return fmt.Sprintf("%s |%s|", t.XclientIP, t.Method)
}

//NewTpackage Constructor
func NewTpackage(sipPackage base.SipMessage, timestamp time.Time) *Tpackage {
	t := new(Tpackage)
	t.SipPackage = sipPackage
	t.Timestamp = timestamp

	if len(sipPackage.Headers("x-clientip")) == 0 {
		t.XclientIP = "x-clientip: None"
	} else {
		t.XclientIP = fmt.Sprint(sipPackage.Headers("x-clientip")[0])
	}

	switch sipPackage.(type) {
	case *base.Request:
		s, _ := sipPackage.(*base.Request)
		t.Method = s.Method
	case *base.Response:
	}

	return t
}

func (t Tpackage) String() string {
	return ""
}
