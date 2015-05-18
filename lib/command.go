package lib

import (
	"encoding/json"
	//	"errors"
	"fmt"
	"github.com/funny/link"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var (
	VFId          int32
	WorkingStatus = true
)

type User struct {
	UserType     string
	Id           string
	WorkTime     string
	CWrongCount  int
	PPutCount    int
	PLastPutTime int64
}

type ResMessage struct {
	Action   string `json:"action"`
	Seq      string `json:"seq"`
	StatInfo MM     `json:"stat_info"`
}

func Process(session *link.Session, req map[string]string) error {
	// panic 异常
	defer func() {
		if err := recover(); err != nil {
			ULogger.Error("Error: panic ", err)
			ULogger.Error(string(debug.Stack()))
		}
	}()

	if !WorkingStatus {
		ULogger.Warn("System has been closed")
		return nil
	}

	command, ok := req["action"]
	if !ok {
		ULogger.Error("client %s %s", session.Conn().RemoteAddr().String(), "bad request ,not found action")
		session.Close()
		return nil

	}
	if (command == "getfile" || command == "answer") && session.State == nil {
		ULogger.Error("client %s %s", session.Conn().RemoteAddr().String(), "c must login frist")
		session.Close()
		return nil
	}

	switch command {
	//p
	case "putfile":
		return putFile(session, req)
	case "reportanswer":
		return reportAnswer(session, req)
	//c
	case "getfile":
		return getFile(session, req)
	case "cstart":
		return cStart(session, req)
	case "answer":
		return answer(session, req)
	case "test001":
		return test001(session, req)
	case "mstart":
		return mStart(session, req)
	case "getstatinfo":
		return getStatInfo(session, req)
	case "conrol":
		return control(session, req)
	default:
		ULogger.Error("client %s %s", session.Conn().RemoteAddr().String(), "not support command")
		session.Close()
		//ULogger.Info("sssss")
	}
	return nil
}

///答题系统控制程序，ip 白名单验证
///
//{
//  "action": "conrol",
//  "event":"start",
//  "seq":"control_001"
//}
func control(session *link.Session, req map[string]string) error {
	remoteAddr := session.Conn().RemoteAddr().String()

	ip := strings.Split(remoteAddr, ":")[0]
	if strings.Index(TCConfig.IpWhiteList, ip) < 0 {
		ULogger.Error("bad request,not in whiteiplist")
		session.Close()
		return nil
	}

	event, _ := req["event"]
	seq := req["seq"]
	stat := map[string]string{"action": "res_control"}
	stat["now"] = strconv.Itoa(int(time.Now().Unix()))
	stat["seq"] = seq
	switch event {
	//开始系统
	case "start":
		fmt.Println("2222")
		WorkingStatus = true
		Exec(`insert into user_activities(user_id,active_time,active_type,user_type,other_info) values(?,now(),'start','system',?)`, session.Conn().LocalAddr().String(), "operation:"+session.Conn().RemoteAddr().String())
		by, _ := json.Marshal(stat)
		session.Send(link.Bytes(by))
		ULogger.Info("send to client %s %s %s", session.Conn().RemoteAddr().String(), "say:", string(by))

	//停止系统
	case "stop":
		//WorkingStatus = false
		Exec(`insert into user_activities(user_id,active_time,active_type,user_type,other_info) values(?,now(),'stop','system',?)`, session.Conn().LocalAddr().String(), "operation:"+session.Conn().RemoteAddr().String())
		by, _ := json.Marshal(stat)
		session.Send(link.Bytes(by))
		ULogger.Info("send to client %s %s %s", session.Conn().RemoteAddr().String(), "say:", string(by))

	case "stat":
		waitcount := QueueInstance.len()
		clientcount := len(VFMapInstance.c_sessions)
		pcount := len(VFMapInstance.p_sessions)
		questioncount := len(VFMapInstance.innerMap)
		stat["productioncount"] = strconv.Itoa(pcount)
		stat["questioncount"] = strconv.Itoa(questioncount)
		stat["waitcount"] = strconv.Itoa(waitcount)
		stat["clientcount"] = strconv.Itoa(clientcount)
		pstatInfo := VFMapInstance.pstatInfo()
		stat["pstatinfo"] = pstatInfo
		cstatInfo := VFMapInstance.cstatInfo()
		stat["cstatinfo"] = cstatInfo
		by, _ := json.Marshal(stat)
		session.Send(link.Bytes(by))
		ULogger.Info("send to client %s %s %s", session.Conn().RemoteAddr().String(), "say:", string(by))
	case "exit":
		session.Close()
	default:
		by, _ := json.Marshal(stat)
		session.Send(link.Bytes(by))
		ULogger.Info("send to client %s %s %s", session.Conn().RemoteAddr().String(), "say:", string(by))

	}

	return nil
}

func test001(session *link.Session, req map[string]string) error {

	data := "1234567890"
	ret := map[string]string{
		"action": "res_test001",
		"seq":    "s_test001",
		"data":   data,
	}
	by, _ := json.Marshal(ret)
	_ = by
	session.Send(link.Bytes(by))
	//session.Send(link.String("123456"))
	//panic("222222")
	return nil //errors.New("test001 error")

}

///======================================
///生产者 ，产生图验图片
func putFile(session *link.Session, req map[string]string) error {
	file := req["file"]
	seq := req["seq"]
	fileid := GetMd5String(file)
	fileAnswer := GetAnswer(fileid)
	if fileAnswer != "" {
		ULogger.Info("file %s %s", fileid, " have in database,direct return")
		//插入题目，状态为9
		vf := &VerifyObj{Id: getId(), P: session, C: nil, FileId: fileid, File: file, Status: 9, Result: "1", Seq: seq, PPutUnix: time.Now().Unix()}
		Exec(`insert into exam(file_id,file_hash,f_status,put_time,answer_result,answer) values(?,?,9,now(),1)`, vf.Id, vf.FileId, fileAnswer)
		//给生产端回应答
		ret := map[string]string{
			"action": "res_putfile",
			"seq":    seq,
			"id":     vf.Id,
			"answer": fileAnswer,
		}
		by, _ := json.Marshal(ret)
		session.Send(link.Bytes(by))
		goto A
	} else {
		vf := &VerifyObj{Id: getId(), P: session, C: nil, FileId: fileid, File: file, Status: 1, Result: "0", Seq: seq, PPutUnix: time.Now().Unix()}
		QueueInstance.Enqueue(vf)
		VFMapInstance.Put(vf)
		ULogger.Info("putfile enqueue %s", vf.String())
	}

A:
	{
		//记录p端的操作
		if session.State == nil {
			userId := session.Conn().RemoteAddr().String()
			user := &User{UserType: "P", Id: userId, WorkTime: time.Now().Format("2006-01-02 15:04:05")}
			session.State = user
			VFMapInstance.AddPSession(session)

			Exec(`insert into user_activities(user_id,active_time,active_type,user_type,other_info) values(?,now(),'begin','production',?)`, userId, userId)
		}
		tmpUser := session.State.(*User)
		tmpUser.PLastPutTime = time.Now().Unix()
		tmpUser.PPutCount++

		return nil
	}
}

func reportAnswer(session *link.Session, req map[string]string) error {
	id := req["id"]
	result := req["result"]

	vf := VFMapInstance.Get(id)
	if vf == nil {
		ULogger.Error("answer,verifyobj not found %s", vf)
		return nil
	}

	if vf.C != nil && !vf.C.IsClosed() {
		if result == "1" {
			vf.C.State.(*User).CWrongCount = 0
		} else {
			vf.C.State.(*User).CWrongCount = vf.C.State.(*User).CWrongCount + 1
			if vf.C.State.(*User).CWrongCount >= 20 {
				ULogger.Error("%s %s %s", vf.C.State.(*User).Id, vf.C.Conn().RemoteAddr().String(), " answer wrong to top limit")
				vf.C.Close()
			}
		}
	}

	vf.Status = 5
	vf.Result = result
	VFMapInstance.Update("p_reportanswer", vf)
	ULogger.Info("reportanswer %s", vf.String())
	return nil
}

///============================
///消费者 获取图片
func getFile(session *link.Session, req map[string]string) error {
	seq := req["seq"]

	vf := QueueInstance.DequeueWithoutPClosed()
	if vf == nil {
		ULogger.Info("queue is nil ,c userid is %s", session.State.(*User).Id)
		ret := map[string]string{
			"action": "res_getfile",
			"seq":    seq,
		}
		by, _ := json.Marshal(ret)
		session.Send(link.Bytes(by))
		ULogger.Info("send to client %s %s %s", session.Conn().RemoteAddr().String(), "say:", string(by))
		return nil
	}
	vf.C = session
	vf.CInfo = session.State.(*User).Id
	vf.Status = 2
	vf.CGetUnix = time.Now().Unix()
	ret := map[string]string{
		"action": "res_getfile",
		"seq":    seq,
		"id":     vf.Id,
		"file":   vf.File,
	}
	by, _ := json.Marshal(ret)
	VFMapInstance.Update("c_getfile", vf)
	session.Send(link.Bytes(by))
	ULogger.Info("res_getfile %s %s %s", session.Conn().RemoteAddr().String(), "say:", vf.String())
	return nil
}

///c端开始答题
func cStart(session *link.Session, req map[string]string) error {

	userid := req["userid"]
	password := req["password"]
	seq := req["seq"]
	ULogger.Info("%s start answer", userid)
	ret := map[string]string{
		"action": "res_cstart",
		"seq":    seq,
		"result": "0",
	}

	if session.State != nil && session.State.(*User).UserType == "C" {
		ULogger.Error("have logined ", userid)
		session.Close()
	}

	if b := Login(userid, password, "c"); !b {
		ULogger.Error("cstart failed ,userid is %s", userid)

		by, _ := json.Marshal(ret)
		session.Send(link.Bytes(by))
		ULogger.Info("send to client %s %s %s", session.Conn().RemoteAddr().String(), "say:", string(by))
		session.Close()
		return nil
	} else {
		user := &User{UserType: "C", Id: userid, WorkTime: time.Now().Format("2006-01-02 15:04:05"), CWrongCount: 0}
		session.State = user

		ret["result"] = "1"

		by, _ := json.Marshal(ret)

		session.Send(link.Bytes(by))
		ULogger.Info("cstart %s %s %s", session.Conn().RemoteAddr().String(), "say:", string(by))
		//c端开始答题
		Exec(`insert into user_activities(user_id,active_time,active_type,user_type,other_info) values(?,now(),'begin','customer',?)`, userid, session.Conn().RemoteAddr().String())
		VFMapInstance.AddCSession(session)
	}

	return nil
}

///m端开登陆
func mStart(session *link.Session, req map[string]string) error {

	userid := req["userid"]
	password := req["password"]
	seq := req["seq"]
	ULogger.Info("m_user start %s", userid)
	ret := map[string]string{
		"action": "res_mstart",
		"seq":    seq,
		"result": "0",
	}

	if session.State != nil && session.State.(*User).UserType == "M" {
		ULogger.Error("m have logined ", userid)
		session.Close()
	}

	if b := Login(userid, password, "m"); !b {
		ULogger.Error("m_start failed ,userid is %s ", userid)

		by, _ := json.Marshal(ret)
		session.Send(link.Bytes(by))
		ULogger.Info("send to client %s %s %s", session.Conn().RemoteAddr().String(), "say:", string(by))
		session.Close()
		return nil
	} else {
		user := &User{UserType: "M", Id: userid, WorkTime: time.Now().Format("2006-01-02 15:04:05"), CWrongCount: 0}
		session.State = user

		ret["result"] = "1"

		by, _ := json.Marshal(ret)

		session.Send(link.Bytes(by))
		ULogger.Info("mstart %s %s %s", session.Conn().RemoteAddr().String(), "say:", string(by))
		//m端登陆
		Exec(`insert into user_activities(user_id,active_time,active_type,user_type,other_info) values(?,now(),'begin','manager',?)`, userid, session.Conn().RemoteAddr().String())
	}

	return nil
}

func getStatInfo(session *link.Session, req map[string]string) error {
	if session.State == nil {
		ULogger.Error("bad request,please login first")
		session.Close()
		return nil
	}
	seq := req["seq"]
	_ = seq
	user := session.State.(*User)
	if user.UserType != "M" {
		ULogger.Error("only M can do it")
		session.Close()
		return nil
	}
	ret := ResMessage{Action: "res_getstatinfo", Seq: seq}
	ret.StatInfo = make(MM, 0)

	userId := user.Id
	_ = userId
	beginTime := req["begin_time"]
	if beginTime == "" {
		beginTime = GetSysLastStartTime()
	} else {
		tm, _ := time.Parse("2006-01-02 15:04:05", beginTime)
		//德国服务器比utc 时间早一个市区
		tm = tm.Add(time.Hour * 2)
		beginTime = tm.Format("2006-01-02 15:04:05")
	}
	uids := GetUsersByManagerId(userId)
	_ = uids
	if uids == nil || len(uids) == 0 {
		by, _ := json.Marshal(ret)
		session.Send(link.Bytes(by))
		return nil
	}
	for _, v := range uids {
		query := `select count(*) from exam where c_userid=? and c_getfile_time > ? and answer is not null`
		finishCount := QueryInt(query, v, beginTime)
		query1 := `select count(*) from exam where c_userid=? and c_getfile_time > ? and answer is not null and answer_result=1`
		rightCount := QueryInt(query1, v, beginTime)
		m := map[string]string{"user_id": v, "finish_count": strconv.Itoa(finishCount), "right_count": strconv.Itoa(rightCount)}
		ret.StatInfo = append(ret.StatInfo, m)
	}
	if len(ret.StatInfo) > 0 {
		sort.Sort(ret.StatInfo)
	}
	by, _ := json.Marshal(ret)
	session.Send(link.Bytes(by))
	ULogger.Info("send to client %s %s %s", session.Conn().RemoteAddr().String(), "say:", string(by))

	return nil
}

//c端回答的问题
func answer(session *link.Session, req map[string]string) error {
	id := req["id"]
	answer := req["answer"]

	vf := VFMapInstance.Get(id)
	if vf == nil {
		ULogger.Error("answer,verifyobj not found,may be timeout,vf is %s", req)

		return nil
	}
	vf.Answer = answer
	vf.Status = 3
	vf.CInfo = session.State.(*User).Id
	//c回答了问题
	VFMapInstance.Update("c_answer", vf)
	if vf.P != nil && !vf.P.IsClosed() {
		//相当于putfile 的应答
		ret := map[string]string{
			"action": "res_putfile",
			"seq":    vf.Seq,
			"id":     vf.Id,
			"answer": answer,
		}
		by, _ := json.Marshal(ret)
		vf.P.Send(link.Bytes(by))
		ULogger.Info("res_putfile %s %s %s", vf.P.Conn().RemoteAddr().String(), "say:", string(by))
		vf.Status = 4
		//给P投递了消息
		VFMapInstance.Update("p_fileanswer", vf)
	} else {
		ULogger.Info("p have closed")
		VFMapInstance.Update("p_close", vf)
	}

	return nil
}

//保证唯一，操作数据库的时候做主键的，go 操作mysql 不熟悉，不知道插入新纪录如何返回主键，有时间研究
func getId() string {
	t := time.Now().Unix()
	VFId := atomic.AddInt32(&VFId, 1)
	return strconv.Itoa(int(t)) + "_" + strconv.Itoa(int(VFId))
}
