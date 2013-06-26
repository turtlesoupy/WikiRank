# coding=utf-8
from __future__ import print_function

import re
import sys
import json
import time
import codecs
import urllib

FORBIDDEN_PAGES = set(["Television film", "Television movie", "Film genre", "Documentary film", "James Bond", "Star Trek", "Star Wars", "Back to the Future",
  "Shoot 'Em up", "V", "Die Fledermaus", "Robin of Locksley", "The Decalogue", "2001 film", "Mad Max", "Carry On (film series)",
  "Detective Conan"
])
FORBIDDEN_RE = re.compile(r"^(Lists? of|Template( talk)?:|\d\d\d\ds? in film|\d\d\d\d film$)")
WIKI_ARTICLES = [
  "List of films: numbers",
  "List of films: A",
  "List of films: B",
  "List of films: C",
  "List of films: D",
  "List of films: E",
  "List of films: F",
  "List of films: G",
  "List of films: H",
  "List of films: I",
  "List of films: J–K",
  "List of films: L",
  "List of films: M",
  "List of films: N–O",
  "List of films: P",
  "List of films: Q–R",
  "List of films: S",
  "List of films: T",
  "List of films: U–W",
  "List of films: X–Z"
]

def orderedUniq(gen):
  seen = set()
  for item in gen:
    if item in seen:
      print("Skipping duplicate %s" % item)
    else:
      seen.add(item)
      yield item

def extractMovies():
  for article in WIKI_ARTICLES:
    print("Fetching %s" % article)
    p = urllib.urlopen("http://en.wikipedia.org/w/api.php?%s" % urllib.urlencode({
      "action": "parse",
      "page": article,
      "format": "json",
    }))

    j = json.loads(p.read())
    for link in linksFromResponse(j):
      yield link
    time.sleep(1)

def linksFromResponse(response):
  p = response["parse"]
  links = p["links"]
  for link in links:
    linkTitle = link["*"]
    if FORBIDDEN_RE.match(linkTitle):
      print("Skipping %s due to regular expression" % linkTitle)
    elif linkTitle in FORBIDDEN_PAGES:
      print("Skipping %s due to forbidden" % linkTitle)
    else:
      yield linkTitle

def main(outputPath):
  with codecs.open(outputPath, "w", "utf-8") as f:
    f.write("\n".join(orderedUniq(extractMovies())))

if __name__ == "__main__":
  if len(sys.argv) != 2:
    print("Missing required argument 'output.txt'")
    sys.exit(1)
  else:
    main(sys.argv[1])
