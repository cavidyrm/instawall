package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"instawall/internal/repository/postgresql"
)

func main() {

	postgresRepo := postgresql.New()

	result, err := postgresRepo.IsPhoneNumberExist("09147786264")
	fmt.Println(result, err)

}
