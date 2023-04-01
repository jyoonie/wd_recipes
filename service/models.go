package service

import (
	"github.com/google/uuid"
)

type Recipe struct {
	RecipeUUID   uuid.UUID           `json:"recipe_uuid,omitempty"`
	UserUUID     uuid.UUID           `json:"user_uuid,omitempty"`
	RecipeName   string              `json:"recipe_name,omitempty"`
	Category     string              `json:"category,omitempty"`
	Ingredients  []RecipeIngredient  `json:"ingredients,omitempty"`  //여기에서 받은건 recipe ingredients 테이블에 저장된당
	Instructions []RecipeInstruction `json:"instructions,omitempty"` //여기에서 받은건 recipe instructions 테이블에 저장된당
	// CreatedAt  time.Time `json:"created_at,omitempty"`
	// UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type SearchRecipes struct {
	UserUUID   uuid.UUID `json:"user_uuid,omitempty"`
	RecipeName string    `json:"recipe_name,omitempty"`
	Category   string    `json:"category,omitempty"`
}

type RecipeIngredient struct {
	IngredientUUID uuid.UUID `json:"ingredient_uuid,omitempty"`
	Amount         int       `json:"amount,omitempty"`
	Unit           string    `json:"unit,omitempty"`
}

type RecipeInstruction struct {
	StepNum     int    `json:"step_num,omitempty"`
	Instruction string `json:"instruction,omitempty"`
}
