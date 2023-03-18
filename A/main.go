package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
)

var bConn net.Conn
var mu sync.Mutex

type handler struct{}

// 实现 http接口
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	r.Write(bConn)
	fmt.Println("read response from B")
	resp, err := http.ReadResponse(bufio.NewReader(bConn), r)
	if err != nil {
		w.Write([]byte("error" + err.Error()))
		return
	}
	io.Copy(w, resp.Body)
}

func main() {
	// 启动 gin 服务
	go func() {
		http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {

		})
		err := http.ListenAndServe(":8080", &handler{})
		if err != nil {
			panic(err)
		}
	}()
	// 在 A 上搭建一个监听 socket 的服务端，等待 B 的连接
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	fmt.Println("Waiting for B to connect...")
	for {
		// 等待 B 的连接
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Println("get conn from B")
		bConn = conn
	}
}
