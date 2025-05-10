package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"instawall/internal/delivery/httpserver"
	"instawall/internal/repository/postgresql"
)

func main() {

	httpserver.Serve()

	postgresRepo := postgresql.New()

	result, err := postgresRepo.IsPhoneNumberExist("09147786264")
	fmt.Println(result, err)

}
