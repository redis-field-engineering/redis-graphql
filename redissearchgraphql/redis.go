package redissearchgraphql

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/alexflint/go-arg"
	"github.com/graphql-go/graphql"
)

func FtSearch(args map[string]interface{}, client *redisearch.Client, c context.Context) []map[string]interface{} {
	var res []map[string]interface{}
	qstring := ""

	if args["raw_query"] == nil {

		for k, v := range args {
			switch v.(type) {
			case string:
				if strings.HasSuffix(k, "_not") {
					qstring += "-@" + strings.TrimSuffix(k, "_not") + ":" + v.(string) + " "
				} else if strings.HasSuffix(k, "_opt") {
					qstring += "~@" + strings.TrimSuffix(k, "_not") + ":" + v.(string) + " "
				} else {
					qstring += "@" + k + ":" + v.(string) + " "
				}

			// this picks up any TAG queries
			case []interface{}:
				myPrefixTags := ""
				myFieldTags := k
				if strings.HasSuffix(k, "_not") {
					myPrefixTags = "-"
					myFieldTags = strings.TrimSuffix(k, "_not")
				} else if strings.HasSuffix(k, "_opt") {
					myPrefixTags = "~"
					myFieldTags = strings.TrimSuffix(k, "_opt")
				}
				joined := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(v.([]interface{}))), "|"), "[]")
				qstring += fmt.Sprintf("%s@%s: {%s} ", myPrefixTags, myFieldTags, joined)

			// this picks up any GEO queries
			case map[string]interface{}:
				myPrefix := ""
				myField := k
				if strings.HasSuffix(k, "_not") {
					myPrefix = "-"
					myField = strings.TrimSuffix(k, "_not")
				} else if strings.HasSuffix(k, "_opt") {
					myPrefix = "~"
					myField = strings.TrimSuffix(k, "_opt")
				}
				qstring += fmt.Sprintf("%s@%s: [%f,%f,%f,%s] ", myPrefix,
					myField, v.(map[string]interface{})["lon"].(float64),
					v.(map[string]interface{})["lat"].(float64),
					v.(map[string]interface{})["radius"].(float64),
					v.(map[string]interface{})["unit"].(string))

			case float64:
				if strings.HasSuffix(k, "_gte") {
					qstring += "@" + strings.TrimSuffix(k, "_gte") +
						":[" + fmt.Sprintf("%f", v.(float64)) + ",+inf] "
				} else if strings.HasSuffix(k, "_lte") {
					qstring += "@" + strings.TrimSuffix(k, "_lte") +
						":[-inf," + fmt.Sprintf("%f", v.(float64)) + "] "
				} else if strings.HasSuffix(k, "_btw") {
					qstring += "@" + strings.TrimSuffix(k, "_btw") +
						":[-inf" + fmt.Sprintf("%f", v.(float64)) + "] "
				} else {
					qstring += "@" + k + ":[" + fmt.Sprintf("%f", v.(float64)) +
						"," + fmt.Sprintf("%f", v.(float64)) + "] "
				}
			}

		}
	} else {
		qstring = args["raw_query"].(string)
	}
	argsMap := c.Value("v").(postVars).v

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
		log.Fatal(err)
	}

	for _, doc := range docs {
		res = append(res, doc.Properties)
	}

	return res
}
