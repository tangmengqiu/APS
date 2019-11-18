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

func CheckAt24(f func()) {
	// check status at 24:00
	go func(){
		for{
			now:=time.Now()
			next:=now.Add(time.Hour*24)
			next=time.Date(next.Year(),next.Month(),next.Day(),0,0,0,0,next.Location())
			t:=time.NewTimer(next.Sub(now))
			<-t.C
			f()
		}
	}()
}

func Check(){
	for _, p := range PersonPipe {
		p.CheckStatusAt24Clock()
	}
}
