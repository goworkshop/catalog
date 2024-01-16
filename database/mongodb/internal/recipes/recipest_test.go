package recipes_test

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/goworkshop/catalog/database/mongodb/internal/recipes"
)

var (
	testDBURL  = "mongodb://localhost:27017"
	testDBName = "test_goworkshop"
)

func setupTestDB(t *testing.T) *mongo.Database {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(testDBURL))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Print connection details for debugging
	t.Log("Connected to MongoDB:", testDBURL)

	db := client.Database(testDBName)

	return db
}

func cleanupTestDB(t *testing.T, db *mongo.Database) {
	if err := db.Drop(context.Background()); err != nil {
		t.Fatalf("Failed to drop test database: %v", err)
	}
}

func TestRecipeRepository_GetAll(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := recipes.NewRecipeRepository(db)

	// Create test recipes
	testRecipes := createTestRecipes()
	for _, recipe := range testRecipes {
		_, err := repo.Create(context.Background(), recipe)
		if err != nil {
			t.Fatalf("Failed to create recipe: %v", err)
		}
	}

	// Test GetAll function
	retrievedRecipes, err := repo.GetAll(context.Background())
	if err != nil {
		t.Fatalf("Failed to get all recipes: %v", err)
	}

	// Assertions
	if len(retrievedRecipes) != len(testRecipes) {
		t.Errorf("Expected %d recipes, got %d", len(testRecipes), len(retrievedRecipes))
	}

	// Additional assertions can be added to compare individual recipes if needed
	for i, retrievedRecipe := range retrievedRecipes {
		assertRecipe(t, testRecipes[i], retrievedRecipe)
	}
}

func TestRecipeRepository_Create(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := recipes.NewRecipeRepository(db)

	testRecipes := createTestRecipes()

	// Test Create function for each test recipe
	for _, recipeToCreate := range testRecipes {
		createdRecipe, err := repo.Create(context.Background(), recipeToCreate)
		if err != nil {
			t.Fatalf("Failed to create recipe: %v", err)
		}

		// Retrieve the created recipe
		retrievedRecipe, err := repo.Get(context.Background(), createdRecipe.ID)
		if err != nil {
			t.Fatalf("Failed to retrieve created recipe: %v", err)
		}

		// Assertions
		assertRecipe(t, createdRecipe, retrievedRecipe)
	}
}

func TestRecipeRepository_Get(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := recipes.NewRecipeRepository(db)

	// Create a test recipe
	createdRecipe, err := repo.Create(context.Background(), &recipes.Recipe{
		Name:        "Test Recipe",
		Description: "This is a test recipe.",
		Favorite:    true,
		Ingredients: []*recipes.Ingredient{
			{Qty: 100, Unit: "g", Name: "Ingredient 1"},
			{Qty: 50, Unit: "ml", Name: "Ingredient 2"},
		},
		Directions: []*recipes.Step{
			{Order: 1, Description: "Step 1"},
			{Order: 2, Description: "Step 2"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	// Test Get function
	retrievedRecipe, err := repo.Get(context.Background(), createdRecipe.ID)
	if err != nil {
		t.Fatalf("Failed to get recipe: %v", err)
	}

	// Assertions
	assertRecipe(t, createdRecipe, retrievedRecipe)
}

func TestRecipeRepository_Update(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := recipes.NewRecipeRepository(db)

	// Create a test recipe
	createdRecipe, err := repo.Create(context.Background(), createTestRecipes()[0])
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	// Modify the created recipe
	createdRecipe.Name = "Updated Test Recipe"
	createdRecipe.Description = "This is an updated test recipe."
	createdRecipe.Favorite = false

	// Test Update function
	err = repo.Update(context.Background(), createdRecipe)
	if err != nil {
		t.Fatalf("Failed to update recipe: %v", err)
	}

	// Retrieve the updated recipe
	updatedRecipe, err := repo.Get(context.Background(), createdRecipe.ID)
	if err != nil {
		t.Fatalf("Failed to get updated recipe: %v", err)
	}

	// Assertions
	assertRecipe(t, createdRecipe, updatedRecipe)
}

func TestRecipeRepository_Delete(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	repo := recipes.NewRecipeRepository(db)

	// Create a test recipe
	testRecipes := createTestRecipes()
	createdRecipe, err := repo.Create(context.Background(), testRecipes[0])
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	// Test Delete function
	err = repo.Delete(context.Background(), createdRecipe.ID)
	if err != nil {
		t.Fatalf("Failed to delete recipe: %v", err)
	}

	// Attempt to retrieve the deleted recipe
	_, err = repo.Get(context.Background(), createdRecipe.ID)

	// Assertions
	if err == nil {
		t.Error("Expected error, but recipe was retrieved after deletion")
	}
}

func createTestRecipes() []*recipes.Recipe {
	return []*recipes.Recipe{
		{
			Name:        "Test Recipe 1",
			Description: "This is a test recipe 1.",
			Favorite:    true,
			Ingredients: []*recipes.Ingredient{
				{Qty: 100, Unit: "g", Name: "Ingredient 1"},
				{Qty: 50, Unit: "ml", Name: "Ingredient 2"},
			},
			Directions: []*recipes.Step{
				{Order: 1, Description: "Step 1"},
				{Order: 2, Description: "Step 2"},
			},
		},
		{
			Name:        "Test Recipe 2",
			Description: "This is a test recipe 2.",
			Ingredients: []*recipes.Ingredient{
				{Qty: 200, Unit: "g", Name: "Ingredient 1"},
				{Qty: 75, Unit: "ml", Name: "Ingredient 2"},
			},
			Directions: []*recipes.Step{
				{Order: 1, Description: "Step 1"},
				{Order: 2, Description: "Step 2"},
				{Order: 3, Description: "Step 3"},
			},
		},
	}
}

func assertRecipe(t *testing.T, recipe1, recipe2 *recipes.Recipe) {
	if recipe1.ID != recipe2.ID {
		t.Errorf("Expected ID %s, got %s", recipe1.ID, recipe2.ID)
	}

	if recipe1.Name != recipe2.Name {
		t.Errorf("Expected Name %s, got %s", recipe1.Name, recipe2.Name)
	}

	if recipe1.Description != recipe2.Description {
		t.Errorf("Expected Description %s, got %s", recipe1.Description, recipe2.Description)
	}

	if recipe1.Favorite != recipe2.Favorite {
		t.Errorf("Expected Favorite %t, got %t", recipe1.Favorite, recipe2.Favorite)
	}

	assertIngredients(t, recipe1.Ingredients, recipe2.Ingredients)
	assertSteps(t, recipe1.Directions, recipe2.Directions)
}

func assertIngredients(t *testing.T, ingredients1, ingredients2 []*recipes.Ingredient) {
	if len(ingredients1) != len(ingredients2) {
		t.Errorf("Expected %d ingredients, got %d", len(ingredients1), len(ingredients2))
	}

	for i, ingredient1 := range ingredients1 {
		ingredient2 := ingredients2[i]
		if ingredient1.Name != ingredient2.Name {
			t.Errorf("Expected ingredient %s, got %s", ingredient1.Name, ingredient2.Name)
		}
		if ingredient1.Qty != ingredient2.Qty {
			t.Errorf("Expected ingredient quantity %f, got %f", ingredient1.Qty, ingredient2.Qty)
		}
		if ingredient1.Unit != ingredient2.Unit {
			t.Errorf("Expected ingredient unit %s, got %s", ingredient1.Unit, ingredient2.Unit)
		}
	}
}

func assertSteps(t *testing.T, steps1, steps2 []*recipes.Step) {
	if len(steps1) != len(steps2) {
		t.Errorf("Expected %d steps, got %d", len(steps1), len(steps2))
	}

	for i, step1 := range steps1 {
		step2 := steps2[i]
		if step1.Description != step2.Description {
			t.Errorf("Expected step description %s, got %s", step1.Description, step2.Description)
		}
		if step1.Order != step2.Order {
			t.Errorf("Expected step order %d, got %d", step1.Order, step2.Order)
		}
	}
}
