package recipes

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RecipeRepository handles CRUD operations for recipes.
type RecipeRepository struct {
	collection *mongo.Collection
}

// NewRecipeRepository creates a new RecipeRepository.
func NewRecipeRepository(db *mongo.Database) *RecipeRepository {
	collection := db.Collection("recipes")
	return &RecipeRepository{collection: collection}
}

// GetAll retrieves a list of recipes from the database.
func (r *RecipeRepository) GetAll(ctx context.Context) ([]*Recipe, error) {
	var recipes []*Recipe

	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to list recipes: %v", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var recipe Recipe
		if err := cursor.Decode(&recipe); err != nil {
			return nil, fmt.Errorf("failed to decode recipe: %v", err)
		}
		recipes = append(recipes, &recipe)
	}

	return recipes, nil
}

// Create adds a new recipe to the database.
func (r *RecipeRepository) Create(ctx context.Context, recipe *Recipe) (*Recipe, error) {
	result, err := r.collection.InsertOne(ctx, recipe)
	if err != nil {
		return nil, fmt.Errorf("failed to insert recipe: %v", err)
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to convert InsertedID to primitive.ObjectID")
	}

	recipe.ID = insertedID.Hex()
	return recipe, nil
}

// Get retrieves a recipe by its ID.
func (r *RecipeRepository) Get(ctx context.Context, recipeID string) (*Recipe, error) {
	var recipe Recipe

	// Convert the recipeID string to an ObjectID
	objectID, err := primitive.ObjectIDFromHex(recipeID)
	if err != nil {
		return nil, fmt.Errorf("invalid recipe ID format: %v", err)
	}

	filter := bson.D{{"_id", objectID}}

	err = r.collection.FindOne(ctx, filter).Decode(&recipe)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("recipe not found")
		}
		return nil, fmt.Errorf("failed to get recipe: %v", err)
	}

	return &recipe, nil
}

// Update updates a recipe in the database.
func (r *RecipeRepository) Update(ctx context.Context, recipe *Recipe) error {
	// Convert the recipe ID string to an ObjectID
	objectID, err := primitive.ObjectIDFromHex(recipe.ID)
	if err != nil {
		return fmt.Errorf("invalid recipe ID format: %v", err)
	}

	// Create a map to represent the fields to be updated
	updateFields := bson.M{
		"name":        recipe.Name,
		"description": recipe.Description,
		"favorite":    recipe.Favorite,
		"ingredients": recipe.Ingredients,
		"directions":  recipe.Directions,
	}

	update := bson.M{"$set": updateFields}

	result, err := r.collection.UpdateByID(ctx, objectID, update)
	if err != nil {
		return fmt.Errorf("failed to update recipe: %v", err)
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("recipe not found")
	}

	return nil
}

// Delete removes a recipe from the database.
func (r *RecipeRepository) Delete(ctx context.Context, recipeID string) error {
	// Convert the recipe ID string to an ObjectID
	objectID, err := primitive.ObjectIDFromHex(recipeID)
	if err != nil {
		return fmt.Errorf("invalid recipe ID format: %v", err)
	}

	filter := bson.D{{"_id", objectID}}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete recipe: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("recipe not found")
	}

	return nil
}
