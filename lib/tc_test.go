package lib

import (
	"encoding/json"
	"os"
	"sort"
	//"database/sql"
	"fmt"
	"testing"
	//	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func TestMain(t *testing.M) {
	fmt.Println("start")

	InitConfig()
	fmt.Println(TCConfig.DBConnectStr)
	code := t.Run()
	os.Exit(code)
}
func Test_02(t *testing.T) {

	fmt.Println("22222")
}
func Test_MD5(t *testing.T) {
	fmt.Println(GetMd5String("123456"))
	str := GetMd5String("123456")
	t.Log(str)
}

func Test_login(t *testing.T) {
	b := Login(`u_001`, `123456`, "c")
	if !b {
		t.Fail()
	}
}

func Test_Exec(t *testing.T) {
	an := GetAnswer("6eb512cb2557734c63f6fa15ef61ef9f")
	fmt.Println(an)
	fmt.Println("wwwwwwwwww")
	fmt.Println(GetSysLastStartTime())
	ret := GetUsersByManagerId("m_0001")
	for _, v := range ret {
		fmt.Println(v)
	}
}

func Test_Json(t *testing.T) {
	m1 := ResMessage{Action: "res_getstatinfo", Seq: "22"}
	m1.StatInfo = make([]map[string]string, 0)
	//m1.StatInfo[0] = map[string]string{"user_id": "w_0001", "finish_count": "22"}
	//m1.StatInfo[1] = map[string]string{"user_id": "w_0002", "finish_count": "22"}
	by, _ := json.Marshal(m1)
	fmt.Println(string(by))
}

func Test_Sort(t *testing.T) {
	mm := make(MM, 5)
	mm[0] = map[string]string{"name": "n1", "finish_count": "20"}
	mm[1] = map[string]string{"name": "n2", "finish_count": "234"}
	mm[2] = map[string]string{"name": "n3", "finish_count": "20"}
	mm[3] = map[string]string{"name": "n4", "finish_count": "1"}
	mm[4] = map[string]string{"name": "n5", "finish_count": "23"}
	sort.Sort(mm)
	fmt.Println(mm)
}
