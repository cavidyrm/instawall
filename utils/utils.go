package utils

func IsPhoneNumberValid(number string) bool {
	if len(number) != 11 || number[:2] != "09" {
		return false
	}
	return true
}
