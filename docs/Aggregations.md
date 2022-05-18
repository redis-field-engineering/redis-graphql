# Aggregations

Aggregations are useful for building reports or calculations on data.

For example group together all records by a field then count by field and sort.

There are 3 different aggregations available and the examples are part of the Fortune 500 dataset in the data directory.

For the examples below, we will be using the RediSearch index Fortune500

## (INDEX_NAME)AggCount

This is a group by and count returned in descending order.

It requires the argument *_agg_groupby* as the field to group by and returns the field *_agg_groupby_count*

```
query {
     Fortune500AggCount(
         hqstate_not: "(ca|ny)"
         tags: ["technology"],
         _agg_groupby: "hqstate",
         )
     {
         hqstate,
         _agg_groupby_count,
     }
 }
```

It will accept the matching query of hqstate not either CA or NY, requires the tag technology and will group and count by state.

The data returned will look like:

```
{
    "data": {
        "Fortune500AggCount": [
            {
                "_agg_groupby_count": 53,
                "hqstate": "ny"
            },
            {
                "_agg_groupby_count": 50,
                "hqstate": "tx"
            },
            {
                "_agg_groupby_count": 36,
                "hqstate": "il"
            },
	    ...
```

## (INDEX_NAME)AggNumGroup

This provides the following numeric functions that can be used to roll up numeric fields

It requires the arguments:

```
_agg_groupby
_agg_num_field
_agg_num_function
```

If using a quantile the following argument is required:

```
_agg_num_quantile
```


|Function|Example|
|--|--|
|Sum of all records|_agg_num_function: "sum"|
|Minimum record|_agg_num_function: "min"|
|Maxiumum record|_agg_num_function: "max"|
|Average record|_agg_num_function: "avg"|
|Standard deviation of all records|_agg_num_function: "stddev"|
|Quantile of all records|_agg_num_function: "quantile", _agg_num_quantile: 0.99|

and returns the following fields

```
_agg_groupby_num
```

Example 1:

```
query {
     Fortune500AggNumGroup(
         _agg_groupby: "hqcity",
         _agg_num_field: "revenues",
         _agg_num_function: "sum",        )
     {
         hqcity,
         _agg_groupby_num,
     }
 }
```

Returns

```
{
    "data": {
        "Fortune500AggNumGroup": [
            {
                "_agg_groupby_num": 1108104,
                "hqcity": "new york"
            },
            {
                "_agg_groupby_num": 485873,
                "hqcity": "bentonville"
            },
            {
                "_agg_groupby_num": 335881,
                "hqcity": "san francisco"
            },

```


## (INDEX_NAME)AggRaw

As with querying we do allow for the possiblity of passing raw aggregations however these are complex

This requires a *raw_agg_plan* argument as an array of strings

Example:

```
    query {
        Fortune500AggRaw(
            hqstate_not: "dc",
            raw_agg_plan: 
	    ["GROUPBY","1","@hqstate","REDUCE",
	    "QUANTILE","2","assets","0.99" "AS",
	    "_agg_groupby_num","SORTBY","2","@_agg_groupby_num","DESC"]
                        )
        {
            hqstate,
            _agg_groupby_num,
        }
    }
```
