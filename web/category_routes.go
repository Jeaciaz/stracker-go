package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/stracker-go/model"
)

func HandleCategoryRoutes() {
	http.HandleFunc("/categories/buttons", categoryButtons)
}

func categoryButtons(w http.ResponseWriter, r *http.Request) {
	err := WriteCategoryButtons(w, 0, 0)
	if err != nil {
    log.Fatal("Failed to write category buttons: ", err)
	}
}

func WriteCategoryButtons(w http.ResponseWriter, month int, year int) error {
	t, err := template.ParseFiles("web/templates/category_buttons.html")
	if err != nil {
		return fmt.Errorf("Failed to parse category_buttons template: %w", err)
	}
	cats, err := model.GetAllCategories()
	if err != nil {
		return fmt.Errorf("Failed to get categories: %w", err)
	}

	type CategoryButtonData struct {
		Category        *model.Category
		SpendingsAmount string
	}
	buttons := make([]*CategoryButtonData, len(cats))
	for i, cat := range cats {
		spendings, err := model.GetSpendingsForCategory(cat.Id, month, year)
		if err != nil {
			return fmt.Errorf("Failed to get spendings for category %s: %w", cat.Name, err)
		}
		var sum float64
		for _, spending := range spendings {
			sum += spending.Amount
		}
		buttons[i] = &CategoryButtonData{Category: cat, SpendingsAmount: strconv.FormatFloat(sum, 'f', 2, 64)}
	}

	err = t.Execute(w, buttons)
	if err != nil {
		return fmt.Errorf("Failed to execute category_buttons template: %w", err)
	}
	return nil
}
