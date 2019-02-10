package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {

	if len(os.Args) < 2 || len(os.Args) > 2 {
		panic("invalid arguments")
	}
	addr := os.Args[1]

	bioIn := bufio.NewReader(os.Stdin)
	bioOut := bufio.NewWriter(os.Stdout)

	var length uint16

	for {
		_ = binary.Read(bioIn, binary.BigEndian, &length)
		buf := make([]byte, length)
		r, err := bioIn.Read(buf)
		if r == 0 || err != nil {
			time.Sleep(25 * time.Millisecond)
			continue
		}

		conn, err := net.Dial("tcp", addr)
		if err != nil {
			_ = log("error", err.Error())
			_ = binary.Write(bioOut, binary.BigEndian, uint16(2))
			_ = binary.Write(bioOut, binary.BigEndian, uint16(0))
			_ = bioOut.Flush()
			continue
		}

		wConn := bufio.NewWriter(conn)
		_ = binary.Write(wConn, binary.BigEndian, &length)
		_ = binary.Write(wConn, binary.BigEndian, buf[:r])
		err = wConn.Flush()
		if err != nil {
			_ = log("error", err.Error())
			_ = binary.Write(bioOut, binary.BigEndian, uint16(2))
			_ = binary.Write(bioOut, binary.BigEndian, uint16(0))
			_ = bioOut.Flush()
			_ = conn.Close()
			continue
		}

		rConn := bufio.NewReader(conn)
		_ = binary.Read(rConn, binary.BigEndian, &length)
		buf = make([]byte, length)
		r, err = rConn.Read(buf)
		if err != nil {
			_ = log("error", err.Error())
			_ = binary.Write(bioOut, binary.BigEndian, uint16(2))
			_ = binary.Write(bioOut, binary.BigEndian, uint16(0))
			_ = bioOut.Flush()
			_ = conn.Close()
			continue
		}

		_ = binary.Write(bioOut, binary.BigEndian, &length)
		_ = binary.Write(bioOut, binary.BigEndian, buf[:r])
		err = bioOut.Flush()
		if err != nil {
			_ = log("error", err.Error())
			_ = binary.Write(bioOut, binary.BigEndian, uint16(2))
			_ = binary.Write(bioOut, binary.BigEndian, uint16(0))
			_ = bioOut.Flush()
			_ = conn.Close()
			continue
		}

		err = conn.Close()
		if err != nil {
			_ = log("error", err.Error())
			continue
		}
	}
}

func log(level string, message string) error {
	_, err := fmt.Fprintf(os.Stderr, "%s [%s] (%s) %s\n", time.Now().Format("15:04:05.999"), level, "authentication", message)
	return err
}
