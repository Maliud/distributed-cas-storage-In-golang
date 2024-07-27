package p2p

import "net"

//RPC, ağdaki iki düğüm arasında her bir aktarım üzerinden gönderilen herhangi bir keyfi veriyi tutar.
type RPC struct {
	From net.Addr
	Payload []byte
	
}