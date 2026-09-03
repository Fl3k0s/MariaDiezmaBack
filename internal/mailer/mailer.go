package mailer

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/smtp"
	"strings"
	"time"

	"mariadiezmaback/internal/config"
	"mariadiezmaback/internal/domain"
)

type Mailer interface {
	SendAppointmentNotification(ctx context.Context, toEmail string, appt *domain.Appointment) error
}

func New(cfg *config.Config, logger *slog.Logger) Mailer {
	if cfg.SMTPHost == "" {
		logger.Warn("SMTP not configured (SMTP_HOST is empty). Using LogMailer (emails will be logged to console)")
		return &LogMailer{logger: logger}
	}

	return &SMTPMailer{
		host:     cfg.SMTPHost,
		port:     cfg.SMTPPort,
		username: cfg.SMTPUser,
		password: cfg.SMTPPassword,
		from:     cfg.SMTPFrom,
		logger:   logger,
	}
}

// LogMailer prints emails to the structured logger for dev/testing
type LogMailer struct {
	logger *slog.Logger
}

func (m *LogMailer) SendAppointmentNotification(ctx context.Context, toEmail string, appt *domain.Appointment) error {
	m.logger.Info("📧 [EMAIL NOTIFICATION DISPATCHED]",
		slog.String("to", toEmail),
		slog.String("subject", fmt.Sprintf("[Maria Diezma] Nueva Cita Solicitada: %s", appt.Name)),
		slog.String("nombre", appt.Name),
		slog.String("email", appt.Email),
		slog.String("telefono", appt.Phone),
		slog.String("fecha", appt.Date),
		slog.String("tramo_horario", appt.TimeSlot),
		slog.String("tipo_cita", appt.Type),
		slog.String("cita_id", appt.ID),
	)
	return nil
}

// SMTPMailer sends emails via standard SMTP server
type SMTPMailer struct {
	host     string
	port     int
	username string
	password string
	from     string
	logger   *slog.Logger
}

func (m *SMTPMailer) SendAppointmentNotification(ctx context.Context, toEmail string, appt *domain.Appointment) error {
	subject := fmt.Sprintf("[Maria Diezma] Nueva Cita Solicitada: %s", appt.Name)

	plainBody := fmt.Sprintf(`Hola,

Se ha recibido una nueva solicitud de cita a través de la web:

--------------------------------------------------
DATOS DE LA CITA SOLICITADA
--------------------------------------------------
- Nombre:         %s
- Email:          %s
- Teléfono:       %s
- Fecha:          %s
- Tramo horario:  %s
- Tipo de cita:   %s
- ID de registro: %s
--------------------------------------------------

Puedes gestionar esta cita directamente en el panel de Backoffice.
`, appt.Name, appt.Email, appt.Phone, appt.Date, appt.TimeSlot, appt.Type, appt.ID)

	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<style>
		body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background-color: #f8fafc; margin: 0; padding: 20px; color: #1e293b; }
		.container { max-width: 600px; margin: 0 auto; background: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 6px -1px rgba(0,0,0,0.1); border: 1px solid #e2e8f0; }
		.header { background: #0f172a; color: #ffffff; padding: 24px; text-align: center; }
		.header h1 { margin: 0; font-size: 20px; font-weight: 600; letter-spacing: 0.5px; }
		.content { padding: 32px 24px; }
		.card { background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 6px; padding: 20px; margin: 20px 0; }
		.row { display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #edf2f7; font-size: 15px; }
		.row:last-child { border-bottom: none; }
		.label { font-weight: 600; color: #64748b; }
		.value { color: #0f172a; font-weight: 500; }
		.footer { background: #f1f5f9; padding: 16px; text-align: center; font-size: 13px; color: #64748b; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>Maria Diezma - Nueva Solicitud de Cita</h1>
		</div>
		<div class="content">
			<p>Has recibido una nueva solicitud de cita desde el formulario de la web:</p>
			<div class="card">
				<div class="row"><span class="label">Nombre:</span><span class="value">%s</span></div>
				<div class="row"><span class="label">Email:</span><span class="value"><a href="mailto:%s">%s</a></span></div>
				<div class="row"><span class="label">Teléfono:</span><span class="value">%s</span></div>
				<div class="row"><span class="label">Fecha solicitada:</span><span class="value"><strong>%s</strong></span></div>
				<div class="row"><span class="label">Tramo horario:</span><span class="value"><strong>%s</strong></span></div>
				<div class="row"><span class="label">Tipo de cita:</span><span class="value">%s</span></div>
				<div class="row"><span class="label">ID de gestión:</span><span class="value"><code>%s</code></span></div>
			</div>
			<p style="font-size: 14px; color: #64748b;">Puedes acceder al Backoffice para revisar, gestionar o responder a esta cita.</p>
		</div>
		<div class="footer">
			Mensaje generado automáticamente por el sistema de MariaDiezmaBack.
		</div>
	</div>
</body>
</html>`, appt.Name, appt.Email, appt.Email, appt.Phone, appt.Date, appt.TimeSlot, appt.Type, appt.ID)

	boundary := fmt.Sprintf("boundary-%d", time.Now().UnixNano())
	msg := []byte(strings.Join([]string{
		fmt.Sprintf("From: %s", m.from),
		fmt.Sprintf("To: %s", toEmail),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"", boundary),
		"",
		fmt.Sprintf("--%s", boundary),
		"Content-Type: text/plain; charset=UTF-8",
		"Content-Transfer-Encoding: 8bit",
		"",
		plainBody,
		"",
		fmt.Sprintf("--%s", boundary),
		"Content-Type: text/html; charset=UTF-8",
		"Content-Transfer-Encoding: 8bit",
		"",
		htmlBody,
		"",
		fmt.Sprintf("--%s--", boundary),
	}, "\r\n"))

	addr := fmt.Sprintf("%s:%d", m.host, m.port)
	var auth smtp.Auth
	if m.username != "" {
		auth = smtp.PlainAuth("", m.username, m.password, m.host)
	}

	// Connect to SMTP server
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to dial SMTP server %s: %w", addr, err)
	}

	client, err := smtp.NewClient(conn, m.host)
	if err != nil {
		_ = conn.Close()
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// Try STARTTLS if available
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{
			ServerName: m.host,
		}
		if err = client.StartTLS(tlsConfig); err != nil {
			m.logger.Warn("STARTTLS failed, continuing with unencrypted connection", slog.Any("error", err))
		}
	}

	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP auth failed: %w", err)
		}
	}

	if err = client.Mail(m.from); err != nil {
		return fmt.Errorf("SMTP MAIL command failed: %w", err)
	}

	if err = client.Rcpt(toEmail); err != nil {
		return fmt.Errorf("SMTP RCPT command failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP DATA command failed: %w", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("SMTP write body failed: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("SMTP close writer failed: %w", err)
	}

	return client.Quit()
}
