### Example searches

Find me all not technology companies raned betwen 100 and 150 with equity greater than 3000 and prefer Financial companies

```
   query {
        ft(
            tags_not: ["technology"],
            rank_bte: [100, 150],
            equity_gte: 3000
            sector_opt: "Financials",
                        )
        {
            company,
            hqcity,
            ticker,
        }
    }
```


### Example aggregations

Count all technology companies not headquarted in CA or NY and group and count by state:

```
   query {
        agg_count(
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

Numerical aggregation

```
   query {
        agg_numgroup(
            _agg_groupby: "hqcity",
            _agg_num_field: "revenues",
            _agg_num_function: "sum",
        )
        {
            hqcity,
            _agg_groupby_num,
        }
    }
```

Raw aggregation

```   
    query {
        agg_raw(
            hqstate_not: "dc",
            raw_agg_plan: ["GROUPBY","1","@hqstate","REDUCE","QUANTILE","2","assets","0.99" "AS","_agg_groupby_num","SORTBY","2","@_agg_groupby_num","DESC"]
                        )
        {
            hqstate,
            _agg_groupby_num,
        }
    }
```
