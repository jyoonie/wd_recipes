package service

import (
	"wd_recipes/store"

	"github.com/google/uuid"
)

func apiRIngr2DBRIngr(r []RecipeIngredient) []store.RecipeIngredient {
	var ingredients []store.RecipeIngredient

	for _, ingr := range r {
		var s store.RecipeIngredient

		s.IngredientUUID = ingr.IngredientUUID
		s.Amount = ingr.Amount
		s.Unit = ingr.Unit
		ingredients = append(ingredients, s)
	}

	return ingredients
}

func apiRInst2DBRInst(r []RecipeInstruction) []store.RecipeInstruction {
	var instructions []store.RecipeInstruction

	for _, inst := range r {
		var s store.RecipeInstruction

		s.StepNum = inst.StepNum
		s.Instruction = inst.Instruction
		instructions = append(instructions, s)
	}

	return instructions
}

func apiRecipe2DBRecipe(r Recipe) store.Recipe {
	return store.Recipe{
		RecipeUUID:   r.RecipeUUID,
		UserUUID:     r.UserUUID,
		RecipeName:   r.RecipeName,
		Category:     r.Category,
		Ingredients:  apiRIngr2DBRIngr(r.Ingredients),
		Instructions: apiRInst2DBRInst(r.Instructions),
	}
}

func dbRIngr2apiRIngr(r []store.RecipeIngredient) []RecipeIngredient {
	var ingredients []RecipeIngredient

	for _, ingr := range r {
		var s RecipeIngredient

		s.IngredientUUID = ingr.IngredientUUID
		s.Amount = ingr.Amount
		s.Unit = ingr.Unit
		ingredients = append(ingredients, s)
	}

	return ingredients
}

func dbRInst2apiRInst(r []store.RecipeInstruction) []RecipeInstruction {
	var instructions []RecipeInstruction

	for _, inst := range r {
		var s RecipeInstruction

		s.StepNum = inst.StepNum
		s.Instruction = inst.Instruction
		instructions = append(instructions, s)
	}

	return instructions
}

func dbRecipe2ApiRecipe(r *store.Recipe) Recipe {
	return Recipe{
		RecipeUUID:   r.RecipeUUID,
		UserUUID:     r.UserUUID,
		RecipeName:   r.RecipeName,
		Category:     r.Category,
		Ingredients:  dbRIngr2apiRIngr(r.Ingredients),
		Instructions: dbRInst2apiRInst(r.Instructions),
	}
}

func apiSearchR2DBSearchR(r SearchRecipes) store.SearchRecipes {
	var out store.SearchRecipes

	if r.UserUUID != uuid.Nil {
		out.UserUUID = &r.UserUUID
	}
	if r.RecipeName != "" {
		out.RecipeName = &r.RecipeName
	}
	if r.Category != "" {
		out.Category = &r.Category
	}

	return out
}
