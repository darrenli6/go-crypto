package mapReduce


// 管理reduce任务

func doReduce(
	jobName string,
	reduceTaksNumber int,

	outFile string,
	nMap int , // 运行map的任务号
	reduceF func(Key string,value []string) string){

}