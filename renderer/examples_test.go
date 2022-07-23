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

func ExampleCidrHost_simple() {
	tmpl := `
{{ cidrHost 16 "10.12.127.0/20" }}
{{ cidrHost 268 "10.12.127.0/20" }}
{{ cidrHost 34 "fd00:fd12:3456:7890:00a2::/72" }}
{{ "10.12.127.0/20" | cidrHost 16 }}
`
	result, err := renderer.New(
		renderer.WithNetFunctions(),
	).Render(tmpl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
	// 10.12.112.16
	// 10.12.113.12
	// fd00:fd12:3456:7890::22
	// 10.12.112.16
}

func ExampleCidrHostEnd_simple() {
	tmpl := `
{{ cidrHostEnd 0 "10.12.127.0/20" }}
{{ cidrHostEnd 268 "10.12.127.0/20" }}
{{ cidrHostEnd 34 "fd00:fd12:3456:7890:00a2::/72" }}
{{ "10.12.127.0/20" | cidrHostEnd 16 }}
`
	result, err := renderer.New(
		renderer.WithNetFunctions(),
	).Render(tmpl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
  // 10.12.127.254
  // 10.12.126.242
  // fd00:fd12:3456:7890:ff:ffff:ffff:ffdc
  // 10.12.127.238
}

func ExampleCidrNetmask_simple() {
	tmpl := `
{{ cidrNetmask "10.0.0.0/12" }}
`
	result, err := renderer.New(
		renderer.WithNetFunctions(),
	).Render(tmpl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
	// 255.240.0.0
}

func ExampleCidrSubnets_simple() {
	tmpl := `
{{ index (cidrSubnets 2 "10.0.0.0/16") 0 }}
{{ index ("10.0.0.0/16" | cidrSubnets 2) 1 }}
{{ range cidrSubnets 3 "10.0.0.0/16" }}
{{ . }}{{ end }}
`
	result, err := renderer.New(
		renderer.WithNetFunctions(),
	).Render(tmpl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
	// 10.0.0.0/18
	// 10.0.64.0/18
	//
	// 10.0.0.0/19
	// 10.0.32.0/19
	// 10.0.64.0/19
	// 10.0.96.0/19
	// 10.0.128.0/19
	// 10.0.160.0/19
	// 10.0.192.0/19
	// 10.0.224.0/19
}

func ExampleCidrSubnetSizes_simple() {
	tmpl := `
{{ range $k, $v := cidrSubnetSizes 4 4 8 4 "10.1.0.0/16" }}
{{ $k }} {{ $v }}{{ end }}
{{ range $k, $v := cidrSubnetSizes 16 16 16 32 "fd00:fd12:3456:7890::/56" }}
{{ $k }} {{ $v }}{{ end }}
`
	result, err := renderer.New(
		renderer.WithNetFunctions(),
	).Render(tmpl)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
	// 0 10.1.0.0/20
	// 1 10.1.16.0/20
	// 2 10.1.32.0/24
	// 3 10.1.48.0/20
	//
	// 0 fd00:fd12:3456:7800::/72
	// 1 fd00:fd12:3456:7800:100::/72
	// 2 fd00:fd12:3456:7800:200::/72
	// 3 fd00:fd12:3456:7800:300::/88
}
