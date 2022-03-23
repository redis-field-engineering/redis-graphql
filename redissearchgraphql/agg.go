package redissearchgraphql

import (
	"context"
	"fmt"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"
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

func FtAggNumGroup(args map[string]interface{}, client *redisearch.Client, c context.Context) ([]map[string]interface{}, error) {
	var res []map[string]interface{}
	q1 := redisearch.NewAggregateQuery()
	argsMap := c.Value("v").(PostVars).Variables

	qstring, err := QueryBuilder(args, argsMap, true)
	if err != nil {
		return res, err
	}

	aggFunction := map[string]redisearch.GroupByReducers{
		"sum":      redisearch.GroupByReducerSum,
		"min":      redisearch.GroupByReducerMin,
		"max":      redisearch.GroupByReducerMax,
		"avg":      redisearch.GroupByReducerAvg,
		"quantile": redisearch.GroupByReducerQuantile,
		"stddev":   redisearch.GroupByReducerStdDev,
	}

	funcArgs := []string{args["_agg_num_field"].(string)}
	if quantile, ok := args["_agg_num_quantile"].(float64); ok {
		funcArgs = append(funcArgs, fmt.Sprintf("%f", quantile))
	}

	if aggFunc, ok := args["_agg_num_function"].(string); ok {

		if q, ok := args["_agg_groupby"]; ok {
			q1 = redisearch.NewAggregateQuery().SetQuery(redisearch.NewQuery(qstring)).
				GroupBy(*redisearch.NewGroupBy().AddFields(fmt.Sprintf("@%s", q)).
					Reduce(*redisearch.NewReducer(aggFunction[aggFunc], funcArgs).SetAlias("_agg_groupby_num"))).
				SortBy([]redisearch.SortingKey{*redisearch.NewSortingKeyDir("@_agg_groupby_num", false)})
		}
	}

	if lim, ok := argsMap["limit"]; ok {
		q1 = q1.Limit(0, int(lim.(float64)))
	}

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

func FtAggRaw(args map[string]interface{}, client *redisearch.Client, c context.Context) ([]map[string]interface{}, error) {
	var res []map[string]interface{}
	var aggPlan redis.Args
	q1 := redisearch.NewAggregateQuery()
	argsMap := c.Value("v").(PostVars).Variables

	qstring, err := QueryBuilder(args, argsMap, true)
	if err != nil {
		return res, err
	}

	q1.Query = redisearch.NewQuery(qstring)

	if lim, ok := argsMap["limit"]; ok {
		q1 = q1.Limit(0, int(lim.(float64)))
	}

	for _, y := range args["raw_agg_plan"].([]interface{}) {
		aggPlan = aggPlan.Add(y)
	}

	q1.AggregatePlan = aggPlan

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
