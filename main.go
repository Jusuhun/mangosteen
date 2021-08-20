package main

//결과 남길 양식 및 폴더 트리 정하기 및 구현
//Line 받아오는 구조 구현 (Client)

//[x] Log가 남는 PC에서 실행할 Line 보내주는 프로그램 구현(Server)

import (
	"container/list"
	"flag"
	"fmt"
	"mangosteena/internal/testa"
	"mangosteena/pkg/info"
)

var VERSION = "dev"
var BUILD = "dev"

func main() {
	//getOption
	version := flag.Bool("version", false, "check version")
	rootPath := flag.String("root", "./", "root path")
	configFile := flag.String("config", "config.yaml", "config file")
	//시작시간을 지정해준다. (지정안하면 실행한 시간부터)
	flag.Parse()

	//check version
	if *version {
		fmt.Println("The version is", VERSION)
		fmt.Println("The Build date is", BUILD)
		return
	}

	//ReadFile
	test := testa.ReadFile(*rootPath + *configFile)

	que := list.New()
	for true /*!test.done*/ {
		if que.Len() == 0 {
			//rest API
			//buffer
			//임시
			que.PushBack(info.DecodeCSV("{time},event,init"))
			que.PushBack(info.DecodeCSV("{time},operation,start-sw,off"))
			que.PushBack(info.DecodeCSV("{time},operation,start-sw,on"))
		}

		if que.Len() == 0 {
			//sleep(100)
		} else {
			test.Valid(que.Front().Value.(info.LogInfo))
			que.Remove(que.Front())
		}

		if true {
			break
		}
	}

	test.WriteReport(*rootPath, VERSION, BUILD)

	fmt.Println("done.")
}
