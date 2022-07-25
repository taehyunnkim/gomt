package net

import (
	"net"
	"io"
	"time"
	"errors"

	"github.com/go-routeros/routeros"
)


func Dial(address string, timeout int) (net.Conn, error) {
	return net.DialTimeout("tcp", address, time.Duration(timeout) * time.Second)
}

func NewRouterOsClient(rwc io.ReadWriteCloser) (*routeros.Client, error) {
	c, err := routeros.NewClient(rwc)
	if err != nil {
		rwc.Close()
		return nil, err
	}
	
	return c, nil
}

func Login(client *routeros.Client, username, password string) (error) {
	err := client.Login(username, password)
	
	if err != nil {
		client.Close()

		if err.Error() == "EOF" {
			return errors.New("Unable to login. Check the api service?")
		}
			
		return err
	}

	return nil
}
