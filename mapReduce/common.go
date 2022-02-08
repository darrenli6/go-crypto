package mapReduce

import (
	"hash/fnv"
	"strconv"
)

// 自定义任务类型

type jobPhase string

const (
	mapPhase jobPhase = "Map"

	reducePhase = "Reduce"
)

type KeyValue struct {
	Key   string
	Value string
}

//reduce task 输出文件名称
func mergeName(jobName string, reduceTask int) string {
	return "mrtmp." + jobName + "-res-" + strconv.Itoa(reduceTask)
}

//生成中间文件名称
func reduceName(jobName string, mapTask int, reduceTask int) string {

	return "mrtmp." + jobName + "-" + strconv.Itoa(mapTask) + "-" + strconv.Itoa(reduceTask)

}

// 哈希函数

func ihash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32() & 0x7fffffff)

}
