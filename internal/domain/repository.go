package domain

type UserRepository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
	UpdatePassword(userID, password string) error
}

type OTPRepository interface {
	Create(otp *OTP) error
	GetByCode(code string) (*OTP, error)
	Delete(id string) error
}