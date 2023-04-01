package mockstore

import (
	"context"
	"time"
	"wd_recipes/store"

	"github.com/google/uuid"
)

var _ store.Store = (*Mockstore)(nil)

type Mockstore struct {
	GetRecipeOverride     func(ctx context.Context, id uuid.UUID) (*store.Recipe, error)
	ListRecipesOverride   func(ctx context.Context, id uuid.UUID) ([]store.Recipe, error)
	SearchRecipesOverride func(ctx context.Context, r store.SearchRecipes) ([]store.Recipe, error)
	CreateRecipeOverride  func(ctx context.Context, r store.Recipe) (*store.Recipe, error)
	UpdateRecipeOverride  func(ctx context.Context, r store.Recipe) (*store.Recipe, error)
	DeleteRecipeOverride  func(ctx context.Context, id uuid.UUID) error
}

func (m *Mockstore) GetRecipe(ctx context.Context, id uuid.UUID) (*store.Recipe, error) {
	if m.GetRecipeOverride != nil {
		return m.GetRecipeOverride(ctx, id)
	}

	return &store.Recipe{
		RecipeUUID: id,
		UserUUID:   uuid.MustParse("2c98fff4-7ccc-4536-8259-67a88380e99c"),
		RecipeName: "kimchi fried rice",
		Category:   "Korean",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

func (m *Mockstore) ListRecipes(ctx context.Context, id uuid.UUID) ([]store.Recipe, error) {
	if m.ListRecipesOverride != nil {
		return m.ListRecipesOverride(ctx, id)
	}

	return []store.Recipe{
		{
			RecipeUUID: uuid.MustParse("ffff7c73-52b0-4e3d-bf3f-0c26785ef972"),
			UserUUID:   id,
			RecipeName: "kimchi jeon",
			Category:   "Korean",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			RecipeUUID: uuid.MustParse("2c98fff4-7ccc-4536-8259-67a88380e99c"),
			UserUUID:   id,
			RecipeName: "kimchi mandu",
			Category:   "Korean",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}, nil
}

func (m *Mockstore) SearchRecipes(ctx context.Context, r store.SearchRecipes) ([]store.Recipe, error) {
	if m.SearchRecipesOverride != nil {
		return m.SearchRecipesOverride(ctx, r)
	}

	return []store.Recipe{
		{
			RecipeUUID: uuid.MustParse("ffff7c73-52b0-4e3d-bf3f-0c26785ef972"),
			UserUUID:   uuid.MustParse("080b5f09-527b-4581-bb56-19adbfe50ebf"),
			RecipeName: "salmon nigiri",
			Category:   "Japanese",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			RecipeUUID: uuid.MustParse("2c98fff4-7ccc-4536-8259-67a88380e99c"),
			UserUUID:   uuid.MustParse("ebe96725-44ef-47ee-979f-8baf823d7283"),
			RecipeName: "salmon nigiri",
			Category:   "Japanese",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}, nil
}

func (m *Mockstore) CreateRecipe(ctx context.Context, r store.Recipe) (*store.Recipe, error) {
	if m.CreateRecipeOverride != nil {
		return m.CreateRecipeOverride(ctx, r)
	}

	recipeID := uuid.New()

	return &store.Recipe{
		RecipeUUID: recipeID,
		UserUUID:   uuid.MustParse("2c98fff4-7ccc-4536-8259-67a88380e99c"),
		RecipeName: "kimchi fried rice",
		Category:   "Korean",
		Ingredients: []store.RecipeIngredient{
			{
				RecipeUUID:     recipeID,                                               //이바부야.. methods.go에서 이미 recipe.RecipeUUID로 채워놨자나..
				IngredientUUID: uuid.MustParse("2c98fff4-7ccc-4536-8259-67a88380e99b"), //니 맘대로 uuid 바꿔도 상관없어~ 숫자, 소문자, 대문자 다 상관없어, as long as they're 0~9,a~f
				Amount:         1,
				Unit:           "kg",
			},
			{
				RecipeUUID:     recipeID,
				IngredientUUID: uuid.MustParse("2c98fff4-7ccc-4536-8259-67a88380e99b"),
				Amount:         1,
				Unit:           "unit",
			},
		},
		Instructions: []store.RecipeInstruction{
			{
				RecipeUUID:  recipeID,
				StepNum:     1,
				Instruction: "Chop kimchi, onion and pork belly",
			},
			{
				RecipeUUID:  recipeID,
				StepNum:     2,
				Instruction: "grill the pan and put some oil on it",
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (m *Mockstore) UpdateRecipe(ctx context.Context, r store.Recipe) (*store.Recipe, error) {
	if m.UpdateRecipeOverride != nil {
		return m.UpdateRecipeOverride(ctx, r)
	}

	r.UpdatedAt = time.Now()

	return &r, nil
}

func (m *Mockstore) DeleteRecipe(ctx context.Context, id uuid.UUID) error {
	if m.DeleteRecipeOverride != nil {
		return m.DeleteRecipeOverride(ctx, id)
	}

	return nil
}
