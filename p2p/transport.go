package p2p


// Peer, uzak düğümü temsil eden bir arayüzdür.

type Peer interface {
	Close() error
}

// Taşıma, iletişimi yöneten her şeydir
// ağdaki düğümler arasında
// bu TCP, UDP, web soketleri gibi biçimlerde olabilir


type Transport interface {
	ListenAndAccept() error
	Consume() <- chan RPC
}