package constant

const (
	// Switching List UI
	DefaultWidth                           int    = 60
	ListHeight                             int    = 11
	DefaultSwitchingContextHeaderMessage   string = "Kubesnap: Context Switcher"
	DefaultSwitchingNamespaceHeaderMessage string = "Kubesnap: Namespace Switcher"
	CurrentContextMarker                   string = " ✓ current-context"
	CurrentNamespaceMarker                 string = " ✓ current-namespace"
	Context                                string = "context"
	Namespace                              string = "namespace"
	DefaultNamespace                       string = "default"
	Bullet                                 string = "✓"

	// Deleting List UI
	DefaultDeletingContextHeaderMessage string = "Select contexts to delete"

	// Color Configuration
	DefaultThemeColor   string = "#4777BA"
	DefaultActiveColor  string = "#00d919"
	DefaultWarningColor string = "#f1c40f"
	DefaultUrlColor     string = "#626262"
	DefaultErrorColor   string = "#d90808"

	// File Paths
	DefaultKubeConfigLocation string = "/.kube/config"
)
