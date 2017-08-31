package order

import (
	"context"
	"fmt"
	"sync"

	"github.com/bakins/kubernetes-envoy-example/api/item"
	"github.com/bakins/kubernetes-envoy-example/api/order"
	"github.com/bakins/kubernetes-envoy-example/util"
	"github.com/satori/go.uuid"
)

// simple memory based storage
type orderStore struct {
	sync.RWMutex
	server     *Server
	items      map[string]*order.Order
	itemClient item.ItemServiceClient
}

func newOrderStore(s *Server, client item.ItemServiceClient) *orderStore {
	return &orderStore{
		server:     s,
		items:      make(map[string]*order.Order),
		itemClient: client,
	}
}

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
	return u.getOrder(req.Id)
}

func (u *orderStore) getOrder(id string) (*order.Order, error) {
	u.RLock()
	defer u.RUnlock()
	item, ok := u.items[id]
	if !ok {
		return nil, util.NewNotFoundError("order", id)
	}
	return item, nil
}

func (u *orderStore) GetOrderDetail(ctx context.Context, req *order.GetOrderDetailRequest) (*order.GetOrderDetailResponse, error) {
	fmt.Println("GetOrderDetail", req.Id)
	o, err := u.getOrder(req.Id)
	if err != nil {
		return nil, err
	}

	res := &order.GetOrderDetailResponse{
		Id:   o.Id,
		User: o.User,
	}

	for _, id := range o.Items {
		fmt.Println("GetOrderDetail item", id)
		i, err := u.itemClient.GetItem(ctx, &item.GetItemRequest{Id: id})
		if err != nil {
			fmt.Println("error getting item", id, err)
			continue
		}
		fmt.Println("GetOrderDetail print item", i)
		res.Items = append(res.Items, i)
	}
	return res, nil
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

func (u *orderStore) LoadSampleData() {
	u.Lock()
	defer u.Unlock()

	for _, o := range util.SampleOrders {
		o := o
		u.items[o.Id] = &o
	}
}
