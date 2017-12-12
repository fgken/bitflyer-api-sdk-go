package bitflyerclient

import (
	"fmt"
)

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

func (client *Client) GetParentOrdersByState(state string) ([]GetParentOrdersResponse, error) {
	param := NewGetParentOrdersParam()
	param.Parent_order_state = state
	return client.GetParentOrders(param)
}

func (client *Client) GetParentOrderState(id string) (string, error) {
	param := NewGetParentOrdersParam()
	for _, count := range []int64{100, 400, 1000} {
		param.Page.Count = count
		orders, err := client.GetParentOrders(param)
		if err != nil {
			return "", err
		}
		for _, order := range orders {
			if order.Parent_order_acceptance_id == id {
				return order.Parent_order_state, nil
			}
		}
		param.Page.Before = orders[len(orders)-1].Id
	}

	return "", fmt.Errorf("not found parent id: %v", id)
}

func (client *Client) GetChildOrdersByChildOrderId(id string) ([]GetChildOrdersResponse, error) {
	param := NewGetChildOrdersParam()
    param.Child_order_id = id
	return client.GetChildOrders(param)
}

