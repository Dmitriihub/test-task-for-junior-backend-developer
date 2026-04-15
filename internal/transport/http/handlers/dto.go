package handlers

import (
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

// ---- request DTOs -----------------------------------------------------------

type taskMutationDTO struct {
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Status      taskdomain.Status       `json:"status"`
	Recurrence  *recurrenceDTO          `json:"recurrence,omitempty"`
}

type recurrenceDTO struct {
	Type      taskdomain.RecurrenceType `json:"type"`
	Interval  int                       `json:"interval,omitempty"`
	MonthDays []int                     `json:"month_days,omitempty"`
	Dates     []time.Time               `json:"dates,omitempty"`
	StartDate time.Time                 `json:"start_date,omitempty"`
	EndDate   *time.Time                `json:"end_date,omitempty"`
}

func (d *recurrenceDTO) toDomain() *taskdomain.Recurrence {
	if d == nil {
		return nil
	}
	return &taskdomain.Recurrence{
		Type:      d.Type,
		Interval:  d.Interval,
		MonthDays: d.MonthDays,
		Dates:     d.Dates,
		StartDate: d.StartDate,
		EndDate:   d.EndDate,
	}
}

// ---- response DTOs ----------------------------------------------------------

type taskDTO struct {
	ID          int64             `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	Recurrence  *recurrenceDTO    `json:"recurrence,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	dto := taskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}

	if task.Recurrence != nil {
		dto.Recurrence = &recurrenceDTO{
			Type:      task.Recurrence.Type,
			Interval:  task.Recurrence.Interval,
			MonthDays: task.Recurrence.MonthDays,
			Dates:     task.Recurrence.Dates,
			StartDate: task.Recurrence.StartDate,
			EndDate:   task.Recurrence.EndDate,
		}
	}

	return dto
}

type occurrencesDTO struct {
	Occurrences []time.Time `json:"occurrences"`
}
