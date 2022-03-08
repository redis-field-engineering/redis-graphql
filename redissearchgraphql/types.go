package redissearchgraphql

type PostVars struct {
	Variables map[string]interface{}
}

type PostData struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
	Variables map[string]interface{} `json:"variables"`
}

type SchemaDocs struct {
	Floats       []string `json:"floats"`
	FloatSuffix  []string `json:"floatSuffix"`
	Strings      []string `json:"strings"`
	StringSuffix []string `json:"stringSuffix"`
	Geos         []string `json:"geos"`
	GeoSuffix    []string `json:"geoSuffix"`
	Tags         []string `json:"tags"`
	TagSuffix    []string `json:"tagSuffix"`
}
