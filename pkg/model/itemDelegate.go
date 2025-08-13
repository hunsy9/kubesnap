package model

type Item string

func (i Item) FilterValue() string { return string(i) }

type ItemDelegate struct{}
