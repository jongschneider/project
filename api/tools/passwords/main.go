package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jongschneider/youtube-project/api/internal/platform/encryption"
)

func main() {
	password := os.Args[1]

	hash, err := encryption.Encrypt(password)
	if err != nil {
		panic(err)
	}

	log.Println(fmt.Sprintf("\t%s  ---->  %s", password, hash))
}
