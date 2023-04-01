package service

import (
	"github.com/google/uuid"
)

func isValidSearchRecipesRequest(r SearchRecipes) bool {
	if r.UserUUID == uuid.Nil && r.RecipeName == "" && r.Category == "" {
		return false
	}

	return true
}

func isValidCreateRecipeRequest(r Recipe) bool {
	switch {
	case r.RecipeUUID != uuid.Nil:
		return false
	case r.UserUUID == uuid.Nil:
		return false
	case r.RecipeName == "":
		return false
	case r.Category == "":
		return false
	case len(r.Ingredients) == 0:
		return false
	case len(r.Instructions) == 0:
		return false
	}

	return true
}

func isValidUpdateRecipeRequest(r Recipe, uidFromPath uuid.UUID) bool {
	switch {
	case uidFromPath != r.RecipeUUID:
		return false
	case r.RecipeUUID == uuid.Nil:
		return false
	case r.UserUUID == uuid.Nil:
		return false
	case r.RecipeName == "":
		return false
	case r.Category == "":
		return false
	case len(r.Ingredients) == 0:
		return false
	case len(r.Instructions) == 0:
		return false
	}

	return true
}
