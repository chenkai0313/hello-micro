package jaeger

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"io"

	"github.com/micro/go-micro/v2/util/log"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

func NewJaegerTracer(serviceName, addr string) (opentracing.Tracer, io.Closer, error) {
	// Sample configuration for testing. Use constant sampling to sample every trace
	// and enable LogSpan to log every span via configured Logger.
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeProbabilistic,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	cfg.ServiceName = serviceName

	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := &jaegerLogger{}
	jMetricsFactory := metrics.NullFactory

	sender, err := jaeger.NewUDPTransport(addr, 0)
	if err != nil {
		log.Errorf("could not initialize jaeger sender: %s", err.Error())
		return nil, nil, err
	}

	reporter := jaeger.NewRemoteReporter(sender)

	return cfg.NewTracer(
		config.Logger(jLogger),
		config.Metrics(jMetricsFactory),
		config.Reporter(reporter),
	)

}

type jaegerLogger struct{}

// Error logs a message at error priority
func (l *jaegerLogger) Error(msg string) {
	log.WithLevel(log.LevelError, msg)
}

// Infof logs a message at info priority
func (l *jaegerLogger) Infof(msg string, args ...interface{}) {
	log.WithLevelf(log.LevelInfo, msg, args...)
}
