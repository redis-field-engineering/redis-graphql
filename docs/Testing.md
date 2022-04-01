### Testing


#### Start a redisearch docker container

```
docker-compose up
```

#### Load some sample data

```
python -m venv .venv
source .venv/bin/activate/
cd data
pip -r install requirements.txt
./dataload.py
```

#### Run the proxy

```
# for help use -h 

go run searchql.go  --redis-host redisnode.com \
                    --redis-port 6380 \
                    --redis-index MySearchIndex
                  
```

#### Read the auto-generated documents with sample queries

[Auto-generated Documentation](http://localhost:8080/docs)


#### Query away!

```
curl -X POST  -H "Content-Type: application/json" \
   --data '{"query": "{ ft(hqstate:\"ca\", hqcity:\"san\", sector: \"Technology\") { company,ceo,sector,hqcity,hqstate } }" }' \
  http://localhost:8080/graphql
```
