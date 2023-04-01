package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"wd_recipes/store"

	"github.com/google/uuid"
)

const defaultTimeout = 5 * time.Second

func (pg *PG) GetRecipe(ctx context.Context, id uuid.UUID) (*store.Recipe, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var recipe store.Recipe

	row := pg.db.QueryRowContext(ctx, sqlGetRecipe, id)
	if err := row.Scan(
		&recipe.RecipeUUID,
		&recipe.UserUUID,
		&recipe.RecipeName,
		&recipe.Category,
		&recipe.CreatedAt,
		&recipe.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, fmt.Errorf("error getting recipe: %w", err)
	}

	ingredients, err := pg.db.QueryContext(ctx, sqlListRecipeIngr, recipe.RecipeUUID)
	if err != nil {
		return nil, fmt.Errorf("error listing recipe ingredients: %w", err)
	}
	defer ingredients.Close()

	for ingredients.Next() {
		var recipeIngr store.RecipeIngredient
		if err := ingredients.Scan(
			&recipeIngr.RecipeUUID,
			&recipeIngr.IngredientUUID,
			&recipeIngr.Amount,
			&recipeIngr.Unit,
		); err != nil {
			return nil, fmt.Errorf("error listing recipe ingredients: %w", err)
		}
		recipe.Ingredients = append(recipe.Ingredients, recipeIngr)
	}

	instructions, err := pg.db.QueryContext(ctx, sqlListRecipeInst, recipe.RecipeUUID)
	if err != nil {
		return nil, fmt.Errorf("error listing recipe instructions: %w", err)
	}
	defer instructions.Close()

	for instructions.Next() {
		var recipeInst store.RecipeInstruction
		if err := instructions.Scan(
			&recipeInst.RecipeUUID,
			&recipeInst.StepNum,
			&recipeInst.Instruction,
		); err != nil {
			return nil, fmt.Errorf("error listing recipe instructions: %w", err)
		}
		recipe.Instructions = append(recipe.Instructions, recipeInst)
	}

	return &recipe, nil
}

func (pg *PG) ListRecipes(ctx context.Context, id uuid.UUID) ([]store.Recipe, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var recipes []store.Recipe

	rows, err := pg.db.QueryContext(ctx, sqlListRecipes, id)
	if err != nil {
		return nil, fmt.Errorf("error listing recipes: %w", err)
	}
	defer rows.Close()

	for rows.Next() { //ListRecipes는 걍 GetRecipe 코드들을 for rows.Next()안에 넣은 형태라 할 수 있음. 레시피가 하나가 아니라 여러개니까 for문 해줘야징
		var recipe store.Recipe
		if err := rows.Scan(
			&recipe.RecipeUUID, //요 부분은 Recipes 테이블에서 불러올 수 있는 정보들이고
			&recipe.UserUUID,
			&recipe.RecipeName,
			&recipe.Category,
			&recipe.CreatedAt,
			&recipe.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error listing recipes: %w", err)
		}

		ingredients, err := pg.db.QueryContext(ctx, sqlListRecipeIngr, recipe.RecipeUUID) //ingredients는 요 안에서 다시 Recipe Ingredients 테이블서 불러줘야지..
		if err != nil {
			return nil, fmt.Errorf("error listing recipe ingredients: %w", err)
		}
		defer ingredients.Close()

		for ingredients.Next() {
			var recipeIngr store.RecipeIngredient
			if err := ingredients.Scan(
				&recipeIngr.RecipeUUID,
				&recipeIngr.IngredientUUID,
				&recipeIngr.Amount,
				&recipeIngr.Unit,
			); err != nil {
				return nil, fmt.Errorf("error listing recipe ingredients: %w", err)
			}
			recipe.Ingredients = append(recipe.Ingredients, recipeIngr)
		}

		instructions, err := pg.db.QueryContext(ctx, sqlListRecipeInst, recipe.RecipeUUID)
		if err != nil {
			return nil, fmt.Errorf("error listing recipe instructions: %w", err)
		}
		defer instructions.Close()

		for instructions.Next() {
			var recipeInst store.RecipeInstruction
			if err := instructions.Scan(
				&recipeInst.RecipeUUID,
				&recipeInst.StepNum,
				&recipeInst.Instruction,
			); err != nil {
				return nil, fmt.Errorf("error listing recipe instructions: %w", err)
			}
			recipe.Instructions = append(recipe.Instructions, recipeInst)
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func (pg *PG) SearchRecipes(ctx context.Context, r store.SearchRecipes) ([]store.Recipe, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	var wheres []string
	var vars []interface{}
	var count int

	if r.UserUUID != nil {
		count++
		wheres = append(wheres, fmt.Sprintf(" user_uuid = $%d", count))
		vars = append(vars, r.UserUUID)
	}

	if r.RecipeName != nil {
		count++
		wheres = append(wheres, fmt.Sprintf(" recipe_name = $%d", count))
		vars = append(vars, r.RecipeName)
	}

	if r.Category != nil {
		count++
		wheres = append(wheres, fmt.Sprintf(" category = $%d", count))
		vars = append(vars, r.Category)
	}

	whereClause := strings.Join(wheres, " AND ")

	var recipes []store.Recipe

	fmt.Println("sql statement is")
	rows, err := pg.db.QueryContext(ctx, sqlsearchRecipes+" WHERE "+whereClause, vars...)
	if err != nil {
		return nil, fmt.Errorf("error searching recipes: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var recipe store.Recipe
		if err = rows.Scan( //여기서 하는 짓은, recipe table에서 찾은 row를 각각 recipe struct의 필드에 맞게 할당해주는 거임.
			&recipe.RecipeUUID,
			&recipe.UserUUID,
			&recipe.RecipeName,
			&recipe.Category,
			&recipe.CreatedAt,
			&recipe.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("error searching recipes: %w", err)
		}

		ingredients, err := pg.db.QueryContext(ctx, sqlListRecipeIngr, recipe.RecipeUUID)
		if err != nil {
			return nil, fmt.Errorf("error listing recipe ingredients: %w", err)
		}
		defer ingredients.Close()

		for ingredients.Next() {
			var recipeIngr store.RecipeIngredient
			if err := ingredients.Scan( //여기서 하는 짓은, recipe_ingredients table에서 찾은 row를 각각 recipe ingredient struct의 필드에 맞게 할당해주는 거임. & 써야하나..(응, 포인터 위치에 값을 복사하는 거니까..)
				&recipeIngr.RecipeUUID,
				&recipeIngr.IngredientUUID,
				&recipeIngr.Amount,
				&recipeIngr.Unit,
			); err != nil {
				return nil, fmt.Errorf("error listing recipe ingredients: %w", err)
			}
			recipe.Ingredients = append(recipe.Ingredients, recipeIngr)
		}

		instructions, err := pg.db.QueryContext(ctx, sqlListRecipeInst, recipe.RecipeUUID)
		if err != nil {
			return nil, fmt.Errorf("error listing recipe instructions: %w", err)
		}
		defer instructions.Close()

		for instructions.Next() {
			var recipeInst store.RecipeInstruction
			if err := instructions.Scan(
				&recipeInst.RecipeUUID,
				&recipeInst.StepNum,
				&recipeInst.Instruction,
			); err != nil {
				return nil, fmt.Errorf("error listing recipe instructions: %w", err)
			}
			recipe.Instructions = append(recipe.Instructions, recipeInst)
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func (pg *PG) CreateRecipe(ctx context.Context, r store.Recipe) (*store.Recipe, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating recipe: %w", err)
	}

	var recipe store.Recipe

	row := tx.QueryRowContext(ctx, sqlCreateRecipe,
		&r.UserUUID,
		&r.RecipeName,
		&r.Category,
	)

	if err = row.Scan(
		&recipe.RecipeUUID,
		&recipe.UserUUID,
		&recipe.RecipeName,
		&recipe.Category,
		&recipe.CreatedAt,
		&recipe.UpdatedAt,
	); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating recipe: %w", err)
	}

	for _, recipeingr := range r.Ingredients { //위의 listRecipes나 searchRecipes에서는 db에서 읽어오는 것이므로 for ingredients.Next()임
		var recipeIngr store.RecipeIngredient

		row := tx.QueryRowContext(ctx, sqlCreateRecipeIngr,
			&recipe.RecipeUUID, //여기 조오심..!! RecipeIngredient db model에 recipe uuid 필드가 있긴 하지만, api model에서 빈채로 요청이 오기 때문에(그리고 db에서 recipe uuid가 만들어지기 때문에) 실제로는 recipe의 recipe_uuid를 recipe_ingredients 테이블에 집어넣어야함.
			&recipeingr.IngredientUUID,
			&recipeingr.Amount,
			&recipeingr.Unit,
		)

		if err = row.Scan(
			&recipeIngr.RecipeUUID,
			&recipeIngr.IngredientUUID,
			&recipeIngr.Amount,
			&recipeIngr.Unit,
		); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error creating recipe ingredients: %w", err)
		}
		recipe.Ingredients = append(recipe.Ingredients, recipeIngr)
	}

	for _, recipeinst := range r.Instructions {
		var recipeInst store.RecipeInstruction

		row := tx.QueryRowContext(ctx, sqlCreateRecipeInst,
			&recipe.RecipeUUID,
			&recipeinst.StepNum,
			&recipeinst.Instruction,
		)

		if err = row.Scan(
			&recipeInst.RecipeUUID,
			&recipeInst.StepNum,
			&recipeInst.Instruction,
		); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error creating recipe instructions: %w", err)
		}
		recipe.Instructions = append(recipe.Instructions, recipeInst)
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating recipe: %w", err)
	}

	return &recipe, nil
}

func (pg *PG) UpdateRecipe(ctx context.Context, r store.Recipe) (*store.Recipe, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error updating recipe: %w", err)
	}

	var recipe store.Recipe

	row := tx.QueryRowContext(ctx, sqlUpdateRecipe,
		r.RecipeName,
		r.Category,
		r.RecipeUUID,
	)

	if err = row.Scan(
		&recipe.RecipeUUID,
		&recipe.UserUUID,
		&recipe.RecipeName,
		&recipe.Category,
		&recipe.CreatedAt,
		&recipe.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tx.Rollback()
			return nil, store.ErrNotFound
		}
		tx.Rollback()
		return nil, fmt.Errorf("error updating recipe: %w", err)
	}

	res, err := tx.ExecContext(ctx, sqlDeleteRecipeIngr,
		r.RecipeUUID,
	)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error deleting recipe ingredients: %w", err) //exec context는 return 값이 없을때 사용하는 거라던데, 이 update recipe는 *store.Recipe를 돌려줘야해서 nil 썼는데, 그럴거면 queryContext를 써야하나..
	}

	if _, affected := res.RowsAffected(); affected != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error deleting recipe ingredients, rows not affected(error: %v)", affected)
	}

	for _, recipeingr := range r.Ingredients { //위에서 삭제하고 다시 create 해주는거징..
		var recipeIngr store.RecipeIngredient

		row := tx.QueryRowContext(ctx, sqlCreateRecipeIngr,
			&recipe.RecipeUUID,
			&recipeingr.IngredientUUID,
			&recipeingr.Amount,
			&recipeingr.Unit,
		)

		if err = row.Scan(
			&recipeIngr.RecipeUUID,
			&recipeIngr.IngredientUUID,
			&recipeIngr.Amount,
			&recipeIngr.Unit,
		); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error creating recipe ingredients: %w", err)
		}
		recipe.Ingredients = append(recipe.Ingredients, recipeIngr)
	}

	res, err = tx.ExecContext(ctx, sqlDeleteRecipeInst,
		r.RecipeUUID,
	)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error deleting recipe instructions: %w", err)
	}

	if _, affected := res.RowsAffected(); affected != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error deleting recipe instructions, rows not affected(error: %v)", affected)
	}

	for _, recipeinst := range r.Instructions {
		var recipeInst store.RecipeInstruction

		row := tx.QueryRowContext(ctx, sqlCreateRecipeInst,
			&recipe.RecipeUUID,
			&recipeinst.StepNum,
			&recipeinst.Instruction,
		)

		if err = row.Scan(
			&recipeInst.RecipeUUID,
			&recipeInst.StepNum,
			&recipeInst.Instruction,
		); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error creating recipe instructions: %w", err)
		}
		recipe.Instructions = append(recipe.Instructions, recipeInst)
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error updating recipe: %w", err)
	}

	return &recipe, nil
}

func (pg *PG) DeleteRecipe(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error deleting recipe: %w", err)
	}

	//deleting certain recipe's everything from recipe_ingredients table
	res, err := tx.ExecContext(ctx, sqlDeleteRecipeIngr, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting recipe ingredients: %w", err)
	}

	if _, affected := res.RowsAffected(); affected != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting recipe ingredients, rows not affected(error: %v)", affected)
	}

	//deleting certain recipe's everything from recipe_instructions table
	res, err = tx.ExecContext(ctx, sqlDeleteRecipeInst, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting recipe instructions: %w", err)
	}

	if _, affected := res.RowsAffected(); affected != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting recipe instructions, rows not affected(error: %v)", affected)
	}

	//deleting certain recipe's everything from recipes table
	res, err = tx.ExecContext(ctx, sqlDeleteRecipe, id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting recipe: %w", err)
	}

	if affected, _ := res.RowsAffected(); affected != 1 {
		tx.Rollback()
		return fmt.Errorf("error deleting recipe, rows affected is %d instead of 1", affected)
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("error deleting recipe: %w", err)
	}

	return nil
}
