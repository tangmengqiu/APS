package src

import (
	vm "APS/src/api/vm"
	httpclient "APS/tools/httpclient"
	"fmt"
)

type Person struct {
	ID              int
	Name            string
	Token           string
	Repo            string
	CommitToday     int
	CommitTotal     int
	DelayNum        int
	ContinuesDayNum int
}

var PersonPipe []Person

func GetUsers() []Person {
	return PersonPipe
}

func NewPerson(_name, _token, _repo string) Person {

	p := Person{
		ID:    len(PersonPipe),
		Name:  _name,
		Token: _token,
		Repo:  _repo,
	}
	PersonPipe = append(PersonPipe, p)
	return p
}
func AddUser(u vm.ReqUser) error {
	return nil
}

func (p Person) GetCommitOfToday() {
	resp, err := httpclient.HTTPGet(p.Repo, p.Token)
	if err != nil {
		log.Error(err.Error())
		return
	}
	fmt.Println(resp)
	//unmarshall the results
}

func (p Person) CheckStatusAt24Clock() {
	if p.CommitToday == 0 {
		//ding ding push
		p.DelayNum++
		p.ContinuesDayNum = 0
	} else {
		// ding ding push
		p.ContinuesDayNum++
	}
}
