package net

import (
	"net"
	"io"
	"time"
	"crypto/tls"

	"github.com/go-routeros/routeros"
)


func Dial(address string, timeout int) (net.Conn, error) {
	return net.DialTimeout("tcp", address, time.Duration(timeout) * time.Second)
}

// DialTLS connects and logs in to a RouterOS device using TLS.
func DialTLS(address, username, password string, tlsConfig *tls.Config) (*routeros.Client, error) {
	conn, err := tls.Dial("tcp", address, tlsConfig)
	if err != nil {
		return nil, err
	}
	return NewClientAndLogin(conn, username, password)
}

func NewClientAndLogin(rwc io.ReadWriteCloser, username, password string) (*routeros.Client, error) {
	c, err := routeros.NewClient(rwc)
	if err != nil {
		rwc.Close()
		return nil, err
	}
	err = c.Login(username, password)
	if err != nil {
		c.Close()
		return nil, err
	}
	return c, nil
}