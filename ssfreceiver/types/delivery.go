package types

type DeliveryMethod string

const (
	DeliveryMethodPush DeliveryMethod = "urn:ietf:rfc:8935"
	DeliveryMethodPoll DeliveryMethod = "urn:ietf:rfc:8936"
)

func IsValidDeliveryMethod(method DeliveryMethod) bool {
	switch method {
	case DeliveryMethodPush, DeliveryMethodPoll:
		return true
	default:
		return false
	}
}
