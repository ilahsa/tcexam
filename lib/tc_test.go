package lib

import (
	//"database/sql"
	"fmt"
	"testing"
	//	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func Test_02(t *testing.T) {

	fmt.Println("22222")
}
func Test_MD5(t *testing.T) {
	fmt.Println(GetMd5String("123456"))
	str := GetMd5String("123456")
	t.Log(str)
}

func Test_login(t *testing.T) {
	b := Login(`u_001`, `123456`)
	if !b {
		t.Fail()
	}
}

func Test_Exec(t *testing.T) {
	//插入题目
	//Exec(`insert into exam(file_id,file_hash,f_status,put_time) values(?,?,1,now())`, `f_id_001`, `f_hash_001`)
	//c端获取问题
	Exec(`update exam set f_status=?,c_userid=? where file_id=?`, 2, `c_user_001`, `f_id_001`)
	//c端回答问题
	Exec(`update exam set f_status=?,answer_time=now(),answer=? where file_id=?`, 3, `answer_0001`, `f_id_001`)

	//给p端下发问题
	Exec(`update exam set f_status=? where file_id=?`, 4, `f_id_001`)
	//p端发送是否正确
	Exec(`update exam set f_status=?,answer_result=? where file_id=?`, 5, 1, `f_id_001`)

	i1 := QueryInt(`select count(*) from exam where c_userid='u_001' and c_getfile_time > '2015-04-17 12:54:31' and answer is not null`)
	fmt.Println(i1)
}
