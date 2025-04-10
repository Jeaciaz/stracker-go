package model

import (
	"database/sql"
	"fmt"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/stracker-go/pkg"
)

type Spending struct {
	Id          string
	Username    string
	Timestamp   int64
	Amount      float64
	Category    string
	Description string
}

type SpendingTimeline struct {
	PeriodStart int64
	Spendings   []*Spending
}

func CreateSpending(categoryId string, username string, amount float64, description string) error {
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("Failed to generate UUID: %w", err)
	}
	db := pkg.GetDb()
	_, err = db.Exec("INSERT INTO spendings (id, username, timestamp, amount, category, description) VALUES (?, ?, ?, ?, ?, ?)",
		id,
		username,
		time.Now().Unix(),
		amount,
		categoryId,
		description,
	)
	if err != nil {
		return fmt.Errorf("Failed to insert spending: %w", err)
	}
	return nil
}

func GetAllSpendings() ([]*Spending, error) {
  db := pkg.GetDb()
  rows, err := db.Query("SELECT * FROM spendings ORDER BY timestamp DESC")
  if err != nil {
    return nil, fmt.Errorf("Failed to load spendings from DB: %w", err)
  }
  var spendings []*Spending
  for rows.Next() {
    s, err := scanSpending(rows)
    if err != nil {
      return nil, fmt.Errorf("Failed to scan spending from SQL response: %w", err)
    }
    spendings = append(spendings, &s)
  }
  defer rows.Close()
  return spendings, nil
}

func GetSpendingById(id string) (*Spending, error) {
  db := pkg.GetDb()
  rows, err := db.Query("SELECT * FROM spendings WHERE id = ?", id)
  if err != nil {
    return nil, fmt.Errorf("Failed to load spending from DB: %w", err)
  }
  defer rows.Close()
  if !rows.Next() {
    return nil, fmt.Errorf("Spending not found")
  }
  s, err := scanSpending(rows)
  if err != nil {
    return nil, fmt.Errorf("Failed to scan spending from SQL response: %w", err)
  }
  return &s, nil
}

func UpdateSpending(id string, amount float64, description string) error {
  db := pkg.GetDb()
  _, err := db.Exec("UPDATE spendings SET amount = ?, description = ? WHERE id = ?", amount, description, id)
  if err != nil {
    return fmt.Errorf("Failed to update spending: %w", err)
  }
  return nil
}

func DeleteSpending (id string) error {
  db := pkg.GetDb()
  _, err := db.Exec("DELETE FROM spendings WHERE id = ?", id)
  if err != nil {
    return fmt.Errorf("Failed to delete spending: %w", err)
  }
  return nil
}

func GetSpendingsForCategory(categoryId string, month int, year int) ([]*Spending, error) {
	db := pkg.GetDb()
  if year == 0 {
    month = int(time.Now().Month())
    year = time.Now().Year()
  }
  timestamp := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Now().Location()).Unix()
	periodStart, periodEnd := pkg.GetPeriod(timestamp, pkg.DefaultPeriodOffset)
	rows, err := db.Query("SELECT * FROM spendings WHERE category = ? AND timestamp > ? AND timestamp < ?", categoryId, periodStart, periodEnd)
	if err != nil {
		return nil, fmt.Errorf("Failed to load spendings from DB: %w", err)
	}
  defer rows.Close()
	var spendings []*Spending
	for rows.Next() {
    s, err := scanSpending(rows)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan spending from SQL response: %w", err)
		}
		spendings = append(spendings, &s)
	}
	return spendings, nil
}

func GetSpendingsTimeline() ([]*SpendingTimeline, error) {
	db := pkg.GetDb()
	rows, err := db.Query("SELECT * FROM spendings")
	if err != nil {
		return nil, fmt.Errorf("Failed to load spendings from DB: %w", err)
	}
  defer rows.Close()

	timelinesMap := make(map[int64]*SpendingTimeline)
	for rows.Next() {
    s, err := scanSpending(rows)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan spending from SQL response: %w", err)
		}
		periodStart, _ := pkg.GetPeriod(s.Timestamp, pkg.DefaultPeriodOffset)
		tl, exists := timelinesMap[periodStart]
		if !exists {
			timelinesMap[periodStart] = &SpendingTimeline{PeriodStart: periodStart, Spendings: []*Spending{&s}}
		} else {
			tl.Spendings = append(tl.Spendings, &s)
		}
	}
	timelines := make([]*SpendingTimeline, 0, len(timelinesMap))
	for _, tl := range timelinesMap {
		timelines = append(timelines, tl)
	}
	slices.SortFunc(timelines, func(i, j *SpendingTimeline) int {
		if i.PeriodStart > j.PeriodStart {
			return -1
		}
		if i.PeriodStart < j.PeriodStart {
			return 1
		}
		return 0
	})

	return timelines, nil
}

func scanSpending(rows *sql.Rows) (Spending, error) {
		var s Spending
    err := rows.Scan(&s.Id, &s.Username, &s.Timestamp, &s.Amount, &s.Category, &s.Description)
    return s, err
}
