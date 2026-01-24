package constant

const (
	// Switching List UI
	DefaultWidth                           int    = 10
	ListHeight                             int    = 15
	DefaultSwitchingContextHeaderMessage   string = "Kubesnap: Context Switcher"
	DefaultSwitchingNamespaceHeaderMessage string = "Kubesnap: Namespace Switcher"
	CurrentContextMarker                   string = " ✓ current-context"
	CurrentNamespaceMarker                 string = " ✓ current-namespace"
	Context                                string = "context"
	Namespace                              string = "namespace"
	DefaultNamespace                       string = "default"
	Bullet                                 string = "✓"

	// Color Configuration
	DefaultThemeColor  string = "#4777BA"
	DefaultActiveColor string = "#00d919"
	DefaultUrlColor    string = "#626262"
	DefaultErrorColor  string = "#d90808"

	// File Paths
	DefaultKubeConfigLocation string = "/.kube/config"
)
