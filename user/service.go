package user

import (
	"context"
)

type Service interface {
	Login(ctx context.Context, userName string, phone string) (*User, error)
	Logout(ctx context.Context, instanceId string, userName string) error
	RegisterAccount(ctx context.Context, userName string, phone string) (*User, error)
	CancelAccount(ctx context.Context, instanceId, userName, phone string) error
	SearchOnlineUsers(ctx context.Context, query map[string]interface{}, withTimeOut bool) ([]*User, error)

	// websocket
	//GetOnlineClient()
}
