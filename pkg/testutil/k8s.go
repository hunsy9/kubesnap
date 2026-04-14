package testutil

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

// MakeNode creates a corev1.Node with the given name and Ready status.
func MakeNode(name string, ready bool) corev1.Node {
	conditionStatus := corev1.ConditionFalse
	if ready {
		conditionStatus = corev1.ConditionTrue
	}

	return corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Status: corev1.NodeStatus{
			Conditions: []corev1.NodeCondition{
				{
					Type:   corev1.NodeReady,
					Status: conditionStatus,
				},
			},
		},
	}
}

// MakePod creates a corev1.Pod with the given name, namespace, and phase.
func MakePod(name, namespace string, phase corev1.PodPhase) corev1.Pod {
	return corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Status: corev1.PodStatus{
			Phase: phase,
		},
	}
}

// MakeEvent creates a corev1.Event with the given name and event type.
func MakeEvent(name, namespace, eventType string) corev1.Event {
	return corev1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Type: eventType,
	}
}

// MakeNamespace creates a corev1.Namespace with the given name.
func MakeNamespace(name string) corev1.Namespace {
	return corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

// CreateFakeClientset creates a fake.Clientset with the given runtime objects.
func CreateFakeClientset(objects ...runtime.Object) *fake.Clientset {
	return fake.NewSimpleClientset(objects...)
}
