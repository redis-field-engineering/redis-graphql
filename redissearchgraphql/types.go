package redissearchgraphql

import "github.com/graphql-go/graphql"

type PostVars struct {
	Variables map[string]interface{}
}

type AllDocs map[string]SchemaDocs

type PostData struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
	Variables map[string]interface{} `json:"variables"`
}

type SchemaDocs struct {
	IndexName    string            `json:"indexName"`
	Floats       []string          `json:"floats"`
	FloatSuffix  []string          `json:"floatSuffix"`
	Strings      []string          `json:"strings"`
	StringSuffix []string          `json:"stringSuffix"`
	Geos         []string          `json:"geos"`
	GeoSuffix    []string          `json:"geoSuffix"`
	Tags         []string          `json:"tags"`
	TagSuffix    []string          `json:"tagSuffix"`
	FieldDocs    map[string]string `json:"fieldDocs"`
}

type Index struct {
	SchemaDocs SchemaDocs     `json:"schemaDocs"`
	Schema     graphql.Schema `json:"schema"`
	SchemaErr  error          `json:"schemaErr"`
}

func NewSchemaDocs() *SchemaDocs {
	return &SchemaDocs{FieldDocs: make(map[string]string)}
}

func (s *AllDocs) GetIndex(indexName string) *SchemaDocs {
	return NewSchemaDocs()
}
