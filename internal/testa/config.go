package testa

import (
	"io/ioutil"
	"mangosteena/pkg/info"

	"github.com/ghodss/yaml"
)

func ReadFile(path string) Test {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var config config
	err = yaml.Unmarshal(yamlFile, &config)

	if err != nil {
		panic(err)
	}

	return config.toTest()
}

type config struct {
	Items []config_item `json:"Items"`
}

type config_item struct {
	Type     string         `json:"Type"`
	Message  string         `json:"Message"`
	Name     string         `json:"Name"`
	MaxCount uint32         `json:"MaxCount"`
	Order    []config_order `json:"Order"`
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

func (c *config_wer) toItem() Conditional {
	var con Conditional
	for i := range c.And {
		con.And = append(con.And, c.And[i].toItem())
	}
	for i := range c.Or {
		con.Or = append(con.Or, c.Or[i].toItem())
	}
	for i := range c.Not {
		con.Not = append(con.Not, c.Not[i].toItem())
	}
	if c.Value != "" {
		con.Value = info.DecodeCSV(c.Value)
	}

	return con
}

func (c *config) toTest() Test {
	var test Test
	test.last = map[string]info.LogInfo{}
	for i := range c.Items {
		var method Method
		method.resultType = c.Items[i].Type
		method.message = c.Items[i].Message
		method.name = c.Items[i].Name
		method.maxCount = c.Items[i].MaxCount
		for _, order := range c.Items[i].Order {
			var method_order Method_order
			method_order.cond = order.Formats.toItem()
			method_order.status_tests = order.Status.toItem()

			method.oders = append(method.oders, method_order)
		}
		test.methods = append(test.methods, method)
	}

	return test
}
