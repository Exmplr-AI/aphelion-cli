package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

var (
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	warningColor = color.New(color.FgYellow, color.Bold)
	infoColor    = color.New(color.FgBlue, color.Bold)
)

func PrintSuccess(format string, args ...interface{}) {
	successColor.Printf("✓ "+format+"\n", args...)
}

func PrintError(format string, args ...interface{}) {
	errorColor.Printf("✗ "+format+"\n", args...)
}

func PrintWarning(format string, args ...interface{}) {
	warningColor.Printf("⚠ "+format+"\n", args...)
}

func PrintInfo(format string, args ...interface{}) {
	infoColor.Printf("ℹ "+format+"\n", args...)
}

func PrintOutput(data interface{}, format string) error {
	switch format {
	case "json":
		return printJSON(data)
	case "yaml":
		return printYAML(data)
	case "table":
		return printTable(data)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

func printJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func printYAML(data interface{}) error {
	encoder := yaml.NewEncoder(os.Stdout)
	defer encoder.Close()
	return encoder.Encode(data)
}

func printTable(data interface{}) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	switch v := data.(type) {
	case map[string]interface{}:
		table.SetHeader([]string{"Key", "Value"})
		for key, value := range v {
			table.Append([]string{key, fmt.Sprintf("%v", value)})
		}
	case []interface{}:
		if len(v) == 0 {
			fmt.Println("No data to display")
			return nil
		}
		
		first := v[0]
		if firstMap, ok := first.(map[string]interface{}); ok {
			var headers []string
			for key := range firstMap {
				headers = append(headers, strings.Title(key))
			}
			table.SetHeader(headers)
			
			for _, item := range v {
				if itemMap, ok := item.(map[string]interface{}); ok {
					var row []string
					for _, header := range headers {
						key := strings.ToLower(header)
						if value, exists := itemMap[key]; exists {
							row = append(row, fmt.Sprintf("%v", value))
						} else {
							row = append(row, "")
						}
					}
					table.Append(row)
				}
			}
		}
	case []map[string]interface{}:
		if len(v) == 0 {
			fmt.Println("No data to display")
			return nil
		}
		
		var headers []string
		for key := range v[0] {
			headers = append(headers, strings.Title(key))
		}
		table.SetHeader(headers)
		
		for _, item := range v {
			var row []string
			for _, header := range headers {
				key := strings.ToLower(header)
				if value, exists := item[key]; exists {
					row = append(row, fmt.Sprintf("%v", value))
				} else {
					row = append(row, "")
				}
			}
			table.Append(row)
		}
	default:
		val := reflect.ValueOf(data)
		if val.Kind() == reflect.Slice {
			if val.Len() == 0 {
				fmt.Println("No data to display")
				return nil
			}
			
			first := val.Index(0)
			if first.Kind() == reflect.Struct {
				typ := first.Type()
				var headers []string
				for i := 0; i < typ.NumField(); i++ {
					field := typ.Field(i)
					if field.IsExported() {
						headers = append(headers, field.Name)
					}
				}
				table.SetHeader(headers)
				
				for i := 0; i < val.Len(); i++ {
					item := val.Index(i)
					var row []string
					for j := 0; j < typ.NumField(); j++ {
						field := typ.Field(j)
						if field.IsExported() {
							value := item.Field(j)
							row = append(row, fmt.Sprintf("%v", value.Interface()))
						}
					}
					table.Append(row)
				}
			}
		} else {
			return printJSON(data)
		}
	}

	table.Render()
	return nil
}

func OutputJSON(data interface{}) error {
	return printJSON(data)
}

func OutputYAML(data interface{}) error {
	return printYAML(data)
}

func OutputTable(data interface{}) error {
	return printTable(data)
}