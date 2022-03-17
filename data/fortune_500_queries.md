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
            hqstate_not: "ca|ny"
            tags: ["technology"],
            _agg_groupby: "hqstate",
            )
        {
            hqstate,
            _agg_groupby_count,
        }
    }
```
