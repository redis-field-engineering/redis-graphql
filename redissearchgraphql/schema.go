package redissearchgraphql

import (
	"fmt"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/graphql-go/graphql"
)

// betweenInput is the input for the numeric between filter
var betweenInput = graphql.NewList(graphql.Float)

// tagInput is the input for the tag filter
var tagInput = graphql.NewList(graphql.String)

// rawAggPlan is an array of strings that we use to build the raw aggregation plan
var rawAggPlan = graphql.NewList(graphql.String)

// geoInputObject is a specialized input object for the geo filter
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

// FtInfo2Schma uses the redisearch-go library to generate a graphql schema
// from the redisearch index.
// see https://redis.io/commands/ft.info/ for more information
//   This adds some extra fields to the schema that are not part of the RediSearch schema
//   All of the extra fields are prefixed with "_"
func FtInfo2Schema(client *redisearch.Client, searchidx string) (graphql.Schema, SchemaDocs, error) {
	idx, err := client.Info()
	var schema graphql.Schema
	var docs SchemaDocs = *NewSchemaDocs()
	docs.IndexName = searchidx

	if err != nil {
		return schema, docs, err
	}

	fields := make(graphql.Fields)
	args := make(graphql.FieldConfigArgument)

	// Add a new argument "raw_query" that will allow us to pass in a raw RediSearch query
	args["raw_query"] = &graphql.ArgumentConfig{
		Type: graphql.String,
	}

	for _, field := range idx.Schema.Fields {

		// For RediSearch TEXT fields we add a new argument and field name that matches
		// the field name in the RediSearch schema
		// Additionally we will add _not and _opt suffixes to the field name
		// to indicate the type of search to be performed
		if field.Type == 0 {
			docs.Strings = append(docs.Strings, field.Name)
			docs.StringSuffix = append(docs.StringSuffix, "not", "opt")
			docs.FieldDocs[field.Name] = "Find records where " + field.Name + " matches STRING"
			docs.FieldDocs[fmt.Sprintf("%s_not", field.Name)] = "Find records where " + field.Name + " does not match STRING"
			docs.FieldDocs[fmt.Sprintf("%s_opt", field.Name)] = "Optionally find records where " + field.Name + " matches STRING"
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

		// For RediSearch NUMERIC fields we add a new argument and field name that matches
		// the field name in the RediSearch schema
		// _gte suffix indicates greater than or equal to
		// _lte suffix indicates less than or equal to
		// _bte suffix indicates between or equal to a list of 2 numbers
		if field.Type == 1 {
			docs.Floats = append(docs.Floats, field.Name)
			docs.FloatSuffix = append(docs.FloatSuffix, "gte", "lte", "bte")
			docs.FieldDocs[field.Name] = "Find records where " + field.Name + " == NUMBER"
			docs.FieldDocs[fmt.Sprintf("%s_gte", field.Name)] = "Find records where " + field.Name + " >=  NUMBER"
			docs.FieldDocs[fmt.Sprintf("%s_lte", field.Name)] = "Find records where " + field.Name + " <= NUMBER"
			docs.FieldDocs[fmt.Sprintf("%s_bte", field.Name)] = "Find records where " + field.Name + " between NUMBER1 and NUMBER2"
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
			args[fmt.Sprintf("%s_bte", field.Name)] = &graphql.ArgumentConfig{
				Type: betweenInput,
			}
		}

		// For RediSearch GEO fields we add a new argument and field name that matches
		// the field name in the RediSearch schema
		// This requires using the geoInputObject for search
		// _not and _opt suffixes indicate the type of search to be performed
		if field.Type == 2 {
			docs.Geos = append(docs.Geos, field.Name)
			docs.GeoSuffix = append(docs.GeoSuffix, "not", "opt")
			docs.FieldDocs[field.Name] = "Find records where " + field.Name + " is within radius of lon,lat"
			docs.FieldDocs[fmt.Sprintf("%s_not", field.Name)] = "Find records where " + field.Name + " is not within radius of lon,lat"
			docs.FieldDocs[fmt.Sprintf("%s_opt", field.Name)] = "Optional find records where " + field.Name + " is optionally within radius of lon,lat"
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

		// For RediSearch GEO fields we add a new argument and field name that matches
		// the field name in the RediSearch schema
		// This requires using the tagInputObject for search
		if field.Type == 3 {
			docs.Tags = append(docs.Tags, field.Name)
			docs.FieldDocs[field.Name] = "Find records where " + field.Name + " == TAG"
			docs.FieldDocs[fmt.Sprintf("%s_not", field.Name)] = "Find records where " + field.Name + " != TAG"
			docs.FieldDocs[fmt.Sprintf("%s_opt", field.Name)] = "Optional find records where " + field.Name + " == TAG"

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

		// Add all of the specialization aggregation result fields
		fields["_agg_groupby_count"] = &graphql.Field{
			Type: graphql.Int,
		}
		fields["_agg_groupby_num"] = &graphql.Field{
			Type: graphql.Float,
		}
		args["_agg_groupby"] = &graphql.ArgumentConfig{
			Type: graphql.String,
		}
		args["_agg_num_field"] = &graphql.ArgumentConfig{
			Type: graphql.String,
		}
		args["_agg_num_function"] = &graphql.ArgumentConfig{
			Type: graphql.String,
		}
		args["_agg_num_quantile"] = &graphql.ArgumentConfig{
			Type: graphql.Float,
		}
		args["raw_agg_plan"] = &graphql.ArgumentConfig{
			Type: rawAggPlan,
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
				// Define the basic FT.SEARCH query - used to be ft
				searchidx: &graphql.Field{
					Type: graphql.NewList(ftType),
					Args: args,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						res, err := FtSearch(p.Args, client, p.Context)
						return res, err
					},
				},
				// Define the basic FT.AGGREGATE and GROUPBY/COUNT query - used to be agg_count
				fmt.Sprintf("%sAggCount", searchidx): &graphql.Field{
					Type: graphql.NewList(ftType),
					Args: args,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						res, err := FtAggCount(p.Args, client, p.Context)
						return res, err
					},
				},
				// Define the basic FT.AGGREGATE with numeric filters - used to be agg_numgroup
				fmt.Sprintf("%sAggNumGroup", searchidx): &graphql.Field{
					Type: graphql.NewList(ftType),
					Args: args,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						res, err := FtAggNumGroup(p.Args, client, p.Context)
						return res, err
					},
				},
				// Define the raw FT.AGGREGATE query - used to be agg_raw
				fmt.Sprintf("%sAggRaw", searchidx): &graphql.Field{
					Type: graphql.NewList(ftType),
					Args: args,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						res, err := FtAggRaw(p.Args, client, p.Context)
						return res, err
					},
				},
			},
		})

	// Set the Schema
	schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryType,
		},
	)
	return schema, docs, nil
}
