package mark

import (
	"context"
	"github.com/chaudhryfaisal/k8s-webhook-pull-policy/internal/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Marker knows how to mark Kubernetes resources.
type Marker interface {
	Mark(ctx context.Context, obj metav1.Object) error
}

// NewLabelMarker returns a new marker that will mark with labels.
func NewLabelMarker(pullPolicy string, logger log.Logger) Marker {
	policy := v1.PullIfNotPresent
	switch pullPolicy {
	case string(v1.PullAlways):
		policy = v1.PullAlways
	case string(v1.PullNever):
		policy = v1.PullNever
	}
	return labelmarker{PullPolicy: policy, logger: logger}
}

type labelmarker struct {
	PullPolicy v1.PullPolicy
	logger     log.Logger
}

func (l labelmarker) Mark(_ context.Context, obj metav1.Object) error {
	pod, ok := obj.(*v1.Pod)
	if ok {
		for i := 0; i < len(pod.Spec.Containers); i++ {
			if pod.Spec.Containers[i].ImagePullPolicy != l.PullPolicy {
				l.logger.Infof("Updated ImagePullPolicy for pod=%s container=%s from=%s to=%s", pod.Name, pod.Spec.Containers[i].Name, pod.Spec.Containers[i].ImagePullPolicy, l.PullPolicy)
				pod.Spec.Containers[i].ImagePullPolicy = l.PullPolicy
			}
		}
	}
	return nil
}

// DummyMarker is a marker that doesn't do anything.
var DummyMarker Marker = dummyMaker(0)

type dummyMaker int

func (dummyMaker) Mark(_ context.Context, obj metav1.Object) error { return nil }
