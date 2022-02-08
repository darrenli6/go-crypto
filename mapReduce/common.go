package mapReduce


// 自定义任务类型

type jobPhase string


const (


	mapPhase jobPhase="Map"

	reducePhase ="Reduce"
)

type KeyValue struct {
	Key string
	Value string
}

