package pkg

import "time"

var DefaultPeriodOffset = 1

func GetPeriod(timestamp int64, dayOffset int) (int64, int64) {
  date := time.Unix(timestamp, 0)
  firstOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
  lastOfMonth := firstOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)
  return firstOfMonth.AddDate(0, 0, dayOffset).Unix(), lastOfMonth.AddDate(0, 0, dayOffset).Unix()
}
