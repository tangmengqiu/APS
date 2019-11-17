package src

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

func CheckCommit() {
	for {

		for _, p := range PersonPipe {
			fmt.Println("i am checking!")
			if err := p.GetCommitOfToday(); err != nil {
				logrus.Info(err.Error())
			}
		}
		time.Sleep(2 * time.Minute)
	}
}

func CheckAt24() {
	// check status at 24:00
}
