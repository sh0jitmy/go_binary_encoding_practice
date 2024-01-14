//go:generate go run generate.go

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"gopkg.in/yaml.v2"
)

type ProtocolInfo struct {
	Protocol struct {
		Path       string `yaml:"path"`
		PacketList []struct {
			Msg    string `yaml:"msg"`
			Format []struct {
				ID       string `yaml:"id"`
				Bytesize int    `yaml:"bytesize"`
			} `yaml:"format"`
		} `yaml:"packetlist"`
	} `yaml:"Protocol"`
}

type GeneratedStruct struct {
	PackageName string
	Structs     []StructInfo
}

type StructInfo struct {
	StructName string
	Fields     []FieldInfo
}

type FieldInfo struct {
	Name      string
	GoType    string
	YAMLField string
}

func main() {
	yamlFile, err := ioutil.ReadFile("protocol.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	protocolInfo := ProtocolInfo{}
	err = yaml.Unmarshal(yamlFile, &protocolInfo)
	if err != nil {
		log.Fatalf("Error parsing YAML: %v", err)
	}

	generatedStructs := generateStructsFromYAML(protocolInfo)
	goCode := generateGoCode(generatedStructs)

	outputFile := "generated_code.go"
	err = ioutil.WriteFile(outputFile, []byte(goCode), 0644)
	if err != nil {
		log.Fatalf("Error writing Go code to file: %v", err)
	}

	fmt.Println("Code generation complete.")
}

func generateStructsFromYAML(protocolInfo ProtocolInfo) []GeneratedStruct {
	var structs []GeneratedStruct

	for _, packet := range protocolInfo.Protocol.PacketList {
		structName := strings.Title(packet.Msg)
		//structName := "C_" + strings.Title(packet.Msg)
		var fields []FieldInfo

		for _, format := range packet.Format {
			field := FieldInfo{
				Name:      strings.Title(format.ID),
				GoType:    determineGoType(format.Bytesize),
				YAMLField: format.ID,
			}
			fields = append(fields, field)
		}

		structs = append(structs, GeneratedStruct{
			PackageName: protocolInfo.Protocol.Path,
			Structs: []StructInfo{
				{
					StructName: structName,
					Fields:     fields,
				},
			},
		})
	}

	return structs
}

func determineGoType(bytesize int) string {
	switch bytesize {
	case 1:
		return "uint8"
	case 2:
		return "uint16"
	case 4:
		return "uint32"
	// Add more cases as needed
	default:
		return "uint8"
	}
}

func generateGoCode(structs []GeneratedStruct) string {
	var code strings.Builder

	code.WriteString(fmt.Sprintf("package %s\n\n", structs[0].PackageName))

	//struct define
	for _, generatedStruct := range structs {
		for _, structInfo := range generatedStruct.Structs {
			code.WriteString(fmt.Sprintf("type %s struct {\n", structInfo.StructName))

			for _, field := range structInfo.Fields {
				code.WriteString(fmt.Sprintf("\t%s %s `yaml:\"%s\"`\n", field.Name, field.GoType, field.YAMLField))
			}

			code.WriteString("}\n\n")
		}
	}
	//function define
	for _, generatedStruct := range structs {
		for _, structInfo := range generatedStruct.Structs {
			code.WriteString(fmt.Sprintf("func (p %s) Encode(b [] byte)(int, error){ \n \n}\n\n",structInfo.StructName))
			code.WriteString(fmt.Sprintf("func (p %s) Decode()([] byte, error){ \n \n}\n\n",structInfo.StructName))
		}
	}	

	return code.String()
}

