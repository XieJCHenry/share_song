package user

type Status string
type SongOpr string

const (
	Unknown        Status = ""
	Online         Status = "online"
	Offline        Status = "offline"
	TimeOutOffline Status = "timeOutOffline"
)

const (
	Add    SongOpr = "add"
	Delete SongOpr = "delete"
)

type SongOperation struct {
	SongId    string
	Operation SongOpr
}
