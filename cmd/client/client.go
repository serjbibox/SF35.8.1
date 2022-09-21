package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	addr  = "localhost:12345"
	proto = "tcp4"
)

func main() {
	conn, err := net.Dial(proto, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	pvCnt := 0
	go func() {
		for {
			pvb, err := reader.ReadBytes('\n')
			if err != nil {
				log.Fatal(err)
			}
			pvCnt++
			msg := strings.Trim(string(pvb), "\n")
			msg = strings.Trim(msg, "\r")
			fmt.Printf("Поговорка №%d: %s\n", pvCnt, msg)
		}
	}()
	fmt.Println("Для закрытия программы введите exit")
	s := ""
	for {
		fmt.Scanln(&s)
		switch s {
		case "exit":
			log.Println("exit now")
			return
		}
	}
}
