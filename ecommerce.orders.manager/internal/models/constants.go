package models

type OrderStatus struct {
	Pending   string
	Confirmed string
	Shipped   string
	Delivered string
	Cancelled string
}

var OrderStatuses = OrderStatus{
	Pending:   "pending",
	Confirmed: "confirmed",
	Shipped:   "shipped",
	Delivered: "delivered",
	Cancelled: "cancelled",
}