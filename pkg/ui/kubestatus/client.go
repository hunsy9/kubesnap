package kubestatus

import (
	"k8s.io/client-go/discovery"
	authv1client "k8s.io/client-go/kubernetes/typed/authentication/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
)

// KubernetesClient is an interface that abstracts the kubernetes.Clientset
// to allow for easier testing with fake clients.
type KubernetesClient interface {
	CoreV1() corev1client.CoreV1Interface
	AuthenticationV1() authv1client.AuthenticationV1Interface
	Discovery() discovery.DiscoveryInterface
}
