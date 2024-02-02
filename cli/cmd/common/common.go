// SPDX-FileCopyrightText: 2024, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/sassoftware/event-provenance-registry/pkg/client"
	"github.com/sassoftware/event-provenance-registry/pkg/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func GetClient(url string) (*client.Client, error) {
	c, err := client.New(url)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func yaml2json(raw []byte) ([]byte, error) {
	var output interface{}
	if err := yaml.Unmarshal(raw, &output); err != nil {
		return nil, err
	}
	content, err := json.Marshal(output)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to yaml : %v", err)
	}
	return content, nil
}

// ToJSON converts to JSON
func ToJSON(data []byte) ([]byte, error) {
	if hasJSONPrefix(data) {
		return data, nil
	}
	return yaml2json(data)
}

func hasJSONPrefix(buf []byte) bool {
	if hasPrefix(buf, []byte("{")) {
		return true
	}
	if hasPrefix(buf, []byte("[")) {
		return true
	}
	return false
}

// IsJSONArray returns true if input is json array, false otherwise
func IsJSONArray(data []byte) bool {
	if len(data) > 0 && hasPrefix(data, []byte("[")) {
		return true
	}
	return false
}

func hasPrefix(buf []byte, prefix []byte) bool {
	trim := bytes.TrimLeftFunc(buf, unicode.IsSpace)
	return bytes.HasPrefix(trim, prefix)
}

func IndentJSON(str string) (string, error) {
	var indentJSON bytes.Buffer
	if err := json.Indent(&indentJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return indentJSON.String(), nil
}

// PrintSearchOutput common routine to output search responses,
// modified by jsonpath expression if given.
func PrintSearchOutput(content, expr string) error {
	cnt := []byte(content)
	v := interface{}(nil)

	var err error
	if expr != "" {
		err := fmt.Errorf("not implemented: %s", expr)
		return err
	}

	err = json.Unmarshal(cnt, &v)
	if err != nil {
		return err
	}

	cnt, err = json.MarshalIndent(v, "", "    ")
	if err != nil {
		return err
	}

	fmt.Println(string(cnt))

	return nil
}

// ProcessSearchFields handles expanding "all" if given,
// or returns defaults or specific list supplied by user.
func ProcessSearchFields(fields []string, i interface{}) ([]string, error) {
	if fields[0] == "all" {
		switch i.(type) {
		case *storage.Event, *storage.EventReceiver, *storage.EventReceiverGroup:
			fields = nil
			var matches []string
			re := regexp.MustCompile("json:\"(.*?)[,\"]")

			e := reflect.ValueOf(i).Elem()
			for i := 0; i < e.NumField(); i++ {
				matches = re.FindStringSubmatch(string(e.Type().Field(i).Tag))
				if len(matches) > 0 {
					if matches[1] != "-" {
						fields = append(fields, matches[1])
					}
				}
			}
		default:
			return nil, fmt.Errorf("ProcessSearchFields only handles Event, EventReceiver, EventReceiverGroup structs, but %T was passed", i)
		}
	}

	var trimmedFields []string
	for _, f := range fields {
		trimmedFields = append(trimmedFields, strings.Trim(f, ","))
	}

	return trimmedFields, nil
}

// BindFlagsE is run under the cobra command as preRunE. It simply binds the flags.
func BindFlagsE(cmd *cobra.Command, _ []string) error {
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}
	return nil
}
