package validation

import (
	"fmt"
	"regexp"
	"strings"
)

type Validator struct {
	errors map[string]string
}

func NewValidator() *Validator {
	return &Validator{
		errors: make(map[string]string),
	}
}

func (v *Validator) Required(field, value string) {
	if strings.TrimSpace(value) == "" {
		v.errors[field] = fmt.Sprintf("%s is required", field)
	}
}

func (v *Validator) MinLength(field, value string, min int) {
	if len(strings.TrimSpace(value)) < min {
		v.errors[field] = fmt.Sprintf("%s must be at least %d characters", field, min)
	}
}

func (v *Validator) MaxLength(field, value string, max int) {
	if len(strings.TrimSpace(value)) > max {
		v.errors[field] = fmt.Sprintf("%s must be no more than %d characters", field, max)
	}
}

func (v *Validator) Email(field, value string) {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(strings.TrimSpace(value)) {
		v.errors[field] = fmt.Sprintf("%s must be a valid email address", field)
	}
}

func (v *Validator) In(field, value string, allowed []string) {
	for _, allowedValue := range allowed {
		if value == allowedValue {
			return
		}
	}
	v.errors[field] = fmt.Sprintf("%s must be one of: %s", field, strings.Join(allowed, ", "))
}

func (v *Validator) PositiveInt(field string, value *int) {
	if value != nil && *value <= 0 {
		v.errors[field] = fmt.Sprintf("%s must be a positive number", field)
	}
}

func (v *Validator) PositiveFloat(field string, value *float64) {
	if value != nil && *value <= 0 {
		v.errors[field] = fmt.Sprintf("%s must be a positive number", field)
	}
}

func (v *Validator) RangeInt(field string, value *int, min, max int) {
	if value != nil && (*value < min || *value > max) {
		v.errors[field] = fmt.Sprintf("%s must be between %d and %d", field, min, max)
	}
}

func (v *Validator) ValidDifficulty(field, value string) {
	valid := []string{"easy", "medium", "hard"}
	v.In(field, value, valid)
}

func (v *Validator) ValidColor(field, value string) {
	colorRegex := regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)
	if value != "" && !colorRegex.MatchString(value) {
		v.errors[field] = fmt.Sprintf("%s must be a valid hex color code", field)
	}
}

func (v *Validator) Custom(field string, condition bool, message string) {
	if !condition {
		v.errors[field] = message
	}
}

func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

func (v *Validator) GetErrors() map[string]string {
	return v.errors
}

func (v *Validator) AddError(field, message string) {
	v.errors[field] = message
}

// Recipe validation functions
func ValidateRecipe(recipe interface{}) *Validator {
	v := NewValidator()

	// Type assertion for recipe struct
	type RecipeStruct struct {
		Title        string  `json:"title"`
		Description  string  `json:"description"`
		Instructions string  `json:"instructions"`
		PrepTime     *int    `json:"prep_time"`
		CookTime     *int    `json:"cook_time"`
		Servings     int     `json:"servings"`
		Difficulty   string  `json:"difficulty"`
		ImageURL     *string `json:"image_url"`
		IsPublic     bool    `json:"is_public"`
	}

	r, ok := recipe.(RecipeStruct)
	if !ok {
		v.AddError("recipe", "Invalid recipe format")
		return v
	}

	// Title validation
	v.Required("title", r.Title)
	v.MinLength("title", r.Title, 3)
	v.MaxLength("title", r.Title, 255)

	// Description validation
	v.MaxLength("description", r.Description, 2000)

	// Instructions validation
	v.Required("instructions", r.Instructions)
	v.MinLength("instructions", r.Instructions, 10)
	v.MaxLength("instructions", r.Instructions, 5000)

	// Time validation
	v.PositiveInt("prep_time", r.PrepTime)
	v.PositiveInt("cook_time", r.CookTime)
	v.RangeInt("prep_time", r.PrepTime, 0, 1440) // Max 24 hours
	v.RangeInt("cook_time", r.CookTime, 0, 1440)

	// Servings validation
	v.Custom("servings", r.Servings >= 1 && r.Servings <= 50, "Servings must be between 1 and 50")

	// Difficulty validation
	if r.Difficulty != "" {
		v.ValidDifficulty("difficulty", r.Difficulty)
	}

	// Image URL validation
	if r.ImageURL != nil && *r.ImageURL != "" {
		v.MaxLength("image_url", *r.ImageURL, 500)
		// Basic URL validation
		if !strings.HasPrefix(*r.ImageURL, "http://") && !strings.HasPrefix(*r.ImageURL, "https://") {
			v.AddError("image_url", "Image URL must start with http:// or https://")
		}
	}

	return v
}

// Ingredient validation functions
func ValidateIngredient(ingredient interface{}) *Validator {
	v := NewValidator()

	type IngredientStruct struct {
		Name     string `json:"name"`
		Category string `json:"category"`
	}

	i, ok := ingredient.(IngredientStruct)
	if !ok {
		v.AddError("ingredient", "Invalid ingredient format")
		return v
	}

	v.Required("name", i.Name)
	v.MinLength("name", i.Name, 2)
	v.MaxLength("name", i.Name, 255)

	if i.Category != "" {
		v.MaxLength("category", i.Category, 100)
	}

	return v
}

// Tag validation functions
func ValidateTag(tag interface{}) *Validator {
	v := NewValidator()

	type TagStruct struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}

	t, ok := tag.(TagStruct)
	if !ok {
		v.AddError("tag", "Invalid tag format")
		return v
	}

	v.Required("name", t.Name)
	v.MinLength("name", t.Name, 2)
	v.MaxLength("name", t.Name, 100)

	v.ValidColor("color", t.Color)

	return v
}
