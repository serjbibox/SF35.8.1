package main

import (
	"bufio"
	"math/rand"
	"net"
	"sync"
	"testing"
	"time"
)

func Test_handleConn(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	srv, cl := net.Pipe()
	type args struct {
		conn     net.Conn
		proverbs []string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "test 1",
			args: args{
				conn: srv,
				proverbs: []string{
					"proverb 1",
					"proverb 2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go handleConn(tt.args.conn, tt.args.proverbs)
			var wg sync.WaitGroup
			wg.Add(1)
			rdCnt := 0
			var res []string
			tn := time.Now()
			go func() {
				for {
					reader := bufio.NewReader(cl)
					b, err := reader.ReadBytes('\n')
					if err != nil {
						t.Log(err)
					}
					res = append(res, string(b))
					if rdCnt < 1 {
						rdCnt++
						continue
					}
					cl.Close()
					wg.Done()
					return
				}
			}()
			wg.Wait()
			if (time.Since(tn).Seconds() < 3) || (len(res) != 2) {
				t.Error("тест не пройден")
			}
			srv.Close()
			cl.Close()
		})
	}
}
