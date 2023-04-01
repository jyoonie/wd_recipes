package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"wd_recipes/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) GetRecipe(c *gin.Context) {
	l := s.l.Named("GetRecipe")

	id := c.Param("id")

	rid, err := uuid.Parse(id)
	if err != nil {
		l.Info("error getting recipe", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	recipe, err := s.db.GetRecipe(context.Background(), rid)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			l.Info("error getting recipe", zap.Error(err))
			c.Status(http.StatusNotFound)
			return
		}
		l.Error("error getting recipe", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dbRecipe2ApiRecipe(recipe))
}

func (s *Service) ListRecipes(c *gin.Context) {
	l := s.l.Named("ListRecipes")

	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		l.Info("error listing recipes", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	recipes, err := s.db.ListRecipes(context.Background(), uid)
	if err != nil {
		l.Error("error listing recipes", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(recipes) == 0 {
		c.Status(http.StatusOK)
		return
	}

	var listRecipesResponse []Recipe

	for _, recipe := range recipes {
		r := dbRecipe2ApiRecipe(&recipe)
		listRecipesResponse = append(listRecipesResponse, r)
	}

	c.JSON(http.StatusOK, listRecipesResponse)
}

func (s *Service) SearchRecipes(c *gin.Context) {
	l := s.l.Named("SearchRecipes")

	var searchRecipesRequest SearchRecipes

	if err := json.NewDecoder(c.Request.Body).Decode(&searchRecipesRequest); err != nil {
		l.Info("error searching recipes", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	if !isValidSearchRecipesRequest(searchRecipesRequest) {
		l.Info("error searching recipes")
		c.Status(http.StatusBadRequest)
		return
	}

	recipes, err := s.db.SearchRecipes(context.Background(), apiSearchR2DBSearchR(searchRecipesRequest))
	if err != nil {
		l.Error("error searching recipes", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(recipes) == 0 {
		c.Status(http.StatusOK)
		return
	}

	var searchRecipesResponse []Recipe

	for _, recipe := range recipes {
		r := dbRecipe2ApiRecipe(&recipe)
		searchRecipesResponse = append(searchRecipesResponse, r)
	}

	c.JSON(http.StatusOK, searchRecipesResponse)
}

func (s *Service) CreateRecipe(c *gin.Context) {
	l := s.l.Named("CreateRecipe")

	var createRecipeRequest Recipe

	if err := json.NewDecoder(c.Request.Body).Decode(&createRecipeRequest); err != nil {
		l.Info("error creating recipe", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	if !isValidCreateRecipeRequest(createRecipeRequest) {
		l.Info("error creating recipe")
		c.Status(http.StatusBadRequest)
		return
	}

	recipe, err := s.db.CreateRecipe(context.Background(), apiRecipe2DBRecipe(createRecipeRequest))
	if err != nil {
		l.Error("error creating recipe", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dbRecipe2ApiRecipe(recipe))
}

func (s *Service) UpdateRecipe(c *gin.Context) {
	l := s.l.Named("UpdateRecipe")

	id := c.Param("id")

	rid, err := uuid.Parse(id)
	if err != nil {
		l.Info("error updating recipe", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	var updateRecipeRequest Recipe

	if err := json.NewDecoder(c.Request.Body).Decode(&updateRecipeRequest); err != nil {
		l.Info("error updating recipe", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	if !isValidUpdateRecipeRequest(updateRecipeRequest, rid) {
		l.Info("error updating recipe", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	recipe, err := s.db.UpdateRecipe(context.Background(), apiRecipe2DBRecipe(updateRecipeRequest))
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			l.Info("error updating recipe", zap.Error(err))
			c.Status(http.StatusNotFound)
			return
		}
		l.Error("error updating recipe", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, dbRecipe2ApiRecipe(recipe))
}

func (s *Service) DeleteRecipe(c *gin.Context) {
	l := s.l.Named("DeleteRecipe")

	id := c.Param("id")

	rid, err := uuid.Parse(id)
	if err != nil {
		l.Info("error deleting recipe", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	if err = s.db.DeleteRecipe(context.Background(), rid); err != nil {
		l.Error("error deleting recipe", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
