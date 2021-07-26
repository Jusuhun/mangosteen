package main

//검사 최종 결과 보완 (Error 허용 개수 지정, Pass 조건 일치율)
//결과 남길 양식 및 폴더 트리 정하기 및 구현
//Line 받아오는 구조 구현 (Client)

//Log가 남는 PC에서 실행할 Line 보내주는 프로그램 구현(Server)

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ghodss/yaml"
)

func main() {
	tests = readConfig("./config/")

	//CSV 파일 읽어서 하는 걸로  && 통신으로 받는 걸로 변경하면 좋음
	array := []format{
		decodeCSV("{time},operation,start-sw,on"),
		decodeCSV("{time},event,init"),
		decodeCSV("{time},operation,start-sw,off"),
	}

	for i := range array {
		onUpdate(array[i])
	}

	fmt.Printf("end")
}

func decodeCSV(msg string) format {
	s := strings.Split(msg, ",")
	return format{s[0], s[1], s[2:]}
}

func onUpdate(info format) {
	for i := range tests {
		tests[i].valid(info)
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
	Name  string        `json:"Name"`
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
	test.last = map[string]format{}
	test.config.name = c.Name
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

type format struct {
	time     string
	kind     string
	elements []string
}

func (f *format) valid() bool {
	if f.time == "" {
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

func (format *format) compare(info format) bool {
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
	config test_config
	result test_result
	last   map[string]format
}
type test_config struct {
	name  string
	items []testItem
}
type test_result struct {
	error_m []string
	warning []string
	message []string
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
	Value format
}

var tests []test

func (test *test) valid(info format) {
	//상태 가 일치 (또는 무상관) & 현 Log가 일치
	for _, t := range test.config.items {
		if t.valid(test.last, info) {
			switch t.resultType {
			case "error":
				test.result.error_m = append(test.result.error_m, t.message)
				break
			case "warning":
				test.result.warning = append(test.result.warning, t.message)
				break
			case "message":
				test.result.message = append(test.result.message, t.message)
				break
			default:
				test.result.error_m = append(test.result.error_m, t.message)
				break
			}
		}
		test.last[info.kind] = info
	}
}

func (item *testItem) valid(last map[string]format, info format) bool {
	return status_va1111(info, item.oders[0].formats) && status_va2222(last, item.oders[0].status_tests)
}

func status_va1111(info format, formats test_Item_wer) bool {
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

func status_va2222(last map[string]format, formats test_Item_wer) bool {
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
