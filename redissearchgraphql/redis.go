package redissearchgraphql

import (
	"context"
	"fmt"
	"strings"

	"github.com/RediSearch/redisearch-go/redisearch"
)

func FtSearch(args map[string]interface{}, client *redisearch.Client, c context.Context) ([]map[string]interface{}, error) {
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

			// this picks up any TAG or between queries
			case []interface{}:
				myPrefixTags := ""
				myFieldTags := k
				if strings.HasSuffix(k, "_not") {
					myPrefixTags = "-"
					myFieldTags = strings.TrimSuffix(k, "_not")
					joined := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(v.([]interface{}))), "|"), "[]")
					qstring += fmt.Sprintf("%s@%s: {%s} ", myPrefixTags, myFieldTags, joined)
				} else if strings.HasSuffix(k, "_opt") {
					myPrefixTags = "~"
					myFieldTags = strings.TrimSuffix(k, "_opt")
					joined := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(v.([]interface{}))), "|"), "[]")
					qstring += fmt.Sprintf("%s@%s: {%s} ", myPrefixTags, myFieldTags, joined)
				} else if strings.HasSuffix(k, "_bte") {
					myFieldTags = strings.TrimSuffix(k, "_bte")
					qstring += fmt.Sprintf("@%s: [%f, %f] ", myFieldTags, v.([]interface{})[0], v.([]interface{})[1])
				}

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
					//  TODO fix the bte
				} else if strings.HasSuffix(k, "_bte") {
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
	argsMap := c.Value("v").(PostVars).Variables

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
		return res, err
	}

	for _, doc := range docs {
		res = append(res, doc.Properties)
	}

	return res, nil
}
