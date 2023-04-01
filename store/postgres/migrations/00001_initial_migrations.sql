-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS wdiet;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wdiet.recipes
(
    recipe_uuid uuid not null default gen_random_uuid()
        constraint recipes_primary_key
            primary key,
    user_uuid            uuid            not null
        constraint user_uuid_fk references wdiet.users,
    recipe_name          varchar(64)     not null,
    category             varchar(64)     not null,
    created_at           timestamp       not null default now(),
    updated_at           timestamp       not null default now()
);

CREATE INDEX ON wdiet.recipes (user_uuid);
-- ALTER TABLE wdiet.recipes ADD FOREIGN KEY user_uuid REFERENCES wdiet.users (user_uuid);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS wdiet;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wdiet.recipe_ingredients
(
    recipe_uuid            uuid            not null
        constraint recipe_uuid_fk references wdiet.recipes,
    ingredient_uuid        uuid            not null
        constraint ingredient_uuid_fk references wdiet.ingredients,
    amount                 integer         not null,
    unit                   varchar(64)     not null,
    PRIMARY KEY (recipe_uuid, ingredient_uuid)
    -- created_at             timestamp       not null default now(), --레시피에 있는 필드니까 안해줘도 상관없음.
    -- updated_at             timestamp       not null default now()
);

-- CREATE INDEX ON wdiets.recipe_ingredients (recipe_uuid);
-- ALTER TABLE wdiet.recipe_ingredients ADD FOREIGN KEY recipe_uuid REFERENCES wdiet.recipes (recipe_uuid);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS wdiet;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wdiet.recipe_instructions
(
    recipe_uuid            uuid            not null
        constraint recipe_uuid_fk references wdiet.recipes,
    step_num               integer         not null,
    instruction            varchar(1024)   not null,
    PRIMARY KEY (recipe_uuid, step_num) --하나의 레시피에 한가지의 단계만 guarantee하기 위함. 1. 양파를 썬다. 1. 육수를 낸다. 1. 고기를 썬다. 이러면 안되니까. 1, 2, 3단계여야 하니까.
    -- created_at             timestamp       not null default now(),
    -- updated_at             timestamp       not null default now()
);

-- CREATE INDEX ON wdiets.recipe_instructions (recipe_uuid); efficient to search the table on this column. Binary tree.
-- ALTER TABLE wdiet.recipe_instructions ADD FOREIGN KEY recipe_uuid REFERENCES wdiet.recipes (recipe_uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wdiet.recipe_ingredients;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS wdiet.recipe_instructions;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS wdiet.recipes;
-- +goose StatementEnd
