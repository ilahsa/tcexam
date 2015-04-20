package lib

/*
放置验证码的队列
*/
import (
	"container/list"
	"errors"
	"sync"
	//	"time"
)

var (
	QueueInstance = newQueue()
	QueuFullError = errors.New("queue full,default size is 500")
)

type Queue struct {
	MaxSize  int
	lst      *list.List
	syncRoot sync.Mutex
	//Ch       chan *VerifyObj
}

func newQueue() *Queue {
	q := &Queue{MaxSize: 500, lst: list.New()}
	//q.startEnChan()
	return q
}

func (q *Queue) Enqueue(v *VerifyObj) error {
	q.syncRoot.Lock()
	defer q.syncRoot.Unlock()
	if q.len() >= 500 {
		return QueuFullError
	}
	q.lst.PushBack(v)

	return nil
}

//排除掉 已经关闭的p 端的问题
func (q *Queue) DequeueWithoutPClosed() *VerifyObj {
	for {
		vf := q.Dequeue()
		if vf == nil {
			return vf
		}
		if vf.P == nil ||
			vf.P.IsClosed() {
			continue
		} else {
			return vf
		}
	}
	return nil
}
func (q *Queue) Dequeue() *VerifyObj {
	q.syncRoot.Lock()
	defer q.syncRoot.Unlock()
	if q.len() <= 0 {
		return nil
	}
	e := q.lst.Front()

	if e != nil {
		q.lst.Remove(e)
		return e.Value.(*VerifyObj)
	} else {
		return nil
	}
}

func (q *Queue) len() int {
	return q.lst.Len()
}

//func (q *Queue) startEnChan() {
//	go func() {
//		for {
//			vf := q.Dequeue()
//			if vf != nil {
//				q.Ch <- vf
//			}
//		}
//	}()
//}

//func (q *Queue) DeChan() *VerifyObj {
//	select {
//	case vf := <-q.Ch:
//		return vf
//	case <-time.After(time.Second * 2):
//		return nil
//	}
//}
