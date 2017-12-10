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

func (client *Client) SendParentOrderStop(side string, price, size float64) (*SendParentOrderResponse, error) {
	param := NewSendParentOrderParam()
	param.Order_method = SIMPLE
	parentOrder := ParentOrder{
		Condition_type: STOP,
		Side:           side,
		Size:           size,
		Trigger_price:  price,
	}
	param.Parameters = append(param.Parameters, parentOrder)
	return client.SendParentOrder(param)
}

func (client *Client) SendParentOrderIFDOCO(conditionType, side string, entry, limit, stop, size float64) (*SendParentOrderResponse, error) {
	param := NewSendParentOrderParam()
	param.Order_method = IFDOCO

	var parentOrder ParentOrder

	/* entry */
	parentOrder = ParentOrder{
		Condition_type: conditionType,
		Side:           side,
		Size:           size,
	}
	switch conditionType {
	case LIMIT:
		parentOrder.Price = entry
	case STOP:
		parentOrder.Trigger_price = entry
	}

	param.Parameters = append(param.Parameters, parentOrder)

	/* limit */
	var oppositeSide string
	switch side {
	case BUY:
		oppositeSide = SELL
	case SELL:
		oppositeSide = BUY
	}

	parentOrder = ParentOrder{
		Condition_type: LIMIT,
		Side:           oppositeSide,
		Size:           size,
		Price:          limit,
	}
	param.Parameters = append(param.Parameters, parentOrder)

	/* stop */
	parentOrder = ParentOrder{
		Condition_type: STOP,
		Side:           oppositeSide,
		Size:           size,
		Trigger_price:  stop,
	}
	param.Parameters = append(param.Parameters, parentOrder)

	return client.SendParentOrder(param)
}
