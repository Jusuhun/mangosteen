package main

//결과 남길 양식 및 폴더 트리 정하기 및 구현
//Line 받아오는 구조 구현 (Client)

//[x] Log가 남는 PC에서 실행할 Line 보내주는 프로그램 구현(Server)

import (
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

	//CSV 파일 읽어서 하는 걸로  && 통신으로 받는 걸로 변경하면 좋음
	array := []info.LogInfo{
		info.DecodeCSV("{time},operation,start-sw,on"),
		info.DecodeCSV("{time},event,init"),
		info.DecodeCSV("{time},operation,start-sw,off"),
	}

	for i := range array {
		test.Valid(array[i])
	}

	test.WriteReport(*rootPath, VERSION, BUILD)

	fmt.Println("done.")
}
