/*
THIS CONTAINS PORTIONS OF SOFTWARE BY:

Copyright (c) 2012, Richard Johnson
All rights reserved.

Redistribution and use in source and binary forms, with or without modification,
are permitted provided that the following conditions are met:

 - Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.
 - Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package ranklib

import (
  "fmt"
  "log"
  "math"
  "strings"
  "container/heap"
)

type TriePage struct {
  Title string
  Rank float64
}

func (t*TriePage) GetRank() float64 {
  return t.Rank
}


func newTriePage(rp *RankedPage) *TriePage {
  t :=  &TriePage {
    Title: rp.Page.Title,
    Rank: rp.Rank,
  }
  return t
}

func NormalizeSuggestionPrefix(prefix string) string {
  return strings.ToUpper(prefix)
}

func CreateTrie(inputFile string, limit int) (trie *Trie, err error) {
  trie = NewTrie()
  rpchan := make(chan *RankedPage, 10000)
  go ReadRankedPages(inputFile, rpchan)
  i := 0
  for page := range rpchan {
    tp := newTriePage(page)
    trie.AddEntry(NormalizeSuggestionPrefix(tp.Title), tp)
    if i % 100000 == 0 && i > 0 {
      log.Printf("Inserted page #%d", i)
      sug, _ := trie.GetTopSuggestions("", 10)
      log.Printf("Suggestions 1: %q", sug)
    } else if i > limit {
      break
    }
    i++
  }

  return trie, nil
}

// The remaining is a ranked autocomplete modification of https://github.com/rjohnsondev/go-trie 

const startLetter = '.'
const endLetter = 'Z'
const noLetters= int(endLetter) - int(startLetter) +1 // optimised for english... with a couple left over

type Rankable interface {
  GetRank() float64
}

type ByRank []Rankable

func (this ByRank) Len() int {
    return len(this)
}
func (this ByRank) Less(i, j int) bool {
    return this[i].GetRank() > this[j].GetRank()
}
func (this ByRank) Swap(i, j int) {
    this[i], this[j] = this[j], this[i]
}

type branch struct {
    Children []*branch
    Value Rankable
    Shortcut []byte
    MaxRank float64 // maximum rank value below this node
}

type Trie struct {
    Tree *branch
    UnicodeMap map[int] int
    NextIndex int
}


func (this *Trie) GetKey(ch byte) int {
    ir := int(ch)
    index := -1
    if ir >= startLetter && ir <= endLetter {
        // we use it's key..
        index = ir - startLetter
    } else {
        mapindex, exists := this.UnicodeMap[ir]
        if !exists {
            index = this.NextIndex
            this.NextIndex++
            this.UnicodeMap[ir] = index
        } else {
            index = mapindex
        }
    }
    return index
}

func (this *Trie) EnsureCapacity(children []*branch, index int) []*branch {
    if len(children) < index+1 {
        for x := len(children); x < index+1; x++ {
            children = append(children, nil)
        }
    }
    return children
}

func (this *Trie) AddEntry(entry string, value Rankable) {
    this.AddToBranch(this.Tree, []byte(entry), value)
}

func (this *Trie) AddToBranch(t *branch, remEntry []byte, value Rankable) {
    oldRank := t.MaxRank
    switch r := value.(type) {
      case Rankable:
        t.MaxRank = math.Max(r.GetRank(), t.MaxRank)
    }

    // can we cheat?
    if t.Shortcut == nil {
        t.Shortcut = remEntry
        t.Value = value
        t.Children = nil // not needed, but it helps think things through
        return
    }

    Shortcut := t.Shortcut

    // are we on the right branch yet?
    if len(remEntry) == 0 && len(Shortcut) == 0 {
        // we are here, set it and forget it
        if t.Value == nil || value.GetRank() > t.Value.GetRank() {
          t.Value = value
        }
        return
    } else {

        // find common prefix
        smallestLen := len(remEntry)
        if smallestLen > len(Shortcut) {
            smallestLen = len(Shortcut)
        }
        var x int
        for x = 0; x < smallestLen && Shortcut[x] == remEntry[x]; x++ {

        }
        commonPrefix := Shortcut[0:x]
        if x < len(Shortcut) {
            // we can assign the t to a child
            ttail := Shortcut[x+1:len(Shortcut)]
            tkey := this.GetKey(Shortcut[x])
            newTBranch := &branch {
                Children: t.Children,
                Value: t.Value,
                MaxRank: oldRank,
                Shortcut: ttail,
            }
            t.Children = make([]*branch, noLetters, noLetters)
            t.Children = this.EnsureCapacity(t.Children, tkey)
            t.Children[tkey] = newTBranch
            t.Shortcut = commonPrefix
            t.Value = nil
        } else {
            // the value of t remains
        }
        if x < len(remEntry) {
            // we can assign the v to a child
            vkey := this.GetKey(remEntry[x])
            vtail := remEntry[x+1:len(remEntry)]
            t.Children = this.EnsureCapacity(t.Children, vkey)
            if t.Children[vkey] == nil {
                newVBranch := &branch {
                    Children: nil,
                    Value: nil,
                    Shortcut: nil,
                }
                t.Children[vkey] = newVBranch
            }
            this.AddToBranch(t.Children[vkey], vtail, value)
        } else {
            // the value of v now takes up the position
            if t.Value == nil || value.GetRank() > t.Value.GetRank() {
              t.Value = value
            }
        }

    }
}

func (this *Trie) DumpTree() {
    fmt.Printf("\n\n")
    this.DumpBranch(this.Tree, 1)
}

func (this *Trie) DumpBranch(t *branch, depth int) {
    // attempt to output a textual view of the Tree
    for x := 0; x < depth; x ++ { fmt.Print("  ") }
    fmt.Printf("- cheat: %s\n", t.Shortcut)
    for x := 0; x < depth; x ++ { fmt.Print("  ") }
    fmt.Printf("- Value: %s\n",t.Value)
    for x := 0; x < depth; x ++ { fmt.Print("  ") }
    fmt.Printf("- MaxRank: %f\n",t.MaxRank)
    if t.Children != nil {
        for x := 0; x < depth; x ++ { fmt.Print("  ") }
        fmt.Printf("- children:\n")
        for y := 0; y < len(t.Children); y++ {
            if t.Children[y] != nil {
                for x := 0; x < depth; x ++ { fmt.Print("  ") }
                charb := make([]byte, 1)
                charb[0] = byte(y)+startLetter
                fmt.Printf(" - %s\n", string(charb))
                this.DumpBranch(t.Children[y], depth+1)
            }
        }
    }
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Rankable

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
    return (*pq[i]).GetRank() < (*pq[j]).GetRank()
}

func (pq PriorityQueue) Swap(i, j int) {
        pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
        item := x.(*Rankable)
        *pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
        old := *pq
        n := len(old)
        item := old[n-1]
        *pq = old[0 : n-1]
        return item
}

func (pq *PriorityQueue) Peek() interface{} {
  return (*pq)[0]
}

type topSuggestionQuery struct {
  suggestions *PriorityQueue
  n int
  watermark float64
}


func (this *Trie) GetTopSuggestions(entry string, n int) ([]Rankable, bool) {
  pq := &PriorityQueue{}
  heap.Init(pq)
  tsq := topSuggestionQuery {
    suggestions: pq,
    n: n,
    watermark: -1,
  }

  tsq, ok := this.getTopSuggestions(entry, tsq)
  nSuggestions := tsq.suggestions.Len()
  ret := make([]Rankable, nSuggestions)
  for i := nSuggestions-1; i >= 0; i-- {
    ret[i] = *(heap.Pop(tsq.suggestions).(*Rankable))
  }
  return ret, ok
}

func (this *branch) topValuesBelow(tsq topSuggestionQuery) (topSuggestionQuery) {
  if this.Value != nil {
    if tsq.suggestions.Len() < tsq.n {
      heap.Push(tsq.suggestions, &this.Value)
      tsq.watermark = (*tsq.suggestions.Peek().(*Rankable)).GetRank()
    } else if tsq.watermark < this.Value.GetRank() {
      heap.Pop(tsq.suggestions)
      heap.Push(tsq.suggestions, &this.Value)
      tsq.watermark = (*tsq.suggestions.Peek().(*Rankable)).GetRank()
    }
  }

  for _, branch := range this.Children {
    if branch != nil && branch.MaxRank > tsq.watermark { // Only explore reasonable subTrees
      tsq = branch.topValuesBelow(tsq)
    }
  }

  return tsq
}

func (this *Trie) getTopSuggestions(entry string, tsq topSuggestionQuery) (topSuggestionQuery, bool) {
    t := this.Tree
    eb := []byte(entry)
    // it's <= here to ensure we get to the cheat comparison nil on a valid path
    for x := 0; x <= len(eb); x++ {
        // if the current branch has a cheat, make sure we match it
        s := t.Shortcut
        var y int
        for y = 0; y < len(s); y++ {
            if x+y >= len(eb) {
                //return nil, true
                return t.topValuesBelow(tsq), true
            }
            if s[y] != eb[x+y] {
                // NO WAY, RETURN NOTHING
                return tsq, false
            }
        }
        x += y


        if x < len(eb) {
            // we got through the cheat!
            index := -1
            // we don't use GetKey here as we don't want reads to pollute our hashmap
            ir := int(eb[x])
            if ir >= startLetter && ir <= endLetter {
                index = ir - startLetter
            } else {
                mapindex, exists := this.UnicodeMap[ir]
                if !exists {
                    return tsq, false
                } else {
                    index = mapindex
                }
            }
            if index > len(t.Children)-1 || t.Children[index] == nil {
                return tsq, false
            }
            t = t.Children[index]
            eb = eb[x:]
            x = 0
        }
    }
    return t.topValuesBelow(tsq), true
}


func (this *Trie) GetEntry(entry string) (value Rankable, validPath bool) {
    t := this.Tree
    eb := []byte(entry)
    // it's <= here to ensure we get to the cheat comparison nil on a valid path
    for x := 0; x <= len(eb); x++ {
        // if the current branch has a cheat, make sure we match it
        s := t.Shortcut
        var y int
        for y = 0; y < len(s); y++ {
            if x+y >= len(eb) {
                return nil, true
            }
            if s[y] != eb[x+y] {
                return nil, false
            }
        }
        x += y
        if x < len(eb) {
            // we got through the cheat!
            index := -1
            // we don't use GetKey here as we don't want reads to pollute our hashmap
            ir := int(eb[x])
            if ir >= startLetter && ir <= endLetter {
                index = ir - startLetter
            } else {
                mapindex, exists := this.UnicodeMap[ir]
                if !exists {
                    return nil, false // no mapping :/
                } else {
                    index = mapindex
                }
            }
            if index > len(t.Children)-1 || t.Children[index] == nil {
                return nil, false
            }
            t = t.Children[index]
            eb = eb[x:]
            x = 0
        }
    }
    return t.Value, true
}

func NewTrie() *Trie {
    t := &Trie {
        Tree: &branch {
            Children: nil,
            Value: nil,
            Shortcut: nil,
        },
        UnicodeMap: make(map[int]int),
        NextIndex: 26,
    }
    return t
}
