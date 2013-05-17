package ranklib

import (
  "testing"
)


func TestSuggestions(t *testing.T) {
  t.Log("HELLO!")
  trie := NewTrie()
  trie.AddEntry("blah", "blah")
  sug, _ := trie.GetSuggestions("bla")

  if len(sug) == 0 || sug[0] != "blah" {
    t.Errorf("Failed to suggest 'blah' for 'bla'")
  }

  sug, _ = trie.GetSuggestions("blahhh")
  if len(sug) > 0 {
    t.Errorf("Too many suggestions for blahhhh")
  }

  sug, _ = trie.GetSuggestions("dsdsdacdvdv")
  if len(sug) > 0 {
    t.Errorf("Too many suggestions for rando")
  }

  trie.AddEntry("bla", "bla")
  sug, _ = trie.GetSuggestions("bl")
  if len(sug) != 2 {
    t.Errorf("Failed to suggest 'blah' and 'bla' for bl")
  }

  trie.AddEntry("bla你", "bla你")
  sug, _ = trie.GetSuggestions("bl")
  if len(sug) != 3 {
    t.Errorf("Failed to suggest unicode for bl")
  }
  sug, _ = trie.GetSuggestions("bla")
  if len(sug) != 3 {
    t.Errorf("Failed to figure out exact match")
  }

  trie.AddEntry("blee", "blee")
  sug, _ = trie.GetSuggestions("bla")
  if len(sug) != 3 {
    t.Errorf("Oversuggested bla")
  }

  trie.DumpTree()
  t.Logf("Suggestions: %q", sug)
}
