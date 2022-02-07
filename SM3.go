package main

import (
	"fmt"
	"github.com/tjfoc/gmsm/sm3"
)

func main12() {

	hash := sm3.New()

	hash.Write([]byte("darren1"))

	result := hash.Sum(nil)

	fmt.Printf("%x\n", result)
}

func main121() {

	result := sm3.Sm3Sum([]byte("darren"))
	fmt.Printf("%x\n", result)
}
