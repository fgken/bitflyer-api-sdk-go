package bitflyerclient

func (client *Client) SendChildOrderMarket(side string, size float64) (*SendChildOrderResponse, error) {
	param := NewSendChildOrderParam()
	param.Child_order_type = MARKET
	param.Side = side
	param.Size = size
	return client.SendChildOrder(param)
}

func (client *Client) SendChildOrderLimit(side string, price, size float64) (*SendChildOrderResponse, error) {
	param := NewSendChildOrderParam()
	param.Child_order_type = LIMIT
	param.Side = side
	param.Price = price
	param.Size = size
	return client.SendChildOrder(param)
}
