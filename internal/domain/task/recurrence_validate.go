package task

import "errors"

// Validate проверяет корректность настроек периодичности.
func (rec *Recurrence) Validate() error {
	if rec == nil {
		return nil
	}

	switch rec.Type {
	case RecurrenceDaily:
		if rec.Interval < 1 {
			return errors.New("recurrence.interval must be >= 1 for daily recurrence")
		}
	case RecurrenceMonthly:
		if len(rec.MonthDays) == 0 {
			return errors.New("recurrence.month_days must not be empty for monthly recurrence")
		}
		for _, d := range rec.MonthDays {
			if d < 1 || d > 30 {
				return errors.New("recurrence.month_days values must be between 1 and 30")
			}
		}
	case RecurrenceSpecificDates:
		if len(rec.Dates) == 0 {
			return errors.New("recurrence.dates must not be empty for specific_dates recurrence")
		}
	case RecurrenceEvenDays, RecurrenceOddDays:
		// дополнительных полей не требуется
	default:
		return errors.New("recurrence.type is invalid; allowed: daily, monthly, specific_dates, even_days, odd_days")
	}

	if rec.EndDate != nil && !rec.StartDate.IsZero() && rec.EndDate.Before(rec.StartDate) {
		return errors.New("recurrence.end_date must not be before start_date")
	}

	return nil
}
