package job

import (
	"strings"
	"os/exec"
	"monitor/global"
	"errors"
	"strconv"
	"golang.org/x/crypto/ssh"
	"fmt"
	"bytes"
)

type Job struct {
	id           int32
	serviceName  string
	sshUser      string
	sshIP        string
	sshPort      string
	startCmd     string
	stopCmd      string
	autoRecover  bool
	mailReceiver string

	rely   []int32
	local  []*CheckCmd
	remote []*CheckCmd
}

func NewJob(id int32, serviceName, sshUser, sshIP, sshPort, startCmd, stopCmd string, autoRecover int, mailReceiver string) *Job {
	return &Job{
		id:           id,
		serviceName:  serviceName,
		sshUser:      sshUser,
		sshIP:        sshIP,
		sshPort:      sshPort,
		startCmd:     startCmd,
		stopCmd:      stopCmd,
		autoRecover:  autoRecover == 1,
		mailReceiver: mailReceiver,
		rely:         make([]int32, 0),
		local:        make([]*CheckCmd, 0),
		remote:       make([]*CheckCmd, 0),
	}
}

func (o *Job) Id() int32 {
	return o.id
}

func (o *Job) ServiceName() string {
	return o.serviceName
}

func (o *Job) IsAutoRecover() bool {
	return o.autoRecover
}

func (o *Job) MailReceiver() string {
	return o.mailReceiver
}

func (o *Job) Rely() []int32 {
	return o.rely
}

func (o *Job) AddRely(relyId int32) {
	o.rely = append(o.rely, relyId)
}

func (o *Job) AddLocalCmd(cmd *CheckCmd) {
	o.local = append(o.local, cmd)
}

func (o *Job) AddRemoteCmd(cmd *CheckCmd) {
	o.local = append(o.remote, cmd)
}


func (o *Job) CheckAll() (healthy bool, err error) {
	// local check
	for _, policy := range o.local {
		healthy, err = o.localCheck(policy)
		if err != nil {
			global.Logger.Error(err.Error())
			global.SendMail(o.mailReceiver, global.AD_TECH_MONITOR, err.Error())
			continue
		}

		if !healthy {
			return
		}
	}

	config := global.GetSSHConfig(o.sshUser)
	if config == nil {
		err = errors.New("get ssh config error")
		return
	}

	// build connection
	var client *ssh.Client
	client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", o.sshIP, o.sshPort), config)
	if err != nil {
		return
	}
	defer client.Close()

	var session *ssh.Session
	session, err = client.NewSession()
	if err != nil {
		return
	}
	defer session.Close()

	// remote check
	for _, policy := range o.remote {
		healthy, err = o.remoteCheck(session, policy)
		if err != nil {
			global.Logger.Error(err.Error())
			global.SendMail(o.mailReceiver, global.AD_TECH_MONITOR, err.Error())
			continue
		}

		if !healthy {
			return
		}
	}
	return
}

func (o *Job) localCheck(check *CheckCmd) (healthy bool, err error) {
	args := strings.Fields(check.checkShell)
	cmd := exec.Command(args[0], args[1:]...)
	var out []byte
	out, err = cmd.Output()
	if err != nil {
		global.Logger.Error(err.Error())
		return
	}

	healthy, err = o.isHealthy(string(out), check.operator, check.checkValue, check.goodMatch)
	return
}

func (o *Job) remoteCheck(session *ssh.Session, check *CheckCmd) (healthy bool, err error) {
	var out bytes.Buffer
	session.Stdout = &out
	if err = session.Run(check.checkShell); err != nil {
		global.Logger.Error(err.Error())
		return
	}

	healthy, err = o.isHealthy(out.String(), check.operator, check.checkValue, check.goodMatch)
	return
}

func (o *Job) isHealthy(left, operator, right string, goodMatch bool) (healthy bool, err error) {
	var match bool
	match, err = o.compare(left, operator, right)
	if err != nil {
		global.Logger.Error(err.Error())
		return
	}

	if goodMatch {
		healthy = match
	} else {
		healthy = !match
	}
	return
}

func (o *Job) compare(left, operator, right string) (match bool, err error) {
	var digitalLeft int
	var digitalRight int
	if _, ok := map[string]interface{}{
		global.LESS_THAN:     nil,
		global.LESS_EQUAL:    nil,
		global.EQUAL:         nil,
		global.GREATOR_EQUAL: nil,
		global.GREATOR_THAN:  nil,
	}[operator]; ok {
		digitalLeft, err = strconv.Atoi(left)
		if err != nil {
			global.Logger.Error(err.Error())
			return
		}

		digitalRight, err = strconv.Atoi(right)
		if err != nil {
			global.Logger.Error(err.Error())
			return
		}
	}

	switch operator {
	case global.LESS_THAN:
		match = digitalLeft < digitalRight
	case global.LESS_EQUAL:
		match = digitalLeft <= digitalRight
	case global.EQUAL:
		match = digitalLeft == digitalRight
	case global.GREATOR_EQUAL:
		match = digitalLeft >= digitalRight
	case global.GREATOR_THAN:
		match = digitalLeft > digitalRight
	case global.IN:
		match = strings.Contains(left, right)
	case global.EX:
		match = !strings.Contains(left, right)
	default:
		match = false
		err = errors.New("unknown operator: " + operator)
	}

	return
}

func (o *Job) stop() (err error) {
	config := global.GetSSHConfig(o.sshUser)
	if config == nil {
		err = errors.New("get ssh config error")
		return
	}

	// build connection
	var client *ssh.Client
	client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", o.sshIP, o.sshPort), config)
	if err != nil {
		return
	}
	defer client.Close()

	var session *ssh.Session
	session, err = client.NewSession()
	if err != nil {
		return
	}
	defer session.Close()

	var out bytes.Buffer
	session.Stdout = &out
	if err = session.Run(o.stopCmd); err != nil {
		global.Logger.Error(err.Error())
		return
	}

	return
}

func (o *Job) start() (err error) {
	config := global.GetSSHConfig(o.sshUser)
	if config == nil {
		err = errors.New("get ssh config error")
		return
	}

	// build connection
	var client *ssh.Client
	client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", o.sshIP, o.sshPort), config)
	if err != nil {
		return
	}
	defer client.Close()

	var session *ssh.Session
	session, err = client.NewSession()
	if err != nil {
		return
	}
	defer session.Close()

	var out bytes.Buffer
	session.Stdout = &out
	if err = session.Run(o.startCmd); err != nil {
		global.Logger.Error(err.Error())
		return
	}

	return
}

func (o *Job) Restart() (err error) {
	config := global.GetSSHConfig(o.sshUser)
	if config == nil {
		err = errors.New("get ssh config error")
		return
	}

	// build connection
	var client *ssh.Client
	client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%s", o.sshIP, o.sshPort), config)
	if err != nil {
		return
	}
	defer client.Close()

	var session *ssh.Session
	session, err = client.NewSession()
	if err != nil {
		return
	}
	defer session.Close()

	var out bytes.Buffer
	session.Stdout = &out
	if err = session.Run(o.stopCmd); err != nil {
		global.Logger.Error(err.Error())
		return
	}

	out.Reset()
	if err = session.Run(o.startCmd); err != nil {
		global.Logger.Error(err.Error())
		return
	}

	return
}
