package postgresql

import "instawall/internal/domain"

func (d DB) IsPhoneNumberExist(phone string) (bool, error) {
	return false, nil
}

func (d DB) IsEmailExist(email string) (bool, error) {
	return false, nil
}

func (d DB) Register(u domain.User) (domain.User, error) {
	return domain.User{}, nil
}
