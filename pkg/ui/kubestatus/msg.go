package kubestatus

import authenticationv1 "k8s.io/api/authentication/v1"

type ApiHealthMsg struct {
	Msg string
}

type AuthInfoMsg struct {
	Msg *authenticationv1.SelfSubjectReview
	Err error
}

type NodeInfoMsg struct {
	TotalNodes       int
	ReadyNodes       int
	NodeAPIStatusErr string
}

type PodInfoMsg struct {
	TotalPods       int
	RunningPods     int
	PodAPIStatusErr string
}

type EventInfoMsg struct {
	WarningEvents     int
	EventAPIStatusErr string
}

type UpdateAvailableMsg string

type AsyncTaskAllDoneMsg struct{}
