package user

type LoginRequest struct {
	UserName string `json:"userName,omitempty"`
	Phone    string `json:"phone,omitempty"`
}
type LoginResponse struct {
	InstanceId string           `json:"instanceId,omitempty"`
	UserName   string           `json:"userName,omitempty"`
	SongList   []*SongOperation `json:"songList,omitempty"`
}

type LogoutRequest struct {
	InstanceId string `json:"instanceId,omitempty"`
	UserName   string `json:"userName,omitempty"`
}
type LogoutResponse struct {
	InstanceId string `json:"instanceId,omitempty"`
}

type RegisterAccountRequest struct {
	UserName string `json:"userName,omitempty"`
	Phone    string `json:"phone,omitempty"`
}
type RegisterAccountResponse struct {
	InstanceId string `json:"instanceId,omitempty"`
	UserName   string `json:"userName,omitempty"`
}

type CancelAccountRequest struct {
	InstanceId string `json:"instanceId,omitempty"`
	UserName   string `json:"userName,omitempty"`
	Phone      string `json:"phone,omitempty"`
}
type CancelAccountResponse struct {
	InstanceId string `json:"instanceId,omitempty"`
}

type SearchOnlineUsersRequest struct {
	Query       map[string]interface{} `json:"query,omitempty"`
	WithTimeOut bool                   `json:"withTimeOut,omitempty"`
}
type SearchOnlineUsersResponse struct {
	WithTimeOut bool `json:"withTimeOut,omitempty"`
	Users       []*User
}
