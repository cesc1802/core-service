package consumer

import (
	"context"
	"github.com/cesc1802/core-service"
	"log"
)

func SendSMSOTP(sc core_service.Service) consumerJob {
	return consumerJob{
		Title: "SEND SMS OTP",
		Handler: func(ctx context.Context, message interface{}) error {
			log.Println("=================== perform send SMS otp to user =======================")
			return nil
		},
	}
}
