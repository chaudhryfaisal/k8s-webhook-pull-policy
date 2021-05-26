package prometheus

import (
	"github.com/chaudhryfaisal/k8s-webhook-pull-policy/internal/http/webhook"
	"github.com/prometheus/client_golang/prometheus"
	gohttpmetrics "github.com/slok/go-http-metrics/metrics"
	gohttpmetricsprometheus "github.com/slok/go-http-metrics/metrics/prometheus"
	whprometheus "github.com/slok/kubewebhook/v2/pkg/metrics/prometheus"

)

// Types used to avoid collisions with the same interface naming.
type httpRecorder = gohttpmetrics.Recorder
type webhookRecorder = whprometheus.Recorder

// Recorder satisfies multiple metrics recording interfaces using a Prometheus backend.
type Recorder struct {
	httpRecorder
	webhookRecorder
}

// NewRecorder returns a new Prometheus Recorder.
func NewRecorder(reg prometheus.Registerer) Recorder {
	// TODO error,
	rec, _ := whprometheus.NewRecorder(whprometheus.RecorderConfig{Registry: reg})

	return Recorder{
		httpRecorder:    gohttpmetricsprometheus.NewRecorder(gohttpmetricsprometheus.Config{Registry: reg}),
		webhookRecorder: *rec,
	}
}

// Interface assertion.
var _ webhook.MetricsRecorder = Recorder{}
