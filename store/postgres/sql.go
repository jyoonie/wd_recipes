package postgres

const sqlGetRecipe = `
	SELECT 	recipe_uuid,
			user_uuid,
			recipe_name,
			category,
			created_at,
			updated_at
	
	FROM 	wdiet.recipes
	
	WHERE	recipe_uuid = $1

	LIMIT 1
	;
`

// recipe 하나당 여러 recipe ingredient이므로 얘는 필연적으로 list recipe ingredients임,, 하나의 recipe ingredient만 불러오는건 의미 ㄴ
const sqlListRecipeIngr = ` 
	SELECT 	recipe_uuid,
			ingredient_uuid,
			amount,
			unit
	
	FROM 	wdiet.recipe_ingredients
	
	WHERE	recipe_uuid = $1
	;

`
const sqlListRecipeInst = `
	SELECT 	recipe_uuid,
			step_num,
			instruction
	
	FROM 	wdiet.recipe_instructions
	
	WHERE	recipe_uuid = $1
	;
`

const sqlListRecipes = `
	SELECT 	recipe_uuid,
			user_uuid,
			recipe_name,
			category,
			created_at,
			updated_at
	
	FROM 	wdiet.recipes
	
	WHERE	user_uuid = $1
	;
`

const sqlsearchRecipes = `
	SELECT 	recipe_uuid,
			user_uuid,
			recipe_name,
			category,
			created_at,
			updated_at

	FROM wdiet.recipes
`

const sqlCreateRecipe = `
	INSERT INTO wdiet.recipes(
		user_uuid,
		recipe_name,
		category
	)
	VALUES(
		$1,
		$2,
		$3
	)
	RETURNING recipe_uuid, user_uuid, recipe_name, category, created_at, updated_at
	;
`

const sqlCreateRecipeIngr = `
	INSERT INTO wdiet.recipe_ingredients(
		recipe_uuid,
		ingredient_uuid,
		amount,
		unit
	)
	VALUES(
		$1,
		$2,
		$3,
		$4
	)
	RETURNING recipe_uuid, ingredient_uuid, amount, unit
	;
`

const sqlCreateRecipeInst = `
	INSERT INTO wdiet.recipe_instructions(
		recipe_uuid,
		step_num,
		instruction
	)
	VALUES(
		$1,
		$2,
		$3
	)
	RETURNING recipe_uuid, step_num, instruction
	;
`

const sqlUpdateRecipe = `
	UPDATE wdiet.recipes
		SET 
			recipe_name = $1,
			category = $2,
			updated_at = now()
	WHERE recipe_uuid = $3
	RETURNING recipe_uuid, user_uuid, recipe_name, category, created_at, updated_at
	;
`

// AND WHERE ingredient_uuid = $2 이 부분 뺌.. 아예 특정 recipe에 해당하는 모든 ingredient들을 싹 다 지울거니까 어차피.. ingredient uuid 필요 ㄴ. fridge ingredient는 한가지 재료를 지우는 것이므로 두 개의 uuid가 필요.
// methods.go UpdateRecipe에서도 for range 안에 이걸 넣을 게 아니라 이건 바깥으로 빼놓음. 돌아가면서 ingredient들을 하나하나 지우는게 아니라 싹 다 지우고 시작.
const sqlDeleteRecipeIngr = `
	DELETE 
	FROM wdiet.recipe_ingredients

	WHERE recipe_uuid = $1
	;
`

const sqlDeleteRecipeInst = `
	DELETE 
	FROM wdiet.recipe_instructions

	WHERE recipe_uuid = $1
	;
`

const sqlDeleteRecipe = `
	DELETE 
		FROM wdiet.recipes

	WHERE recipe_uuid = $1
	;
`
