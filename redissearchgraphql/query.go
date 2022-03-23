package redissearchgraphql

import (
	"fmt"
	"strings"
)

func QueryBuilder(args, argsMap map[string]interface{}, wildcard bool) (string, error) {
	qstring := ""
	query_conditions := []string{}
	if args["raw_query"] == nil {

		for k, v := range args {
			if strings.HasPrefix(k, "_agg_") {
				// Exclude anything that starts with _agg_
				continue
			} else {
				switch v.(type) {
				case string:
					if strings.HasSuffix(k, "_not") {
						query_conditions = append(query_conditions, "-@"+strings.TrimSuffix(k, "_not")+":"+v.(string))
					} else if strings.HasSuffix(k, "_opt") {
						query_conditions = append(query_conditions, "~@"+strings.TrimSuffix(k, "_not")+":"+v.(string))
					} else {
						query_conditions = append(query_conditions, "@"+k+":"+v.(string))
					}

				// this picks up any TAG or between queries
				case []interface{}:
					myPrefixTags := ""
					myFieldTags := k
					// we don't want to pick up the raw_sggregation plan here
					if k != "raw_agg_plan" {
						if strings.HasSuffix(k, "_not") {
							myPrefixTags = "-"
							myFieldTags = strings.TrimSuffix(k, "_not")
							joined := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(v.([]interface{}))), "|"), "[]")
							query_conditions = append(query_conditions, fmt.Sprintf("%s@%s:{%s}", myPrefixTags, myFieldTags, joined))
						} else if strings.HasSuffix(k, "_opt") {
							myPrefixTags = "~"
							myFieldTags = strings.TrimSuffix(k, "_opt")
							joined := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(v.([]interface{}))), "|"), "[]")
							query_conditions = append(query_conditions, fmt.Sprintf("%s@%s:{%s}", myPrefixTags, myFieldTags, joined))
						} else if strings.HasSuffix(k, "_bte") {
							myFieldTags = strings.TrimSuffix(k, "_bte")
							query_conditions = append(query_conditions, fmt.Sprintf("@%s:[%f, %f]", myFieldTags, v.([]interface{})[0], v.([]interface{})[1]))
						} else {
							joined := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(v.([]interface{}))), "|"), "[]")
							query_conditions = append(query_conditions, fmt.Sprintf("%s@%s:{%s}", myPrefixTags, myFieldTags, joined))
						}
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
					query_conditions = append(query_conditions, fmt.Sprintf("%s@%s: [%f,%f,%f,%s]", myPrefix,
						myField, v.(map[string]interface{})["lon"].(float64),
						v.(map[string]interface{})["lat"].(float64),
						v.(map[string]interface{})["radius"].(float64),
						v.(map[string]interface{})["unit"].(string)))

				case float64:
					if strings.HasSuffix(k, "_gte") {
						query_conditions = append(query_conditions, "@"+strings.TrimSuffix(k, "_gte")+
							":["+fmt.Sprintf("%f", v.(float64))+",+inf]")
					} else if strings.HasSuffix(k, "_lte") {
						query_conditions = append(query_conditions, "@"+strings.TrimSuffix(k, "_lte")+
							":[-inf,"+fmt.Sprintf("%f", v.(float64))+"]")
					} else {
						query_conditions = append(query_conditions, "@"+k+":["+fmt.Sprintf("%f", v.(float64))+
							","+fmt.Sprintf("%f", v.(float64))+"]")
					}
				}

			}

			qstring = strings.Join(query_conditions, " ")

			if ormatch, ok := argsMap["ormatch"]; ok {
				if ormatch.(bool) {
					qstring = "(" + strings.Join(query_conditions, " | ") + ")"
				}
			}

		}
	} else {
		qstring = args["raw_query"].(string)
	}
	if wildcard && len(query_conditions) < 1 {
		qstring = "*"
	}
	return qstring, nil
}
