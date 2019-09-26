package eureka

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/eureka"
	"github.com/go-kit/kit/sd/lb"
	"github.com/hudl/fargo"
	"os"
	"time"
)

var (
	fargoConnection fargo.EurekaConnection
	logger          = log.NewLogfmtLogger(os.Stdout)
)

func init() {
	var fargoConfig fargo.Config
	fargoConfig.Eureka.ServiceUrls = []string{"http://eureka-server-9001:9001/eureka", "http://eureka-server-9002:9002/eureka", "http://eureka-server-9003:9003/eureka"}
	fargoConfig.Eureka.PollIntervalSeconds = 1
	fargoConnection = fargo.NewConnFromConfig(fargoConfig)
}

func getInstancer(serviceName string) *eureka.Instancer {
	instancer := eureka.NewInstancer(&fargoConnection, serviceName, logger)
	return instancer
}

func NewDiscoverEndpoint(ctx context.Context, logger log.Logger) endpoint.Endpoint {
	serviceName := "user-service-consumer"
	duration := 500 * time.Millisecond
	//基于eureka客户端、服务名称等信息，
	instancer := getInstancer(serviceName)
	//针对calculate接口创建sd.Factory
	factory := arithmeticFactory(ctx, "POST", "/api/calculate")

	//使用ceureka连接实例（发现服务系统）、factory创建sd.Factory
	endpointer := sd.NewEndpointer(instancer, factory, logger)

	//创建RoundRibbon负载均衡器
	balancer := lb.NewRoundRobin(endpointer)

	//为负载均衡器增加重试功能，同时该对象为endpoint.Endpoint
	retry := lb.Retry(1, duration, balancer)

	return retry
}
