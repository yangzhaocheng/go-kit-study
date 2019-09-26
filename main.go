package main

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	endpoint2 "micro-service/endpoint"
	"micro-service/eureka"
	"micro-service/service"
	"micro-service/service/impl"
	"micro-service/transport"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	errChan := make(chan error)

	var svc service.Service = impl.ArithmeticService{}
	endpoint := endpoint2.NewArithmeticEndpoint(svc)
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	r := transport.NewHttpHandler(ctx, endpoint, logger)
	register := eureka.Register(9000, "arithmetic", "127.0.0.1")
	go func() {
		fmt.Println("http server start at port:9000")
		register.Register()
		handler := r
		errChan <- http.ListenAndServe(":9000", handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println(<-errChan)
	register.Deregister()
}
