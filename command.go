package main

import (
	"encoding/json"
	"github.com/funny/link"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	VFId int32
)

type User struct {
	UserType string
	Id       string
	WorkTime string
}

func Process(session *link.Session, req map[string]string) error {
	command, ok := req["action"]
	if !ok {
		ULogger.Error("client", session.Conn().RemoteAddr().String(), "bad request ,not found action")
		session.Close()
		return nil

	}
	if (command == "getfile" || command == "answer") && session.State == nil {
		ULogger.Error("client", session.Conn().RemoteAddr().String(), "c must login frist")
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
	default:
		ULogger.Error("client", session.Conn().RemoteAddr().String(), "not support command")
		session.Close()
		//ULogger.Info("sssss")
	}
	return nil
}

///======================================
///生产者 ，产生图验图片
func putFile(session *link.Session, req map[string]string) error {
	file := req["file"]
	seq := req["seq"]
	fileid := GetMd5String(file)
	vf := &VerifyObj{Id: getId(), P: session, C: nil, FileId: fileid, File: file, Status: 1, Result: "0", Seq: seq, PPutUnix: time.Now().Unix()}
	QueueInstance.Enqueue(vf)
	VFMapInstance.Put(vf)
	ULogger.Infof("putfile 进队列 %v\n", vf)
	//记录p端的操作
	if session.State == nil {
		userId := session.Conn().RemoteAddr().String()
		user := &User{UserType: "P", Id: userId, WorkTime: time.Now().Format("2006-01-02 15:04:05")}
		session.State = user
		Exec(`insert into user_activities(user_id,active_time,active_type,user_type,other_info) values(?,now(),'begin','production',?)`, userId, userId)
	}
	return nil
}

func reportAnswer(session *link.Session, req map[string]string) error {
	id := req["id"]
	result := req["result"]

	vf := VFMapInstance.Get(id)
	if vf == nil {
		ULogger.Errorf("answer,verifyobj not found,%v\n")
		return nil
	}
	vf.Result = result
	VFMapInstance.Update("p_reportanswer", vf)

	return nil
}

///============================
///消费者 获取图片
func getFile(session *link.Session, req map[string]string) error {
	seq := req["seq"]

	//死等
	vf := QueueInstance.DeChan()
	if vf == nil {
		ULogger.Error("getfile time out,sessioninfo is %s", session.State.(*User).Id)
		ret := map[string]string{
			"action": "res_getfile",
			"seq":    seq,
		}
		by, _ := json.Marshal(ret)
		session.Send(link.Bytes(by))
		ULogger.Info("send to client", session.Conn().RemoteAddr().String(), "say:", string(by))
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
	ULogger.Info("send to client", session.Conn().RemoteAddr().String(), "say:", string(by))
	return nil
}

///c端开始答题
func cStart(session *link.Session, req map[string]string) error {
	userid := req["userid"]
	password := req["password"]
	seq := req["seq"]
	ULogger.Infof("user %s start answer", userid)
	ret := map[string]string{
		"action": "res_cstart",
		"seq":    seq,
		"result": "0",
	}
	if session.State != nil && session.State.(*User).UserType == "C" {
		ULogger.Error("have logined ", userid)
		session.Close()
	}

	if b := Login(userid, password); !b {
		ULogger.Errorf("cstart failed ,userid is %s password is %s", userid, password)

		by, _ := json.Marshal(ret)
		session.Send(link.Bytes(by))
		ULogger.Info("send to client", session.Conn().RemoteAddr().String(), "say:", string(by))
		session.Close()
		return nil
	} else {
		user := &User{UserType: "C", Id: userid, WorkTime: time.Now().Format("2006-01-02 15:04:05")}
		session.State = user

		ret["result"] = "1"

		by, _ := json.Marshal(ret)

		session.Send(link.Bytes(by))
		ULogger.Info("send to client", session.Conn().RemoteAddr().String(), "say:", string(by))
		//c端开始答题
		Exec(`insert into user_activities(user_id,active_time,active_type,user_type,other_info) values(?,now(),'begin','customer',?)`, userid, session.Conn().RemoteAddr().String())
		VFMapInstance.AddSession(session)
	}

	return nil
}

//c端回答的问题
func answer(session *link.Session, req map[string]string) error {
	id := req["id"]
	answer := req["answer"]

	vf := VFMapInstance.Get(id)
	if vf == nil {
		ULogger.Errorf("answer,verifyobj not found,may be timeout%v\n")
		Exec(`update exam set f_status=8,answer_time=now(),answer=? where file_id=?`, answer, id)
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
		ULogger.Info("send to client", vf.P.Conn().RemoteAddr().String(), "say:", string(by))
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
