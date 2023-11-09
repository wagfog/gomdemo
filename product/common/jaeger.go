package common

import (
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// create chain tracing
// 这段函数是一个用于创建Jaeger分布式追踪器的工厂函数
// 接收参数：函数接收两个参数，分别是serviceName和addr，分别代表服务的名称和地址。
func NewTracer(serviceName string, addr string) (opentracing.Tracer, io.Closer, error) {
	//创建配置：通过config.Configuration结构体创建了配置cfg，其中包括了服务名称、采样器配置和报告器配置。
	cfg := config.Configuration{
		ServiceName: serviceName,
		//采样器配置
		Sampler: &config.SamplerConfig{
			//这里使用了常量采样器jaeger.SamplerTypeConst，参数为1，表示对所有的span进行采样。
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		//报告器配置：
		Reporter: &config.ReporterConfig{
			BufferFlushInterval: 1 * time.Second, //报告器配置中指定了缓冲刷新间隔
			LogSpans:            true,            //是否记录span日志
			LocalAgentHostPort:  addr,            //本地代理主机端口
		},
	}
	return cfg.NewTracer() //返回追踪器：最后通过cfg.NewTracer()创建并返回了一个Jaeger追踪器实例。
}
