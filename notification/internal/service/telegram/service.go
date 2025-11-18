package telegram

import (
	"bytes"
	"context"
	"embed"
	"text/template"

	"go.uber.org/zap"

	"github.com/max-kriv0s/go-microservices-edu/notification/internal/client/http"
	"github.com/max-kriv0s/go-microservices-edu/notification/internal/config"
	"github.com/max-kriv0s/go-microservices-edu/notification/internal/model"
	"github.com/max-kriv0s/go-microservices-edu/platform/pkg/logger"
)

//go:embed templates/assembled_notification.tmpl
var assembledTemplateFS embed.FS

var assembledTemplate = template.Must(template.ParseFS(assembledTemplateFS, "templates/assembled_notification.tmpl"))

//go:embed templates/paid_notification.tmpl
var paidTemplateFS embed.FS

var paidTemplate = template.Must(template.ParseFS(paidTemplateFS, "templates/paid_notification.tmpl"))

type assembledTemlateData struct {
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64
}

type paidTemplateData struct {
	OrderUUID     string
	UserUUID      string
	PaymentMethod string
}

type service struct {
	telegramClient http.TelegramClient
	chatID         int64
}

func NewService(telegramClient http.TelegramClient) *service {
	return &service{
		telegramClient: telegramClient,
		chatID:         config.AppConfig().TelegramBot.ChatID(),
	}
}

func (s *service) SendOrderPaidNotification(ctx context.Context, event model.OrderPaidEvent) error {
	message, err := s.buildOrderPaidMessage(event)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, s.chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int64("chat_id", s.chatID), zap.String("message", message))
	return nil
}

func (s *service) buildOrderPaidMessage(event model.OrderPaidEvent) (string, error) {
	data := paidTemplateData{
		OrderUUID:     event.OrderUUID,
		UserUUID:      event.UserUUID,
		PaymentMethod: event.PaymentMethod,
	}

	var buf bytes.Buffer
	err := paidTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *service) SendOrderAssembledNotification(ctx context.Context, event model.ShipAssembledEvent) error {
	message, err := s.buildOrderAssembledMessage(event)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, s.chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int64("chat_id", s.chatID), zap.String("message", message))
	return nil
}

func (s *service) buildOrderAssembledMessage(event model.ShipAssembledEvent) (string, error) {
	data := assembledTemlateData{
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	var buf bytes.Buffer
	err := assembledTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
