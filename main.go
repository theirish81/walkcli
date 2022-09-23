package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TylerBrock/colorjson"
	"github.com/pborman/getopt/v2"
	"github.com/theirish81/gowalker"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"strings"
)

// main function
func main() {
	template, beautifyJSON, beautifyYAML, color, err := ParseParams()
	if err != nil {
		getopt.PrintUsage(os.Stdout)
		os.Exit(1)
	}
	template, subTemplates, err := TransformTemplate(template)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: error loading template ", err)
		os.Exit(1)
	}

	// loading from piped standard input
	if data, err := ReadPipe(); err == nil {
		if output, err := Render(template, subTemplates, beautifyJSON, beautifyYAML, color, data); err == nil {
			fmt.Println(string(output))
		} else {
			_, _ = fmt.Fprintln(os.Stderr, "ERROR: ", err)
		}
	} else {
		_, _ = fmt.Fprintln(os.Stderr, "ERROR: ", err)
		os.Exit(1)
	}
}

// ParseParams parses the command line string
func ParseParams() (string, bool, bool, bool, error) {
	template := getopt.StringLong("template", 't', "", "The template")
	beautifyJSON := getopt.BoolLong("json-beautify", 'j', "false", "Beautify JSON output")
	beautifyYAML := getopt.BoolLong("beautify-yaml", 'y', "false", "Beautify YAML output")
	color := getopt.BoolLong("color", 'c', "false", "Colored output (JSON)")
	getopt.Parse()
	if *template == "" {
		return "", false, false, false, errors.New("template is absent")
	}
	return *template, *beautifyJSON, *beautifyYAML, *color, nil
}

// TransformTemplate will decide whether the provided string is a template or a file. If it's a template, then it will
// simply return it. If it's a file, it will load it as main template, load all other files in the directory as
// sub-templates, and return everything
func TransformTemplate(template string) (string, gowalker.SubTemplates, error) {
	subTemplates := gowalker.NewSubTemplates()
	if strings.HasPrefix(template, "file://") {
		// transforming the string to a path by removing the file:// prefix
		template = strings.Replace(template, "file://", "", 1)
		var err error
		// loading all sub-templates
		if template, subTemplates, err = gowalker.LoadTemplatesFromDisk(template); err != nil {
			return "", subTemplates, err
		}
	}
	return template, subTemplates, nil
}

// Render will use the templates and transform the input accordingly
func Render(template string, subTemplates gowalker.SubTemplates, beautifyJSON bool, beautifyYAML bool, color bool, data []byte) ([]byte, error) {
	var output []byte
	var str interface{}
	err := yaml.Unmarshal(data, &str)
	if err != nil {
		return nil, errors.New("unable to parse input: " + err.Error())
	}
	// Render everything
	if res, err := gowalker.RenderAll(context.Background(), template, subTemplates, str, gowalker.NewFunctions()); err == nil {
		output = []byte(res)
		// if beautifyJSON, then we try and beautify the output as a JSON
		if beautifyJSON {
			err = json.Unmarshal([]byte(res), &str)
			if err != nil {
				return nil, errors.New("output string is not valid JSON: " + err.Error())
			}
			// if color is true, format the output in JSON-color
			if color {
				f := colorjson.NewFormatter()
				f.Indent = 2
				output, err = f.Marshal(str)
			} else {
				// otherwise marshal in JSON without colors
				output, err = json.MarshalIndent(str, "", "  ")
			}
			if err != nil {
				return nil, errors.New("could not unmarshal valid JSON: " + err.Error())
			}
		}
		// if beautifyYAML, then we try and beautify the output as a YAML
		if beautifyYAML {
			err = yaml.Unmarshal([]byte(res), &str)
			if err != nil {
				return nil, errors.New("output string is not valid YAML: " + err.Error())
			}
			output, err = yaml.Marshal(str)
			if err != nil {
				return nil, errors.New("could not unmarshal valid YAML: " + err.Error())
			}
		}
		return output, nil
	} else {
		return nil, err
	}
}

// ReadPipe reads the input from the pipe and returns a byte array
func ReadPipe() ([]byte, error) {
	nBytes, nChunks := int64(0), int64(0)
	r := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0, 4*1024)
	data := make([]byte, 0)
	for {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			return data, err
		}
		nChunks++
		nBytes += int64(len(buf))
		data = append(data, buf...)
		if err != nil && err != io.EOF {
			return data, err
		}
	}
	return data, nil
}
