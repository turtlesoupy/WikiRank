import json

newjson = []
for u in json.load(open("ranked.json")):
  newjson.append({
    "title": u["Title"],
    "rank": u["Rank"]
  })

json.dump(newjson, open("simplerank.json", "w"))
