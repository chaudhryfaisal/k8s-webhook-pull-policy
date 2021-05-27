package prometheus_test

import (
	"context"
	"github.com/chaudhryfaisal/k8s-webhook-pull-policy/internal/mutation/prometheus"
	"testing"
	"time"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceMonitorSafer(t *testing.T) {
	tests := map[string]struct {
		minScrapeInterval time.Duration
		servMon           *monitoringv1.ServiceMonitor
		expServMon        *monitoringv1.ServiceMonitor
	}{
		"Having a correct scrape interval should not mutate the service monitor.": {
			minScrapeInterval: 10 * time.Second,
			servMon: &monitoringv1.ServiceMonitor{
				Spec: monitoringv1.ServiceMonitorSpec{
					Endpoints: []monitoringv1.Endpoint{
						{Interval: "15s"},
						{Interval: "20s"},
						{Interval: "30s"},
					},
				},
			},
			expServMon: &monitoringv1.ServiceMonitor{
				Spec: monitoringv1.ServiceMonitorSpec{
					Endpoints: []monitoringv1.Endpoint{
						{Interval: "15s"},
						{Interval: "20s"},
						{Interval: "30s"},
					},
				},
			},
		},

		"Having a incorrect scrape interval should not mutate the service monitor.": {
			minScrapeInterval: 16 * time.Second,
			servMon: &monitoringv1.ServiceMonitor{
				Spec: monitoringv1.ServiceMonitorSpec{
					Endpoints: []monitoringv1.Endpoint{
						{Interval: "30s"},
						{Interval: "15s"},
						{Interval: "20s"},
					},
				},
			},
			expServMon: &monitoringv1.ServiceMonitor{
				Spec: monitoringv1.ServiceMonitorSpec{
					Endpoints: []monitoringv1.Endpoint{
						{Interval: "30s"},
						{Interval: "16s"},
						{Interval: "20s"},
					},
				},
			},
		},

		"Having a service monitor without interval should set the minimum one.": {
			minScrapeInterval: 11 * time.Second,
			servMon: &monitoringv1.ServiceMonitor{
				Spec: monitoringv1.ServiceMonitorSpec{
					Endpoints: []monitoringv1.Endpoint{
						{Interval: "30s"},
						{Interval: "15s"},
						{Interval: "20s"},
						{},
					},
				},
			},
			expServMon: &monitoringv1.ServiceMonitor{
				Spec: monitoringv1.ServiceMonitorSpec{
					Endpoints: []monitoringv1.Endpoint{
						{Interval: "30s"},
						{Interval: "15s"},
						{Interval: "20s"},
						{Interval: "11s"},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			s := prometheus.NewServiceMonitorSafer(test.minScrapeInterval)
			err := s.EnsureSafety(context.TODO(), test.servMon)
			require.NoError(err)

			assert.Equal(test.expServMon, test.servMon)
		})
	}
}
