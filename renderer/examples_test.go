package renderer_test

import (
	"fmt"

	"github.com/VirtusLab/render/renderer"
	"github.com/VirtusLab/render/renderer/parameters"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func ExampleJSONPath_simple() {
	json := `{
	"welcome":{
		"message":["Good Morning", "Hello World!"]
	}
}`
	expression := "{$.welcome.message[1]}"

	params := parameters.Parameters{
		"json":       json,
		"expression": expression,
	}

	tmpl := `{{ .json | fromJson | jsonPath .expression }}`

	result, err := renderer.New(
		renderer.WithParameters(params),
		renderer.WithExtraFunctions(),
	).Render(tmpl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

	// Output:
	// Hello World!
}

func ExampleJSONPath_array() {
	json := `["Good Morning", "Hello World!"]`
	expression := "{$[1]}"

	params := parameters.Parameters{
		"json":       json,
		"expression": expression,
	}

	tmpl := `
{{ .json | fromJson | jsonPath .expression }}
`

	result, err := renderer.New(
		renderer.WithParameters(params),
		renderer.WithExtraFunctions(),
	).Render(tmpl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

	// Output:
	// Hello World!
}

func ExampleJSONPath_wildcard() {
	json := `{
	"welcome":{
		"message":["Good Morning", "Hello World!"]
	}
}`
	expression := "{$.welcome.message[*]}"

	params := parameters.Parameters{
		"json":       json,
		"expression": expression,
	}

	tmpl := `
{{- range $m := .json | fromJson | jsonPath .expression }}
{{ $m }} 
{{- end }}
`

	result, err := renderer.New(
		renderer.WithParameters(params),
		renderer.WithExtraFunctions(),
	).Render(tmpl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

	// Output:
	// Good Morning
	// Hello World!
}

func ExampleJSONPath_yaml() {
	yaml := `---
welcome:
  message:
    - "Good Morning"
    - "Hello World!"
`
	expression := "{$.welcome.message[1]}"

	params := parameters.Parameters{
		"yaml":       yaml,
		"expression": expression,
	}

	tmpl := `{{ .yaml | fromYaml | jsonPath .expression }}`

	result, err := renderer.New(
		renderer.WithParameters(params),
		renderer.WithExtraFunctions(),
	).Render(tmpl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

	// Output:
	// Hello World!
}

func ExampleJSONPath_multi() {
	yaml := `---
data:
  en: "Hello World!"
---
data:
  pl: "Witaj Świecie!"
`
	expression := "{$[*].data}"

	params := parameters.Parameters{
		"yaml":       yaml,
		"expression": expression,
	}

	tmpl := `
{{- range $m := .yaml | fromYaml | jsonPath .expression }}
{{ range $k, $v := $m }}{{ $k }} {{ $v }}{{ end }}
{{- end }}
`

	result, err := renderer.New(
		renderer.WithParameters(params),
		renderer.WithExtraFunctions(),
	).Render(tmpl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

	// Output:
	// en Hello World!
	// pl Witaj Świecie!
}

func ExampleJSONPath_multi_nested() {
	yamlWithJson := `---
data:
  en.json: |2
    {
      "welcome":{
        "message":["Good Morning", "Hello World!"]
      }
    }
---
data:
  pl.json: |2
    {
      "welcome":{
        "message":["Dzień dobry", "Witaj Świecie!"]
      }
    }
`
	expression := "{$[*].data.*}"
	expression2 := "{$.welcome.message[1]}"

	params := parameters.Parameters{
		"yamlWithJson": yamlWithJson,
		"expression":   expression,
		"expression2":  expression2,
	}

	tmpl := `
{{- range $r := .yamlWithJson | fromYaml | jsonPath .expression }}
{{ $r | fromJson | jsonPath $.expression2 }}
{{- end }}
`
	result, err := renderer.New(
		renderer.WithParameters(params),
		renderer.WithExtraFunctions(),
	).Render(tmpl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

	// Output:
	// Hello World!
	// Witaj Świecie!
}

func ExampleN_simple() {
	tmpl := `
{{ range $i := n 0 10 }}{{ $i }} {{ end }}
`
	result, err := renderer.New(
		renderer.WithExtraFunctions(),
	).Render(tmpl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

	// Output:
	// 0 1 2 3 4 5 6 7 8 9 10
}

func ExampleN_empty() {
	tmpl := `
{{ range $i := n 0 0 }}{{ $i }} {{ end }}
`
	result, err := renderer.New(
		renderer.WithExtraFunctions(),
	).Render(tmpl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)

	// Output:
	// 0
}
