package main

import (
	"github.com/funny/link"
)

type VerifyObj struct {
	Id     string
	P      *link.Session //生产者session
	C      *link.Session //消费者session
	CInfo  string
	Status int //状态 1 生产者入队列 2 消费者获取 3 消费者回复答案 4 给生产者回发
	//5 p 端回复是否正确  6 问题被回收 7 p端已经关闭 8 c 端超时 回答问题
	FileId   string
	File     string
	Answer   string
	Result   string //是否正确 1 正确 0 错误
	Seq      string //上传问题时的seq
	CGetUnix int64  //c端获取的时间戳 5分钟超时
	PPutUnix int64  //p端放入的时间戳 15分钟超时
}
