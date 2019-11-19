package src

import (
	vm "APS/src/api/vm"
	"APS/tools"
	ding "APS/tools/ding"
	"APS/tools/storage"
	"context"
	"encoding/json"
	"errors"

	"github.com/google/go-github/v28/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type Person struct {
	ID              int
	UUID            string
	Name            string
	Token           string
	Repo            string
	CommitToday     int
	CommitTotal     int
	DelayNum        int
	ContinuesDayNum int
}

var PersonPipe []*Person

var MDataBase storage.Storage

func GetUsers() []*Person {
	return PersonPipe
}

func NewPerson(_name, _token, _repo string) Person {

	p := Person{
		ID:    len(PersonPipe),
		UUID:  tools.StringUUID(),
		Name:  _name,
		Token: _token,
		Repo:  _repo,
	}
	PersonPipe = append(PersonPipe, &p)
	return p
}

func AddUser(u vm.ReqUser) error {
	p := NewPerson(u.Name, u.Token, u.Repo)
	//put in storage
	if err := MDataBase.AddOrUpdate(p.UUID, p); err != nil {
		logrus.WithField("event", "add user").Error(err.Error())
		return err
	}
	return nil
}

func DeleteUser(name string) error {
	for idx, p := range PersonPipe {
		if p.Name == name {
			//delete
			PersonPipe = append(PersonPipe[:idx], PersonPipe[idx+1:]...)
			return nil
		}
	}
	return errors.New("no such user nameed: " + name)
}

func (p *Person) GetCommitOfToday() error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: p.Token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)
	// list all repositories for the authenticated user
	commits, _, err := client.Repositories.ListCommits(ctx, p.Name, p.Repo, nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"event": "Get github commits",
		}).Error(err.Error())
		return err
	}
	numOfCommits := len(commits)
	if p.CommitTotal == 0 {
		//first add this user
		//push welcome msg
		var req ding.Req
		req.MakeMessage(p.Name, GlobalConfig.DingUrl, "欢迎新加入的朋友", p.CommitToday, p.CommitTotal, p.ContinuesDayNum)
		req.DingDing()
		p.CommitTotal = numOfCommits
		if err := MDataBase.AddOrUpdate(p.UUID, p); err != nil {
			logrus.WithField("event", "add user").Error(err.Error())
			return err
		}
		return nil
	}
	//not first time
	numOfCommitsOfToday := numOfCommits - p.CommitTotal
	p.CommitTotal = numOfCommits
	if numOfCommitsOfToday != 0 {
		if p.CommitToday == 0 {
			p.CommitToday = numOfCommitsOfToday
		} else {
			p.CommitToday += numOfCommitsOfToday
		}
		if err := MDataBase.AddOrUpdate(p.UUID, p); err != nil {
			logrus.WithField("event", "add user").Error(err.Error())
			return err
		}
		var req ding.Req
		req.MakeMessage(p.Name, GlobalConfig.DingUrl, p.Name+" 有新提交了", p.CommitToday, p.CommitTotal, p.ContinuesDayNum)
		req.DingDing()
	}
	return nil
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
	if err := MDataBase.AddOrUpdate(p.UUID, p); err != nil {
		logrus.WithField("event", "add user").Error(err.Error())
	}
	var req ding.Req
	req.MakeMessage(p.Name, GlobalConfig.DingUrl, "每日监督: "+p.Name, p.CommitToday, p.CommitTotal, p.ContinuesDayNum)
	req.DingDing()
}

func SyncMemoryToUsers() {
	allUsers := MDataBase.GetAll()
	for _, dataByte := range allUsers {
		p := new(Person)
		err := json.Unmarshal(dataByte, p)
		if err != nil {
			logrus.Error("SyncMemoryToUsers Unmarshal err", err)
			continue
		}
		PersonPipe = append(PersonPipe, p)
	}
	logrus.Info("Sync Memory To Db Ok")
}
