package redissearchgraphql

import (
	"context"
	"fmt"

	"github.com/RediSearch/redisearch-go/redisearch"
)

func FtAggCount(args map[string]interface{}, client *redisearch.Client, c context.Context) ([]map[string]interface{}, error) {
	var res []map[string]interface{}
	q1 := redisearch.NewAggregateQuery()
	argsMap := c.Value("v").(PostVars).Variables

	qstring, err := QueryBuilder(args, argsMap, true)
	if err != nil {
		return res, err
	}

	if q, ok := args["_agg_groupby"]; ok {
		q1 = redisearch.NewAggregateQuery().SetQuery(redisearch.NewQuery(qstring)).
			GroupBy(*redisearch.NewGroupBy().AddFields(fmt.Sprintf("@%s", q)).
				Reduce(*redisearch.NewReducer(redisearch.GroupByReducerCount, []string{}).SetAlias("_agg_groupby_count"))).
			SortBy([]redisearch.SortingKey{*redisearch.NewSortingKeyDir("@_agg_groupby_count", false)})
	}

	if lim, ok := argsMap["limit"]; ok {
		q1 = q1.Limit(0, int(lim.(float64)))
	}

	//if verbatim, ok := argsMap["verbatim"]; ok {
	//	if verbatim.(bool) {
	//		q1 = q1.SetFlags(redisearch.QueryVerbatim)
	//	}
	//}

	docs, _, err := client.Aggregate(q1)

	if err != nil {
		return res, err
	}

	for _, doc := range docs {
		if len(doc) == 4 {
			res = append(res, map[string]interface{}{doc[0]: doc[1], doc[2]: doc[3]})
		}
	}

	return res, nil
}
