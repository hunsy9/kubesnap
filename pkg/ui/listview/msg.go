package listview

type OperationResultMsg struct {
	Err error
}

type RenameResultMsg struct {
	OriginalName string
	NewName      string
	Err          error
}

type DeletionResultMsg struct {
	Err     error
	Targets []string
}
