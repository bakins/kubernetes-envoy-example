package util

import (
	"github.com/bakins/kubernetes-envoy-example/api/order"
	"github.com/bakins/kubernetes-envoy-example/api/user"
)

// SampleUsers is sample user data for testing.
var SampleUsers = []user.User{
	user.User{
		Id:      "0a30cbc7-91f1-4637-832f-9ddb9978e1c5",
		Name:    "User One",
		Address: "1234 5th Street",
	},
	user.User{
		Id:      "d6dd61c4-75a6-4424-91b2-ab740501c94f",
		Name:    "User Two",
		Address: "6th Ave",
	},
	user.User{
		Id:      "32cc174b-3586-4263-9c67-fce988d246ca",
		Name:    "User Three",
		Address: "Down by the river",
	},
}

// SampleOrders is sample order data for testing.
var SampleOrders = []order.Order{
	order.Order{
		Id:   "22f4ee36-2c3e-450f-bb65-36e52656b183",
		User: "0a30cbc7-91f1-4637-832f-9ddb9978e1c5",
		Items: []string{
			"6ab9e0c2-e7be-4120-a3e9-62c39b7dbfd7",
			"4415fede-7462-4f12-b87f-ede596ec6ee2",
		},
	},
	order.Order{
		Id:   "034735ae-b5c8-47dd-9b59-30eacc62473a",
		User: "0a30cbc7-91f1-4637-832f-9ddb9978e1c5",
		Items: []string{
			"6962f4ff-b752-4103-b90c-1f9bcec30913",
			"dff79aa1-6b13-4aeb-8dca-22a45322a293",
		},
	},
	order.Order{
		Id:   "4b2f25c3-f0fb-4640-9fc8-a3066a8022f6",
		User: "d6dd61c4-75a6-4424-91b2-ab740501c94f",
		Items: []string{
			"5d210689-cca7-4e81-8437-05b20f658ad0",
		},
	},
}

/*
items from item service
   "6ab9e0c2-e7be-4120-a3e9-62c39b7dbfd7"
    "4415fede-7462-4f12-b87f-ede596ec6ee2"
    "5d210689-cca7-4e81-8437-05b20f658ad0"
    "dff79aa1-6b13-4aeb-8dca-22a45322a293"
	"6962f4ff-b752-4103-b90c-1f9bcec30913"

*/
