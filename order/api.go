package order

import (
	"context"
	"sync"

	"github.com/bakins/kubernetes-envoy-example/api/order"
	"github.com/bakins/kubernetes-envoy-example/util"
	"github.com/satori/go.uuid"
)

// simple memory based storage
type orderStore struct {
	sync.RWMutex
	server *Server
	items  map[string]*order.Order
}

func newOrderStore(s *Server) *orderStore {
	return &orderStore{
		server: s,
		items:  make(map[string]*order.Order),
	}
}

// TODO: populate stub data

func (u *orderStore) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.Order, error) {
	u.Lock()
	defer u.Unlock()

	item := &order.Order{
		Id:   uuid.NewV4().String(),
		User: req.User,
		// we'd really make sure these are valid, I guess.
		Items: req.Items,
	}

	u.items[item.Id] = item

	return item, nil
}

func (u *orderStore) GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.Order, error) {
	u.RLock()
	defer u.RUnlock()

	item, ok := u.items[req.Id]
	if !ok {
		return nil, util.NewNotFoundError("order", req.Id)
	}

	return item, nil
}

func (u *orderStore) ListOrders(ctx context.Context, req *order.ListOrdersRequest) (*order.ListOrdersResponse, error) {
	u.RLock()
	defer u.RUnlock()

	res := &order.ListOrdersResponse{
		Orders: make([]*order.Order, 0),
	}

	for _, v := range u.items {
		if req.User != "" && req.User != v.User {
			continue
		}
		res.Orders = append(res.Orders, v)
	}

	return res, nil
}

func (u *orderStore) DeleteOrder(ctx context.Context, req *order.DeleteOrderRequest) (*order.Order, error) {
	u.Lock()
	defer u.Unlock()

	item, ok := u.items[req.Id]
	if !ok {
		return nil, util.NewNotFoundError("order", req.Id)
	}

	delete(u.items, req.Id)

	return item, nil
}

func (u *orderStore) UpdateOrder(ctx context.Context, req *order.Order) (*order.Order, error) {
	u.Lock()
	defer u.Unlock()

	_, ok := u.items[req.Id]
	if !ok {
		return nil, util.NewNotFoundError("order", req.Id)
	}

	// in a real system we may merge/patch
	u.items[req.Id] = req

	return req, nil
}
