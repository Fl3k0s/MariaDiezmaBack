package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/mailer"
	"mariadiezmaback/internal/repository"
	"mariadiezmaback/pkg/validator"
)

type AppointmentService struct {
	reqRepo           repository.RequestRepository
	mailer            mailer.Mailer
	notificationEmail string
	logger            *slog.Logger
}

func NewAppointmentService(
	reqRepo repository.RequestRepository,
	mailService mailer.Mailer,
	notificationEmail string,
	logger *slog.Logger,
) *AppointmentService {
	return &AppointmentService{
		reqRepo:           reqRepo,
		mailer:            mailService,
		notificationEmail: notificationEmail,
		logger:            logger,
	}
}

func (s *AppointmentService) Create(ctx context.Context, input domain.CreateAppointmentInput) (*domain.Appointment, error) {
	v := validator.New()
	v.Required(input.Name, "nombre")
	v.Required(input.Email, "email")
	v.Email(input.Email, "email")
	v.Required(input.Phone, "telefono")
	v.Required(input.Date, "fecha")
	v.Required(input.TimeSlot, "tramo_horario")
	v.Required(input.Type, "tipo_cita")

	if !v.Valid() {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, v.Errors)
	}

	id := uuid.NewString()
	now := time.Now().UTC()

	// 1. Persist as backoffice RequestItem so it can be managed from the dashboard
	reqItem := &domain.RequestItem{
		ID:          id,
		Type:        "appointment",
		Status:      domain.StatusPending,
		Priority:    domain.PriorityHigh,
		SenderName:  input.Name,
		SenderEmail: input.Email,
		SenderPhone: input.Phone,
		Subject:     fmt.Sprintf("Solicitud de Cita: %s - %s", input.Type, input.Date),
		Message: fmt.Sprintf(
			"Nueva cita solicitada para el día %s en el tramo horario %s.\nTipo de cita: %s.\nContacto: %s (%s).",
			input.Date, input.TimeSlot, input.Type, input.Name, input.Phone,
		),
		Metadata: map[string]any{
			"fecha":         input.Date,
			"tramo_horario": input.TimeSlot,
			"tipo_cita":     input.Type,
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.reqRepo.Create(ctx, reqItem); err != nil {
		return nil, fmt.Errorf("failed to save appointment request: %w", err)
	}

	appt := &domain.Appointment{
		ID:        id,
		Name:      input.Name,
		Email:     input.Email,
		Phone:     input.Phone,
		Date:      input.Date,
		TimeSlot:  input.TimeSlot,
		Type:      input.Type,
		Status:    string(domain.StatusPending),
		CreatedAt: now,
	}

	// 2. Dispatch email notification to team/management
	go func(notificationAppt *domain.Appointment) {
		mailCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := s.mailer.SendAppointmentNotification(mailCtx, s.notificationEmail, notificationAppt); err != nil {
			s.logger.Error("Failed to send appointment notification email",
				slog.String("appointment_id", notificationAppt.ID),
				slog.String("to", s.notificationEmail),
				slog.Any("error", err),
			)
		} else {
			s.logger.Info("Appointment notification email sent successfully",
				slog.String("appointment_id", notificationAppt.ID),
				slog.String("to", s.notificationEmail),
			)
		}
	}(appt)

	return appt, nil
}
