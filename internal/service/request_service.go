package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/repository"
	"mariadiezmaback/pkg/validator"
)

type RequestService struct {
	repo repository.RequestRepository
}

func NewRequestService(repo repository.RequestRepository) *RequestService {
	return &RequestService{repo: repo}
}

func (s *RequestService) Create(ctx context.Context, input domain.CreateRequestInput) (*domain.RequestItem, error) {
	v := validator.New()
	v.Required(input.SenderName, "sender_name")
	v.Required(input.SenderEmail, "sender_email")
	v.Email(input.SenderEmail, "sender_email")
	v.Required(input.Subject, "subject")
	v.Required(input.Message, "message")

	if input.Type == "" {
		input.Type = "general"
	}

	if input.Priority == "" {
		input.Priority = domain.PriorityMedium
	} else {
		v.In(string(input.Priority), []string{
			string(domain.PriorityLow),
			string(domain.PriorityMedium),
			string(domain.PriorityHigh),
			string(domain.PriorityUrgent),
		}, "priority")
	}

	if !v.Valid() {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, v.Errors)
	}

	now := time.Now().UTC()
	item := &domain.RequestItem{
		ID:          uuid.NewString(),
		Type:        input.Type,
		Status:      domain.StatusPending,
		Priority:    input.Priority,
		SenderName:  input.SenderName,
		SenderEmail: input.SenderEmail,
		SenderPhone: input.SenderPhone,
		Subject:     input.Subject,
		Message:     input.Message,
		Metadata:    input.Metadata,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *RequestService) GetByID(ctx context.Context, id string) (*domain.RequestItem, error) {
	if id == "" {
		return nil, domain.ErrInvalidInput
	}
	return s.repo.GetByID(ctx, id)
}

func (s *RequestService) List(ctx context.Context, filter domain.RequestFilter) ([]domain.RequestItem, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 || filter.PerPage > 100 {
		filter.PerPage = 20
	}
	return s.repo.List(ctx, filter)
}

func (s *RequestService) UpdateStatus(ctx context.Context, id string, input domain.UpdateStatusInput) (*domain.RequestItem, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	v := validator.New()
	v.In(string(input.Status), []string{
		string(domain.StatusPending),
		string(domain.StatusInProgress),
		string(domain.StatusResolved),
		string(domain.StatusCancelled),
	}, "status")

	if !v.Valid() {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, v.Errors)
	}

	item.Status = input.Status
	if input.InternalNote != "" {
		item.InternalNote = input.InternalNote
	}
	item.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *RequestService) Update(ctx context.Context, id string, input domain.UpdateRequestInput) (*domain.RequestItem, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Type != nil && *input.Type != "" {
		item.Type = *input.Type
	}
	if input.Status != nil {
		item.Status = *input.Status
	}
	if input.Priority != nil {
		item.Priority = *input.Priority
	}
	if input.Subject != nil && *input.Subject != "" {
		item.Subject = *input.Subject
	}
	if input.Message != nil && *input.Message != "" {
		item.Message = *input.Message
	}
	if input.InternalNote != nil {
		item.InternalNote = *input.InternalNote
	}
	if input.AssignedTo != nil {
		item.AssignedTo = input.AssignedTo
	}
	item.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *RequestService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrInvalidInput
	}
	return s.repo.Delete(ctx, id)
}
