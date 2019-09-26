package consumer

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"micro-service/eureka"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {



	//创建日志组件
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	ctx := context.Background()

	//创建Endpoint
	discoverEndpoint := eureka.NewDiscoverEndpoint(ctx, logger)

	//创建传输层
	r := eureka.NewHttpHandler(discoverEndpoint)

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	//开始监听
	go func() {
		logger.Log("transport", "HTTP", "addr", "7001")
		errc <- http.ListenAndServe(":7001", r)
	}()

	// 开始运行，等待结束
	logger.Log("exit", <-errc)
}
