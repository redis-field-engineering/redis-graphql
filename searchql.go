package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/alexflint/go-arg"
	"github.com/graphql-go/graphql"
)

/*****************************************************************************/
/* Shared data variables to allow dynamic reloads
/*****************************************************************************/

type postVars struct {
	v map[string]interface{}
}

type postData struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
	Variables map[string]interface{} `json:"variables"`
}

var schema graphql.Schema

var args struct {
	Addr          string `help:"where to listen for websocket requests" default:"localhost:8080" arg:"env:LISTEN"`
	RedisServer   string `help:"Redis to connect to" default:"localhost" arg:"--redis-host, -s, env:REDIS_SERVER"`
	RedisPort     int    `help:"Redis port to connect to" default:"6379" arg:"--redis-port, -p, env:REDIS_PORT"`
	RedisPassword string `help:"Redis password" default:"" arg:"--redis-password, -a, env:REDIS_PASSWORD"`
	RedisIndex    string `help:"RediSearch Index" default:"idx" arg:"--redis-index, -i, env:REDIS_INDEX"`
}

var TagScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Tag",
	Description: "Redisearch Tag",
})

func ftSearch(args map[string]interface{}, client *redisearch.Client, c context.Context) []map[string]interface{} {
	var res []map[string]interface{}
	fmt.Printf("%+v\n", c.Value("v"))

	qstring := ""

	for k, v := range args {
		qstring += "@" + k + ":" + v.(string) + " "
	}
	docs, _, err := client.Search(redisearch.NewQuery(qstring))

	if err != nil {
		log.Fatal(err)
	}

	for _, doc := range docs {
		res = append(res, doc.Properties)
	}

	return res
}

func FtInfo2Schema(client *redisearch.Client) error {
	idx, err := client.Info()

	if err != nil {
		log.Fatal("cannot do info on index:"+args.RedisIndex, " Error: ", err)
	}

	fields := make(graphql.Fields)
	args := make(graphql.FieldConfigArgument)

	for _, field := range idx.Schema.Fields {
		//fmt.Printf("%+v\n", field)
		if field.Type == 0 {
			fields[field.Name] = &graphql.Field{
				Type: graphql.String,
			}
			args[field.Name] = &graphql.ArgumentConfig{
				Type: graphql.String,
			}
		}

		if field.Type == 1 {
			fields[field.Name] = &graphql.Field{
				Type: graphql.Float,
			}
			args[field.Name] = &graphql.ArgumentConfig{
				Type: graphql.Float,
			}
		}

		//if field.Type == 3 {
		//	fields[field.Name] = &graphql.Field{
		//		Type: TagScalar,
		//	}
		//	args[field.Name] = &graphql.ArgumentConfig{
		//		Type: graphql.String,
		//	}
		//}

	}

	var ftType = graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "FT",
			Fields: fields,
		},
	)

	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"ft": &graphql.Field{
					Type: graphql.NewList(ftType),
					Args: args,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return ftSearch(p.Args, client, p.Context), nil
					},
				},
			},
		})
	schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryType,
		},
	)
	return nil
}

func main() {
	arg.MustParse(&args)

	searchClient := redisearch.NewClient(
		fmt.Sprintf("%s:%d", args.RedisServer, args.RedisPort),
		args.RedisIndex,
	)
	nerr := FtInfo2Schema(searchClient)
	if nerr != nil {
		log.Fatal(nerr)
	}

	http.HandleFunc("/graphql", func(w http.ResponseWriter, req *http.Request) {
		var p postData
		if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
			w.WriteHeader(400)
			return
		}
		c := context.Background()
		z := postVars{v: p.Variables}
		c = context.WithValue(c, "v", z)
		result := graphql.Do(graphql.Params{
			Context:        c,
			Schema:         schema,
			RequestString:  p.Query,
			VariableValues: p.Variables,
			OperationName:  p.Operation,
		})
		if err := json.NewEncoder(w).Encode(result); err != nil {
			fmt.Printf("could not write result to response: %s", err)
		}
	})

	fmt.Println("Now server is running on " + args.Addr)
	fmt.Println(`Example:  curl -X POST  -H "Content-Type: application/json"  --data '{ "variables": {"foo": 1}, "query": "{ ft(hqstate:\"ca\", hqcity:\"san\", sector: \"Technology\") { company,ceo,sector,hqcity,hqstate } }" }' http://localhost:8080/graphql`)
	http.ListenAndServe(args.Addr, nil)
}
