package cloud

import (
	"errors"
	"net"
	"sync"
)

const (
	NetworkProtocol = "tcp"
	DefaultBuffer   = 2048
)

type Cloud struct {
	alias  string
	local  *net.IPAddr
	remote *net.IPAddr
	conn   *net.IPConn
	mux    sync.Mutex
}


func (c *Cloud) Init(alias string, local, remote *net.IPAddr) {
	c.alias, c.local, c.remote = alias, local, remote
}

func resolveLocalRemoteIPAddr(NetworkProtocol, localAddr, remoteAddr string) (*net.IPAddr, *net.IPAddr, error) {
	local, err := net.ResolveIPAddr(NetworkProtocol, localAddr)
	if err != nil {
		return nil, nil, err
	}

	remote, err := net.ResolveIPAddr(NetworkProtocol, remoteAddr)
	if err != nil {
		return nil, nil, err
	}

	return local, remote, err
}

var cloudConnectionPool = sync.Pool{
	New: func() any {
		return new(*net.IPConn)
	},
}

func (c *Cloud) CloudEstablishConnection(local_addr, remote_addr string) (*net.IPConn, error) {
	local, remote, err := resolveLocalRemoteIPAddr(NetworkProtocol, local_addr, remote_addr)
	if err != nil {
		return nil, nil
	}

	c.mux.Lock()
	conn, err := net.DialIP(NetworkProtocol, local, remote)
	if err != nil {
		c.mux.Unlock()
		return nil, err
	}
	c.conn = conn
	cloudConnectionPool.Put(c.conn)

	c.mux.Unlock()

	return conn, err
}

var cloudTunnel = make(chan []byte, DefaultBuffer)

func (*Cloud) Connect() (*net.IPConn, error) {
	connInterface := cloudConnectionPool.Get()
	if connInterface == nil {
		return nil, errors.New("no available connections in the pool")
	}

	conn, ok := connInterface.(*net.IPConn)
	if !ok {
		return nil, errors.New("failed to convert pool item to *net.IPConn")
	}

	return conn, nil
}

