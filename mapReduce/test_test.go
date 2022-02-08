package mapReduce

import (
	"testing"
	"fmt"
	"os"
	"log"
	"bufio"
	"strings"
)

const (
	nNumber =100
)

// 创建一个N个编号的输入文件
// 通过mapreduce进行处理
// 检查最终输出文件中是否包含N个编号


// 自定义map分割处理函数
func MapFunc(file string,value string) (res []KeyValue){

   words := strings.Fields(value)

   for _, w :=range words{
   	  kv :=KeyValue{w,""}
   	  res =append(res,kv)
   }

   return


}

// 聚合函数

func ReduceFunc(Key string,value []string) string {

  return ""

}

func TestSequentialSignal(t *testing.T) {
	Sequential("test",
		makeInputs(1),
		1,
		MapFunc ,ReduceFunc)

}



func TestSequentialMany(t *testing.T) {
	Sequential("test",
		makeInputs(5),
		3,
		MapFunc,ReduceFunc)


}




// 创建输入文件
// 根据指定的数量,返回创建好的文件名列表
// 写入相应的数量
// count 创建文件的数量
func makeInputs(num int) []string{


	var names []string

	var i=0
	for f:=0;f<num;f++{
		// 文件命名方式,根据课程命名
		names=append(names,fmt.Sprintf("824-mrinput-%d.txt",f))
		// 创建文件
		file,err :=os.Create(names[f])
		if err!=nil{
			log.Fatalf("create input file [%s] failed  error : ",file,err)
		}

		w:=bufio.NewWriter(file)

		for i<(f+1) *(nNumber/num){
			// 写入i到w中
			fmt.Fprintf(w,"%d\n",i)
			i++
		}
		// 把buffer写到文件中
		w.Flush()
		file.Close()
	}

	return names
}