## Sample Data and tools

### Fortune500

This will load the Fortune 500 data into an index called fortune500_v1

[This](./fortune_500_queries.md) contains some sample queries on the data that can just be pasted into [Postman](https://www.postman.com/downloads/)

To load the data

```
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
dataload.py
```

### Game matchmaking

This helps to generate some data to load for Game Matchmaking
