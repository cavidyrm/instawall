package otp

import (
	"fmt"
	"math/rand"
	"time"
)

type OTPService interface {
	Generate() string
	Send(email, code string) error
}

type otpService struct{}

func NewOTPService() OTPService {
	return &otpService{}
}

func (s *otpService) Generate() string {
	// Generate a 6-digit OTP
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (s *otpService) Send(email, code string) error {
	// In a real application, this would send an email with the OTP
	// For this example, we'll just print it to the console
	fmt.Printf("Sending OTP %s to %s\n", code, email)
	return nil
}