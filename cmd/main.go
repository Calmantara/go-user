package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	ginrouter "github.com/Calmantara/go-user/common/infra/gin/router"
)

func main() {
	// call dependencies injection
	conf, server, err := BuildInRuntime()
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		server.SERVE(ginrouter.WithPort(fmt.Sprintf("%v", conf["httpport"])))
		wg.Done()
	}()

	fmt.Printf("runtime:%v\napplication-name:%s\napplication-port:%v\napplication-env:%v\n",
		time.Now().Format(time.RFC850), conf["name"], conf["httpport"], os.Getenv("ENV"))
	wg.Wait()
}
