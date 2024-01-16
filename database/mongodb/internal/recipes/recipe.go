package recipes

import (
	"fmt"
)

// Recipe represents a recipe.
type Recipe struct {
	ID          string        `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Description string        `json:"description" bson:"description"`
	Favorite    bool          `json:"favorite" bson:"favorite"`
	Ingredients []*Ingredient `json:"ingredients" bson:"ingredients"`
	Directions  []*Step       `json:"directions" bson:"directions"`
}

// String returns a human-readable and formatted representation of the Recipe.
func (r *Recipe) String() string {
	var ingredientsString string
	for _, i := range r.Ingredients {
		ingredientsString += fmt.Sprintf("  - %s\n", i)
	}

	var directionsString string
	for _, d := range r.Directions {
		directionsString += fmt.Sprintf("  - %s\n", d)
	}

	return fmt.Sprintf("Recipe ID: %s\nName: %s\nDescription: %s\nFavorite: %t\nIngredients:\n%sDirections:\n%s",
		r.ID, r.Name, r.Description, r.Favorite, ingredientsString, directionsString)
}

// Ingredient represents an ingredient in a recipe.
type Ingredient struct {
	Qty  float64 `json:"qty" bson:"qty"`
	Unit string  `json:"unit" bson:"unit"`
	Name string  `json:"name" bson:"name"`
}

// String returns a human-readable representation of the Ingredient.
func (i *Ingredient) String() string {
	return fmt.Sprintf("Qty: %f %s, Name: %s", i.Qty, i.Unit, i.Name)
}

// Step represents a step in the recipe procedure.
type Step struct {
	Order       int    `json:"order" bson:"order"`
	Description string `json:"description" bson:"description"`
}

// String returns a human-readable representation of the Step.
func (s *Step) String() string {
	return fmt.Sprintf("Order: %d, Description: %s", s.Order, s.Description)
}
