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
命令行参数信息
*/
type CommandInfo struct {
	name        string
	valueType   string
	value       string
	description string
}

/*
根据配置文件激活不同的模块，获取命令行信息
*/
func activeInfos() (map[string]bool, []CommandInfo) {
	var oneAgentModules = make(map[string]bool)
	var commandInfos []CommandInfo
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
			jsonCommandInfos := jsonConfigInfos.GetJsonArray("command")
			for _, jsonCommandInfo := range jsonCommandInfos {
				commandInfos = append(commandInfos, CommandInfo{
					name:        jsonCommandInfo.GetString("name"),
					valueType:   jsonCommandInfo.GetString("type"),
					value:       jsonCommandInfo.GetString("value"),
					description: jsonCommandInfo.GetString("description"),
				})
			}
		}
	}
	for moduleName := range oneAgentModules {
		fmt.Printf("激活模块%s\n", moduleName)
	}
	return oneAgentModules, commandInfos
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
	oneAgentModules, oneAgentCommandLines := activeInfos()
	if oneAgentModules[Module_File] {
		go runFileBeat(oneAgentModules, done[0], oneAgentCommandLines)
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
func addNodeExporterFlag(runFlags *pflag.FlagSet, commandInfos []CommandInfo) {
	for _, commandInfo := range commandInfos {
		findFlag := flag.CommandLine.Lookup(commandInfo.name)
		if findFlag == nil {
			switch commandInfo.valueType {
			case "string":
				{
					flag.String(commandInfo.name, commandInfo.value, commandInfo.description)
					break
				}
			case "boolean":
				{
					flag.Bool(commandInfo.name, "true" == commandInfo.value, commandInfo.description)
					break
				}
			}

		}
		runFlags.AddGoFlag(flag.CommandLine.Lookup(commandInfo.name))
	}
}

// func addNodeExporterConfigFileFlag(runFlags *pflag.FlagSet) {
// 	findFlag := flag.CommandLine.Lookup("web.config.file")
// 	if findFlag == nil {
// 		flag.String("web.config.file", "", "[EXPERIMENTAL] Path to configuration file that can enable TLS or authentication. See: https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md")
// 	}
// 	runFlags.AddGoFlag(flag.CommandLine.Lookup("web.config.file"))
// }

// func addNodeExporterListenAddressFlag(runFlags *pflag.FlagSet) {
// 	findFlag := flag.CommandLine.Lookup("web.listen-address")
// 	if findFlag == nil {
// 		flag.String("web.listen-address", ":9100", "Addresses on which to expose metrics and web interface. Repeatable for multiple addresses.")
// 	}
// 	runFlags.AddGoFlag(flag.CommandLine.Lookup("web.listen-address"))
// }

/*
运行FileBeat
*/
func runFileBeat(activeModules map[string]bool, doneChain chan<- string, commandInfos []CommandInfo) {
	settings := fileBeatCmd.FilebeatSettings()
	if activeModules[Module_Node] {
		addNodeExporterFlag(settings.RunFlags, commandInfos)
	}
	if err := fileBeatCmd.Filebeat(inputs.Init, settings).Execute(); err != nil {
		doneChain <- " file error " + err.Error()
	} else {
		doneChain <- " file done "
	}
}
