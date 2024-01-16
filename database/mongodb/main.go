package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/goworkshop/catalog/database/mongodb/internal/recipes"
)

func main() {
	// Connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Create a new MongoDB database
	db := client.Database("goworkshop")

	// Initialize the repository
	recipeRepo := recipes.NewRecipeRepository(db)

	// Some recipes to play with
	recipe1 := &recipes.Recipe{
		Name:        "Pasta Carbonara",
		Description: "Classic Italian pasta dish with eggs, cheese, pancetta, and black pepper.",
		Favorite:    true,
		Ingredients: []*recipes.Ingredient{
			{Qty: 200, Unit: "g", Name: "Spaghetti"},
			{Qty: 100, Unit: "g", Name: "Pancetta"},
			{Qty: 2, Unit: "", Name: "Eggs"},
			{Qty: 50, Unit: "g", Name: "Parmesan Cheese"},
			{Qty: 50, Unit: "g", Name: "Pecorino Cheese"},
			{Qty: 1, Unit: "", Name: "Black Pepper"},
		},
		Directions: []*recipes.Step{
			{Order: 1, Description: "Boil spaghetti until al dente."},
			{Order: 2, Description: "Fry pancetta until crispy."},
			{Order: 3, Description: "Mix eggs with grated cheeses."},
			{Order: 4, Description: "Combine everything and add black pepper."},
		},
	}

	recipe2 := &recipes.Recipe{
		Name:        "Caprese Salad",
		Description: "Simple and delicious salad with fresh tomatoes, mozzarella, and basil.",
		Ingredients: []*recipes.Ingredient{
			{Qty: 4, Unit: "", Name: "Tomatoes"},
			{Qty: 1, Unit: "", Name: "Mozzarella"},
			{Qty: 1, Unit: "bunch", Name: "Basil"},
			{Qty: 2, Unit: "tbsp", Name: "Olive Oil"},
			{Qty: 1, Unit: "tbsp", Name: "Balsamic Vinegar"},
		},
		Directions: []*recipes.Step{
			{Order: 1, Description: "Slice tomatoes and mozzarella."},
			{Order: 2, Description: "Arrange slices on a plate with basil leaves."},
			{Order: 3, Description: "Drizzle with olive oil and balsamic vinegar."},
		},
	}

	recipe3 := &recipes.Recipe{
		Name:        "Chocolate Cake",
		Description: "Decadent chocolate cake with rich frosting.",
		Favorite:    true,
		Ingredients: []*recipes.Ingredient{
			{Qty: 200, Unit: "g", Name: "Flour"},
			{Qty: 50, Unit: "g", Name: "Cocoa Powder"},
			{Qty: 200, Unit: "g", Name: "Sugar"},
			{Qty: 200, Unit: "g", Name: "Butter"},
			{Qty: 4, Unit: "", Name: "Eggs"},
		},
		Directions: []*recipes.Step{
			{Order: 1, Description: "Mix dry ingredients (flour, cocoa powder, sugar)."},
			{Order: 2, Description: "Cream butter and sugar, then add eggs one by one."},
			{Order: 3, Description: "Fold in the dry ingredients."},
			{Order: 4, Description: "Bake in preheated oven."},
		},
	}

	// Create them in the database
	createdRecipe1, err := recipeRepo.Create(ctx, recipe1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created Recipe 1:\n%+v\n", createdRecipe1)

	createdRecipe2, err := recipeRepo.Create(ctx, recipe2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created Recipe 2:\n%+v\n", createdRecipe2)

	createdRecipe3, err := recipeRepo.Create(ctx, recipe3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created Recipe 3:\n%+v\n", createdRecipe3)

	// List all
	allRecipes, err := recipeRepo.GetAll(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("All Recipes:")
	for _, r := range allRecipes {
		fmt.Printf("%+v\n", r)
	}

	// Get one by ID
	recipeByID, err := recipeRepo.Get(ctx, createdRecipe1.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Recipe by ID:\n%+v\n", recipeByID)

	// Update one
	recipeToUpdate := createdRecipe1
	recipeToUpdate.Favorite = false
	err = recipeRepo.Update(ctx, recipeToUpdate)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Updated Recipe")

	// Delete one
	err = recipeRepo.Delete(ctx, createdRecipe3.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted Recipe 3")

	// List all again
	allRecipes, err = recipeRepo.GetAll(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("All Recipes after deletion:")
	for _, r := range allRecipes {
		fmt.Printf("%+v\n", r)
	}
}
