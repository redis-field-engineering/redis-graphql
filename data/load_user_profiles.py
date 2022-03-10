#!/usr/bin/env python

import csv
import sys
from redisearch import Client, IndexDefinition, TextField, NumericField, TagField, GeoField


def load_data(redis_server, redis_port, redis_password):

    load_client = Client(
      'matchmaking',
      host=redis_server,
      password=redis_password,
      port=redis_port
   )

    definition = IndexDefinition(
      prefix=['user:'],
      language='English',
    )

    load_client.create_index(
      (
          TextField('id'),
          TextField('email'),
          TextField('username'),
          NumericField('mmr', sortable=True),
          NumericField('experience', sortable=True),
          NumericField('rating', sortable=True),
          TagField('group_tags'),
          TagField('secondary_group_tags'),
          TagField('blacklist_tags'),
          TagField('play_style_tags'),
          GeoField('location'),
      ),
      definition=definition)
  
  
    pipe = load_client.redis.pipeline(transaction=False)
  
    if len(sys.argv) > 1:
        load_file = sys.argv[1]
    else:
        load_file = 'users.csv'

    with open(load_file, newline='') as csvfile:
        reader = csv.DictReader(csvfile)
        row_count = 0
        for row in reader:
            pipe.hset("user:%s" %(row['username']), mapping = row)
            row_count += 1
            if row_count % 500 == 0:
                pipe.execute()

    pipe.execute()
  
if __name__ == "__main__":
   load_data(
      redis_server='localhost',
      redis_port=6379,
      redis_password=None
   )
