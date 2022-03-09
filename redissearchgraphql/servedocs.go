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
    <style>
    pre { background-color: black; color: white; margin-left: 80px; }
    p { margin-left: 65px; }
    h2 { color: maroon; margin-left: 40px; }
    h3 { color: green; margin-left: 60px; }
    </style>
    <body>

    <h2>Numeric Queries</h2>
    {{ range $val := .Floats }}
	<h3>{{ $val }}</h3>
   <p>Find all documents where the {{ $val }} field is equal to 3.14</p>
   <p>Example Query:</p>
    <pre>
    query {
	ft(
	    {{ $val }}: 3.14,
	    )
	{
    	{{- range $x := $.Strings }}
	    {{ $x }},
    	{{- end }}
    	{{- range $y := $.Floats }}
	    {{ $y }},
    	{{- end }}
    	{{- range $z := $.Geos }}
	    {{ $z }},
    	{{- end }}
	}
    }
</pre>
    	{{ range $val2 := $.FloatSuffix }} <!-- need $. here instead of . -->
	    <h3>{{ $val }}_{{ $val2 }}</h3>
	{{ end }}
    {{ end }}

    <h2>String Queries</h2>
    {{ range $val := .Strings }}
	<h3>{{ $val }}</h3>
   <p>Example Query:</p>
   <p>Find all documents where the {{ $val }} field matches my Value</p>
    <pre>
    query {
	ft(
	    {{ $val }}: "myValue",
	    )
	{
    	{{- range $x := $.Strings }}
	    {{ $x }},
    	{{- end }}
    	{{- range $y := $.Floats }}
	    {{ $y }},
    	{{- end }}
    	{{- range $z := $.Geos }}
	    {{ $z }},
    	{{- end }}
	}
    }
</pre>
    	{{ range $val2 := $.StringSuffix }} <!-- need $. here instead of . -->
	    <h3>{{ $val }}_{{ $val2 }}</h3>
	{{ end }}
    {{ end }}

    <h2>Geo Queries</h2>
    {{ range $val := .Geos }}
	<h3>{{ $val }}</h3>
   <p>Find all documents where the {{ $val }} is within a 10km radius of the point (lon,lat)</p>
   <p>Example Query:</p>
    <pre>
    query {
	ft(
	    {{ $val }}: {lat: 37.377658, lon: -122.064228, radius: 10, unit: "km"}},
	    )
	{
    	{{- range $x := $.Strings }}
	    {{ $x }},
    	{{- end }}
    	{{- range $y := $.Floats }}
	    {{ $y }},
    	{{- end }}
    	{{- range $z := $.Geos }}
	    {{ $z }},
    	{{- end }}
	}
    }
</pre>
    	{{ range $val2 := $.GeoSuffix }} <!-- need $. here instead of . -->
	    <h3>{{ $val }}_{{ $val2 }}</h3>
	{{ end }}
    {{ end }}

    <h3>Tag Queries</h3>
    {{ range $val := .Tags }}
	<h3>{{ $val }}</h3>
    	<p>Query based on tags</p>
   <p>Find all documents where the tag names {{ $val }} is present</p>
    <p>Example Query:</p>
    <pre>
    query {
	ft(
	    {{ $val }}: ["tag1", "tag2"]
	    )
	{
    	{{- range $x := $.Strings }}
	    {{ $x }},
    	{{- end }}
    	{{- range $y := $.Floats }}
	    {{ $y }},
    	{{- end }}
    	{{- range $z := $.Geos }}
	    {{ $z }},
    	{{- end }}
	}
    }
</pre>

    	{{ range $val2 := $.TagSuffix }} <!-- need $. here instead of . -->
	    <h3>{{ $val }}_{{ $val2 }}</h3>
	{{ end }}

    {{ end }}

    <h2>Raw Queries</h2>
    <p>This allows you to run raw queries using the <a href="https://oss.redis.com/redisearch/Query_Syntax/">RediSearch</a> syntax.</p>

    <p>Example Query:</p>
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
    	{{- range $val := .Geos }}
	    {{ $val }},
    	{{- end }}
    	{{- range $val := .Tags }}
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
