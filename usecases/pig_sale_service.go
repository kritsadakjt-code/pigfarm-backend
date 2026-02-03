package usecases

type PigSaleService struct {
	saleRepo PigSaleRepository
}

func NewPigSaleService(saleRepo PigSaleRepository) *PigSaleService {
	return &PigSaleService{saleRepo: saleRepo}
}
