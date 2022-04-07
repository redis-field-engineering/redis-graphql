## Schema Tips and Tricks

Redisearch allows for several [field types](https://redis.io/commands/ft.create/#field-types)

Knowing when to use which field type is key to developing a useful schema.

### Text

Use this type when you are looking to do a full text search and the text field does not contain special characters.

If you require an exact match or have special characters like an email address use a TAG

Text allows prefixes such that searching "jen*" will match "jennifer82" and also allows fuzzy searching such that "%bill%" will match "will", "bill" and "ball"

### Numeric

Numeric types are useful when you want to match a numeric range such as 12 <= X <= 21.

They are also useful in aggregations where you might want to sum a number over a search like tallying up a shopping cart total.

### Geo

Useful for searching within a radius of a longitude and latitude such as find all ramen shops within a 20km distance.

Also useful in aggregations to sort by the closest to a given location.

### Tags

Tags are very efficient and should be used whenever there is a low cardinality or when an exact match is required.

Special characters need to be escaped for example chris@example.com should be

```
query {
     ft( email: ["chris\\@example\\.com"],)
     { username, }
 }
```

