package protocol

type Protocol struct {
	Path string                 `json:"path,omitempty"`
	Body map[string]interface{} `json:"body,omitempty"`
}
