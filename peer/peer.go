package peer

import (
	"fmt"
	"net"
)

type Peer struct {
	IP   net.IP
	Port uint16
}

func (p Peer) String() string {
	return fmt.Sprintf("%s:%d", p.IP.String(), p.Port)
}
