package user

import (
	"context"
	"sync"

	"github.com/bakins/kubernetes-envoy-example/api/user"
	"github.com/bakins/kubernetes-envoy-example/util"
	"github.com/satori/go.uuid"
)

// simple memory based storage
type userStore struct {
	sync.RWMutex
	server *Server
	items  map[string]*user.User
}

func newUserStore(s *Server) *userStore {
	return &userStore{
		server: s,
		items:  make(map[string]*user.User),
	}
}

// TODO: populate stub data

func (u *userStore) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.User, error) {
	u.Lock()
	defer u.Unlock()

	item := &user.User{
		Id:      uuid.NewV4().String(),
		Name:    req.Name,
		Address: req.Address,
		Email:   req.Address,
	}

	u.items[item.Id] = item

	return item, nil
}

func (u *userStore) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.User, error) {
	u.RLock()
	defer u.RUnlock()

	item, ok := u.items[req.Id]
	if !ok {
		return nil, util.NewNotFoundError("user", req.Id)
	}

	return item, nil
}

func (u *userStore) ListUsers(ctx context.Context, req *user.ListUsersRequest) (*user.ListUsersResponse, error) {
	u.RLock()
	defer u.RUnlock()

	res := &user.ListUsersResponse{
		Users: make([]*user.User, 0, len(u.items)),
	}
	for _, v := range u.items {
		res.Users = append(res.Users, v)
	}

	return res, nil
}

func (u *userStore) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.User, error) {
	u.Lock()
	defer u.Unlock()

	item, ok := u.items[req.Id]
	if !ok {
		return nil, util.NewNotFoundError("user", req.Id)
	}

	delete(u.items, req.Id)

	return item, nil
}

func (u *userStore) UpdateUser(ctx context.Context, req *user.User) (*user.User, error) {
	u.Lock()
	defer u.Unlock()

	_, ok := u.items[req.Id]
	if !ok {
		return nil, util.NewNotFoundError("user", req.Id)
	}

	// in a real system we may merge/patch
	u.items[req.Id] = req

	return req, nil
}

func (u *userStore) LoadSampleData() {
	u.Lock()
	defer u.Unlock()

	for _, o := range util.SampleUsers {
		o := o
		u.items[o.Id] = &o
	}
}
