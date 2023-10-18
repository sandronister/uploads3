package main

import (
	"fmt"
	"os"
)

func main() {
	qtd := 5000

	for i := 0; i < qtd; i++ {
		f, err := os.Create(fmt.Sprintf("./tmp/file%d.txt", i))
		if err != nil {
			fmt.Printf("Ocorreu um error %d", err)
		}
		defer f.Close()
		f.WriteString("Vai que é tua")
	}
}
