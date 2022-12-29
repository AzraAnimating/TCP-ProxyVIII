package main

import (
	"net"
)

func openProxy() {

	tcpAddress, addressError := net.ResolveTCPAddr("tcp", "0.0.0.0:25565")
	if addressError != nil {
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddress)
	if err != nil {
		return
	}

	for {
		clientConnection, connectionErr := listener.Accept()
		if connectionErr != nil {
			continue
		}

		openBackend(clientConnection)
	}

}

func openBackend(clientConnection net.Conn) {
	backendConnection, isConnected := openTCPConnection("45.137.68.201:25561")

	if !isConnected {
		return
	}

	go initProxy(clientConnection, backendConnection)
}

func initProxy(clientConnection net.Conn, backendConnection net.Conn) {
	var buffer = make([]byte, 128)
	var serverBuffer = make([]byte, 128)

	var readFromClient = func() {
		defer clientConnection.Close()
		defer backendConnection.Close()
		for {
			readBytes, clientConnectionReadError := clientConnection.Read(buffer)
			if clientConnectionReadError != nil {
				clientConnection.Close()
				backendConnection.Close()
				//print("Read From Client Error ")
				//println(clientConnectionReadError.Error())
				return
			}

			if readBytes > 0 {
				_, backendConnectionWriteError := backendConnection.Write(format(buffer, readBytes))
				buffer = make([]byte, 128)
				if backendConnectionWriteError != nil {
					clientConnection.Close()
					backendConnection.Close()
					//print("Write to Server Error ")
					//println(backendConnectionWriteError.Error())
					return
				}
			} else {
				clientConnection.Close()
				backendConnection.Close()
				return
			}
		}
	}

	var readFromServer = func() {
		defer backendConnection.Close()
		defer clientConnection.Close()
		for {

			readBytes, backendConnectionReadError := backendConnection.Read(serverBuffer)
			if backendConnectionReadError != nil {
				clientConnection.Close()
				backendConnection.Close()
				//print("Read From Server Error ")
				//println(backendConnectionReadError.Error())
				return
			}

			if readBytes > 0 {
				_, clientWriteConnectionError := clientConnection.Write(format(serverBuffer, readBytes))
				serverBuffer = make([]byte, 128)
				if clientWriteConnectionError != nil {
					clientConnection.Close()
					backendConnection.Close()
					//print("Write to Client Error ")
					//println(clientWriteConnectionError.Error())
					return
				}
			} else {
				clientConnection.Close()
				backendConnection.Close()
				return
			}
		}
	}

	go readFromServer()
	go readFromClient()
}

func format(rawBuffer []byte, readBytes int) []byte {
	return rawBuffer[:readBytes]
}
