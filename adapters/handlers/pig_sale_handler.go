package handlers

import "backend/usecases"

type HttpPigSale struct {
	saleService *usecases.PigSaleService
}

func NewHttpPigSale(saleService *usecases.PigSaleService) *HttpPigSale {
	return &HttpPigSale{saleService: saleService}
}
