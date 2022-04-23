package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/alexflint/go-arg"
	"github.com/graphql-go/graphql"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	rsq "github.com/redis-field-engineering/RediSearchGraphQL/redissearchgraphql"
)

var args struct {
	Addr          string `help:"where to listen for websocket requests" default:"localhost:8080" arg:"env:LISTEN"`
	RedisServer   string `help:"Redis to connect to" default:"localhost" arg:"--redis-host, -s, env:REDIS_SERVER"`
	RedisPort     int    `help:"Redis port to connect to" default:"6379" arg:"--redis-port, -p, env:REDIS_PORT"`
	RedisPassword string `help:"Redis password" default:"" arg:"--redis-password, -a, env:REDIS_PASSWORD"`
	RedisIndex    string `help:"RediSearch Index" default:"idx" arg:"--redis-index, -i, env:REDIS_INDEX"`
}

func main() {
	// Parse the command line arguments
	arg.MustParse(&args)

	// Initialize Prometheus histogram and summary metrics
	rsq.InitPrometheus()

	searchClient := redisearch.NewClient(
		fmt.Sprintf("%s:%d", args.RedisServer, args.RedisPort),
		args.RedisIndex,
	)
	schema, docs, nerr := rsq.FtInfo2Schema(searchClient, args.RedisIndex)
	if nerr != nil {
		log.Fatal(nerr)
	}

	http.HandleFunc("/docs", docs.ServeDocs)

	http.HandleFunc("/graphql", func(w http.ResponseWriter, req *http.Request) {
		var p rsq.PostData
		if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
			rsq.IncrPromPostErrors()
			w.WriteHeader(400)
			return
		}
		c := context.Background()
		z := rsq.PostVars{Variables: p.Variables}
		c = context.WithValue(c, "v", z)
		start := time.Now()
		result := graphql.Do(graphql.Params{
			Context:        c,
			Schema:         schema,
			RequestString:  p.Query,
			VariableValues: p.Variables,
			OperationName:  p.Operation,
		})
		rsq.ObserveGraphqlDuration(time.Since(start).Milliseconds())
		if result.Errors != nil {
			rsq.IncrQueryErrors()
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			fmt.Printf("could not write result to response: %s", err)
		}
	})

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`OK`))
	})

	fmt.Println("Server is running on " + args.Addr + " and providing data from index: " + args.RedisIndex)
	http.ListenAndServe(args.Addr, nil)
}
