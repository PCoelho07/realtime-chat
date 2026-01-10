package client

import "fmt"

func Present(res <-chan string) {
	for e := range res {
		fmt.Print(e)
	}
}
