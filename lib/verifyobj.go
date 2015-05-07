package lib

import (
	"bytes"
	"github.com/funny/link"
	"strconv"
	"time"
)

type VerifyObj struct {
	Id     string
	P      *link.Session //生产者session
	C      *link.Session //消费者session
	CInfo  string
	Status int //状态 1 生产者入队列 2 消费者获取 3 消费者回复答案 4 给生产者回发
	//5 p 端回复是否正确  6 问题被回收 7 p端已经关闭 8 c 端超时 回答问题 9 从服务端直接获取答案
	FileId   string
	File     string
	Answer   string
	Result   string //是否正确 1 正确 0 错误
	Seq      string //上传问题时的seq
	CGetUnix int64  //c端获取的时间戳 5分钟超时
	PPutUnix int64  //p端放入的时间戳 15分钟超时
}

func (v *VerifyObj) String() string {
	bys := bytes.Buffer{}
	bys.WriteString("id:" + v.Id)
	bys.WriteString(";md5:" + v.FileId)
	bys.WriteString(";puttime:" + strconv.Itoa(int(v.PPutUnix)))
	bys.WriteString(";gettime:" + strconv.Itoa(int(v.CGetUnix)))
	bys.WriteString(";nowtime:" + strconv.Itoa(int(time.Now().Unix())))
	pinfo := ""
	if v.P != nil && !v.P.IsClosed() {
		pinfo = v.P.Conn().RemoteAddr().String()
	}
	cinfo := ""
	if v.C != nil && !v.C.IsClosed() {
		cinfo = v.C.Conn().RemoteAddr().String()
	}
	bys.WriteString(";p:" + pinfo)
	bys.WriteString(";c:" + cinfo)
	bys.WriteString(";answer:" + v.Answer)
	bys.WriteString(";p_seq:" + v.Seq)
	bys.WriteString(";status:" + strconv.Itoa(v.Status))
	return bys.String()
}
