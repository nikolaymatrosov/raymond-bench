package main

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"ses-templates/pkg/raymond"

	"github.com/rodaine/table"
)

type Template struct {
	TemplateName string `json:"TemplateName"`
	SubjectPart  string `json:"SubjectPart"`
	HtmlPart     string `json:"HtmlPart"`
	TextPart     string `json:"TextPart"`
}

type Example struct {
	Template Template `json:"Template"`
	Data     map[string]interface{}
}

func main() {
	// Open the JSON file
	//jsonFile, err := os.Open("example.json")
	//if err != nil {
	//	log.Fatalf("Failed to open file: %s", err)
	//}
	//defer jsonFile.Close()
	//
	//// Read the file content
	//byteValue, err := io.ReadAll(jsonFile)
	//if err != nil {
	//	log.Fatalf("Failed to read file: %s", err)
	//}
	//
	//// Unmarshal the JSON data
	//var example Example
	//if err := json.Unmarshal(byteValue, &example); err != nil {
	//	log.Fatalf("Failed to unmarshal JSON: %s", err)
	//}
	//
	//// Print the content
	//fmt.Printf("Template Name: %s\n", example.Template.TemplateName)
	//fmt.Printf("Subject Part: %s\n", example.Template.SubjectPart)
	//fmt.Printf("HTML Part: %s\n", example.Template.HtmlPart)
	//fmt.Printf("Text Part: %s\n", example.Template.TextPart)
	//fmt.Printf("\n\n")

	// Find inline partials
	//partials, content := findInlinePartials(example.Template.HtmlPart)
	//content := "{{#each names}}{{text}}{{/each}}"
	//tpl := raymond.MustParse(content)
	//for _, t := range partials {
	//	raymond.RegisterPartial(t.Name, t.Content)
	//}
	//tpl.SetMaxOutputSize(10)

	// Render the template
	//_, err = tpl.Exec(example.Data)
	//if err != nil {
	//	log.Fatalf("Failed to render template: %s", err)
	//}

	//fmt.Printf("Html:\n%s", result)
	fmt.Print("check is template is simple examples\n")
	templates := []string{
		"{{#each tags}}{{#each tags}}{{#each tags}}{{this}}{{/each}}{{/each}}{{/each}}",
		"{{#if foo}}{{#if bar}}{{foo}}{{/if}}{{/if}}",
		"{{foo}} {{bar}}",
		strings.Repeat("{{foo}} {{bar}}", 100),
	}
	for _, t := range templates {
		tpl := raymond.MustParse(t)
		fmt.Printf("%v	%s\n", tpl.IsSimple(), t)
	}
	fmt.Printf("\n\n")

	tbl := table.New("Duration", "Length", "Template")
	for _, t := range templates {
		start := time.Now()
		res, _ := render(t, map[string]any{
			"tags": numbers,
			"foo":  "foo",
			"bar":  "bar",
		})
		tbl.AddRow(time.Since(start), len(res), t)
	}
	tbl.Print()

}

func render(template string, data map[string]interface{}) (string, error) {
	// Find inline partials
	partials, content := findInlinePartials(template)
	tpl := raymond.MustParse(content)
	for _, t := range partials {
		raymond.RegisterPartial(t.Name, t.Content)
	}
	tpl.SetMaxOutputSize(10 * 1024 * 1024) // 10MB
	// Render the template
	result, err := tpl.Exec(data)
	if err != nil {
		return "", err
	}

	return result, nil
}

var numbers = genNumbers(100)

func genNumbers(n int) []string {
	var res []string
	for i := 0; i < n; i++ {
		res = append(res, strconv.Itoa(i))
	}
	return res
}

var longText = times(10, strings.Repeat("a", 1024))

func times(n int, s string) []map[string]any {
	result := make([]map[string]any, n)
	for i := range result {
		result[i] = map[string]any{"text": s}
	}
	return result
}

type LimitedBuffer struct {
	buf    *bytes.Buffer
	maxCap int
}

func NewLimitedBuffer(maxCap int) *LimitedBuffer {
	return &LimitedBuffer{
		buf:    bytes.NewBuffer(make([]byte, 0, maxCap)),
		maxCap: maxCap,
	}
}

func (lb *LimitedBuffer) Write(p []byte) (n int, err error) {
	if lb.buf.Len()+len(p) > lb.maxCap {
		return 0, errors.New("buffer capacity exceeded")
	}
	return lb.buf.Write(p)
}

func (lb *LimitedBuffer) WriteString(s string) (n int, err error) {
	if lb.buf.Len()+len(s) > lb.maxCap {
		return 0, errors.New("buffer capacity exceeded")
	}
	return lb.buf.WriteString(s)
}

func (lb *LimitedBuffer) String() string {
	return lb.buf.String()
}

// Find inline partials
// {{#\* inline "partialName"}}.*{{/inline~?}}
// rest of the content
type Partial struct {
	Name    string
	Content string
}

func findInlinePartials(content string) ([]Partial, string) {
	// Find all inline partials
	re := regexp.MustCompile(`{{#\*\s+inline\s+"(?P<name>.*?)"\s?}}(?P<content>.*?){{/\s*inline\s*~?}}\n`)
	matches := re.FindAllStringSubmatch(content, -1)

	// Create a slice of partials
	partials := make([]Partial, 0, len(matches))
	for _, match := range matches {
		partials = append(partials, Partial{
			Name:    match[1],
			Content: match[2],
		})
	}
	// Remove inline partials from the content
	content = re.ReplaceAllString(content, "")

	return partials, content
}
