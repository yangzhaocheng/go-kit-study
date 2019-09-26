package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {

	flag.Parse()
	port := 7002
	ctx := context.Background()
	errChan := make(chan error)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	fieldKeys := []string{"method"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "raysonxin",
		Subsystem: "arithmetic_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "raysonxin",
		Subsystem: "arithemetic_service",
		Name:      "request_latency",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	var svc Service
	svc = ArithmeticService{}

	//----- service ----
	// add logging middleware
	svc = LoggingMiddleware(logger)(svc)
	svc = Metrics(requestCount, requestLatency)(svc)
	endpoint := MakeArithmeticEndpoint(svc)
	//---- endpoint -----
	ratebucket := rate.NewLimiter(rate.Every(time.Second*1), 100)
	endpoint = NewTokenBucketLimitterWithBuildIn(ratebucket)(endpoint)
	//创建健康检查的Endpoint，未增加限流
	healthEndpoint := MakeHealthCheckEndpoint(svc)

	//把算术运算Endpoint和健康检查Endpoint封装至ArithmeticEndpoints
	endpts := ArithmeticEndpoints{
		ArithmeticEndpoint:  endpoint,
		HealthCheckEndpoint: healthEndpoint,
	}

	//创建http.Handler
	r := MakeHttpHandler(ctx, endpts, logger)

	//创建注册对象
	register := Register(port, "arithmetic", "127.0.0.1")

	go func() {
		fmt.Println("Http Server start at port:" + strconv.Itoa(port))
		//启动前执行注册
		register.Register()
		handler := r
		errChan <- http.ListenAndServe(":"+strconv.Itoa(port), handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	//服务退出取消注册
	register.Deregister()
	fmt.Println(error)
}
