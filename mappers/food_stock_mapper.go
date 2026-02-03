package mappers

import (
	"backend/entities"
	"backend/models"
	"backend/utils"

	"gorm.io/gorm"
)

func FoodStockModelToEntity(model models.FoodStock) entities.FoodStock {
	return entities.FoodStock{
		ID:           utils.UintToString(model.ID),
		FoodTypeID:   utils.UintToString(model.FoodTypeID),
		Amount:       model.Amount,
		DateTime:     model.DateTime,
		Note:         model.Note,
		FoodTypeName: model.FoodType.Name,
		CreatedBy:    model.Creator.FullName,
		CreatedRole:  model.Creator.Role,
		UpdatedBy:    model.Updater.FullName,
		UpdatedRole:  model.Updater.Role,

		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func FoodStockEntityToModel(entity entities.FoodStock) models.FoodStock {
	return models.FoodStock{
		Model: gorm.Model{
			ID:        utils.StringToUint(entity.ID),
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		FoodTypeID: utils.StringToUint(entity.FoodTypeID),
		Amount:     entity.Amount,
		DateTime:   entity.DateTime,
		Note:       entity.Note,
		CreatedBy:  utils.StringToUint(entity.CreatedBy),
		UpdatedBy:  utils.StringToUint(entity.UpdatedBy),
	}
}

func FoodStockToEntities(model []models.FoodStock) []entities.FoodStock {
	entities := make([]entities.FoodStock, len(model))
	for i, m := range model {
		entities[i] = FoodStockModelToEntity(m)
	}
	return entities
}
