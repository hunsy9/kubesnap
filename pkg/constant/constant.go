package constant

const (
	// UI Message
	DefaultSwitchingContextHeaderMessage   string = "Kubesnap: Context Switcher"
	DefaultSwitchingNamespaceHeaderMessage string = "Kubesnap: Namespace Switcher"
	CurrentContextMarker                   string = " ✓ current-context"
	CurrentNamespaceMarker                 string = " ✓ current-namespace"
	Context                                string = "context"
	Namespace                              string = "namespace"

	// UI Configuration
	DefaultWidth                int    = 10
	ListHeight                  int    = 15
	DefaultThemeColor           string = "#8b11f5ff"
	CurrentContextMarkerColor   string = "#10bb00"
	CurrentNamespaceMarkerColor string = "#10bb00"

	// File Paths
	DefaultKubeConfigLocation string = "/.kube/config"
)
