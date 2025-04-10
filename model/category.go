package model

import (
	"fmt"

	"github.com/stracker-go/pkg"
)

type Category struct {
  Id string
  Name string
  Emoji string
}

func GetAllCategories() ([]*Category, error) {
  db := pkg.GetDb()
  rows, err := db.Query("SELECT * FROM categories")
  if err != nil {
    return nil, fmt.Errorf("Failed to load categories from DB: %w", err)
  }

  categories := []*Category{}
  for rows.Next() {
    var c Category
    err = rows.Scan(&c.Id, &c.Name, &c.Emoji)
    if err != nil {
      return nil, fmt.Errorf("Failed to scan category from SQL response: %w", err)
    }
    categories = append(categories, &c)
  }
  return categories, nil
}

func GetCategoryById(id string) (*Category, error) {
  db := pkg.GetDb()
  rows, err := db.Query("SELECT * FROM categories WHERE id = ?", id)
  if err != nil {
    return nil, fmt.Errorf("Failed to load category from DB: %w", err)
  }
  if !rows.Next() {
    return nil, fmt.Errorf("Category not found")
  }
  var c Category
  err = rows.Scan(&c.Id, &c.Name, &c.Emoji)
  if err != nil {
    return nil, fmt.Errorf("Failed to scan category from SQL response: %w", err)
  }
  return &c, nil
}
