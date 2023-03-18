package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func createConnection() (conn net.Conn, err error) {
	// 连接到 A 的 socket
	connA, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		return nil, err
	}
	return connA, nil
}

func main() {
	go func() {
		for {
			connA, err := createConnection()
			if err != nil {
				fmt.Println("err", err)
				time.Sleep(time.Second)
				continue
			}
			err = handleConnection(connA)
			if err != nil {
				fmt.Println("err", err)
				connA.Close()
				time.Sleep(time.Second)
			}
		}
	}()
	http.HandleFunc("/hello", hello)
	err := http.ListenAndServe(":8070", nil)
	if err != nil {
		panic(err)
	}

}

func hellodata(r *http.Request) string {
	return time.Now().String() + "hello world"
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handle in B")
	w.Write([]byte(hellodata(r) + " IN B"))
}

func handleConnection(conn net.Conn) (err error) {
	fmt.Println("handle form A")
	// 读取HTTP请求
	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		// 处理读取请求失败的情况
		fmt.Println("err", err)
		return
	}
	// 处理HTTP请求
	var bodyBuf bytes.Buffer
	bodyBuf.WriteString(hellodata(req) + "FROM A")
	// 将HTTP响应写回到连接中
	resp := &http.Response{
		StatusCode: http.StatusOK,
		ProtoMajor: req.ProtoMajor,
		ProtoMinor: req.ProtoMinor,
		Header:     make(http.Header),
	}
	resp.Header.Set("Content-Type", "text/plain")
	reader := bufio.NewReader(&bodyBuf)
	resp.ContentLength = int64(bodyBuf.Len())
	resp.Body = ioutil.NopCloser(reader)
	resp.Write(conn)
	return nil
}
