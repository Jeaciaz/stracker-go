package web

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/stracker-go/model"
)

func HandleSpendingRoutes() {
	http.HandleFunc("/spending/create-modal", submitSpendingModal)
	http.HandleFunc("/spending/create", createSpending)
	http.HandleFunc("/spending/edit-modal", editSpendingModal)
	http.HandleFunc("/spending/edit", editSpending)
	http.HandleFunc("/spending/delete-modal", deleteSpendingModal)
	http.HandleFunc("/spending/delete", deleteSpending)
	http.HandleFunc("/spending/timeline", handleSpendingTimeline)
	http.HandleFunc("/spending/list", handleSpendingList)
}

func submitSpendingModal(w http.ResponseWriter, r *http.Request) {
	category, err := model.GetCategoryById(r.URL.Query().Get("category"))
	if err != nil {
		log.Fatal("Failed to get category: ", err)
	}
	t, err := template.ParseFiles("web/templates/submit-spending-modal.html")
	if err != nil {
		log.Fatal("Failed to parse submit-spending-modal template: ", err)
	}
	t.Execute(w, category)
}

func createSpending(w http.ResponseWriter, r *http.Request) {
	category, err := model.GetCategoryById(r.URL.Query().Get("category"))
	if err != nil {
		log.Fatal("Failed to get category: ", err)
	}

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		log.Fatal("Failed to parse amount: ", err)
	}
	err = model.CreateSpending(r.URL.Query().Get("category"), "Test user", amount, r.FormValue("description"))
	if err != nil {
		log.Fatal("Failed to create spending: ", err)
	}

	err = WriteConfetti(w, category.Emoji)
	if err != nil {
		log.Fatal("Failed to write confetti: ", err)
	}
	err = WriteCloseModal(w)
	if err != nil {
		log.Fatal("Failed to write close modal: ", err)
	}
	err = WriteSpendingData(w, 0, 0)
	if err != nil {
		log.Fatal("Failed to write spending data: ", err)
	}
	err = WriteCategoryButtons(w)
	if err != nil {
		log.Fatal("Failed to write category buttons: ", err)
	}
}

func editSpendingModal(w http.ResponseWriter, r *http.Request) {
	spending, err := model.GetSpendingById(r.URL.Query().Get("spendingId"))
	if err != nil {
		log.Fatal("Failed to get spending: ", err)
	}
	t, err := template.ParseFiles("web/templates/edit-spending-modal.html")
	if err != nil {
		log.Fatal("Failed to parse edit-spending-modal template: ", err)
	}
	err = t.Execute(w, spending)
	if err != nil {
		log.Fatal("Failed to execute edit-spending-modal template: ", err)
	}
}

func editSpending(w http.ResponseWriter, r *http.Request) {
	spending, err := model.GetSpendingById(r.URL.Query().Get("spendingId"))
	if err != nil {
		log.Fatal("Failed to get spending: ", err)
	}
	category, err := model.GetCategoryById(spending.Category)
	if err != nil {
		log.Fatal("Failed to get category: ", err)
	}
	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		log.Fatal("Failed to parse amount: ", err)
	}
	err = model.UpdateSpending(r.URL.Query().Get("spendingId"), amount, r.FormValue("description"))
	if err != nil {
		log.Fatal("Failed to update spending: ", err)
	}

	err = WriteCloseModal(w)
	if err != nil {
		log.Fatal("Failed to write close modal: ", err)
	}
	err = WriteConfetti(w, category.Emoji)
	if err != nil {
		log.Fatal("Failed to write confetti: ", err)
	}
	err = WriteSpendingData(w, 0, 0)
	if err != nil {
		log.Fatal("Failed to write spending data: ", err)
	}
}

func deleteSpendingModal(w http.ResponseWriter, r *http.Request) {
	spending, err := model.GetSpendingById(r.URL.Query().Get("spendingId"))
	if err != nil {
		log.Fatal("Failed to get spending: ", err)
	}
	t, err := template.ParseFiles("web/templates/delete-spending-modal.html")
	if err != nil {
		log.Fatal("Failed to parse delete-spending-modal template: ", err)
	}
	err = t.Execute(w, spending)
	if err != nil {
		log.Fatal("Failed to execute delete-spending-modal template: ", err)
	}
}

func deleteSpending(w http.ResponseWriter, r *http.Request) {
	spending, err := model.GetSpendingById(r.URL.Query().Get("spendingId"))
	if err != nil {
		log.Fatal("Failed to get spending: ", err)
	}
	category, err := model.GetCategoryById(spending.Category)
	if err != nil {
		log.Fatal("Failed to get category: ", err)
	}
	err = model.DeleteSpending(r.URL.Query().Get("spendingId"))
	if err != nil {
		log.Fatal("Failed to delete spending: ", err)
	}

	err = WriteCloseModal(w)
	if err != nil {
		log.Fatal("Failed to write close modal: ", err)
	}
	err = WriteConfetti(w, category.Emoji)
	if err != nil {
		log.Fatal("Failed to write confetti: ", err)
	}
	err = WriteSpendingData(w, 0, 0)
	if err != nil {
		log.Fatal("Failed to write spending data: ", err)
	}
}

func handleSpendingTimeline(w http.ResponseWriter, r *http.Request) {
	monthQ := r.URL.Query().Get("month")
	yearQ := r.URL.Query().Get("year")
	year, err := strconv.ParseInt(yearQ, 0, 0)
	month, err := strconv.ParseInt(monthQ, 0, 0)
	if err != nil {
		year = 0
	}
	err = writeSpendingTimeline(w, int(month), int(year))
	if err != nil {
		log.Fatal("Failed to write spending timeline: ", err)
	}
}

type SpendingListItem struct {
	Id           string
	Amount       float64
	CategoryName string
	Description  string
	Datetime     string
}

func handleSpendingList(w http.ResponseWriter, r *http.Request) {
	err := writeSpendingList(w)
	if err != nil {
		log.Fatal("Failed to write spending list: ", err)
	}
}

func WriteSpendingData(w http.ResponseWriter, month int, year int) error {
	err := WriteCategoryButtons(w, month, year)
	if err != nil {
		return fmt.Errorf("Failed to write category buttons: %w", err)
	}
	err = writeSpendingList(w)
	if err != nil {
		return fmt.Errorf("Failed to write spending list: %w", err)
	}
	err = writeSpendingTimeline(w, 0, 0)
	if err != nil {
		return fmt.Errorf("Failed to write spending timeline: %w", err)
	}
	return nil
}

func writeSpendingList(w http.ResponseWriter) error {
	spendings, err := model.GetAllSpendings()
	if err != nil {
		return err
	}
	t, err := template.ParseFiles("web/templates/spending-list.html")
	if err != nil {
		return err
	}
	spendingItems := make([]*SpendingListItem, 0, len(spendings))
	for _, spending := range spendings {
		category, err := model.GetCategoryById(spending.Category)
		if err != nil {
			return err
		}
		spendingItems = append(spendingItems, &SpendingListItem{
			Id:           spending.Id,
			Amount:       spending.Amount,
			CategoryName: category.Name,
			Description:  spending.Description,
			Datetime:     time.Unix(spending.Timestamp, 0).Format("2006-01-02 15:04:05"),
		})
	}
	err = t.Execute(w, spendingItems)
	if err != nil {
		return err
	}
	return nil
}

func writeSpendingTimeline(w http.ResponseWriter, month int, year int) error {
	timeline, err := model.GetSpendingsTimeline()
	if err != nil {
		return err
	}
	t, err := template.ParseFiles("web/templates/spending-timeline.html")
	if err != nil {
		return err
	}

	type TimelineDisplay struct {
		SpendingsSum  float64
		PeriodLabel   string
		PeriodMonth   int
		PeriodYear    int
		IsHighlighted bool
	}
	timelineDisplay := make([]*TimelineDisplay, 0, len(timeline))
	for i, timelineEntry := range timeline {
		periodStart := time.Unix(timelineEntry.PeriodStart, 0)
		var sum float64
		for _, spending := range timelineEntry.Spendings {
			sum += spending.Amount
		}
		periodMonth := int(periodStart.Month())
		periodYear := periodStart.Year()
		isHighlighted := false
		if periodMonth == month && periodYear == year && year > 0 {
			isHighlighted = true
		}
		if year == 0 && i == 0 {
			isHighlighted = true
		}
		timelineDisplay = append(timelineDisplay, &TimelineDisplay{
			SpendingsSum:  sum,
			PeriodLabel:   periodStart.Format("Jan '2006"),
			PeriodMonth:   periodMonth,
			PeriodYear:    periodYear,
			IsHighlighted: isHighlighted,
		})
	}

	err = t.Execute(w, timelineDisplay)
	if err != nil {
		return err
	}

  if year > 0 {
    err = WriteCategoryButtons(w, month, year)
    if err != nil {
      return err
    }
  }
	return nil
}
