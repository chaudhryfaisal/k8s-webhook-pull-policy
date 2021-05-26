package mark_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/chaudhryfaisal/k8s-webhook-pull-policy/internal/mutation/mark"
)

func TestLabelMarkerMark(t *testing.T) {
	tests := map[string]struct {
		marks  string
		obj    metav1.Object
		expObj metav1.Object
	}{
		"Should set ImagePullPolicy to Always from PullNever": {
			marks: "Always",
			obj: &corev1.Pod{
				Spec: corev1.PodSpec{Containers: []corev1.Container{{ImagePullPolicy: corev1.PullNever}}},
				ObjectMeta: metav1.ObjectMeta{
					Name: string(corev1.PullAlways),
				},
			},
			expObj: &corev1.Pod{
				Spec: corev1.PodSpec{Containers: []corev1.Container{{ImagePullPolicy: corev1.PullAlways}}},
				ObjectMeta: metav1.ObjectMeta{
					Name: string(corev1.PullAlways),
				},
			},
		},
		"Should set ImagePullPolicy to IfNotPresent from PullAlways": {
			marks: "IfNotPresent",
			obj: &corev1.Pod{
				Spec: corev1.PodSpec{Containers: []corev1.Container{{ImagePullPolicy: corev1.PullNever}}},
				ObjectMeta: metav1.ObjectMeta{
					Name: string(corev1.PullIfNotPresent),
				},
			},
			expObj: &corev1.Pod{
				Spec: corev1.PodSpec{Containers: []corev1.Container{{ImagePullPolicy: corev1.PullIfNotPresent}}},
				ObjectMeta: metav1.ObjectMeta{
					Name: string(corev1.PullIfNotPresent),
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			m := mark.NewLabelMarker(test.marks)

			err := m.Mark(context.TODO(), test.obj)
			require.NoError(err)

			assert.Equal(test.expObj, test.obj)
		})
	}
}
