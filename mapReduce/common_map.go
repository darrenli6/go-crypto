package mapReduce


// 实现一个 map管理函数 从 input file 中读取内容
// 将输出分为指定数量的中间文件
// 自定义分区标准


func doMap(
	jobName string,
	mapTaskNumber int,
	inFile string,
	nReduce int,
	mapF func(file string,content string)  []KeyValue ){


}