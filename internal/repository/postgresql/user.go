package postgresql

import (
	"fmt"
	"instawall/internal/domain"
)

func (d *DB) IsPhoneNumberExist(phone string) (bool, error) {
	result, err := d.db.Exec("SELECT FROM users WHERE phone = $1", phone)
	fmt.Println(result)
	return true, err
}

func (d DB) IsEmailExist(email string) (bool, error) {
	return false, nil
}

func (d DB) Register(u domain.User) (domain.User, error) {
	return domain.User{}, nil
}
