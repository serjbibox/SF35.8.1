package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	proverbsUrl = "https://go-proverbs.github.io/"
	addr        = "0.0.0.0:12345"
	proto       = "tcp4"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	pv, err := getProverbs(proverbsUrl)
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.Listen(proto, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal(err)
			}
			log.Println("установлено новое соединение, клиент:", conn.RemoteAddr())
			go handleConn(conn, pv)
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

// Обработчик. Вызывается для каждого соединения.
func handleConn(conn net.Conn, proverbs []string) {
	defer conn.Close()
	for {
		conn.Write([]byte(proverbs[rand.Intn(len(proverbs))] + "\n\r"))
		time.Sleep(3 * time.Second)
	}
}

// Загружает поговорки по заданному URL
func getProverbs(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	data := parseHtml(string(body))
	return data, nil
}

func parseHtml(text string) (data []string) {
	tkn := html.NewTokenizer(strings.NewReader(text))
	var vals []string
	var isH3, isProverb bool

	for {
		tt := tkn.Next()
		switch {
		case tt == html.ErrorToken:
			return vals
		case tt == html.StartTagToken:
			t := tkn.Token()
			if !isH3 {
				isH3 = t.Data == "h3"
			} else {
				isProverb = t.Data == "a"
			}
		case tt == html.TextToken:
			t := tkn.Token()
			if isProverb {
				vals = append(vals, t.Data)
			}
			isH3 = false
			isProverb = false
		}
	}
}
