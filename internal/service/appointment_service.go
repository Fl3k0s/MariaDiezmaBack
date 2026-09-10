package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
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
	v.Required(input.Name, "nombre_apellidos")
	v.Required(input.Phone, "telefono_contacto")
	v.Required(input.Date, "fecha")
	v.Required(input.TimeSlot, "franja_horaria")
	v.Required(input.Type, "tipo_cita")

	if input.Email != "" {
		v.Email(input.Email, "email")
	}

	if !v.Valid() {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidInput, v.Errors)
	}

	id := uuid.NewString()
	now := time.Now().UTC()

	// 1. Persist as backoffice RequestItem so it can be managed from the dashboard
	metadata := map[string]any{
		"tipo_cita":         input.Type,
		"fecha":             input.Date,
		"franja_horaria":    input.TimeSlot,
		"tramo_horario":     input.TimeSlot,
		"nombre_apellidos":  input.Name,
		"telefono_contacto": input.Phone,
	}
	if input.EstimatedDate != nil && *input.EstimatedDate != "" {
		metadata["fecha_estimada"] = *input.EstimatedDate
	}
	if input.Details != "" {
		metadata["detalles"] = input.Details
	}

	msgLines := []string{
		fmt.Sprintf("Nueva cita solicitada para el día %s en el tramo horario %s.", input.Date, input.TimeSlot),
		fmt.Sprintf("Tipo de cita: %s.", input.Type),
		fmt.Sprintf("Contacto: %s (%s).", input.Name, input.Phone),
	}
	if input.Email != "" {
		msgLines = append(msgLines, fmt.Sprintf("Email: %s.", input.Email))
	}
	if input.EstimatedDate != nil && *input.EstimatedDate != "" {
		msgLines = append(msgLines, fmt.Sprintf("Fecha estimada: %s.", *input.EstimatedDate))
	}
	if input.Details != "" {
		msgLines = append(msgLines, fmt.Sprintf("Detalles: %s", input.Details))
	}

	reqItem := &domain.RequestItem{
		ID:          id,
		Type:        "appointment",
		Status:      domain.StatusPending,
		Priority:    domain.PriorityHigh,
		SenderName:  input.Name,
		SenderEmail: input.Email,
		SenderPhone: input.Phone,
		Subject:     fmt.Sprintf("Solicitud de Cita: %s - %s", input.Type, input.Date),
		Message:     strings.Join(msgLines, "\n"),
		Metadata:    metadata,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.reqRepo.Create(ctx, reqItem); err != nil {
		return nil, fmt.Errorf("failed to save appointment request: %w", err)
	}

	appt := &domain.Appointment{
		ID:            id,
		Name:          input.Name,
		Email:         input.Email,
		Phone:         input.Phone,
		Date:          input.Date,
		TimeSlot:      input.TimeSlot,
		Type:          input.Type,
		EstimatedDate: input.EstimatedDate,
		Details:       input.Details,
		Status:        string(domain.StatusPending),
		CreatedAt:     now,
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
