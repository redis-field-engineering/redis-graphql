package redissearchgraphql

import (
	"context"

	"github.com/RediSearch/redisearch-go/redisearch"
)

// FtSearch is the the most basic search query.
// Given a list of fields detailed from the docs page http://localhost:8080/docs
// it will return a JSON array of results matching those queries
// see https://redis.io/commands/ft.search for more information
func FtSearch(args map[string]interface{}, clients map[string]*redisearch.Client, index string, c context.Context) ([]map[string]interface{}, error) {
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
		if limOffset, ok := argsMap["limit_offset"]; ok {
			q = q.Limit(int(limOffset.(float64)), int(lim.(float64)))
		}
	}

	if verbatim, ok := argsMap["verbatim"]; ok {
		if verbatim.(bool) {
			q = q.SetFlags(redisearch.QueryVerbatim)
		}
	}

	client := clients[index]

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
