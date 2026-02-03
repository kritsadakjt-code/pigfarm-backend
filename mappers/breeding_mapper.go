package mappers

import (
	"backend/entities"
	"backend/models"
	"backend/utils"

	"gorm.io/gorm"
)

func BreedingModelToEntity(model models.Breeding) entities.Breeding {
	// ถ้าไม่มีค่าเเล้วอยากให้เป็น - ซึ่งไม่มีก็ได้
	fatherCodeName := "-"
	if model.Father.ID != 0 {
		fatherCodeName = model.Father.CodeName
	}
	motherCodeName := "-"
	if model.Mother.ID != 0 {
		motherCodeName = model.Mother.CodeName
	}

	createdName := "-"
	createdRole := "-"
	if model.Creator.ID != 0 {
		createdName = model.Creator.FullName
		createdRole = model.Creator.Role
	}

	updatedName := "-"
	updatedRole := "-"
	if model.Updater.ID != 0 {
		updatedName = model.Updater.FullName
		updatedRole = model.Updater.Role
	}

	return entities.Breeding{
		ID:        utils.UintToString(model.ID),
		FatherID:  utils.UintToString(model.FatherID),
		MotherID:  utils.UintToString(model.MotherID),
		CreatedBy: utils.UintToString(model.CreatedBy),
		UpdatedBy: utils.UintToString(model.UpdatedBy),

		FatherCodename: fatherCodeName,
		MotherCodename: motherCodeName,
		CreatedName:    createdName,
		CreatedRole:    createdRole,
		UpdatedName:    updatedName,
		UpdatedRole:    updatedRole,

		BreedingDate:  model.BreedingDate,
		ExpectedBirth: model.ExpectedBirth,
		Status:        model.Status,
		Result:        model.Result,
		Note:          model.Note,
		//time
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}

}

func BreedingEntityToModel(e entities.Breeding) models.Breeding {
	return models.Breeding{
		Model: gorm.Model{
			ID:        utils.StringToUint(e.ID),
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		},
		FatherID:      utils.StringToUint(e.FatherID),
		MotherID:      utils.StringToUint(e.MotherID),
		CreatedBy:     utils.StringToUint(e.CreatedBy),
		UpdatedBy:     utils.StringToUint(e.UpdatedBy),
		BreedingDate:  e.BreedingDate,
		ExpectedBirth: e.ExpectedBirth,
		Status:        e.Status,
		Result:        e.Result,
		Note:          e.Note,
	}
}

// กรณี return หลายตัว
func BreedingToEntities(models []models.Breeding) []entities.Breeding {
	entities := make([]entities.Breeding, len(models))
	for i, m := range models {
		entities[i] = BreedingModelToEntity(m)
	}
	return entities
}
