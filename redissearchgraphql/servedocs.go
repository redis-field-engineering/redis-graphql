package redissearchgraphql

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var dataTempl = template.Must(template.New("").Parse(dataHTML))

// dataHTML is the HTML template for the docs page
const dataHTML = `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>GraphQL Documentation</title>
    </head>
    <style>
    pre { background-color: black; color: white; margin-left: 80px; width: 65%; }
    p { margin-left: 65px; }
    h1 { color: blue; margin-left: 25px; }
    h2 { color: maroon; margin-left: 40px; }
    h3 { color: green; margin-left: 60px; }
    h4 { color: black; margin-left: 70px; }
    table {
	border-collapse: collapse;
	width: 50%;
	margin-left: 80px;
      }
      
      th, td {
	text-align: left;
	padding: 8px;
      }
      
      tr:nth-child(even) {
	background-color: #D6EEEE;
      }
    </style>
    <body>

    <h1> {{ .IndexName }} GraphQL Documentation</h1>

    <h2>Fields</h2>
    <p>Below is a list of all the fields in the index and their types that are availalbe to query or fetch.</p>
    <table>
    <tr><th>Field</th><th>Type</th></tr>
    {{ range $val := .Floats }}
    <tr><td>{{ $val }}</td><td>Numeric</td></tr>
    {{ end }}
    {{ range $val := .Strings }}
    <tr><td>{{ $val }}</td><td>String</td></tr>
    {{ end }}
    {{ range $val := .Geos }}
    <tr><td>{{ $val }}</td><td>Geo</td></tr>
    {{ end }}
    {{ range $val := .Tags }}
    <tr><td>{{ $val }}</td><td>Tag</td></tr>
    {{ end }}

    </table>

    <h2>Numeric Queries</h2>
    {{ range $val := .Floats }}
	<h3>{{ $val }}</h3>
   	<h4>{{ index $.FieldDocs $val  }}<h4>
   <p>Example Query:</p>
    <pre>
    query {
	{{ $.IndexName }}(
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
    	{{- range $z := $.Tags }}
	    {{ $z }},
    	{{- end }}
	}
    }
</pre>
    	{{ range $val2 := $.FloatSuffix }} <!-- need $. here instead of . -->
	    <h3>{{ $val }}_{{ $val2 }}</h3>
	    {{ $combined := (printf "%s_%s" $val $val2 ) }}
	    <h4>{{ index $.FieldDocs $combined  }}<h4>
   <p>Example Query:</p>
    <pre>
    query {
	{{ $.IndexName }}(
	    {{- if eq $val2 "bte" }}
	    {{ $combined }}: [10, 20],
	    {{- else }}
	    {{ $combined }}: 3.14,
	    {{- end }}
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
    	{{- range $z := $.Tags }}
	    {{ $z }},
    	{{- end }}
	}
    }
</pre>
	{{ end }}
    {{ end }}

    <h2>String Queries</h2>
    {{ range $val := .Strings }}
	<h3>{{ $val }}</h3>
   	<h4>{{ index $.FieldDocs $val  }}<h4>
   <p>Example Query:</p>
    <pre>
    query {
	{{ $.IndexName }}(
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
    	{{- range $z := $.Tags }}
	    {{ $z }},
    	{{- end }}
	}
    }
</pre>
    	{{ range $val2 := $.StringSuffix }} <!-- need $. here instead of . -->
	    <h3>{{ $val }}_{{ $val2 }}</h3>
	    {{ $combined := (printf "%s_%s" $val $val2 ) }}
	    <h4>{{ index $.FieldDocs $combined  }}<h4>
   <p>Example Query:</p>
    <pre>
    query {
	{{ $.IndexName }}(
	    {{ $combined }}: "myValue",
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
    	{{- range $z := $.Tags }}
	    {{ $z }},
    	{{- end }}
	}
    }
</pre>
	{{ end }}
    {{ end }}

    <h2>Geo Queries</h2>
    {{ range $val := .Geos }}
	<h3>{{ $val }}</h3>
   	<h4>{{ index $.FieldDocs $val  }}<h4>
   <p>Example Query:</p>
    <pre>
    query {
	{{ $.IndexName }}(
	    {{ $val }}: {lat: 37.377658, lon: -122.064228, radius: 10, unit: "km"},
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
    	{{- range $z := $.Tags }}
	    {{ $z }},
    	{{- end }}
	}
    }
</pre>
    	{{ range $val2 := $.GeoSuffix }} <!-- need $. here instead of . -->
	    <h3>{{ $val }}_{{ $val2 }}</h3>
	    {{ $combined := (printf "%s_%s" $val $val2 ) }}
	    <h4>{{ index $.FieldDocs $combined  }}<h4>
   <p>Example Query:</p>
    <pre>
    query {
	{{ $.IndexName }}(
	    {{ $combined }}: {lat: 37.377658, lon: -122.064228, radius: 10, unit: "km"}},
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
    	{{- range $z := $.Tags }}
	    {{ $z }},
    	{{- end }}
	}
    }
</pre>
	{{ end }}
    {{ end }}

    <h2>Tag Queries</h2>
    	<p>Query based on tags</p>
    {{ range $val := .Tags }}
	<h3>{{ $val }}</h3>
   	<h4>{{ index $.FieldDocs $val  }}<h4>
    <p>Example Query:</p>
    <pre>
    query {
	{{ $.IndexName }}(
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
    	{{- range $z := $.Tags }}
	    {{ $z }},
    	{{- end }}
	}
    }
</pre>

    	{{ range $val2 := $.TagSuffix }} <!-- need $. here instead of . -->
	    <h3>{{ $val }}_{{ $val2 }}</h3>
	        {{ $combined := (printf "%s_%s" $val $val2 ) }}
   		<h4>{{ index $.FieldDocs $combined  }}<h4>
    <p>Example Query:</p>
    <pre>
    query {
	{{ $.IndexName }}(
	    {{ $combined }}: ["tag1", "tag2"]
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
    	{{- range $z := $.Tags }}
	    {{ $z }},
    	{{- end }}
	}
    }
</pre>
	{{ end }}

    {{ end }}

    <h2>Raw Queries</h2>
    <h3>This allows you to run raw queries using the <a href="https://oss.redis.com/redisearch/Query_Syntax/">RediSearch</a> syntax.</h3>

    <p>Example Query:</p>
    <pre>
    query {
	{{ $.IndexName }}(
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

    <p>This document is auto-generated from the RediSearch Schema</p>
    <p>For more information see <a href="https://redisearch.io/docs/schema/">RediSearch Schema Docs</a></p>
    <p>For full documentation see <a href="https://github.com/redis-field-engineering/RediSearchGraphQL">RediSearch GraphQL Docs</a></p>


    </body>
</html>
`

// ServeDocs generates the documentation for the schema and displays using the template above
func (alldocs AllDocs) ServeDocs(w http.ResponseWriter, r *http.Request) {
	promDocsViewCount.Inc()
	vars := mux.Vars(r)
	idx := vars["index"]
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var v interface{} = alldocs[idx]
	err := dataTempl.Execute(w, &v)
	if err != nil {
		fmt.Println(err)
	}
}
