package main

import (
	"errors"
	"fmt"
	"io"
	"net"

	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)

	go func() {
		defer close(lines)
		defer f.Close()

		line := ""
		for {
			buffer := make([]byte, 8)
			n, err := f.Read(buffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					fmt.Println("End of file reached")
					break
				}

				fmt.Println("Error reading file:", err)
				break
			}

			parts := strings.Split(string(buffer[:n]), "\n")
			line += parts[0]
			if len(parts) > 1 {
				lines <- line
				line = parts[1]
			}
		}
	}()

	return lines
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			return
		}

		ch := getLinesChannel(conn)
		for line := range ch {
			fmt.Printf("read: %s\n", line)
		}
	}
}
