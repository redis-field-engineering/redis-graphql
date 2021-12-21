--- wrk -c 10 -d 60s -s example.lua http://localhost:8080/graphql
wrk.method = "POST"
wrk.body = '{"query": "{ ft(hqstate:\\\"ca\\\", hqcity:\\\"san\\\", sector: \\\"Technology\\\") { company,ceo,sector,hqcity,hqstate } }" }'
wrk.headers["Content-Type"] = "application/json"
