package redissearchgraphql

import (
	"fmt"
	"log"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/graphql-go/graphql"
)

var tagInput = graphql.NewList(graphql.String)
var geoInputObject = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "geo",
	Fields: graphql.InputObjectConfigFieldMap{
		"unit": &graphql.InputObjectFieldConfig{
			Type:         graphql.String,
			DefaultValue: "km",
		},
		"lat": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.Float),
		},
		"lon": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.Float),
		},
		"radius": &graphql.InputObjectFieldConfig{
			Type: graphql.NewNonNull(graphql.Float),
		},
	}})

func FtInfo2Schema(client *redisearch.Client, searchidx string) (graphql.Schema, error) {
	idx, err := client.Info()
	var schema graphql.Schema

	if err != nil {
		log.Fatal("cannot do info on index:"+searchidx, " Error: ", err)
	}

	fields := make(graphql.Fields)
	args := make(graphql.FieldConfigArgument)

	// Handle the case of a raw query
	args["raw_query"] = &graphql.ArgumentConfig{
		Type: graphql.String,
	}

	for _, field := range idx.Schema.Fields {
		if field.Type == 0 {
			fields[field.Name] = &graphql.Field{
				Type: graphql.String,
			}
			args[field.Name] = &graphql.ArgumentConfig{
				Type: graphql.String,
			}
			args[fmt.Sprintf("%s_not", field.Name)] = &graphql.ArgumentConfig{
				Type: graphql.String,
			}
			args[fmt.Sprintf("%s_opt", field.Name)] = &graphql.ArgumentConfig{
				Type: graphql.String,
			}
		}

		if field.Type == 1 {
			fields[field.Name] = &graphql.Field{
				Type: graphql.Float,
			}
			args[field.Name] = &graphql.ArgumentConfig{
				Type: graphql.Float,
			}
			args[fmt.Sprintf("%s_gte", field.Name)] = &graphql.ArgumentConfig{
				Type: graphql.Float,
			}
			args[fmt.Sprintf("%s_lte", field.Name)] = &graphql.ArgumentConfig{
				Type: graphql.Float,
			}
			// TODO: handle between!
			args[fmt.Sprintf("%s_btw", field.Name)] = &graphql.ArgumentConfig{
				Type: graphql.String,
			}
		}

		// GEO TYPE
		if field.Type == 2 {
			fields[field.Name] = &graphql.Field{
				Type: graphql.String,
			}
			args[field.Name] = &graphql.ArgumentConfig{
				Type: geoInputObject,
			}
			args[fmt.Sprintf("%s_not", field.Name)] = &graphql.ArgumentConfig{
				Type: geoInputObject,
			}
			args[fmt.Sprintf("%s_opt", field.Name)] = &graphql.ArgumentConfig{
				Type: geoInputObject,
			}
		}

		// TAGS
		if field.Type == 3 {
			fields[field.Name] = &graphql.Field{
				Type: graphql.String,
			}
			args[field.Name] = &graphql.ArgumentConfig{
				Type: tagInput,
			}
			args[fmt.Sprintf("%s_not", field.Name)] = &graphql.ArgumentConfig{
				Type: tagInput,
			}
			args[fmt.Sprintf("%s_opt", field.Name)] = &graphql.ArgumentConfig{
				Type: tagInput,
			}
		}

	}

	var ftType = graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "FT",
			Fields: fields,
		},
	)

	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"ft": &graphql.Field{
					Type: graphql.NewList(ftType),
					Args: args,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return FtSearch(p.Args, client, p.Context), nil
					},
				},
				"raw": &graphql.Field{
					Type: graphql.NewList(ftType),
					Args: args,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return FtSearch(p.Args, client, p.Context), nil
					},
				},
			},
		})
	schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryType,
		},
	)
	return schema, nil
}
