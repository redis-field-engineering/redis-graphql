package redissearchgraphql

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

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
			Name:    "redisgraphql_graphql_duration_microseconds",
			Help:    "The amount of time it takes to process a GraphQL request",
			Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000, 25000, 50000, 100000, 250000, 500000, 1000000},
		})
)

func IncrPromPostErrors() {
	promPostErrorCount.Inc()
}

func IncrQueryErrors() {
	promQueryErrorCount.Inc()
}

func ObserveGraphqlDuration(dur int64) {
	promGraphqlHistogram.Observe(float64(dur))
}

func InitPrometheus() {
	prometheus.MustRegister(promGraphqlHistogram)
}
