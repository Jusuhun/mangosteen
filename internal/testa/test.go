package testa

import (
	"io/ioutil"
	"log"
	"mangosteena/pkg/idgen"
	"mangosteena/pkg/info"
	"os"

	"github.com/ghodss/yaml"
)

type Test struct {
	methods []Method
	last    map[string]info.LogInfo
}

func (test *Test) Valid(info info.LogInfo) {
	//상태 가 일치 (또는 무상관) & 현 Log가 일치
	for i := range test.methods {
		if test.methods[i].valid(test.last, info) {
			test.methods[i].addHistory(test.last, info)
		}
		test.last[info.Kind] = info
	}
}

func (test *Test) WriteReport(rootPath string, version string, build string) {
	var IDGen idgen.IdGenerater
	var report struct {
		DateTime         string   `json:DateTime`
		CheckerVersion   string   `json:CheckerVersion`
		CheckerBuildDate string   `json:CheckerBuildDate`
		DetailPath       []string `json:_detail-path`
	}
	for _, method := range test.methods {
		var HistoryIDGen idgen.IdGenerater
		type subReport_detail struct {
			Id       string `json:Id`
			Type     string `json:Type`
			Message  string `json:Message`
			FilePath string `json:FilePath`
		}
		var subReport struct {
			Name     string             `json:Name`
			Id       string             `json:Id`
			DateTime string             `json:DateTime`
			Done     bool               `json:Done`
			Detail   []subReport_detail `json:_detail`
		}
		subReport.Id = IDGen.GenerateID()
		subReport.Name = method.name

		for i := range method.historys {
			id := HistoryIDGen.GenerateID()
			detail := subReport_detail{
				Id:       id,
				Type:     method.historys[i].resultType,
				Message:  method.historys[i].message,
				FilePath: "./report-" + subReport.Name + "-" + subReport.Id + "/history-" + id + ".yaml",
			}

			subReport.Detail = append(subReport.Detail, detail)

			var history struct {
				Id          string `json:Id`
				ResultType  string
				Message     string
				MaxCount    uint32
				Match_count uint32
				Info        info.LogInfo
				History     map[string]info.LogInfo
			}

			history.Id = detail.Id
			history.ResultType = method.historys[i].resultType
			history.Message = method.historys[i].message
			history.MaxCount = method.historys[i].maxCount
			history.Match_count = method.historys[i].match_count
			history.Info = method.historys[i].info
			history.History = method.historys[i].History
			//더욱 자세한 내용을 구조체로 만든다.
			saveYamlFile(rootPath+"report-"+subReport.Name+"-"+subReport.Id+"/", "history-"+id+".yaml", history)
		}

		subReport.Done = false
		saveYamlFile(rootPath, "report-"+subReport.Name+"-"+subReport.Id+".yaml", subReport)

		report.DetailPath = append(report.DetailPath, "./report-"+subReport.Name+"-"+subReport.Id+".yaml")
	}

	report.CheckerVersion = version
	report.CheckerBuildDate = build
	report.DateTime = "2021-07-29 02:10:41 +0900 KST"
	saveYamlFile(rootPath, "report.yaml", report)
}

func saveYamlFile(filePath string, fileName string, o interface{}) {
	//Create a folder/directory at a full qualified path
	err := os.MkdirAll(filePath, 0755)
	if err != nil {
		log.Fatal(err)
	}

	b, err := yaml.Marshal(o)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(filePath+fileName, b, 0644)
	if err != nil {
		return
	}
}
