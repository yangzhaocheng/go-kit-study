package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {

	// 创建环境变量
	flag.Parse()
	port := 7001
	//创建日志组件
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	ctx := context.Background()

	//创建Endpoint
	discoverEndpoint := MakeDiscoverEndpoint(ctx, logger)

	//创建传输层
	r := MakeHttpHandler(discoverEndpoint)

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	//开始监听
	go func() {
		logger.Log("transport", "HTTP", "addr", port)
		errc <- http.ListenAndServe(":"+strconv.Itoa(port), r)
	}()

	// 开始运行，等待结束
	logger.Log("exit", <-errc)
}
