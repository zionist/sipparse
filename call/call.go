package call

import (
	"fmt"

	"github.com/zionist/sipparse/tpackage"
)

// Call struct. Earch call equals SCenario
type Call struct {
	packages []*tpackage.Tpackage
	callID   string
}

//GetPackages packages getter
func (c *Call) GetPackages() []*tpackage.Tpackage {
	return c.packages
}

func (c Call) String() string {
	s := fmt.Sprintf("======= %s ======= \n", c.callID)
	for p := range c.packages {
		s = fmt.Sprintf("%s %s \n", s, c.packages[p])
	}
	return s
}

//NewCall Call constructor
func NewCall(pack *tpackage.Tpackage) *Call {
	c := new(Call)
	c.callID = fmt.Sprint(pack.SipPackage.Headers("Call-Id")[0])
	c.packages = make([]*tpackage.Tpackage, 0)
	return c
}

// AddPackage to the Call
func (c *Call) AddPackage(pack *tpackage.Tpackage) {
	c.packages = append(c.packages, pack)
}

//CheckPackageInCall does Tpackage belongs to Call
func (c *Call) CheckPackageInCall(pack *tpackage.Tpackage) bool {
	//fmt.Println(pack.SipPackage.Headers("Call-Id")[0])
	//sfmt.Println(pack.SipPackage.Headers("Call-Id")[0] == c.callID)

	return fmt.Sprint(pack.SipPackage.Headers("Call-Id")[0]) == c.callID

	//for i := range c.packages {
	//fmt.Println(c.packages[i] == nil)
	//	c := (c.packages[i].SipPackage.Headers("Call-Id")[0])
	//	p := pack.SipPackage.Headers("Call-Id")[0]
	//	if c == p {
	//		fmt.Println("Yes")
	//	} else {
	//		fmt.Println("No")
	//	}

	//}
}
