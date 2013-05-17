package ranklib

import (
  "testing"
)


type RankString struct {
  s string
  r float64
}

func (this RankString) GetRank() float64 {
  return this.r
}

func TestRankedSuggestions(t *testing.T) {
  trie := NewTrie()
  trie.AddEntry("hello1", RankString{"hello1", 1.0})
  trie.DumpTree()
  trie.AddEntry("hello2", RankString{"hello2", 2.0})
  trie.DumpTree()
  trie.AddEntry("hello3", RankString{"hello3", 3.0})
  sug, _ := trie.GetTopSuggestions("h", 1)
  trie.DumpTree()
  t.Logf("Suggestions: %q", sug)
  trie.AddEntry("hello12", RankString{"hello12", 12.0})
  sug, _ = trie.GetTopSuggestions("h", 1)
  trie.DumpTree()
  t.Logf("Suggestions: %q", sug)
}

func TestUnrankedSuggestions(t *testing.T) {
  trie := NewTrie()
  trie.AddEntry("blah", RankString{"blah", 1})
  sug, _ := trie.GetTopSuggestions("bla", 10)

  if len(sug) == 0 {
    t.Errorf("Failed to suggest 'blah' for 'bla'")
  }

  sug, _ = trie.GetTopSuggestions("blahhh", 10)
  if len(sug) > 0 {
    t.Errorf("Too many suggestions for blahhhh")
  }

  sug, _ = trie.GetTopSuggestions("dsdsdacdvdv", 10)
  if len(sug) > 0 {
    t.Errorf("Too many suggestions for rando")
  }

  trie.AddEntry("bla", RankString{"bla", 2})
  sug, _ = trie.GetTopSuggestions("bl", 10)
  if len(sug) != 2 {
    t.Errorf("Failed to suggest 'blah' and 'bla' for bl")
  }

  trie.AddEntry("bla你", RankString{"bla你", 2})
  sug, _ = trie.GetTopSuggestions("bl", 10)
  if len(sug) != 3 {
    t.Errorf("Failed to suggest unicode for bl")
  }
  sug, _ = trie.GetTopSuggestions("bla", 10)
  if len(sug) != 3 {
    t.Errorf("Failed to figure out exact match")
  }

  trie.AddEntry("blee", RankString{"blee", 2})
  sug, _ = trie.GetTopSuggestions("bla", 10)
  if len(sug) != 3 {
    t.Errorf("Oversuggested bla")
  }

  trie.DumpTree()
  t.Logf("Suggestions: %q", sug)
}
