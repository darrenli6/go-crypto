package mapReduce

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// 实现一个 map管理函数 从 input file 中读取内容
// 将输出分为指定数量的中间文件
// 自定义分区标准

func doMap(
	jobName string,
	mapTaskNumber int,
	inFile string,
	nReduce int,
	mapF func(file string, content string) []KeyValue) {

	//从输入文件inFile中读取内容
	f, err := os.Open(inFile)
	if err != nil {
		log.Printf("open file %s failed %v \n", inFile, err)
	}
	defer f.Close()
	// 从输入文件读取内容
	content, err := ioutil.ReadAll(f)
	if nil != err {
		log.Printf("read the content of file failed ! %v\n", err)
	}

	// 通过调用mapF对内容进行处理，分割map任务输出

	kvs := mapF(inFile, string(content))

	// 生成一个编码对象
	// 创建每一个map 生成nreduce 中间文件对象
	encoder := make([]*json.Encoder, nReduce)
	// 创建nReduce 个中间结果文件
	for i := 0; i < nReduce; i++ {
		//生成中间文件的名称

		file_name := reduceName(jobName, mapTaskNumber, i)
		f, err := os.Create(file_name)
		if nil != err {
			log.Printf("unable to create [%s]: %v \n", file_name, err)
		}
		defer f.Close()
		encoder[i] = json.NewEncoder(f)

	}

	// 将kvs的内容存入前面生成的中间文件中去
	for _, v := range kvs {
		// 自定义规则对k进行分类
		// 此处以编号对nReduce取余分类

		index := ihash(v.Key) % nReduce
		if err := encoder[index].Encode(&v); nil != err {
			log.Printf("Unable to write file \n ")
		}

	}
}
