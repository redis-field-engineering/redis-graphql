package redissearchgraphql

import (
	"context"

	"github.com/RediSearch/redisearch-go/redisearch"
)

func FtSearch(args map[string]interface{}, client *redisearch.Client, c context.Context) ([]map[string]interface{}, error) {
	promFtSearchCount.Inc()
	var res []map[string]interface{}
	qstring := ""
	argsMap := c.Value("v").(PostVars).Variables

	qstring, err := QueryBuilder(args, argsMap, false)
	if err != nil {
		return res, err
	}

	q := redisearch.NewQuery(qstring)

	if lim, ok := argsMap["limit"]; ok {
		q = q.Limit(0, int(lim.(float64)))
	}

	if verbatim, ok := argsMap["verbatim"]; ok {
		if verbatim.(bool) {
			q = q.SetFlags(redisearch.QueryVerbatim)
		}
	}

	docs, _, err := client.Search(q)

	if err != nil {
		promPostErrorCount.Inc()
		return res, err
	}

	for _, doc := range docs {
		res = append(res, doc.Properties)
	}

	return res, nil
}
