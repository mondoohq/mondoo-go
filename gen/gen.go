// This implementation is inspired from https://github.com/shurcooL/githubv4/blob/master/gen.go

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/shurcooL/graphql/ident"
)

func main() {
	flag.Parse()

	token, ok := os.LookupEnv("MONDOO_API_TOKEN")
	if !ok {
		log.Fatalln(fmt.Errorf("MONDOO_API_TOKEN environment variable not set"))
	}
	err := generateSchema(token, ".")
	if err != nil {
		log.Fatalln(err)
	}
}

// generateSchema generates the mondoogql package in basePath.
func generateSchema(token string, basePath string) error {
	// fetch the graphql schema
	schema, err := loadSchema(token)
	if err != nil {
		return err
	}

	// write the generated files
	for filename, t := range templates {
		var buf bytes.Buffer
		// render go file
		err := t.Execute(&buf, schema)
		if err != nil {
			return err
		}

		// format go file with gofmt
		out, err := format.Source(buf.Bytes())
		if err != nil {
			log.Println(err)
			out = []byte("// gofmt error: " + err.Error() + "\n\n" + buf.String())
		}

		// write generated go file
		outfile := filepath.Join(basePath, filename)
		fmt.Println("writing", outfile)
		err = os.WriteFile(outfile, out, 0o644)
		if err != nil {
			return err
		}
	}

	return nil
}

// loadSchema loads the GraphQL schema from the Mondoo API.
func loadSchema(token string) (schema interface{}, err error) {
	apiHost := "us.api.mondoo.com"
	introspection := `
{
  __schema {
    queryType {
      name
    }
    mutationType {
      name
    }
    subscriptionType {
      name
    }
    types {
      ...FullType
    }
    directives {
      name
      description
      locations
      args {
        ...InputValue
      }
    }
  }
}

fragment FullType on __Type {
  kind
  name
  description
  fields(includeDeprecated: true) {
    name
    description
    args {
      ...InputValue
    }
    type {
      ...TypeRef
    }
    isDeprecated
    deprecationReason
  }
  inputFields {
    ...InputValue
  }
  interfaces {
    ...TypeRef
  }
  enumValues(includeDeprecated: true) {
    name
    description
    isDeprecated
    deprecationReason
  }
  possibleTypes {
    ...TypeRef
  }
}

fragment InputValue on __InputValue {
  name
  description
  type {
    ...TypeRef
  }
  defaultValue
}

fragment TypeRef on __Type {
  kind
  name
  ofType {
    kind
    name
    ofType {
      kind
      name
      ofType {
        kind
        name
        ofType {
          kind
          name
          ofType {
            kind
            name
            ofType {
              kind
              name
              ofType {
                kind
                name
              }
            }
          }
        }
      }
    }
  }
}
`
	// do introspection query
	req, err := http.NewRequest(
		"POST",
		"https://"+apiHost+"/query",
		strings.NewReader(`{"query":`+strconv.Quote(introspection)+`}`),
	)
	if err != nil {
		return nil, err
	}

	// set headers
	req.Header.Set("Authorization", "bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Host", apiHost)
	req.Header.Set("Origin", "https://"+apiHost)

	// send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-200 OK status code: %v body: %q", resp.Status, body)
	}

	// decode graphql schema
	err = json.NewDecoder(resp.Body).Decode(&schema)
	return schema, err
}

// templates maps generated filenames with the go template to use for them.
var templates = map[string]*template.Template{
	// generate all enum types
	"enum.go": parseTemplate(`// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: MPL-2.0
// 
// Code generated by gen.go; DO NOT EDIT.

package mondoogql
{{range .data.__schema.types | sortByName}}{{if and (eq .kind "ENUM") (not (internal .name))}}
{{template "enum" .}}
{{end}}{{end}}

{{- define "enum" -}}
// {{.name}}
type {{.name}} string

const ({{range .enumValues}}
	{{enumIdentifier $.name .name}} {{$.name}} = {{.name | quote}} {{end}}
)
{{- end -}}
`),

	// generate all input types
	"input.go": parseTemplate(`// Copyright (c) Mondoo, Inc.
// SPDX-License-Identifier: MPL-2.0
//
// Code generated by gen.go; DO NOT EDIT.

package mondoogql

// Input represents one of the Input structs:
//
// {{join (inputObjects .data.__schema.types) ", "}}.
type Input interface{}
{{range .data.__schema.types | sortByName}}{{if eq .kind "INPUT_OBJECT"}}
{{template "inputObject" .}}
{{end}}{{end}}


{{- define "inputObject" -}}
// {{.name}}
type {{.name}} struct {{"{"}}{{range .inputFields}}{{if eq .type.kind "NON_NULL"}}
	// (Required.)
	{{.name | identifier}} {{.type | type}} ` + "`" + `json:"{{.name}}"` + "`" + `{{end}}{{end}}
{{range .inputFields}}{{if ne .type.kind "NON_NULL"}}
	// (Optional.)
	{{.name | identifier}} {{.type | type}} ` + "`" + `json:"{{.name}},omitempty"` + "`" + `{{end}}{{end}}
}
{{- end -}}
`),
}

// parseTemplate takes a text template and returns a parsed template.
// This function is
func parseTemplate(text string) *template.Template {
	// typeString returns a string representation of GraphQL type t.
	var typeString func(t map[string]interface{}) string

	// typeString returns a string representation of GraphQL type t.
	typeString = func(t map[string]interface{}) string {
		switch t["kind"] {
		case "NON_NULL":
			s := typeString(t["ofType"].(map[string]interface{}))
			if !strings.HasPrefix(s, "*") {
				panic(fmt.Errorf("nullable type %q doesn't begin with '*'", s))
			}
			return s[1:] // Strip star from nullable type to make it non-null.
		case "LIST":
			return "*[]" + typeString(t["ofType"].(map[string]interface{}))
		default:
			return "*" + t["name"].(string)
		}
	}

	// compile template with custom functions
	return template.Must(template.New("").Funcs(template.FuncMap{
		"internal": func(s string) bool { return strings.HasPrefix(s, "__") },
		"quote":    strconv.Quote,
		"join":     strings.Join,
		"sortByName": func(types []interface{}) []interface{} {
			sort.Slice(types, func(i, j int) bool {
				ni := types[i].(map[string]interface{})["name"].(string)
				nj := types[j].(map[string]interface{})["name"].(string)
				return ni < nj
			})
			return types
		},
		"inputObjects": func(types []interface{}) []string {
			var names []string
			for _, t := range types {
				t := t.(map[string]interface{})
				if t["kind"].(string) != "INPUT_OBJECT" {
					continue
				}
				names = append(names, t["name"].(string))
			}
			sort.Strings(names)
			return names
		},
		"identifier": func(name string) string { return ident.ParseLowerCamelCase(name).ToMixedCaps() },
		"enumIdentifier": func(enum, value string) string {
			return enum + ident.ParseScreamingSnakeCase(value).ToMixedCaps()
		},
		"type":  typeString,
		"clean": func(s string) string { return strings.Join(strings.Fields(s), " ") },
		"endSentence": func(s string) string {
			s = strings.ToLower(s[0:1]) + s[1:]
			switch {
			default:
				s = "represents " + s
			case strings.HasPrefix(s, "autogenerated "):
				s = "is an " + s
			case strings.HasPrefix(s, "specifies "):
				// Do nothing.
			}
			if !strings.HasSuffix(s, ".") {
				s += "."
			}
			return s
		},
		"fullSentence": func(s string) string {
			if !strings.HasSuffix(s, ".") {
				s += "."
			}
			return s
		},
	}).Parse(text))
}
