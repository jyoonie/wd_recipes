package store

import (
	"time"

	"github.com/google/uuid"
)

type Recipe struct {
	RecipeUUID   uuid.UUID
	UserUUID     uuid.UUID
	RecipeName   string
	Category     string
	Ingredients  []RecipeIngredient
	Instructions []RecipeInstruction
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type SearchRecipes struct {
	UserUUID   *uuid.UUID
	RecipeName *string
	Category   *string
}

type RecipeIngredient struct { //your db model always matches your table. 그래서 여기에서는 init migrate up에 있는 모든 필드 다 있음.
	RecipeUUID     uuid.UUID
	IngredientUUID uuid.UUID
	Amount         int
	Unit           string
}

type RecipeInstruction struct {
	RecipeUUID  uuid.UUID
	StepNum     int
	Instruction string
}
