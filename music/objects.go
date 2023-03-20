package music

type FromType string

const (
	FromUser   FromType = "user"
	FromSystem FromType = "system"
)

type Song struct {
	InstanceId string   `json:"instanceId,omitempty"`
	Name       string   `json:"name,omitempty"`
	Artists    []string `json:"artists,omitempty"`
	Length     uint32   `json:"length,omitempty"`
	Lyric      string   `json:"lyric,omitempty"`
	Source     string   `json:"source,omitempty"`
	Cover      string   `json:"cover,omitempty"`
	DeleteFlag bool     `json:"deleteFlag,omitempty"`
	From       FromType `json:"from,omitempty"`
	BinarySize uint64   `json:"binarySize,omitempty"`
}
