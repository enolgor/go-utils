package main

import (
	"fmt"

	"github.com/enolgor/go-utils/examples"
)

type Key string

func main() {
	fmt.Printf("%+v\n", examples.PORT)
	fmt.Printf("%+v\n", examples.HOST)
	fmt.Printf("%+v\n", examples.TIMEZONE.String())
	fmt.Printf("%+v\n", examples.TEST)
	fmt.Printf("%+v\n", examples.LANG)
}
