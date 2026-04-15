package task_test

import (
	"testing"
	"time"

	"example.com/taskservice/internal/domain/task"
)

// helpers
func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func ptr(t time.Time) *time.Time { return &t }

// ---- Validate ---------------------------------------------------------------

func TestRecurrence_Validate(t *testing.T) {
	tests := []struct {
		name    string
		rec     *task.Recurrence
		wantErr bool
	}{
		{
			name:    "nil recurrence is valid",
			rec:     nil,
			wantErr: false,
		},
		{
			name:    "daily with interval 1 is valid",
			rec:     &task.Recurrence{Type: task.RecurrenceDaily, Interval: 1},
			wantErr: false,
		},
		{
			name:    "daily with interval 0 is invalid",
			rec:     &task.Recurrence{Type: task.RecurrenceDaily, Interval: 0},
			wantErr: true,
		},
		{
			name:    "daily with negative interval is invalid",
			rec:     &task.Recurrence{Type: task.RecurrenceDaily, Interval: -1},
			wantErr: true,
		},
		{
			name:    "monthly with valid days is valid",
			rec:     &task.Recurrence{Type: task.RecurrenceMonthly, MonthDays: []int{1, 15, 30}},
			wantErr: false,
		},
		{
			name:    "monthly with empty days is invalid",
			rec:     &task.Recurrence{Type: task.RecurrenceMonthly, MonthDays: []int{}},
			wantErr: true,
		},
		{
			name:    "monthly with day 31 is invalid",
			rec:     &task.Recurrence{Type: task.RecurrenceMonthly, MonthDays: []int{31}},
			wantErr: true,
		},
		{
			name:    "monthly with day 0 is invalid",
			rec:     &task.Recurrence{Type: task.RecurrenceMonthly, MonthDays: []int{0}},
			wantErr: true,
		},
		{
			name:    "specific_dates with dates is valid",
			rec:     &task.Recurrence{Type: task.RecurrenceSpecificDates, Dates: []time.Time{date(2025, 1, 1)}},
			wantErr: false,
		},
		{
			name:    "specific_dates with empty dates is invalid",
			rec:     &task.Recurrence{Type: task.RecurrenceSpecificDates, Dates: []time.Time{}},
			wantErr: true,
		},
		{
			name:    "even_days is valid",
			rec:     &task.Recurrence{Type: task.RecurrenceEvenDays},
			wantErr: false,
		},
		{
			name:    "odd_days is valid",
			rec:     &task.Recurrence{Type: task.RecurrenceOddDays},
			wantErr: false,
		},
		{
			name:    "unknown type is invalid",
			rec:     &task.Recurrence{Type: "weekly"},
			wantErr: true,
		},
		{
			name: "end_date before start_date is invalid",
			rec: &task.Recurrence{
				Type:      task.RecurrenceEvenDays,
				StartDate: date(2025, 6, 1),
				EndDate:   ptr(date(2025, 1, 1)),
			},
			wantErr: true,
		},
		{
			name: "end_date equal to start_date is valid",
			rec: &task.Recurrence{
				Type:      task.RecurrenceEvenDays,
				StartDate: date(2025, 6, 2),
				EndDate:   ptr(date(2025, 6, 2)),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rec.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

// ---- OccursOn ---------------------------------------------------------------

func TestRecurrence_OccursOn_Daily(t *testing.T) {
	start := date(2025, 1, 1)

	rec := &task.Recurrence{
		Type:      task.RecurrenceDaily,
		Interval:  3,
		StartDate: start,
	}

	tests := []struct {
		date time.Time
		want bool
	}{
		{date(2025, 1, 1), true},  // день 0 — совпадает
		{date(2025, 1, 2), false}, // день 1
		{date(2025, 1, 3), false}, // день 2
		{date(2025, 1, 4), true},  // день 3 — совпадает
		{date(2025, 1, 7), true},  // день 6 — совпадает
		{date(2024, 12, 31), false}, // до start_date
	}

	for _, tt := range tests {
		got := rec.OccursOn(tt.date)
		if got != tt.want {
			t.Errorf("OccursOn(%s) = %v, want %v", tt.date.Format("2006-01-02"), got, tt.want)
		}
	}
}

func TestRecurrence_OccursOn_Monthly(t *testing.T) {
	rec := &task.Recurrence{
		Type:      task.RecurrenceMonthly,
		MonthDays: []int{1, 15},
		StartDate: date(2025, 1, 1),
	}

	tests := []struct {
		date time.Time
		want bool
	}{
		{date(2025, 1, 1), true},
		{date(2025, 1, 15), true},
		{date(2025, 1, 10), false},
		{date(2025, 2, 1), true},
		{date(2025, 2, 15), true},
		{date(2025, 2, 14), false},
	}

	for _, tt := range tests {
		got := rec.OccursOn(tt.date)
		if got != tt.want {
			t.Errorf("OccursOn(%s) = %v, want %v", tt.date.Format("2006-01-02"), got, tt.want)
		}
	}
}

func TestRecurrence_OccursOn_SpecificDates(t *testing.T) {
	rec := &task.Recurrence{
		Type:  task.RecurrenceSpecificDates,
		Dates: []time.Time{date(2025, 3, 10), date(2025, 4, 22)},
	}

	tests := []struct {
		date time.Time
		want bool
	}{
		{date(2025, 3, 10), true},
		{date(2025, 4, 22), true},
		{date(2025, 3, 11), false},
		{date(2025, 4, 21), false},
	}

	for _, tt := range tests {
		got := rec.OccursOn(tt.date)
		if got != tt.want {
			t.Errorf("OccursOn(%s) = %v, want %v", tt.date.Format("2006-01-02"), got, tt.want)
		}
	}
}

func TestRecurrence_OccursOn_EvenDays(t *testing.T) {
	rec := &task.Recurrence{
		Type:      task.RecurrenceEvenDays,
		StartDate: date(2025, 1, 1),
	}

	tests := []struct {
		date time.Time
		want bool
	}{
		{date(2025, 1, 2), true},
		{date(2025, 1, 4), true},
		{date(2025, 1, 30), true},
		{date(2025, 1, 1), false},
		{date(2025, 1, 3), false},
		{date(2025, 1, 15), false},
	}

	for _, tt := range tests {
		got := rec.OccursOn(tt.date)
		if got != tt.want {
			t.Errorf("OccursOn(%s) = %v, want %v", tt.date.Format("2006-01-02"), got, tt.want)
		}
	}
}

func TestRecurrence_OccursOn_OddDays(t *testing.T) {
	rec := &task.Recurrence{
		Type:      task.RecurrenceOddDays,
		StartDate: date(2025, 1, 1),
	}

	tests := []struct {
		date time.Time
		want bool
	}{
		{date(2025, 1, 1), true},
		{date(2025, 1, 3), true},
		{date(2025, 1, 15), true},
		{date(2025, 1, 2), false},
		{date(2025, 1, 4), false},
		{date(2025, 1, 30), false},
	}

	for _, tt := range tests {
		got := rec.OccursOn(tt.date)
		if got != tt.want {
			t.Errorf("OccursOn(%s) = %v, want %v", tt.date.Format("2006-01-02"), got, tt.want)
		}
	}
}

func TestRecurrence_OccursOn_EndDate(t *testing.T) {
	end := date(2025, 1, 10)
	rec := &task.Recurrence{
		Type:      task.RecurrenceEvenDays,
		StartDate: date(2025, 1, 1),
		EndDate:   &end,
	}

	tests := []struct {
		date time.Time
		want bool
	}{
		{date(2025, 1, 2), true},
		{date(2025, 1, 10), true},  // граница включительно
		{date(2025, 1, 12), false}, // после end_date
	}

	for _, tt := range tests {
		got := rec.OccursOn(tt.date)
		if got != tt.want {
			t.Errorf("OccursOn(%s) = %v, want %v", tt.date.Format("2006-01-02"), got, tt.want)
		}
	}
}

// ---- NextOccurrences --------------------------------------------------------

func TestRecurrence_NextOccurrences(t *testing.T) {
	rec := &task.Recurrence{
		Type:      task.RecurrenceMonthly,
		MonthDays: []int{1},
		StartDate: date(2025, 1, 1),
	}

	from := date(2025, 1, 1)
	results := rec.NextOccurrences(from, 3)

	if len(results) != 3 {
		t.Fatalf("expected 3 occurrences, got %d", len(results))
	}
	if !results[0].Equal(date(2025, 1, 1)) {
		t.Errorf("first occurrence: got %s, want 2025-01-01", results[0].Format("2006-01-02"))
	}
	if !results[1].Equal(date(2025, 2, 1)) {
		t.Errorf("second occurrence: got %s, want 2025-02-01", results[1].Format("2006-01-02"))
	}
	if !results[2].Equal(date(2025, 3, 1)) {
		t.Errorf("third occurrence: got %s, want 2025-03-01", results[2].Format("2006-01-02"))
	}
}

func TestRecurrence_NextOccurrences_NilReturnsNil(t *testing.T) {
	var rec *task.Recurrence
	results := rec.NextOccurrences(time.Now(), 5)
	if results != nil {
		t.Errorf("expected nil, got %v", results)
	}
}

func TestRecurrence_NextOccurrences_RespectsEndDate(t *testing.T) {
	end := date(2025, 1, 6)
	rec := &task.Recurrence{
		Type:      task.RecurrenceEvenDays,
		StartDate: date(2025, 1, 1),
		EndDate:   &end,
	}

	results := rec.NextOccurrences(date(2025, 1, 1), 10)
	// Чётные числа до 6-го включительно: 2, 4, 6
	if len(results) != 3 {
		t.Errorf("expected 3 occurrences, got %d: %v", len(results), results)
	}
}
