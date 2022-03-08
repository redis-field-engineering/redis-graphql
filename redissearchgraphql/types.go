package redissearchgraphql

type PostVars struct {
	Variables map[string]interface{}
}

type PostData struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
	Variables map[string]interface{} `json:"variables"`
}
