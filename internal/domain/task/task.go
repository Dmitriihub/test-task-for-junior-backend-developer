package task

import "time"

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

// RecurrenceType — тип периодичности задачи.
type RecurrenceType string

const (
	RecurrenceDaily         RecurrenceType = "daily"          // каждые N дней
	RecurrenceMonthly       RecurrenceType = "monthly"        // в конкретные числа месяца
	RecurrenceSpecificDates RecurrenceType = "specific_dates" // только указанные даты
	RecurrenceEvenDays      RecurrenceType = "even_days"      // чётные числа месяца
	RecurrenceOddDays       RecurrenceType = "odd_days"       // нечётные числа месяца
)

// Recurrence — настройки периодичности задачи.
//
// Заполнять нужно только поля, соответствующие выбранному Type:
//
//	daily          → Interval (>= 1, каждые N дней)
//	monthly        → MonthDays (числа от 1 до 30)
//	specific_dates → Dates (список конкретных дат)
//	even_days      → (доп. полей нет)
//	odd_days       → (доп. полей нет)
//
// StartDate и EndDate применяются ко всем типам.
type Recurrence struct {
	Type      RecurrenceType `json:"type"`
	Interval  int            `json:"interval,omitempty"`   // только для daily
	MonthDays []int          `json:"month_days,omitempty"` // только для monthly
	Dates     []time.Time    `json:"dates,omitempty"`      // только для specific_dates
	StartDate time.Time      `json:"start_date,omitempty"` // начало, включительно
	EndDate   *time.Time     `json:"end_date,omitempty"`   // конец, включительно; nil = бессрочно
}

// OccursOn сообщает, должна ли задача быть активна в указанную дату.
func (rec *Recurrence) OccursOn(date time.Time) bool {
	if rec == nil {
		return false
	}

	date = date.UTC().Truncate(24 * time.Hour)

	if !rec.StartDate.IsZero() {
		if date.Before(rec.StartDate.UTC().Truncate(24 * time.Hour)) {
			return false
		}
	}
	if rec.EndDate != nil {
		if date.After(rec.EndDate.UTC().Truncate(24 * time.Hour)) {
			return false
		}
	}

	day := date.Day()

	switch rec.Type {
	case RecurrenceDaily:
		interval := rec.Interval
		if interval < 1 {
			interval = 1
		}
		start := rec.StartDate.UTC().Truncate(24 * time.Hour)
		diff := int(date.Sub(start).Hours()) / 24
		return diff%interval == 0

	case RecurrenceMonthly:
		for _, d := range rec.MonthDays {
			if d == day {
				return true
			}
		}
		return false

	case RecurrenceSpecificDates:
		for _, d := range rec.Dates {
			if d.UTC().Truncate(24*time.Hour).Equal(date) {
				return true
			}
		}
		return false

	case RecurrenceEvenDays:
		return day%2 == 0

	case RecurrenceOddDays:
		return day%2 != 0
	}

	return false
}

// NextOccurrences возвращает до n ближайших дат (начиная с from),
// когда задача активна по расписанию. Горизонт поиска — 365 дней.
func (rec *Recurrence) NextOccurrences(from time.Time, n int) []time.Time {
	if rec == nil || n <= 0 {
		return nil
	}

	var results []time.Time
	current := from.UTC().Truncate(24 * time.Hour)
	horizon := current.AddDate(0, 0, 365)

	if rec.EndDate != nil {
		end := rec.EndDate.UTC().Truncate(24 * time.Hour)
		if end.Before(horizon) {
			horizon = end
		}
	}

	for !current.After(horizon) && len(results) < n {
		if rec.OccursOn(current) {
			results = append(results, current)
		}
		current = current.AddDate(0, 0, 1)
	}
	return results
}

// Task — основная доменная модель задачи.
type Task struct {
	ID          int64       `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Status      Status      `json:"status"`
	Recurrence  *Recurrence `json:"recurrence,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusNew, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}
