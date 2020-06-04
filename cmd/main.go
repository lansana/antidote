package main

import (
	"fmt"
	"log"

	"github.com/lansana/antidote"
)

func main() {
	a := antidote.New()
	a.Mix(&antidote.Ingredients{URL: "https://fujifish.com"})

	html, err := a.Cure()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(html)
}
