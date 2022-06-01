package redissearchgraphql

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Create variables for Prometheus metrics
var (
	promDocsViewCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "redisgraphql_docs_count",
		Help: "The total number of document requests processed",
	})
	promFtSearchCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "redisgraphql_ft_searches_count",
		Help: "The total number of ft searches processed",
	})
	promAggCountCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "redisgraphql_agg_count_count",
		Help: "The total number of count aggregations processed",
	})
	promAggNumgroupCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "redisgraphql_agg_numgroup_counts",
		Help: "The total number of numgroup aggregations processed",
	})
	promAggRawCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "redisgraphql_agg_raw_counts",
		Help: "The total number of raw aggregations processed",
	})
	promQueryErrorCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "redisgraphql_query_errrors",
		Help: "The total number of query errors",
	})
	promPostErrorCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "redisgraphql_post_errrors",
		Help: "The total number of times the post could not be JSON decoded",
	})
	promGraphqlHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "redisgraphql_graphql_duration_milliseconds",
			Help:    "The amount of time it takes to process a GraphQL request in ms",
			Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000, 25000, 50000, 100000, 250000, 500000, 1000000},
		})
	promRedisHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "redisgraphql_redis_duration_milliseconds",
			Help:    "The amount of time it takes to process a GraphQL request in ms",
			Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000, 25000, 50000, 100000, 250000, 500000, 1000000},
		})
	promPoolActive = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "redisgraphql_pool_active",
		Help: "The number of active connections in the pool",
	})
	promPoolIdle = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "redisgraphql_pool_idle",
		Help: "The number of idle connections in the pool",
	})
	promPoolWait = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "redisgraphql_pool_wait",
		Help: "The number of wait connections in the pool",
	})
)

// IncrPromPostErrors increments the post error counter
// It's exported to allow for top levels availability
func IncrPromPostErrors() {
	promPostErrorCount.Inc()
}

// IncrQueryErrors increments the post error counter
// It's exported to allow for top levels availability
func IncrQueryErrors() {
	promQueryErrorCount.Inc()
}

// ObserveGraphqlDuration observes the duration of a GraphQL request for latency purposes
// It's exported to allow for top levels availability
func ObserveGraphqlDuration(dur int64) {
	promGraphqlHistogram.Observe(float64(dur))
}

// ObserveGraphqlDuration observes the duration of a Redisearch request for latency purposes
// It's exported to allow for top levels availability
func ObserveRedisDuration(dur int64) {
	promRedisHistogram.Observe(float64(dur))
}

// InitPrometheus initializes the Prometheus metrics which is necessary for
// for histogram and summary types
func InitPrometheus() {
	prometheus.MustRegister(promGraphqlHistogram)
	prometheus.MustRegister(promRedisHistogram)
}
