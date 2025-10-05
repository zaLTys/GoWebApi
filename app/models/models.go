package models

import (
	"database/sql/driver"
	"errors"
)

// Color as a string-based enum for easy JSON + DB storage
type Color string

const (
	Red   Color = "Red"
	Green Color = "Green"
	Blue  Color = "Blue"
)

func (c Color) IsValid() bool {
	switch c {
	case Red, Green, Blue:
		return true
	}
	return false
}

// Implement driver.Valuer for DB write
func (c Color) Value() (driver.Value, error) {
	if c == "" {
		return nil, nil
	}
	return string(c), nil
}

// Implement sql.Scanner for DB read
func (c *Color) Scan(value interface{}) error {
	if value == nil {
		*c = ""
		return nil
	}
	v, ok := value.(string)
	if !ok {
		return errors.New("invalid Color type")
	}
	*c = Color(v)
	return nil
}

// Book model (GORM automatically creates table 'books')
type Book struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Author string `json:"author"`
	Title  string `json:"title"`
	Pages  int    `json:"pages"`
	Color  *Color `json:"color,omitempty"`
}

func (book *Book) Update(update Book) {
	if update.Author != "" {
		book.Author = update.Author
	}
	if update.Title != "" {
		book.Title = update.Title
	}
	if update.Pages != 0 {
		book.Pages = update.Pages
	}
	if update.Color != nil {
		book.Color = update.Color
	}
}
