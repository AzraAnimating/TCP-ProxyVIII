package main

import (
	"github.com/docker/go-connections/sockets"
)

func openProxy() {
	listener, err := sockets.NewTCPSocket("0.0.0.0:12345", nil)
	if err != nil {
		return
	}

	for {
		clientConnection, connectionErr := listener.Accept()
		if connectionErr != nil {
			continue
		}

		backendConnection, isConnected := openTCPConnection("mc.fuxelbau.de:25565")

		if !isConnected {
			continue
		}

		var buffer = make([]byte, 2048)
		var serverBuffer = make([]byte, 2048)

		var readFromClient = func() {
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
					buffer = make([]byte, 2048)
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
				}
			}
		}

		var readFromServer = func() {
			for {
				readBytes, backendConnectionReadError := backendConnection.Read(serverBuffer)
				if err != backendConnectionReadError {
					clientConnection.Close()
					backendConnection.Close()
					//print("Read From Server Error ")
					//println(backendConnectionReadError.Error())
					return
				}

				if readBytes > 0 {
					_, clientWriteConnectionError := clientConnection.Write(format(serverBuffer, readBytes))
					serverBuffer = make([]byte, 2048)
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
				}
			}
		}

		go readFromServer()
		go readFromClient()

	}

}

func format(rawBuffer []byte, readBytes int) []byte {
	var sendPacket []byte
	for i := 0; i < readBytes; i++ {
		sendPacket = append(sendPacket, rawBuffer[i])
	}
	return sendPacket
}
