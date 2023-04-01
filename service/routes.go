package service

func (s *Service) registerRoutes() {
	authorized := s.r.Group("/")
	//authorized.Use(s.ValidateToken)
	{
		authorized.GET("/recipes/:id", s.GetRecipe)
		authorized.GET("/users/:id/recipes", s.ListRecipes)
		authorized.POST("/recipes/search", s.SearchRecipes)
		authorized.POST("/recipes", s.CreateRecipe)
		authorized.POST("/recipes/:id", s.UpdateRecipe)
		authorized.DELETE("/recipes/:id", s.DeleteRecipe)
	}
}
