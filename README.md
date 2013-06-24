#WikiRank 

This is an investigation into performing PageRank on Wikipedia articles 
in order to quantify _influence_ of different items (e.g. movies or universities).
It is a bit of a mess, but you can roughly follow the code from
* `ranklib/extractgraph.go`: Extracts links from a Wikipedia dump
* `ranklib/ranker.go`: Wrapper for calculating PageRank
* `ranklib/pagerank.go`: Performs PageRank via a power iteration method and outputs 
  a rank vector for all pages
* `wikirank.go`: Main executable / wrapper for the pipeline

A more detailed description and example output is available on my blog: 
http://blog.argteam.com/coding/university-ranking-wikipedia/
