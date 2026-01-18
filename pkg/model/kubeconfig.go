package model

type KubeConfig struct {
	APIVersion     string    `yaml:"apiVersion"`
	Kind           string    `yaml:"kind"`
	CurrentContext string    `yaml:"current-context"`
	Contexts       []Context `yaml:"contexts"`
	Clusters       []Cluster `yaml:"clusters"`
	Users          []User    `yaml:"users"`
}

type Context struct {
	Name    string        `yaml:"name"`
	Context ContextDetail `yaml:"context"`
}

type ContextDetail struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

type Cluster struct {
	Name    string        `yaml:"name"`
	Cluster ClusterDetail `yaml:"cluster"`
}

type ClusterDetail struct {
	Server string `yaml:"server"`
}

type User struct {
	Name string     `yaml:"name"`
	User UserDetail `yaml:"user"`
}

type UserDetail struct {
	Token string `yaml:"token,omitempty"`
}
