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
	"github.com/gomodule/redigo/redis"
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

	// Initialize RediSearch client
	pool := &redis.Pool{Dial: func() (redis.Conn, error) {
		return redis.Dial(
			"tcp",
			fmt.Sprintf("%s:%d", args.RedisServer, args.RedisPort),
			redis.DialPassword(args.RedisPassword))
	}}

	// Build the Redis Client for searching
	searchClient := redisearch.NewClientFromPool(
		pool,
		args.RedisIndex,
	)

	// Build the graphql schema from the RediSearch Index
	// https://redis.io/commands/ft.info/ details the index serch schema
	// that we will map to a graphql schema
	schema, docs, nerr := rsq.FtInfo2Schema(searchClient, args.RedisIndex)
	if nerr != nil {
		log.Fatal(nerr)
	}

	// Serve the auto-generated graphql schema docs
	http.HandleFunc("/docs", docs.ServeDocs)

	// Perform all graphql queries against the schema
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

		// Update Prometheus metrics allowing us to track the response time in ms
		rsq.ObserveGraphqlDuration(time.Since(start).Milliseconds())

		// If we get a bad query make sure to update the metrics
		if result.Errors != nil {
			rsq.IncrQueryErrors()
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			fmt.Printf("could not write result to response: %s", err)
		}
	})

	// Serve the prometheus metrics
	http.Handle("/metrics", promhttp.Handler())

	// Return a 200 OK response for use as a health check
	// TODO: Add a health check using Redis PING command
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`OK`))
	})

	fmt.Println("Server is running on " + args.Addr + " and providing data from index: " + args.RedisIndex)
	http.ListenAndServe(args.Addr, nil)
}
