#!/usr/bin/env python

import csv
import random
from faker import Faker

fake=Faker()

headers = ['id','username', 'email', 'mmr', 'experience', 'location', 'group_tags', 'secondary_group_tags', 'blacklist_tags', 'play_style_tags', 'rating']

# generate 1000 grouptags
group_types = [
    'band', 'bunch', 'clan', 'crew', 'clique', 'club', 'cabal', 'combine', 'collective', 'community',
    'faction', 'gang', 'league', 'mob', 'pack', 'posse', 'squad', 'group', 'squad', 'syndicate', 'team', 'trust'
    ]

play_style_tags = {
    0 : ['high_mobile', 'med_mobile', 'med_stationary', 'stationary'],
    1 : ['sniper', 'selective', 'random', 'sprayer'],
    2 : ['high_accuracy', 'med_accuracy', 'low_accuracy'],
}

user_set = set({})

groups = []

for _ in range(1000):
    groups.append("{}_{}".format(fake.color_name().lower(), group_types[random.randint(0,len(group_types)-1)]))


seed_users = []
for x in range(40000):
    username = fake.user_name()
    while username in user_set:
        username = "{}{}".format(fake.user_name(), x % 4096)
    seed_users.append(username)
    user_set.add(username)

with open('users.csv', 'w', newline='') as csvfile:
    writer = csv.DictWriter(csvfile, fieldnames=headers)
    writer.writeheader()
    for usr in range(40000):
        loc = fake.location_on_land(coords_only=True)
        writer.writerow({
            'id' : usr,
            'username' : seed_users[usr],
            'email' : fake.email(),
            'mmr' : random.randint(0,10000),
            'experience' : random.randint(0,10000),
            'location' : "{},{}".format(loc[0], loc[1]),
            'group_tags' : groups[usr % 1000],
            'secondary_group_tags' : groups[(usr + 1) % 1000],
            'blacklist_tags' : groups[(usr + 2) % 1000],
            'play_style_tags' : play_style_tags[random.randint(0,2)][random.randint(0,2)],
            'rating' : random.randint(0,10000)
        })
        email = fake.ascii_email()
        gt = []
        for _ in range(random.randint(0,3)):
            gt.append(groups[random.randint(0, len(groups)-1)])
        sgt = []
        if len(gt) > 0:
            for _ in range(random.randint(0,2)):
                sgt.append(groups[random.randint(0, len(groups)-1)])
        bt = []
        for _ in range(random.randint(0,3)):
            bt.append(seed_users[random.randint(0, len(seed_users)-1)])
        pst = []
        for x in range(random.randint(0,2)):
            pst.append(play_style_tags[x][random.randint(0, len(play_style_tags[x])-1)])

        if usr < len(seed_users):
            username=seed_users[usr]
        else:
            username = fake.user_name()
            # ensure we don't have duplicate user names
            while username in user_set:
                username = "{}{}".format(fake.user_name(), usr % 10240)

        user_set.add(username)

        writer.writerow({
            'id': str(usr).zfill(12),
            'username': username,
            'email': email,
            'mmr': random.randint(1000, 4000),
            'experience': random.randint(50, 750) ,
            'location': "{},{}".format(loc[1], loc[0]),
            'group_tags': ",".join(gt),
            'secondary_group_tags': ",".join(sgt),
            'blacklist_tags': ",".join(bt),
            'play_style_tags': ",".join(pst),
            'rating': random.randint(1, 10),
        })
