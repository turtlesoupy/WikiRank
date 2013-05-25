from __future__ import print_function

import sys
import time
import urllib
import codecs
from BeautifulSoup import BeautifulSoup

URLS = ["http://www.4icu.org/reviews/index%d.htm" % i for i in xrange(1,28)]

def extract_universities(htmlDoc):
  soup = BeautifulSoup(htmlDoc, convertEntities=BeautifulSoup.HTML_ENTITIES)
  return [link.string.strip() for link in soup("a") \
    if link.string and link['href'] and link['href'].startswith("http://www.4icu.org/reviews/")]

def main(outputPath):
  with codecs.open(outputPath, "w", "utf-8") as f:
    for url in URLS:
      p = urllib.urlopen(url)
      f.write("\n".join(extract_universities(p.read())))
      f.write("\n")
      f.flush()
      time.sleep(1)

if __name__ == "__main__":
  if len(sys.argv) != 2:
    print("Missing required argument 'output.txt'")
  else:
    main(sys.argv[1])
