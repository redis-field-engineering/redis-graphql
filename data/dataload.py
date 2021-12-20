#! /usr/bin/env python

from redisearch import Client, IndexDefinition, TextField, NumericField, TagField
import csv

def load_data(redis_server, redis_port, redis_password):
   load_client = Client(
      'fortune500-v1',
      host=redis_server,
      password=redis_password,
      port=redis_port
   )
   
   definition = IndexDefinition(
           prefix=['fortune500:'],
           language='English',
           score_field='title',
           score=0.5
           )
   load_client.create_index(
           (
               TextField("title", weight=5.0),
               TextField('website'),
               TextField('company'),
               NumericField('employees', sortable=True),
               TextField('industry', sortable=True),
               TextField('sector', sortable=True),
               TextField('hqcity', sortable=True),
               TextField('hqstate', sortable=True),
               TextField('ceo'),
               TextField('ceoTitle'),
               NumericField('rank', sortable=True),
               NumericField('assets', sortable=True),
               NumericField('revenues', sortable=True),
               NumericField('profits', sortable=True),
               NumericField('equity', sortable=True),
               TagField('tags'),
               TextField('ticker')
               ),        
       definition=definition)

   with open('./fortune500.csv', encoding='utf-8') as csv_file:
      csv_reader = csv.reader(csv_file, delimiter=',')
      line_count = 0
      for row in csv_reader:
         if line_count > 0:
            load_client.redis.hset(
                    "fortune500:%s" %(row[1].replace(" ", '')),
                    mapping = {
                        'title': row[1],
                        'company': row[1],
                        'rank': row[0],
                        'website': row[2],
                        'employees': row[3],
                        'sector': row[4],
                        'tags': ",".join(row[4].replace('&', '').replace(',', '').replace('  ', ' ').split()).lower(),
                        'industry': row[5],
                        'hqcity': row[8],
                        'hqstate': row[9],
                        'ceo': row[12],
                        'ceoTitle': row[13],
                        'ticker': row[15],
                        'revenues': row[17],
                        'profits': row[19],
                        'assets': row[21],
                        'equity': row[22]

               })
         line_count += 1
   # Finally Create the alias
   load_client.aliasadd("idx")


if __name__ == "__main__":
   load_data(
      redis_server='localhost',
      redis_port=6379,
      redis_password=None
   )
