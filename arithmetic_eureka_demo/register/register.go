package main

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/eureka"
	"github.com/hudl/fargo"
	"os"
	"strconv"
)

func Register(port int, app string, ip string) *eureka.Registrar {
	logger := log.NewLogfmtLogger(os.Stdout)
	var fargoConfig fargo.Config
	fargoConfig.Eureka.ServiceUrls = []string{"http://eureka-server-9001:9001/eureka", "http://eureka-server-9002:9002/eureka", "http://eureka-server-9003:9003/eureka"}
	fargoConfig.Eureka.PollIntervalSeconds = 1
	instance := &fargo.Instance{
		XMLName:          struct{}{},
		HostName:         ip+":"+app+":"+strconv.Itoa(port),
		App:              app,
		IPAddr:           ip,
		VipAddress:       "",
		SecureVipAddress: "",
		Status:           fargo.UP,
		Overriddenstatus: "",
		Port:             port,
		PortJ:            fargo.Port{},
		SecurePort:       0,
		SecurePortJ:      fargo.Port{},
		HomePageUrl:      "",
		StatusPageUrl:    "",
		HealthCheckUrl:   "",
		CountryId:        0,
		DataCenterInfo:   fargo.DataCenterInfo{Name: fargo.MyOwn},
		LeaseInfo:        fargo.LeaseInfo{RenewalIntervalInSecs: 1},
		Metadata:         fargo.InstanceMetadata{},
		UniqueID:         genUid,
	}
	fargoConnection := fargo.NewConnFromConfig(fargoConfig)
	register := eureka.NewRegistrar(&fargoConnection, instance, logger)
	return register
}

func genUid(i fargo.Instance) string {
	return i.IPAddr+":"+i.App+":"+strconv.Itoa(i.Port)
}