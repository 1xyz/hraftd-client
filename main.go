package main

import (
	"github.com/1xyz/hraftd-client/cmd"
	"github.com/1xyz/hraftd-client/config"
	"time"
)

func main() {
	cfg := &config.Config{
		URL:     "http://localhost:11001",
		Timeout: time.Second * 10,
	}

	cmd.Execute(cfg)

	//c := client.NewHttpClient("http://localhost:11001", time.Second * 30)
	//if err := c.Put("hello1", "world"); err != nil {
	//	fmt.Printf("Put error = %v", err)
	//}
	//
	//result, err := c.Get("hello1")
	//if err != nil {
	//	fmt.Printf("Get error = %v", err)
	//	return
	//}
	//fmt.Printf("Get Result = %v", result)
}
