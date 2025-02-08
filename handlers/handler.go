package handlers

type Handler struct {
	payment *Payment
}

func NewHandle(payment *Payment) *Handler {
	return &Handler{payment: payment}
}
