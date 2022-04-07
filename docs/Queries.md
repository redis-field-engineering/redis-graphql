# Queries

## FT

### Matching 
There is one query type "ft", the full text search.  

This can be called as follows:

```
    ft(
        field1: "value1",
        field2: "value2",
    )
    {
        field3,
        field4,
        field5
    }
}
```

By default this will search for up to 10 records where field1 matches value1 AND field2 matches value2.

Fields 3,4,5 will be returned.

The default of the AND the limit of the records returned can be modified with GraphQL variables

| QUERY | GRAPHQL Variables |
|--|--|
| ```     ft( field1: "value1", field2: "value2",) { field3, field4, field5 } }```| ```{"limit": 2, "ormatch": true}```|

The above will now return up to 2 records where field1 matches value1 OR field2 matches value2


### Matching options 

Based on the type of field additional arguments can be created to allow for better matching

| Field Type | Suffix | Explanation | Example | Details |
| -- | -- | -- | -- | -- |
| Text | _not | Negative Match | username_not: "tohokujunko" | Return records where username does not match "tohokujunko" |
| Text | _opt | Optional Match | username_opt: "tohokujunko" | Return records where username optionally matches "tohokujunko" |
| Numeric | _gte | Greater than or equal to | rating_gte: 5 | Return records where rating >= 5 | 
| Numeric | _lte | Less than or equal to | rating_lte: 2 | Return records where rating <= 2 | 
| Numeric | _bte | Between or equal to | rating_bte: [2, 7] | Return records where 2 <= rating <= 7 | 
| Geo | _not | Negative Match | location_not: {lat: 37.377658, lon: -122.064228, radius: 10, unit: "km"}} | Return records not in a 10km radius of lon/lat |
| Geo | _opt | Optional Match | location_opt: {lat: 37.377658, lon: -122.064228, radius: 10, unit: "km"}} | Return records optionally in a 10km radius of lon/lat |
| Tag | _not | Negative Match |  group_tags_not: ["tag1", "tag2"] | Return record does not have tag1 or tag2 |
| Tag | _opt | Optional Match | group_tags_opt: ["tag1", "tag2"] | Return records optionally if has tag1 or tag2 |

### Tag Matching

When matching tags you will need to escape any special characters or spaces and only exact matches are available


```
query {
     ft( email: ["chris\\@example\\.com"],)
     { username }
 }
```

### Raw Queries

The Redisearch query language can be very powerful and with the options and variables will cover the vast majority of use cases, it is also possible to pass a very complex matching query using raw_query

```
   query {
        ft(
            raw_query: "(@field1:value1 @field2:[1, 100])|(@field2:val2 ~@field7:val7)",
        )
        {
            field1,
            field10,
            field11,
            field12
        }
    }

```