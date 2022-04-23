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
)
