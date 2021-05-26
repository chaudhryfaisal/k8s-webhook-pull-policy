package mark

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Marker knows how to mark Kubernetes resources.
type Marker interface {
	Mark(ctx context.Context, obj metav1.Object) error
}

// NewLabelMarker returns a new marker that will mark with labels.
func NewLabelMarker(pullPolicy string) Marker {
	policy := v1.PullIfNotPresent
	switch pullPolicy {
	case string(v1.PullAlways):
		policy = v1.PullAlways
	case string(v1.PullNever):
		policy = v1.PullNever
	}
	return labelmarker{PullPolicy: policy}
}

type labelmarker struct {
	PullPolicy v1.PullPolicy
}

func (l labelmarker) Mark(_ context.Context, obj metav1.Object) error {
	pod, ok := obj.(*v1.Pod)
	if ok {
		for i := 0; i < len(pod.Spec.Containers); i++ {
			pod.Spec.Containers[i].ImagePullPolicy = l.PullPolicy
		}
	}
	return nil
}

// DummyMarker is a marker that doesn't do anything.
var DummyMarker Marker = dummyMaker(0)

type dummyMaker int

func (dummyMaker) Mark(_ context.Context, obj metav1.Object) error { return nil }
