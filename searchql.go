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
	rsq "github.com/redis-field-engineering/RediSearchGraphQL/redissearchgraphql"
)

/*****************************************************************************/
/* Shared data variables to allow dynamic reloads
/*****************************************************************************/

var schema graphql.Schema

var args struct {
	Addr          string `help:"where to listen for websocket requests" default:"localhost:8080" arg:"env:LISTEN"`
	RedisServer   string `help:"Redis to connect to" default:"localhost" arg:"--redis-host, -s, env:REDIS_SERVER"`
	RedisPort     int    `help:"Redis port to connect to" default:"6379" arg:"--redis-port, -p, env:REDIS_PORT"`
	RedisPassword string `help:"Redis password" default:"" arg:"--redis-password, -a, env:REDIS_PASSWORD"`
	RedisIndex    string `help:"RediSearch Index" default:"idx" arg:"--redis-index, -i, env:REDIS_INDEX"`
}

func main() {
	arg.MustParse(&args)

	searchClient := redisearch.NewClient(
		fmt.Sprintf("%s:%d", args.RedisServer, args.RedisPort),
		args.RedisIndex,
	)
	schema, nerr := rsq.FtInfo2Schema(searchClient, args.RedisIndex)
	if nerr != nil {
		log.Fatal(nerr)
	}

	http.HandleFunc("/docs", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "this is where we get the docs")
	})

	http.HandleFunc("/graphql", func(w http.ResponseWriter, req *http.Request) {
		var p rsq.PostData
		if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
			w.WriteHeader(400)
			return
		}
		c := context.Background()
		z := rsq.PostVars{Variables: p.Variables}
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
	fmt.Println(`Example:  curl -X POST  -H "Content-Type: application/json"  --data '{ "variables": {"limit": 29, "verbatim": true}, "query": "{ ft(hqstate:\"ca\", hqcity:\"san\", sector: \"Technology\") { company,ceo,sector,hqcity,hqstate } }" }' http://localhost:8080/graphql`)
	http.ListenAndServe(args.Addr, nil)
}
