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
				ID   string `yaml:"id"`
				Type string `yaml:"type"`
				Tag  string `yaml:"tag"`
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
	Tag       string
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
		var fields []FieldInfo

		for _, format := range packet.Format {
			field := FieldInfo{
				Name:      strings.Title(format.ID),
				GoType:    determineGoType(format.Type),
				YAMLField: format.ID,
				Tag:       format.Tag,
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

func determineGoType(typeString string) string {
	switch typeString {
	case "uint16":
		return "uint16"
	case "[]byte":
		return "[]byte"
	case "byte_array":
		return "[]byte"
	// Add more cases as needed
	default:
		return "uint8"
	}
}

func generateGoCode(structs []GeneratedStruct) string {
	var code strings.Builder

	code.WriteString(fmt.Sprintf("package %s\n\n", structs[0].PackageName))
	code.WriteString("import  github.com/ghostiam/binstruct \n\n")

	for _, generatedStruct := range structs {
		for _, structInfo := range generatedStruct.Structs {
			code.WriteString(fmt.Sprintf("type %s struct {\n", structInfo.StructName))

			for _, field := range structInfo.Fields {
				code.WriteString(fmt.Sprintf("\t%s %s `yaml:\"%s\" %s`\n", field.Name, field.GoType, field.YAMLField, field.Tag))
			}

			code.WriteString("}\n\n")
			code.WriteString(generateEncodeDecodeCode(structInfo.StructName))
		}
	}

	return code.String()
}


func generateEncodeDecodeCode(structname string) string {
	var code strings.Builder
	code.WriteString(fmt.Sprintf("func (*%s p) Encode()([]byte,error) {\n", structname))
	code.WriteString("\tvar bin []byte(\n")
	code.WriteString("\terr := binary.Write(bin,binary.BigEndian,p)\n")
	code.WriteString("}\n\n")
	code.WriteString(fmt.Sprintf("func Decode_%s(bindata []byte)(%s,error) {\n", structname,structname))
	code.WriteString(fmt.Sprintf("\tvar st %s\n",structname))
	code.WriteString("\terr := binstruct.UnmarshalBE(bindata,&st)\n")
	code.WriteString("\treturn st,err\n")
	code.WriteString("}\n\n")
	return code.String()

}
