package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8080/httprpc", strings.NewReader(`{"method":"HelloService.Hello", "params":["星期三"], "id":0}`))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println(err)
	if resp.Body != nil {
		defer func() {
			_ = resp.Body.Close()
		}()
	}
	bs, _ := io.ReadAll(resp.Body)
	fmt.Println(string(bs))
}
