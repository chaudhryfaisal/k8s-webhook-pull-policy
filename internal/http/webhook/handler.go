package webhook

import (
	"context"
	"fmt"
	"net/http"

	kwhhttp "github.com/slok/kubewebhook/v2/pkg/http"
	kwhlog "github.com/slok/kubewebhook/v2/pkg/log"
	kwhmodel "github.com/slok/kubewebhook/v2/pkg/model"
	kwhwebhook "github.com/slok/kubewebhook/v2/pkg/webhook"
	kwhmutating "github.com/slok/kubewebhook/v2/pkg/webhook/mutating"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/chaudhryfaisal/k8s-webhook-pull-policy/internal/log"
)

// kubewebhookLogger is a small proxy to use our logger with Kubewebhook.
type kubewebhookLogger struct {
	log.Logger
}

func (l kubewebhookLogger) WithValues(kv map[string]interface{}) kwhlog.Logger {
	return kubewebhookLogger{Logger: l.Logger.WithKV(kv)}
}
func (l kubewebhookLogger) WithCtxValues(ctx context.Context) kwhlog.Logger {
	return l.WithValues(kwhlog.ValuesFromCtx(ctx))
}
func (l kubewebhookLogger) SetValuesOnCtx(parent context.Context, values map[string]interface{}) context.Context {
	return kwhlog.CtxWithValues(parent, values)
}

// allmark sets up the webhook handler for marking all kubernetes resources using Kubewebhook library.
func (h handler) allMark() (http.Handler, error) {
	mt := kwhmutating.MutatorFunc(func(ctx context.Context, ar *kwhmodel.AdmissionReview, obj metav1.Object) (*kwhmutating.MutatorResult, error) {
		err := h.marker.Mark(ctx, obj)
		if err != nil {
			return nil, fmt.Errorf("could not mark the resource: %w", err)
		}

		return &kwhmutating.MutatorResult{
			MutatedObject: obj,
			Warnings:      []string{"Resource marked with custom labels"},
		}, nil
	})

	logger := kubewebhookLogger{Logger: h.logger.WithKV(log.KV{"lib": "kubewebhook", "webhook": "allMark"})}
	wh, err := kwhmutating.NewWebhook(kwhmutating.WebhookConfig{
		ID:      "allMark",
		Logger:  logger,
		Mutator: mt,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create webhook: %w", err)
	}
	whHandler, err := kwhhttp.HandlerFor(kwhhttp.HandlerConfig{
		Webhook: kwhwebhook.NewMeasuredWebhook(h.metrics, wh),
		Logger:  logger,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create handler from webhook: %w", err)
	}

	return whHandler, nil
}