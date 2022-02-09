package mapReduce

// 任务调度的函数

// master结构
type Master struct {
}

func Sequential(
	jobName string, // 任务名称
	files []string, // 输出文件 待处理文件
	nReduce int,    // 文件的分区数量

	mapF func(string, string) []KeyValue,
	reduceF func(string, []string) string) {

	// 执行分配的任务
	// 在mapReduce 中,任务要分为map任务和reduce任务

	mr := newMaster()
	mr.run(jobName, files, nReduce, func(phase jobPhase) {

		switch phase {
		case mapPhase:

			// 执行map任务
			// map任务的调用次数由输入文件个数决定
			for i, f := range files {
				doMap(jobName, i, f, nReduce, mapF)
			}

			break
		case reducePhase:

			// reduce任务的调用次数由nreduce 大小来决定
			for i := 0; i < nReduce; i++ {

				doReduce(jobName, i, mergeName(jobName, i), len(files), reduceF)
			}

			break

		}
	})

	//
	//doMap("",mapF)
	//doReduce("",reduceF)

}

// 初始化
func newMaster() *Master {

	return nil
}

// 实际上执行函数
// 执行给定的任务

func (mr *Master) run(jobName string, files []string, nreduce int, schedule func(phase jobPhase)) {




	// 顺序执行map任务

	schedule(mapPhase)

	// 执行reduce任务

	schedule(reducePhase)
}
