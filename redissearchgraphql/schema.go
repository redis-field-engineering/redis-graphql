package redissearchgraphql

import (
	"fmt"
	"log"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/graphql-go/graphql"
)

var betweenInput = graphql.NewList(graphql.Float)
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

func FtInfo2Schema(client *redisearch.Client, searchidx string) (graphql.Schema, SchemaDocs, error) {
	idx, err := client.Info()
	var schema graphql.Schema
	var docs SchemaDocs = *NewSchemaDocs()
	docs.IndexName = searchidx

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

		// Strings
		if field.Type == 0 {
			docs.Strings = append(docs.Strings, field.Name)
			docs.StringSuffix = append(docs.StringSuffix, "not", "opt")
			docs.FieldDocs[field.Name] = "Find documents where " + field.Name + " == STRING"
			docs.FieldDocs[fmt.Sprintf("%s_not", field.Name)] = "Find documents where " + field.Name + " != STRING"
			docs.FieldDocs[fmt.Sprintf("%s_opt", field.Name)] = "Optionally find documents where " + field.Name + " == STRING"
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

		// Numeric
		if field.Type == 1 {
			docs.Floats = append(docs.Floats, field.Name)
			docs.FloatSuffix = append(docs.FloatSuffix, "gte", "lte", "bte")
			docs.FieldDocs[field.Name] = "Find documents where " + field.Name + " == NUMBER"
			docs.FieldDocs[fmt.Sprintf("%s_gte", field.Name)] = "Find documents where " + field.Name + " >=  NUMBER"
			docs.FieldDocs[fmt.Sprintf("%s_lte", field.Name)] = "Find documents where " + field.Name + " <= NUMBER"
			docs.FieldDocs[fmt.Sprintf("%s_bte", field.Name)] = "Find documents where " + field.Name + " between NUMBER1 and NUMBER2"
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
			args[fmt.Sprintf("%s_bte", field.Name)] = &graphql.ArgumentConfig{
				Type: betweenInput,
			}
		}

		// GEO TYPE
		if field.Type == 2 {
			docs.Geos = append(docs.Geos, field.Name)
			docs.GeoSuffix = append(docs.GeoSuffix, "not", "opt")
			docs.FieldDocs[field.Name] = "Find documents where " + field.Name + " is within radius of lon,lat"
			docs.FieldDocs[fmt.Sprintf("%s_not", field.Name)] = "Find documents where " + field.Name + " is not within radius of lon,lat"
			docs.FieldDocs[fmt.Sprintf("%s_opt", field.Name)] = "Optional find documents where " + field.Name + " is optionally within radius of lon,lat"
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
			docs.Tags = append(docs.Tags, field.Name)
			docs.FieldDocs[field.Name] = "Find documents where " + field.Name + " == TAG"
			docs.FieldDocs[fmt.Sprintf("%s_not", field.Name)] = "Find documents where " + field.Name + " != TAG"
			docs.FieldDocs[fmt.Sprintf("%s_opt", field.Name)] = "Optional find documents where " + field.Name + " == TAG"

			docs.TagSuffix = append(docs.TagSuffix, "not", "opt")
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
						res, err := FtSearch(p.Args, client, p.Context)
						return res, err
					},
				},
			},
		})
	schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryType,
		},
	)
	return schema, docs, nil
}
