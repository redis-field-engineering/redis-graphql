package redissearchgraphql

import (
	"fmt"
	"html/template"
	"net/http"
)

var dataTempl = template.Must(template.New("").Parse(dataHTML))

const dataHTML = `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>GraphQL Documentation</title>
    </head>
    <body>

    <h2>Numeric Queries</h2>
    {{ range $val := .Floats }}
	<p>{{ $val }}</p>
    	{{ range $val2 := $.FloatSuffix }} <!-- need $. here instead of . -->
	    <p>{{ $val }}_{{ $val2 }}</p>
	{{ end }}
    {{ end }}

    <h2>String Queries</h2>
    {{ range $val := .Strings }}
	<p>{{ $val }}</p>
    	{{ range $val2 := $.StringSuffix }} <!-- need $. here instead of . -->
	    <p>{{ $val }}_{{ $val2 }}</p>
	{{ end }}
    {{ end }}

    <h2>Geo Queries</h2>
    {{ range $val := .Geos }}
	<p>{{ $val }}</p>
    	{{ range $val2 := $.GeoSuffix }} <!-- need $. here instead of . -->
	    <p>{{ $val }}_{{ $val2 }}</p>
	{{ end }}
    {{ end }}

    <h2>Tag Queries</h2>
    {{ range $val := .Tags }}
	<p>{{ $val }}</p>
    	{{ range $val2 := $.TagSuffix }} <!-- need $. here instead of . -->
	    <p>{{ $val }}_{{ $val2 }}</p>
	{{ end }}
    {{ end }}

    <h2>raw Queries</h2>
    This allows you to run raw queries using the Redisearch Syntax
    <pre>
    query {
	ft(
	    raw_query: "*",
	    )
	{
    	{{- range $val := .Strings }}
	    {{ $val }},
    	{{- end }}
    	{{- range $val := .Floats }}
	    {{ $val }},
    	{{- end }}
	}
    }
    </pre>



    </body>
</html>
`

func (d *SchemaDocs) ServeDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var v interface{} = d
	err := dataTempl.Execute(w, &v)
	if err != nil {
		fmt.Println(err)
	}
}

func ServDocs() string {
	html := "<html>manamana</html>"
	return html
}
