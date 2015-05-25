package lib

import (
	"strconv"
)

type MM []map[string]string

func (list MM) Len() int {
	return len(list)
}

func (list MM) Less(i, j int) bool {
	t1, _ := strconv.Atoi(list[i]["finish_count"])
	t2, _ := strconv.Atoi(list[j]["finish_count"])
	if t1 > t2 {
		return true
	} else {
		return false
	}
}

func (list MM) Swap(i, j int) {
	var temp map[string]string = list[i]
	list[i] = list[j]
	list[j] = temp
}
