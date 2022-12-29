package main

import "net"

func openTCPConnection(address string) (net.Conn, bool) {
	tcpAddress, resolveError := net.ResolveTCPAddr("tcp", address)
	if resolveError != nil {
		return nil, false
	}

	connection, connectError := net.DialTCP("tcp", nil, tcpAddress)
	if connectError != nil {
		return nil, false
	}

	return connection, true
}
