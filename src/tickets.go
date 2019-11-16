package src

import (
	"fmt"
	"time"
)

func CheckCommit() {
	for {
		fmt.Println("i am checking!")
		time.Sleep(5 * time.Second)
	}
}

func CheckAt24() {
	// check status at 24:00
}
