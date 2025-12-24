package network

import (
	"fmt"
	"time"
	"net"
	"luna_l2/shared"
)

// Memory layout:
	// 0 - 9: command registers
	// 10 - 13: outbound ptr
	// 14 - 17: inbound ptr

var MemoryNetwork[22]byte
var Connections []net.Conn

func waitForAcceptBit() {
	for {
		if MemoryNetwork[2] == 0x00 {
			time.Sleep(15 * time.Millisecond)
		} else {
			MemoryNetwork[2] = 0x00
			break
		}
	}
}

func waitForAcceptBitWithCheck(connection net.Conn) bool {
	for {
		if MemoryNetwork[0] == 0 {
			connection.Close()
			return true
		}
		if MemoryNetwork[2] == 0x00 {
			time.Sleep(15 * time.Millisecond)
		} else {
			MemoryNetwork[2] = 0x00
			break
		}
	}
	return false
}

func NetHandleConn(connection net.Conn) {
	buf := make([]byte, 2048)

	termed := waitForAcceptBitWithCheck(connection)
	if termed == true {
		return
	}

	top:
	command := MemoryNetwork[5]

	if MemoryNetwork[0] == 0x00 {
		connection.Close()
		return
	}

	switch command {
	case 0x00:
		connection.Close()
		goto done
	case 0x01:	
		n, err := connection.Read(buf)
		if err == nil && n > 0 {
			output := uint32(MemoryNetwork[10]) << 24 | uint32(MemoryNetwork[11]) << 16 | uint32(MemoryNetwork[12]) << 8 | uint32(MemoryNetwork[13])
			for i := 0; i < n; i++ {
				shared.MapperWrite(output + uint32(i), buf[i])
			}
		} else {
			fmt.Println("luna-l2: network controller error: ", err)
			return
		}	
		MemoryNetwork[3] = 0x01
		termed := waitForAcceptBitWithCheck(connection)
		if termed == true {
			return
		}
		MemoryNetwork[3] = 0x00
	case 0x02:
		input := uint32(MemoryNetwork[14]) << 24 | uint32(MemoryNetwork[15]) << 16 | uint32(MemoryNetwork[16]) << 8 | uint32(MemoryNetwork[17])
		size := uint32(MemoryNetwork[18]) << 24 | uint32(MemoryNetwork[19]) << 16 | uint32(MemoryNetwork[20]) << 8 | uint32(MemoryNetwork[21])
		buf := make([]byte, size)
		j := 0
		for i := input; i < input + size; i++ {
			buf[j] = shared.Mapper(i)
			j++
		}

		_, err := connection.Write(buf)
		if err != nil {
			fmt.Println("luna-l2: network controller error: ", err)
			return	
		}

		MemoryNetwork[3] = 0x01
		termed := waitForAcceptBitWithCheck(connection)
		if termed == true {
			return
		}
		MemoryNetwork[3] = 0x00	
	}
	
	goto top
	done:	
}

func NetController() {
	// copy(MemoryNetwork[4107:], []byte("enp0s0")) // Interface name
	for {
		if MemoryNetwork[0] != 0x00 {	
			connType := ""
			switch MemoryNetwork[1] {
			case 0x00:
				connType = "tcp" // TCP client
			case 0x01:
				connType = "tcp" // TCP server
			default:
				fmt.Println("luna-l2: invalid NIC controller mode")
			}
			
			if MemoryNetwork[1] == 0x00 {	
				MemoryNetwork[0] = 0x00
				first := fmt.Sprintf("%d", uint8(MemoryNetwork[2]))
				second := fmt.Sprintf("%d", uint8(MemoryNetwork[3]))
				third := fmt.Sprintf("%d", uint8(MemoryNetwork[4]))
				fourth := fmt.Sprintf("%d", uint8(MemoryNetwork[5]))
				port := fmt.Sprintf("%d", uint16(MemoryNetwork[6]) << 8 | uint16(MemoryNetwork[7]))
				timeout := uint16(MemoryNetwork[8]) << 8 | uint16(MemoryNetwork[9])

				addr := first + "." + second + "." + third + "." + fourth + ":" + port	
				switch connType {
				case "tcp":
					conn, err := net.Dial(connType, addr)	
					if err != nil {
						fmt.Println("luna-l2: network controller error: ", err)	
						break
					}
					conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
					defer conn.Close()
				
					input := uint32(MemoryNetwork[14]) << 24 | uint32(MemoryNetwork[15]) << 16 | uint32(MemoryNetwork[16]) << 8 | uint32(MemoryNetwork[17])
					size := uint32(MemoryNetwork[18]) << 24 | uint32(MemoryNetwork[19]) << 16 | uint32(MemoryNetwork[20]) << 8 | uint32(MemoryNetwork[21])
					buf := make([]byte, size)
					j := 0
					for i := input; i < input + size; i++ {
						buf[j] = shared.Mapper(i)
						j++
					}
					_, err = conn.Write(buf)
					if err != nil {
						fmt.Println("luna-l2: network controller error: ", err)	
						break
					}

					buf = make([]byte, 2048)
					output := uint32(MemoryNetwork[10]) << 24 | uint32(MemoryNetwork[11]) << 16 | uint32(MemoryNetwork[12]) << 8 | uint32(MemoryNetwork[13])
					for i := 0; i < len(buf); i++ {
						shared.MapperWrite(output + uint32(i), buf[i])
					}

					conn.Close()	
				}
			} else if MemoryNetwork[1] == 0x01 {	
				port := fmt.Sprintf("%d", uint16(MemoryNetwork[6]) << 8 | uint16(MemoryNetwork[7]))

				listener, err := net.Listen("tcp", ":" + port)
				if err != nil {
					fmt.Println("luna-l2: network controller error: ", err)
					return
				}
				defer listener.Close()

				for {
					if MemoryNetwork[0] == 0x00 {
						break
					}	
					conn, err := listener.Accept()	
					MemoryNetwork[4] = 0x01
					waitForAcceptBit()
					MemoryNetwork[4] = 0x00	
					if err != nil {
						fmt.Println("luna-l2: network controller error: ", err)
						return
					}
					NetHandleConn(conn)
				}	
				listener.Close()
				MemoryNetwork[0] = 0x00
			}
		}
		time.Sleep(15 * time.Millisecond)
	}
}
