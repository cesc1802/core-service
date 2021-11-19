package consumer

import (
	"context"
	"github.com/cesc1802/core-service"
	"log"
)

func SendEmailOTP(sc core_service.Service) consumerJob {
	return consumerJob{
		Title: "SEND EMAIL OTP",
		Handler: func(ctx context.Context, message interface{}) error {
			log.Println("=================== perform send EMAIL otp to user =======================")
			return nil
		},
	}
}
