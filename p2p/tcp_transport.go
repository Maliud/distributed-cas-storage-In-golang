package p2p

import (
	"fmt"
	"net"
)

// TCPPeer, TCP ile kurulmuş bir bağlantı üzerinden uzak düğümü temsil eder.
type TCPPeer struct {

	// conn eşin temel bağlantısıdır
	conn net.Conn
	//eğer bir bağlantı çevirir ve alırsak => outbound == true
	// eğer bir bağlantıyı kabul eder ve alırsak => outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// Close, Peer arayüzünü uygular.
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC),
	}
}

// Consume, ağdaki başka bir eşten alınan gelen mesajları okumak için salt okunur kanal döndürecek olan Tranport arayüzünü uygular.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP kabul hatası: %s\n", err)
		}

		fmt.Printf("yeni gelen bağlantı:%+v\n", conn)

		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {

	var err error

	defer func()  {
		fmt.Printf("dropping peer connection: %s\n", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)

	if err = t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil{
			return
		}
	}

	//Loop'u oku
	rpc := RPC{}
	//buf := make([]byte, 2000)
	for {
		err = t.Decoder.Decode(conn, &rpc)

		if err != nil {
			return
		}
		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
	}
}
