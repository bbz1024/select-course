package promet

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"
	"select-course/demo7/src/constant/config"
)

func ExtractContext(ctx context.Context) prometheus.Labels {
	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		return prometheus.Labels{
			"traceID": span.TraceID().String(),
			"spanID":  span.SpanID().String(),
			"host":    config.EnvCfg.BaseHost,
		}
	}
	return nil
}
