package main

//결과 남길 양식 및 폴더 트리 정하기 및 구현
//Line 받아오는 구조 구현 (Client)

//[x] Log가 남는 PC에서 실행할 Line 보내주는 프로그램 구현(Server)

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ghodss/yaml"
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
	yamlFile, err := ioutil.ReadFile(*rootPath + *configFile)
	if err != nil {
		panic(err)
	}

	var config config
	err = yaml.Unmarshal(yamlFile, &config)

	if err != nil {
		panic(err)
	}

	test1 = config.toTest()

	//CSV 파일 읽어서 하는 걸로  && 통신으로 받는 걸로 변경하면 좋음
	array := []logInfo{
		decodeCSV("{time},operation,start-sw,on"),
		decodeCSV("{time},event,init"),
		decodeCSV("{time},operation,start-sw,off"),
	}

	for i := range array {
		test1.valid(array[i])
	}

	test1.toReport(*rootPath)
}

func decodeCSV(msg string) logInfo {
	s := strings.Split(msg, ",")
	return logInfo{
		unixDate:       s[0],
		utcNanoseconds: 0,
		index:          0,
		kind:           s[1],
		elements:       s[2:],
	}
}

func readConfig(path string) []test {
	var result []test
	//통신으로 받는 걸로 변경하면 좋음
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return result
	}
	for _, file := range files {
		if strings.Contains(file.Name(), ".yml") || strings.Contains(file.Name(), ".yaml") {
			yamlFile, err := ioutil.ReadFile(path + file.Name())

			var config config
			err = yaml.Unmarshal(yamlFile, &config)

			if err != nil {
				panic(err)
			}

			result = append(result, config.toTest())
		}
	}

	return result
}

type config struct {
	Items []config_item `json:"Items"`
}

type config_item struct {
	Type    string         `json:"Type"`
	Message string         `json:"Message"`
	Count   uint32         `json:"Count"`
	Order   []config_order `json:"Order"`
}

type config_order struct {
	Formats        config_wer `json:"Formats"`
	Status         config_wer `json:"Status"`
	Timeout_second uint16     `json:"Timeout_second"`
}

type config_wer struct {
	And   []config_wer `json:"And"`
	Or    []config_wer `json:"Or"`
	Not   []config_wer `json:"Not"`
	Value string       `json:"Value"`
}

func (c *config_wer) toItem() test_Item_wer {
	var test test_Item_wer
	for i := range c.And {
		test.And = append(test.And, c.And[i].toItem())
	}
	for i := range c.Or {
		test.Or = append(test.Or, c.Or[i].toItem())
	}
	for i := range c.Not {
		test.Not = append(test.Not, c.Not[i].toItem())
	}
	if c.Value != "" {
		test.Value = decodeCSV(c.Value)
	}

	return test
}

func (c *config) toTest() test {
	var test test
	test.last = map[string]logInfo{}
	for i := range c.Items {
		var item testItem
		item.resultType = c.Items[i].Type
		item.message = c.Items[i].Message
		item.count = c.Items[i].Count
		for _, order := range c.Items[i].Order {
			var testItem_order testItem_order
			testItem_order.formats = order.Formats.toItem()
			testItem_order.status_tests = order.Status.toItem()

			item.oders = append(item.oders, testItem_order)
		}
		test.config.items = append(test.config.items, item)
	}

	return test
}

type logInfo struct {
	unixDate       string
	utcNanoseconds int64
	index          uint32
	kind           string
	elements       []string
}

func (f *logInfo) valid() bool {
	if f.unixDate == "" {
		return false
	}
	if f.kind == "" {
		return false
	}
	if len(f.elements) == 0 {
		return false
	}
	for _, element := range f.elements {
		if element == "" {
			return false
		}
	}
	return true
}

func (format *logInfo) compare(info logInfo) bool {
	if format.valid() == false {
		return false
	}
	if info.valid() == false {
		return false
	}
	if format.kind != info.kind {
		return false
	}
	for i := range format.elements {
		if format.elements[i] != info.elements[i] {
			return false
		}
	}
	return true
}

type test struct {
	config   test_config
	result   test_result
	report   string
	reports  []string
	historys []string
	last     map[string]logInfo
}

//
// DateTime: 2021-07-29 02:10:41 +0900 KST #Unix Date Time Format
// _detail-path:
//   - "./SinarioName-{id}/TestName/report-testName-{id}.yaml"
//   - "./SinarioName-{id}/TestName/report-testName-{id}.yaml"

// report-testName-{id}.yaml
// Name: testName
// Id: "e22ghex"
// CheckerVersion: V1 #필요에 따라 여러종류의 Checker가 실행될 수 있음
// DateTime: 2021-07-29 02:10:41 +0900 KST #Unix Date Time Format
// Done: true
// _detail:
//   -
//     id: "28bac"
//     Type: Error
//     Message: abdcdd
//     filePath: "./SinarioName-{id}/TestName/report-검사방식이름/history-{id}.yaml"
// Error: [] #없으면 생략
// Warning: [""] #없으면 생략
// Message: ["",""] #없으면 생략
type test_config struct {
	items []testItem
}
type test_result struct {
	Error_m []string
	Warning []string
	Message []string
}

type testItem struct {
	//상태 들
	resultType  string
	message     string
	count       uint32
	match_count uint32
	oders       []testItem_order
}
type testItem_order struct {
	status_tests test_Item_wer
	formats      test_Item_wer
}

type test_Item_wer struct {
	And   []test_Item_wer
	Or    []test_Item_wer
	Not   []test_Item_wer
	Value logInfo
}

var test1 test

func (test *test) valid(info logInfo) {
	//상태 가 일치 (또는 무상관) & 현 Log가 일치
	for _, t := range test.config.items {
		if t.valid(test.last, info) {
			switch t.resultType {
			case "error":
				test.result.Error_m = append(test.result.Error_m, t.message)
				break
			case "warning":
				test.result.Warning = append(test.result.Warning, t.message)
				break
			case "message":
				test.result.Message = append(test.result.Message, t.message)
				break
			default:
				test.result.Error_m = append(test.result.Error_m, t.message)
				break
			}
		}
		test.last[info.kind] = info
	}
}

func (test *test) toReport(rootPath string) {
	// rootPath
	// ├── config.yaml
	// ├── report.yaml
	// ├── report-testName-{id}.yaml
	// ├── report-testName-{id}
	// │   ├── history-{id}.yaml
	// │   ├── history-{id}.yaml
	// │   ├── history-{id}.yaml
	// │   ├── history-{id}.yaml
	// ├── report-testName-{id}.yaml
	// ├── report-testName-{id}
	// │   ├── history-{id}.yaml
	// │   ├── history-{id}.yaml
	// │   ├── history-{id}.yaml
	// │   ├── history-{id}.yaml
	saveYamlFile(rootPath+"report.yaml", test.result)
	// saveYamlFile(rootPath+"Report-testName-{id}.yaml", test.result)
}

func saveYamlFile(filePath string, o interface{}) {
	b, err := yaml.Marshal(o)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(filePath, b, 0644)
	if err != nil {
		return
	}
}

func (item *testItem) valid(last map[string]logInfo, info logInfo) bool {
	return status_va1111(info, item.oders[0].formats) && status_va2222(last, item.oders[0].status_tests)
}

func status_va1111(info logInfo, formats test_Item_wer) bool {
	if len(formats.And) != 0 {
		for i := range formats.And {
			if status_va1111(info, formats.And[i]) == false {
				return false
			}
		}
		return true
	} else if len(formats.Or) != 0 {
		for i := range formats.Or {
			if status_va1111(info, formats.Or[i]) == true {
				return true
			}
		}
		return false
	} else if len(formats.Not) != 0 {
		for i := range formats.Not {
			if status_va1111(info, formats.Not[i]) == true {
				return false
			}
		}
		return true
	} else {
		if formats.Value.valid() == false {
			return false
		}

		return formats.Value.compare(info)
	}
}

func status_va2222(last map[string]logInfo, formats test_Item_wer) bool {
	if len(formats.And) != 0 {
		for i := range formats.And {
			if status_va2222(last, formats.And[i]) == false {
				return false
			}
		}
		return true
	} else if len(formats.Or) != 0 {
		for i := range formats.Or {
			if status_va2222(last, formats.Or[i]) == true {
				return true
			}
		}
		return false
	} else if len(formats.Not) != 0 {
		for i := range formats.Not {
			if status_va2222(last, formats.Not[i]) == true {
				return false
			}
		}
		return true
	} else {
		if formats.Value.valid() == false {
			return false
		}

		if kkk, ok := last[formats.Value.kind]; ok {
			return formats.Value.compare(kkk)
		}
		return false
	}
}
