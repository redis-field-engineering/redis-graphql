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
    pre { background-color: black; color: white; margin-left: 80px; width: 65%; }
    p { margin-left: 65px; }
    h2 { color: maroon; margin-left: 40px; }
    h3 { color: green; margin-left: 60px; }
    h4 { color: black; margin-left: 70px; }
    </style>
    <body>

    <h2>Numeric Queries</h2>
    {{ range $val := .Floats }}
	<h3>{{ $val }}</h3>
   	<h4>{{ index $.FieldDocs $val  }}<h4>
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
	    {{ $combined := (printf "%s_%s" $val $val2 ) }}
	    <h4>{{ index $.FieldDocs $combined  }}<h4>
   <p>Example Query:</p>
    <pre>
    query {
	ft(
	    {{ $combined }}: 3.14,
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
	{{ end }}
    {{ end }}

    <h2>String Queries</h2>
    {{ range $val := .Strings }}
	<h3>{{ $val }}</h3>
   	<h4>{{ index $.FieldDocs $val  }}<h4>
   <p>Example Query:</p>
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
	    {{ $combined := (printf "%s_%s" $val $val2 ) }}
	    <h4>{{ index $.FieldDocs $combined  }}<h4>
   <p>Example Query:</p>
    <pre>
    query {
	ft(
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
	    {{ $combined := (printf "%s_%s" $val $val2 ) }}
	    <h4>{{ index $.FieldDocs $combined  }}<h4>
   <p>Example Query:</p>
    <pre>
    query {
	ft(
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
	        {{ $combined := (printf "%s_%s" $val $val2 ) }}
   		<h4>{{ index $.FieldDocs $combined  }}<h4>
    <p>Example Query:</p>
    <pre>
    query {
	ft(
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

    <p>This document is auto-generated from the RediSearch Schema</p>
    <p>For more information see <a href="https://redisearch.io/docs/schema/">RediSearch Schema Docs</a></p>
    <p>For full documentation see <a href="https://github.com/redis-field-engineering/RediSearchGraphQL">RediSearch GraphQL Docs</a></p>


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
