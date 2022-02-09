package mapReduce

import (
	"os"
	"log"
	"encoding/json"
)

// 管理reduce任务
//合并操作
func doReduce(
	jobName string,
	reduceTaksNumber int,  // 任务编号

	outFile string,
	nMap int , // 运行map的任务号
	reduceF func(Key string,value []string) string){

		var result map[string][]string=make(map[string][]string)

		// 打开每一个中间文件
		for i:=0;i<nMap;i++{

			interFile := reduceName(jobName,i,reduceTaksNumber)

			f,err :=os.Open(interFile)

			if err!=nil{
				log.Printf("read content from file [%s] filed ! %v \n",interFile,err)
			}

			defer  f.Close()

			decoder :=json.NewDecoder(f)

			var kv KeyValue

			for ; decoder.More();{
				//
				err :=decoder.Decode(&kv)
				if err!=nil{
					log.Printf("json decode failed %v",err)
				}
				// 将具有相同key合并

				result[kv.Key] =append(result[kv.Key],kv.Value)


			}
		}



		// 获取内容
		// 把相同key的内容合并
		// 生成最终文件



		var keys []string
		for key,_ :=range result{
			keys=append(keys, key)
		}

		// 新建输出文件

		out_file,err :=os.Create(outFile)
		if nil!=err{
			log.Printf("create outfile failed  %v",err)
		}

		defer out_file.Close()

		encoder :=json.NewEncoder(out_file)
		for _, key :=range  keys{
			encoder.Encode(KeyValue{key,reduceF(key,result[key])})
		}


}