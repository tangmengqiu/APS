package src

import (
	vm "APS/src/api/vm"
	"APS/tools"
	ding "APS/tools/ding"
	"APS/tools/storage"
	"context"
	"encoding/json"
	"errors"
	"time"

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
	UpdateAt        string
	LastCommitSHA   string
}

var PersonPipe []*Person

var MDataBase storage.Storage

func GetUsers() []*Person {
	return PersonPipe
}

func NewPerson(_name, _token, _repo string) *Person {

	p := Person{
		ID:    len(PersonPipe),
		UUID:  tools.StringUUID(),
		Name:  _name,
		Token: _token,
		Repo:  _repo,
	}
	PersonPipe = append(PersonPipe, &p)
	return &p
}

func AddUser(u vm.ReqUser) error {
	p := NewPerson(u.Name, u.Token, u.Repo)

	if err := p.InitUserCommitInfo(); err != nil {
		return err
	}
	return nil
}

func DeleteUser(name string) error {
	for idx, p := range PersonPipe {
		if p.Name == name {
			//delete
			PersonPipe = append(PersonPipe[:idx], PersonPipe[idx+1:]...)
			if !MDataBase.Delete(p.UUID) {
				return errors.New("delete user in bolt failed")
			}
			return nil
		}
	}
	return errors.New("no such user nameed: " + name)
}

func (p *Person) InitUserCommitInfo() error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: p.Token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	var lop = github.ListOptions{
		Page:    1,
		PerPage: 500,
	}
	opt := github.CommitsListOptions{
		ListOptions: lop,
	}
	client := github.NewClient(tc)
	commits, _, err := client.Repositories.ListCommits(ctx, p.Name, p.Repo, &opt)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"event": "Get github commits",
		}).Error(err.Error())
		return err
	}
	numOfCommits := len(commits)
	if numOfCommits == 0 {
		return errors.New("禁止空仓库加入")
	}
	latestCommit := commits[0]
	p.CommitTotal = numOfCommits
	p.CommitToday = 0
	p.ContinuesDayNum = 0
	p.LastCommitSHA = *(latestCommit.SHA)

	var req ding.Req
	req.MakeMessage(p.Name, GlobalConfig.DingUrl, "欢迎新加入的朋友", p.CommitToday, p.CommitTotal, p.ContinuesDayNum)
	req.DingDing()
	p.UpdateAt = time.Now().Format("2006-01-02 15:04:05")
	if err := MDataBase.AddOrUpdate(p.UUID, p); err != nil {
		logrus.WithField("event", "add user").Error(err.Error())
		return err
	}
	return nil
}

func (p *Person) GetCommitOfToday() error {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: p.Token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	var lop = github.ListOptions{
		Page:    1,
		PerPage: 3,
	}
	opt := github.CommitsListOptions{
		ListOptions: lop,
	}
	client := github.NewClient(tc)
	commits, _, err := client.Repositories.ListCommits(ctx, p.Name, p.Repo, &opt)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"event": "Get github commits",
		}).Error(err.Error())
		return err
	}
	numOfCommits := len(commits)
	if numOfCommits == 0 {
		return nil
	}
	latestCommit := commits[0]

	// new commits
	if *latestCommit.SHA != p.LastCommitSHA {
		//新提交了
		cc, _, err := client.Repositories.CompareCommits(ctx, p.Name, p.Repo, p.LastCommitSHA, *latestCommit.SHA)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"event": "Get github commits compare",
			}).Error(err.Error())
			return err
		}
		numOfCommitsOfToday := *cc.TotalCommits
		if numOfCommitsOfToday != 0 {
			p.CommitToday += numOfCommitsOfToday
			p.LastCommitSHA = *latestCommit.SHA
			p.CommitTotal += numOfCommitsOfToday
			p.UpdateAt = time.Now().Format("2006-01-02 15:04:05")
			if err := MDataBase.AddOrUpdate(p.UUID, p); err != nil {
				logrus.WithField("event", "add user").Error(err.Error())
				return err
			}
			var req ding.Req
			req.MakeMessage(p.Name, GlobalConfig.DingUrl, p.Name+" 有新提交了", p.CommitToday, p.CommitTotal, p.ContinuesDayNum)
			req.DingDing()
		}

	}
	return nil
}

func (p *Person) CheckStatusAt24Clock() {
	if p.CommitToday == 0 {
		//ding ding push
		p.DelayNum++
		p.ContinuesDayNum = 0
	} else {
		// ding ding push
		p.ContinuesDayNum++
	}
	p.CommitToday = 0
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
	logrus.Info("Sync Memory From Db Ok")
}
