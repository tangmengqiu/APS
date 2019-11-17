package ding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Req struct {
	Secret     string `json:"secret"`
	AppKey     string `json:"app_key"`
	TemplateId string `json:"template_id"`
	Url        string `json:"url"`
	Data       data   `json:"data"`
}

func (req *Req) MakeMessage(_name, _url, _type string, ct, cl, cd int) *Req {
	req.Url = _url
	req.Data.Title = _type
	req.Data.Name = _name
	req.Data.CommitToday = ct
	req.Data.CommitTotoal = cl
	req.Data.ContinuesDayNum = cd
	return req
}

type data struct {
	Title           string `json:"title"`
	Name            string `json:"name"`
	CommitToday     int    `json:"commit_today"`
	CommitTotoal    int    `json:"commit_total"`
	ContinuesDayNum int    `json:"continues_day_num`
}

type MarkDownMsg struct {
	MsgType   string   `json:"msgtype"`
	MarkDownT MarkDown `json:"markdown"`
	At        At       `json:"at"`
}
type MarkDown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

func (req *Req) DingDing() {
	var md MarkDownMsg
	md.MsgType = "markdown"
	md.MarkDownT.Title = req.Data.Title
	if req.Data.CommitTotoal == 0 {
		//fist time ,welcome new users
		md.MarkDownT.Text = fmt.Sprintf("%s%s%s", "欢迎: ", req.Data.Name, " 选手")
	} else {
		md.MarkDownT.Text = fmt.Sprintf("%s%s%s%s",
			req.Data.Name+" 选手",
			"\n今日刷题数: "+strconv.Itoa(req.Data.CommitToday),
			"\n累计刷题数: "+strconv.Itoa(req.Data.CommitTotoal),
			"\n连续刷题天数: "+strconv.Itoa(req.Data.ContinuesDayNum))
	}

	fmt.Print(md.MarkDownT.Text)
	mdByte, err := json.Marshal(&md)
	if err != nil {
		logrus.Error("dingding marshal err", err)
		return
	}
	fmt.Println(string(mdByte))
	dingURL := req.Url
	dingRequest, err := http.NewRequest("POST",
		dingURL,
		bytes.NewReader(mdByte))
	//dingRequest, err := http.NewRequest("POST", "https://oapi.dingtalk.com/robot/send?access_token=ab6893afba86b066067cb898d0f5df44ccc56395e02887ac12b159acbb6a74c5", bytes.NewReader(mdByte))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"event": "new ding talk requst",
		}).Error(err)
	}
	dingRequest.Header.Set("Content-Type", "application/json")
	dingHTTPClient := http.Client{}
	dingResp, err := dingHTTPClient.Do(dingRequest)
	_, err = ioutil.ReadAll(dingResp.Body)
	if err != nil {
		logrus.WithField("event", "push ding talk").Error(err)
	}
	//logrus.Info(string(readAll))
	return
}
