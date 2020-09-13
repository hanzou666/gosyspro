package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"
)

func main() {
	implType := os.Getenv("IMPL_TYPE")
	// TCPでlocalhostの8888を使用することを宣言
	listener, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	fmt.Println("Server is running at localhost:8888")
	for {
		// 待ち受け状態にする。net.Connを受け取るとセッション成立
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		}
		switch implType {
		case "1":
			go receive1(&conn)
		case "2":
			continue
		case "3":
			continue
		case "4":
			go receive4(conn)
		case "5":
			// keep aliveなしでパイプライニングはするバージョンを書くと勉強になりそう
			continue
		default:
			go receive(&conn)
		}
	}
}

// HTTP/1.0 の基本形
func receive(conn *net.Conn) {
	fmt.Printf("Accept %v\n", (*conn).RemoteAddr())
	// リクエストの解析はサボる
	request, err := http.ReadRequest(
		bufio.NewReader(*conn))
	if err != nil {
		fmt.Println(err)
	}
	//デバッグ用
	dump, err := httputil.DumpRequest(request, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(dump))
	response := http.Response{
		StatusCode: 200,
		ProtoMajor: 1,
		ProtoMinor: 0,
		Header:     nil,
		Body:       ioutil.NopCloser(strings.NewReader("Hello, World\n")),
	}
	if err := response.Write(*conn); err != nil {
		fmt.Println(err)
	}
	if err := (*conn).Close(); err != nil {
		fmt.Println(err)
	}
}

// 速度改善1: Keep-Alive対応
func receive1(conn *net.Conn) {
	defer (*conn).Close()
	fmt.Printf("Accept %v\n", (*conn).RemoteAddr())
	for {
		// タイムアウトの設定
		(*conn).SetReadDeadline(time.Now().Add(5 * time.Second))
		request, err := http.ReadRequest(bufio.NewReader(*conn))
		// タイムアウトかソケットクローズ時は終了する
		if err != nil {
			neterr, ok := err.(net.Error) // ダウンキャスト
			if ok && neterr.Timeout() {
				fmt.Println("Timout")
				break
			} else if err == io.EOF {
				break
			}
			fmt.Println(err)
		}

		dump, err := httputil.DumpRequest(request, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(dump))
		// もしHTTP/1.1より古いのをセットしていると、writeを呼んだときにcloseヘッダを付与してしまう
		response := http.Response{
			StatusCode: 200,
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     nil,
			Body:       ioutil.NopCloser(strings.NewReader("Hello, World\n")),
		}
		if err := response.Write(*conn); err != nil {
			fmt.Println(err)
		}
	}
	// リクエストの解析はサボる
	//デバッグ用

}
