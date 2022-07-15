package dto

type (
	Video struct {
		Caption string `bson:"caption" json:"caption" validate:"required"`
	}
)
