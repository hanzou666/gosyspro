package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

// 速度改善版4: パイプライニング
// 状態を変更しないGETやHEADは並列に処理
// レスポンスの順番はリクエストの順番に
func receive4(conn net.Conn) {
	fmt.Printf("Accept %v\n", conn.RemoteAddr())
	sessionResponses := make(chan chan *http.Response, 50)
	defer close(sessionResponses)

	go writeToConn(sessionResponses, conn)
	reader := bufio.NewReader(conn)
	for {
		//レスポンスを受け取って、セッションのキューに詰める
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		request, err := http.ReadRequest(reader)
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
		sessionResponse := make(chan *http.Response)
		sessionResponses <- sessionResponse
		go handleRequest(request, sessionResponse)
	}
}

func handleRequest(r *http.Request, resultReceiver chan *http.Response) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(dump))
	content := "Hello, World\n"
	response := &http.Response{
		StatusCode:    200,
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: int64(len(content)),
		Body:          ioutil.NopCloser(strings.NewReader(content)),
	}
	resultReceiver <- response
}

func writeToConn(sessionResponses chan chan *http.Response, conn net.Conn) {
	defer conn.Close()
	for sessionResponse := range sessionResponses {
		response := <-sessionResponse
		response.Write(conn)
		close(sessionResponse)
	}

}
