package main

import (
	"flag"
	"fmt"
	"os"

	node_exporter_main "github.com/chaolihf/node_exporter"
	jjson "github.com/chaolihf/udpgo/json"
	fileBeatCmd "github.com/elastic/beats/v7/filebeat/cmd"
	inputs "github.com/elastic/beats/v7/filebeat/input/default-inputs"
	"github.com/spf13/pflag"
)

const (
	Module_File = "filebeat"
	Module_Node = "node_exporter"
)

/*
根据配置文件激活不同的模块
*/
func activeModules() map[string]bool {
	var oneAgentModules = make(map[string]bool)
	oneAgentModules[Module_File] = false
	oneAgentModules[Module_Node] = false
	filePath := "config.json"
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("读取文件%s出错:%s\n", filePath, err.Error())
	} else {
		jsonConfigInfos, err := jjson.NewJsonObject([]byte(content))
		if err != nil {
			fmt.Printf("文件%sJson格式出错:%s\n", filePath, err.Error())
		} else {
			jsonModuleInfos := jsonConfigInfos.GetJsonArray("module")
			for _, jsonModuleInfo := range jsonModuleInfos {
				oneAgentModules[jsonModuleInfo.GetStringValue()] = true
			}
		}
	}
	for moduleName := range oneAgentModules {
		fmt.Printf("激活模块%s\n", moduleName)
	}
	return oneAgentModules
}

/*
开始应用
*/
func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("程序退出原因:", r)
		}
	}()
	done := make([]chan string, 1)
	oneAgentModules := activeModules()
	if oneAgentModules[Module_File] {
		go runFileBeat(oneAgentModules, done[0])
	}
	if !oneAgentModules[Module_Node] {
		node_exporter_main.Main()
	} else {
		for _, moduleDone := range done {
			v := <-moduleDone
			fmt.Println("run finish ", v)
		}
	}
}

/*
创建包含Node_Export相关的参数
*/
func addNodeExporterConfigFileFlag(runFlags *pflag.FlagSet) {
	findFlag := flag.CommandLine.Lookup("web.config.file")
	if findFlag == nil {
		flag.String("web.config.file", "", "[EXPERIMENTAL] Path to configuration file that can enable TLS or authentication. See: https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md")
	}
	runFlags.AddGoFlag(flag.CommandLine.Lookup("web.config.file"))
}

func addNodeExporterListenAddressFlag(runFlags *pflag.FlagSet) {
	findFlag := flag.CommandLine.Lookup("web.listen-address")
	if findFlag == nil {
		flag.String("web.listen-address", ":9100", "Addresses on which to expose metrics and web interface. Repeatable for multiple addresses.")
	}
	runFlags.AddGoFlag(flag.CommandLine.Lookup("web.listen-address"))
}

/*
运行FileBeat
*/
func runFileBeat(activeModules map[string]bool, doneChain chan<- string) {
	settings := fileBeatCmd.FilebeatSettings()
	if activeModules[Module_Node] {
		addNodeExporterConfigFileFlag(settings.RunFlags)
		addNodeExporterListenAddressFlag(settings.RunFlags)
	}
	if err := fileBeatCmd.Filebeat(inputs.Init, settings).Execute(); err != nil {
		doneChain <- " file error " + err.Error()
	} else {
		doneChain <- " file done "
	}
}
