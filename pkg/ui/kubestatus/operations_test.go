package kubestatus

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/hunsy9/kubesnap/pkg/testutil"
)

func TestGetNodeInfo(t *testing.T) {
	tests := []struct {
		name       string
		nodes      []corev1.Node
		wantTotal  int
		wantReady  int
		wantErrMsg string
	}{
		{
			name: "all nodes ready",
			nodes: []corev1.Node{
				testutil.MakeNode("node1", true),
				testutil.MakeNode("node2", true),
				testutil.MakeNode("node3", true),
			},
			wantTotal:  3,
			wantReady:  3,
			wantErrMsg: "",
		},
		{
			name: "some nodes not ready",
			nodes: []corev1.Node{
				testutil.MakeNode("node1", true),
				testutil.MakeNode("node2", false),
				testutil.MakeNode("node3", true),
			},
			wantTotal:  3,
			wantReady:  2,
			wantErrMsg: "",
		},
		{
			name: "all nodes not ready",
			nodes: []corev1.Node{
				testutil.MakeNode("node1", false),
				testutil.MakeNode("node2", false),
			},
			wantTotal:  2,
			wantReady:  0,
			wantErrMsg: "",
		},
		{
			name:       "no nodes",
			nodes:      []corev1.Node{},
			wantTotal:  0,
			wantReady:  0,
			wantErrMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build runtime objects from nodes
			var objects []runtime.Object
			for i := range tt.nodes {
				objects = append(objects, &tt.nodes[i])
			}

			fakeClient := testutil.CreateFakeClientset(objects...)
			model := NewStatusModel(fakeClient, "test-cluster", "default", "https://localhost:6443")

			msg := model.getNodeInfo()
			nodeMsg, ok := msg.(NodeInfoMsg)
			if !ok {
				t.Fatalf("expected NodeInfoMsg, got %T", msg)
			}

			if nodeMsg.TotalNodes != tt.wantTotal {
				t.Errorf("TotalNodes = %v, want %v", nodeMsg.TotalNodes, tt.wantTotal)
			}
			if nodeMsg.ReadyNodes != tt.wantReady {
				t.Errorf("ReadyNodes = %v, want %v", nodeMsg.ReadyNodes, tt.wantReady)
			}
			if nodeMsg.NodeAPIStatusErr != tt.wantErrMsg {
				t.Errorf("NodeAPIStatusErr = %v, want %v", nodeMsg.NodeAPIStatusErr, tt.wantErrMsg)
			}
		})
	}
}

func TestGetPodInfo(t *testing.T) {
	tests := []struct {
		name        string
		pods        []corev1.Pod
		wantTotal   int
		wantRunning int
		wantErrMsg  string
	}{
		{
			name: "all pods running",
			pods: []corev1.Pod{
				testutil.MakePod("pod1", "default", corev1.PodRunning),
				testutil.MakePod("pod2", "default", corev1.PodRunning),
				testutil.MakePod("pod3", "kube-system", corev1.PodRunning),
			},
			wantTotal:   3,
			wantRunning: 3,
			wantErrMsg:  "",
		},
		{
			name: "mixed pod phases",
			pods: []corev1.Pod{
				testutil.MakePod("pod1", "default", corev1.PodRunning),
				testutil.MakePod("pod2", "default", corev1.PodPending),
				testutil.MakePod("pod3", "kube-system", corev1.PodFailed),
				testutil.MakePod("pod4", "kube-system", corev1.PodRunning),
			},
			wantTotal:   4,
			wantRunning: 2,
			wantErrMsg:  "",
		},
		{
			name: "no running pods",
			pods: []corev1.Pod{
				testutil.MakePod("pod1", "default", corev1.PodPending),
				testutil.MakePod("pod2", "default", corev1.PodFailed),
			},
			wantTotal:   2,
			wantRunning: 0,
			wantErrMsg:  "",
		},
		{
			name:        "no pods",
			pods:        []corev1.Pod{},
			wantTotal:   0,
			wantRunning: 0,
			wantErrMsg:  "",
		},
		{
			name: "pods in multiple namespaces",
			pods: []corev1.Pod{
				testutil.MakePod("pod1", "default", corev1.PodRunning),
				testutil.MakePod("pod2", "kube-system", corev1.PodRunning),
				testutil.MakePod("pod3", "monitoring", corev1.PodSucceeded),
			},
			wantTotal:   3,
			wantRunning: 2,
			wantErrMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var objects []runtime.Object
			for i := range tt.pods {
				objects = append(objects, &tt.pods[i])
			}

			fakeClient := testutil.CreateFakeClientset(objects...)
			model := NewStatusModel(fakeClient, "test-cluster", "default", "https://localhost:6443")

			msg := model.getPodInfo()
			podMsg, ok := msg.(PodInfoMsg)
			if !ok {
				t.Fatalf("expected PodInfoMsg, got %T", msg)
			}

			if podMsg.TotalPods != tt.wantTotal {
				t.Errorf("TotalPods = %v, want %v", podMsg.TotalPods, tt.wantTotal)
			}
			if podMsg.RunningPods != tt.wantRunning {
				t.Errorf("RunningPods = %v, want %v", podMsg.RunningPods, tt.wantRunning)
			}
			if podMsg.PodAPIStatusErr != tt.wantErrMsg {
				t.Errorf("PodAPIStatusErr = %v, want %v", podMsg.PodAPIStatusErr, tt.wantErrMsg)
			}
		})
	}
}

func TestGetEventInfo(t *testing.T) {
	tests := []struct {
		name         string
		events       []corev1.Event
		wantWarnings int
		wantErrMsg   string
	}{
		{
			name: "no warning events",
			events: []corev1.Event{
				testutil.MakeEvent("event1", "default", "Normal"),
				testutil.MakeEvent("event2", "default", "Normal"),
			},
			wantWarnings: 0,
			wantErrMsg:   "",
		},
		{
			name: "some warning events",
			events: []corev1.Event{
				testutil.MakeEvent("event1", "default", "Normal"),
				testutil.MakeEvent("event2", "default", "Warning"),
				testutil.MakeEvent("event3", "kube-system", "Warning"),
			},
			wantWarnings: 2,
			wantErrMsg:   "",
		},
		{
			name: "all warning events",
			events: []corev1.Event{
				testutil.MakeEvent("event1", "default", "Warning"),
				testutil.MakeEvent("event2", "kube-system", "Warning"),
			},
			wantWarnings: 2,
			wantErrMsg:   "",
		},
		{
			name:         "no events",
			events:       []corev1.Event{},
			wantWarnings: 0,
			wantErrMsg:   "",
		},
		{
			name: "events in multiple namespaces",
			events: []corev1.Event{
				testutil.MakeEvent("event1", "default", "Warning"),
				testutil.MakeEvent("event2", "kube-system", "Normal"),
				testutil.MakeEvent("event3", "monitoring", "Warning"),
			},
			wantWarnings: 2,
			wantErrMsg:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var objects []runtime.Object
			for i := range tt.events {
				objects = append(objects, &tt.events[i])
			}

			fakeClient := testutil.CreateFakeClientset(objects...)
			model := NewStatusModel(fakeClient, "test-cluster", "default", "https://localhost:6443")

			msg := model.getEventInfo()
			eventMsg, ok := msg.(EventInfoMsg)
			if !ok {
				t.Fatalf("expected EventInfoMsg, got %T", msg)
			}

			if eventMsg.WarningEvents != tt.wantWarnings {
				t.Errorf("WarningEvents = %v, want %v", eventMsg.WarningEvents, tt.wantWarnings)
			}
			if eventMsg.EventAPIStatusErr != tt.wantErrMsg {
				t.Errorf("EventAPIStatusErr = %v, want %v", eventMsg.EventAPIStatusErr, tt.wantErrMsg)
			}
		})
	}
}

func TestGetAuthInfo(t *testing.T) {
	// Note: Testing getAuthInfo with fake client is limited because
	// fake.Clientset doesn't fully implement SelfSubjectReviews.Create behavior.
	// The fake client will return an empty SelfSubjectReview without proper UserInfo.
	// This test verifies that the function doesn't panic and returns the expected message type.

	t.Run("fake client returns empty auth info", func(t *testing.T) {
		fakeClient := testutil.CreateFakeClientset()
		model := NewStatusModel(fakeClient, "test-cluster", "default", "https://localhost:6443")

		msg := model.getAuthInfo()
		authMsg, ok := msg.(AuthInfoMsg)
		if !ok {
			t.Fatalf("expected AuthInfoMsg, got %T", msg)
		}

		// fake client returns empty result without error
		if authMsg.Err != nil {
			t.Errorf("unexpected error: %v", authMsg.Err)
		}
	})
}

// Note: Testing checkApiHealth with fake client is not possible because
// fake.Clientset's Discovery().RESTClient() returns nil, causing a panic.
// The checkApiHealth function requires a real REST client to make HTTP requests
// to the /healthz endpoint. This would require integration tests with a real
// or mock HTTP server.
//
// For unit testing checkApiHealth, consider:
// 1. Using httptest.Server to mock the API server
// 2. Creating a custom mock that implements the required interfaces
// 3. Running integration tests against a real cluster

func TestCheckForUpdate(t *testing.T) {
	// Note: checkForUpdate calls external version checking which we don't want
	// to do in unit tests. This test just verifies the function returns the correct type.

	t.Run("returns update available message type", func(t *testing.T) {
		fakeClient := testutil.CreateFakeClientset()
		model := NewStatusModel(fakeClient, "test-cluster", "default", "https://localhost:6443")

		msg := model.checkForUpdate()
		_, ok := msg.(UpdateAvailableMsg)
		if !ok {
			t.Fatalf("expected UpdateAvailableMsg, got %T", msg)
		}
	})
}

func TestNewStatusModel(t *testing.T) {
	tests := []struct {
		name        string
		clusterName string
		namespace   string
		endpoint    string
	}{
		{
			name:        "basic model creation",
			clusterName: "test-cluster",
			namespace:   "default",
			endpoint:    "https://localhost:6443",
		},
		{
			name:        "custom namespace",
			clusterName: "prod-cluster",
			namespace:   "kube-system",
			endpoint:    "https://10.0.0.100:6443",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := testutil.CreateFakeClientset()
			model := NewStatusModel(fakeClient, tt.clusterName, tt.namespace, tt.endpoint)

			if model == nil {
				t.Fatal("expected non-nil model")
			}
			if model.clusterName != tt.clusterName {
				t.Errorf("clusterName = %v, want %v", model.clusterName, tt.clusterName)
			}
			if model.namespace != tt.namespace {
				t.Errorf("namespace = %v, want %v", model.namespace, tt.namespace)
			}
			if model.endpoint != tt.endpoint {
				t.Errorf("endpoint = %v, want %v", model.endpoint, tt.endpoint)
			}
			if !model.isHealthLoading {
				t.Error("expected isHealthLoading to be true")
			}
			if !model.isAuthLoading {
				t.Error("expected isAuthLoading to be true")
			}
			if !model.isNodeLoading {
				t.Error("expected isNodeLoading to be true")
			}
			if !model.isPodLoading {
				t.Error("expected isPodLoading to be true")
			}
			if !model.isEventLoading {
				t.Error("expected isEventLoading to be true")
			}
			if model.totalAsyncTasks != 6 {
				t.Errorf("totalAsyncTasks = %v, want 6", model.totalAsyncTasks)
			}
		})
	}
}
