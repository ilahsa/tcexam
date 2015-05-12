package lib

import (
	"bytes"
	//"fmt"
	"strconv"
	//	"strconv"
	"encoding/json"
	"github.com/funny/link"
	"sync"
	"time"
)

var (
	VFMapInstance = newVFMapInstance()
)

func init() {
	go func() {
		t1 := time.NewTimer(time.Second * 60)
		t2 := time.NewTimer(time.Second * 30)
		for {
			select {
			case <-t1.C:
				nowUnix := time.Now().Unix()
				for k, v := range VFMapInstance.innerMap {
					//p 端超时
					if v.PPutUnix+900 < nowUnix {
						ULogger.Info(v.Id, "file timeover 15 minutes,puttime ", v.PPutUnix, " now is ", time.Now().Unix())
						//						delete(VFMapInstance.innerMap, k)
						//						if v.P != nil && v.P.IsClosed() {
						//							v.P.Close()
						//						}
					}
					//被c 端获取且超时
					if v.Status == 2 && (v.CGetUnix+300) < nowUnix {
						ULogger.Info(v.Id, " recycle")
						VFMapInstance.Update("recycle", v)
						delete(VFMapInstance.innerMap, k)
						uid := ""
						if v.C != nil {
							uid = v.C.State.(*User).Id
						}
						Exec(`insert into recycle_log(c_id,file_id,file_hash,put_time,c_getfile_time,recycle_time) 
values(?,?,?,?,?,?)`, uid, v.Id, v.FileId, v.PPutUnix, v.CGetUnix, time.Now().Unix())
						vf := &VerifyObj{Id: getId(), P: v.P, C: nil, FileId: v.FileId, File: v.File, Status: 1, Result: "0", Seq: v.Seq, PPutUnix: v.PPutUnix}
						QueueInstance.Enqueue(vf)
						VFMapInstance.Put(vf)
						ULogger.Info("recycle putfile enqueue %s", vf.String())
					}
				}
				t1.Reset(time.Second * 60)
			case <-t2.C:
				//fmt.Println("ss")
				stat := map[string]string{"action": "stat"}
				stat["now"] = strconv.Itoa(int(time.Now().Unix()))
				//TCServer.
				pstatInfo := VFMapInstance.pstatInfo()
				for k, v := range VFMapInstance.c_sessions {
					if v == nil || v.State == nil {
						delete(VFMapInstance.c_sessions, k)
					}
					u := v.State.(*User)
					query := `select count(*) from exam where c_userid=? and c_getfile_time > ? and answer is not null`
					//c端回答的问题总数
					canswer := QueryInt(query, u.Id, u.WorkTime)
					//ULogger.Info(query+, u.Id, u.WorkTime)

					query1 := `select count(*) from exam where c_userid=? and c_getfile_time > ? and answer is not null and answer_result=1`
					canswerrigth := QueryInt(query1, u.Id, u.WorkTime)
					waitcount := QueueInstance.len()
					clientcount := len(VFMapInstance.c_sessions)
					pcount := len(VFMapInstance.p_sessions)
					questioncount := len(VFMapInstance.innerMap)

					stat["productioncount"] = strconv.Itoa(pcount)
					stat["questioncount"] = strconv.Itoa(questioncount)
					stat["finishcount"] = strconv.Itoa(canswer)
					stat["rightcount"] = strconv.Itoa(canswerrigth)
					stat["waitcount"] = strconv.Itoa(waitcount)
					stat["clientcount"] = strconv.Itoa(clientcount)
					stat["pstatinfo"] = pstatInfo
					by, _ := json.Marshal(stat)
					v.Send(link.Bytes(by))
					ULogger.Info("send to client %s %s %s", v.Conn().RemoteAddr().String(), "say:", string(by))
				}
				t2.Reset(time.Second * 30)
			}
		}
	}()
}

type VFMap struct {
	innerMap   map[string]*VerifyObj
	c_sessions map[uint64]*link.Session
	p_sessions map[uint64]*link.Session
	syncRoot   sync.Mutex
}

func newVFMapInstance() *VFMap {
	return &VFMap{innerMap: make(map[string]*VerifyObj, 500), c_sessions: map[uint64]*link.Session{}, p_sessions: map[uint64]*link.Session{}}
}

func (m *VFMap) Put(vf *VerifyObj) {
	m.syncRoot.Lock()
	defer m.syncRoot.Unlock()
	_, ok := m.innerMap[vf.Id]
	if ok {
		ULogger.Error("key exists")
		return
	}

	m.innerMap[vf.Id] = vf
	//添加数据库的insert
	//插入题目
	Exec(`insert into exam(file_id,file_hash,f_status,put_time) values(?,?,1,now())`, vf.Id, vf.FileId)
}

//添加c 端的session
//同一账户只能登陆一个
func (m *VFMap) AddCSession(s *link.Session) {
	m.syncRoot.Lock()
	defer m.syncRoot.Unlock()

	u := s.State.(*User)
	for k, v := range m.c_sessions {
		tmpU := v.State.(*User)
		if u.Id == tmpU.Id && !v.IsClosed() {
			ULogger.Info("click c, cinfo is %s", v.Conn().RemoteAddr().String())
			notify := `{
	"action": "notify_kick",
	"seq": "s_0001"
}`
			v.Send(link.String(notify))

			delete(m.c_sessions, k)
			v.Close()
		}
	}

	m.c_sessions[s.Id()] = s
}

//p 端的统计信息
func (m *VFMap) pstatInfo() string {
	m.syncRoot.Lock()
	defer m.syncRoot.Unlock()

	bufs := bytes.Buffer{}
	for k, v := range m.p_sessions {
		if v == nil {
			delete(m.p_sessions, k)
		}
		tmpU := v.State.(*User)

		bufs.WriteString("p_addr:" + v.Conn().RemoteAddr().String() + ";pputcount:" + strconv.Itoa(tmpU.PPutCount) + ";plastputtime:" + strconv.Itoa(int(tmpU.PLastPutTime)) + "\r\n")
	}
	return bufs.String()
}

//添加p 端的session
func (m *VFMap) AddPSession(s *link.Session) {
	m.syncRoot.Lock()
	defer m.syncRoot.Unlock()

	m.p_sessions[s.Id()] = s
}

func (m *VFMap) Get(id string) *VerifyObj {
	m.syncRoot.Lock()
	defer m.syncRoot.Unlock()
	_, ok := m.innerMap[id]
	if !ok {
		//ULogger.Error("verifyobj not found")
		return nil
	}
	return m.innerMap[id]
}

//c 端关闭
func (m *VFMap) DelSessionByC(s *link.Session) {
	//	m.syncRoot.Lock()
	//	defer m.syncRoot.Unlock()

	delete(m.c_sessions, s.Id())
	//回收c 的任务
	for k, v := range m.innerMap {
		if v.Status == 2 && v.C != nil && v.C == s {
			ULogger.Info("c is closed, recycle file,fileid is %s", v.Id)
			VFMapInstance.Update("recycle", v)
			delete(VFMapInstance.innerMap, k)
			Exec(`insert into recycle_log(c_id,file_id,file_hash,put_time,c_getfile_time,recycle_time) 
values(?,?,?,?,?,?)`, s.State.(*User).Id, v.Id, v.FileId, v.PPutUnix, v.CGetUnix, time.Now().Unix())

			vf := &VerifyObj{Id: getId(), P: v.P, C: nil, FileId: v.FileId, File: v.File, Status: 1, Result: "0", Seq: v.Seq, PPutUnix: v.PPutUnix}
			QueueInstance.Enqueue(vf)
			VFMapInstance.Put(vf)
			ULogger.Info("recycle putfile enqueue %s", vf.String())
		}
	}
}

//p 端关闭，清楚掉所有的session
func (m *VFMap) DelSessionByP(s *link.Session) {
	//	m.syncRoot.Lock()
	//	defer m.syncRoot.Unlock()

	delete(m.p_sessions, s.Id())
	for k, v := range m.innerMap {
		if v.P == s {
			ULogger.Info("p is closed %s", s.Conn().RemoteAddr().String())
			m.Update("p_closed", v)
			delete(m.innerMap, k)
		}
	}
}

///更新map
func (m *VFMap) Update(action string, vf *VerifyObj) {
	//	m.syncRoot.Lock()
	//	defer m.syncRoot.Unlock()
	_, ok := m.innerMap[vf.Id]
	if !ok {
		ULogger.Error("key not exists")
		return
	}
	//更新数据库
	switch action {
	case "c_getfile":
		//c端获取问题
		Exec(`update exam set f_status=?,c_userid=?,c_getfile_time=now() where file_id=?`, 2, vf.CInfo, vf.Id)
	case "c_answer":
		//c端回答问题
		Exec(`update exam set f_status=?,answer_time=now(),answer=? where file_id=?`, 3, vf.Answer, vf.Id)
	case "p_fileanswer":
		//给p端下发问题
		Exec(`update exam set f_status=? where file_id=?`, 4, vf.Id)
	case "p_reportanswer":
		//p端发送是否正确
		Exec(`update exam set f_status=?,answer_result=? where file_id=?`, 5, vf.Result, vf.Id)
		delete(m.innerMap, vf.Id)
	case "p_closed":
		//p端已经关闭
		Exec(`update exam set f_status=7 where file_id=?`, vf.Id)
	case "recycle":
		//问题被回收
		Exec(`update exam set f_status=6,recycle_time=now() where file_id=?`, vf.Id)

	}
}
