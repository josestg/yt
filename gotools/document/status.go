package document

//go:generate stringer -type=Status

type Status int

const (
	StatusSignRequested Status = iota
	StatusSigned
	StatusRejected
	StatusVerified
)
