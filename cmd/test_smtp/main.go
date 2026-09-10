package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"mariadiezmaback/internal/config"
	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/mailer"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	fmt.Printf("Config loaded:\n")
	fmt.Printf("SMTP_HOST: %s\n", cfg.SMTPHost)
	fmt.Printf("SMTP_PORT: %d\n", cfg.SMTPPort)
	fmt.Printf("SMTP_USER: %s\n", cfg.SMTPUser)
	fmt.Printf("SMTP_FROM: %s\n", cfg.SMTPFrom)
	fmt.Printf("NOTIFICATION_EMAIL: %s\n", cfg.NotificationEmail)
	hasPassword := cfg.SMTPPassword != ""
	fmt.Printf("SMTP_PASSWORD set: %v\n", hasPassword)

	m := mailer.New(cfg, logger)

	testAppt := &domain.Appointment{
		ID:        "test-smtp-001",
		Name:      "Prueba Sistema",
		Email:     "prueba@mariadiezma.es",
		Phone:     "+34 600 000 000",
		Date:      "2026-10-10",
		TimeSlot:  "12:00 - 13:00",
		Type:      "Prueba de Envío",
		Details:   "Mensaje de prueba para verificar SMTP de Google",
	}

	ctx := context.Background()
	err := m.SendAppointmentNotification(ctx, cfg.NotificationEmail, testAppt)
	if err != nil {
		fmt.Printf("Error sending email: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Success! Email sent successfully.")
}
