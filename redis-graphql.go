package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/alexflint/go-arg"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	rsq "github.com/redis-field-engineering/RediSearchGraphQL/redissearchgraphql"
	"go.uber.org/zap"
)

var args struct {
	Addr          string `help:"where to listen for websocket requests" default:"localhost:8080" arg:"env:LISTEN"`
	RedisServer   string `help:"Redis to connect to" default:"localhost" arg:"--redis-host, -s, env:REDIS_SERVER"`
	RedisPort     int    `help:"Redis port to connect to" default:"6379" arg:"--redis-port, -p, env:REDIS_PORT"`
	RedisPassword string `help:"Redis password" default:"" arg:"--redis-password, -a, env:REDIS_PASSWORD"`
}

func main() {
	// Parse the command line arguments
	arg.MustParse(&args)

	// Initialize the Logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

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
		"", // we don't need this to run the list command
	)

	indicies, ierr := searchClient.List()
	if ierr != nil {
		sugar.Fatal("Error getting index list", zap.Error(ierr))
	}

	searchIndices := make(map[string]*redisearch.Client, len(indicies))

	for i, x := range indicies {
		searchIndices[x] = redisearch.NewClientFromPool(pool, indicies[i])
	}

	// Build the graphql schema from the RediSearch Index
	// https://redis.io/commands/ft.info/ details the index serch schema
	// that we will map to a graphql schema
	schema, docs, nerr := rsq.FtInfo2Schema(searchIndices)
	if nerr != nil {
		sugar.Fatalw("Failed to build schema", "error", nerr)
	}

	// Setup our mux router for handling multiple Redisearch Indices
	router := mux.NewRouter()

	// Serve the auto-generated graphql schema docs
	router.HandleFunc("/docs", docs.ServeAllDocs)
	router.HandleFunc("/docs/{index}", docs.ServeDocs)

	// Perform all graphql queries against the schema
	router.HandleFunc("/graphql", func(w http.ResponseWriter, req *http.Request) {
		var p rsq.PostData
		if err := json.NewDecoder(req.Body).Decode(&p); err != nil {

			// Log the error
			sugar.Infow(
				"graphql_decode_error",
				"ip", req.RemoteAddr,
				"path", req.URL.Path,
				"error", err,
			)
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

		// Log the request
		sugar.Infow(
			"graphql_query",
			"duration_ms", time.Since(start).Milliseconds(),
			"ip", req.RemoteAddr,
			"path", req.URL.Path,
		)

		// Update Prometheus metrics allowing us to track the response time in ms
		rsq.ObserveGraphqlDuration(time.Since(start).Milliseconds())

		// If we get a bad query make sure to update the metrics
		if result.Errors != nil {
			rsq.IncrQueryErrors()
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			sugar.Errorw("Failed to encode response", "error", err)
		}
	})

	// Serve the prometheus metrics
	router.Handle("/metrics", promhttp.Handler())

	// Return a 200 OK response for use as a health check
	// TODO: Add a health check using Redis PING command
	router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`OK`))
	})

	sugar.Infow(
		"server_started",
		"addr", args.Addr,
	)

	http.ListenAndServe(args.Addr, router)
}
