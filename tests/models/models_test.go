package models_test

import (
	"books-api/app/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColor_IsValid(t *testing.T) {
	tests := []struct {
		name  string
		color models.Color
		want  bool
	}{
		{"Red is valid", models.Red, true},
		{"Green is valid", models.Green, true},
		{"Blue is valid", models.Blue, true},
		{"Purple is invalid", models.Color("Purple"), false},
		{"Empty is invalid", models.Color(""), false},
		{"Random is invalid", models.Color("Random"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.color.IsValid())
		})
	}
}

func TestColor_Value(t *testing.T) {
	color := models.Red
	value, err := color.Value()
	assert.NoError(t, err)
	assert.Equal(t, "Red", value)

	emptyColor := models.Color("")
	value, err = emptyColor.Value()
	assert.NoError(t, err)
	assert.Nil(t, value)
}

func TestColor_Scan(t *testing.T) {
	var color models.Color
	
	// Test scanning a string
	err := color.Scan("Blue")
	assert.NoError(t, err)
	assert.Equal(t, models.Blue, color)

	// Test scanning nil
	err = color.Scan(nil)
	assert.NoError(t, err)
	assert.Equal(t, models.Color(""), color)

	// Test scanning invalid type
	err = color.Scan(123)
	assert.Error(t, err)
}

func TestBook_Update(t *testing.T) {
	red := models.Red
	blue := models.Blue
	
	original := &models.Book{
		ID:     1,
		Title:  "Original Title",
		Author: "Original Author",
		Pages:  100,
		Color:  &red,
	}

	update := models.Book{
		Title: "Updated Title",
		Pages: 200,
		Color: &blue,
	}

	original.Update(update)

	assert.Equal(t, "Updated Title", original.Title)
	assert.Equal(t, "Original Author", original.Author) // Should not change
	assert.Equal(t, 200, original.Pages)
	assert.Equal(t, &blue, original.Color)
}

func TestBook_Update_PartialUpdate(t *testing.T) {
	red := models.Red
	
	original := &models.Book{
		ID:     1,
		Title:  "Original Title",
		Author: "Original Author",
		Pages:  100,
		Color:  &red,
	}

	// Only update title
	update := models.Book{
		Title: "New Title",
	}

	original.Update(update)

	assert.Equal(t, "New Title", original.Title)
	assert.Equal(t, "Original Author", original.Author)
	assert.Equal(t, 100, original.Pages) // Should not change (0 is ignored)
	assert.Equal(t, &red, original.Color) // Should not change
}