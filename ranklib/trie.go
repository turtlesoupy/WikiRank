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
)

type TrieNode struct {
  Prefix string
  Data *TriePage
  Children []TrieNode
}

type TriePage struct {
  Title string
  Rank float64
}


func newTriePage(rp *RankedPage) *TriePage {
  t :=  &TriePage {
    Title: rp.Page.Title,
    Rank: rp.Rank,
  }
  return t
}

func CreateAndWriteTrie(inputFile string, outputFile string) (err error) {
  trie := NewTrie()

  rpchan := make(chan *RankedPage, 10000)
  go ReadRankedPages(inputFile, rpchan)
  i := 0
  for page := range rpchan {
    tp := newTriePage(page)

    trie.AddEntry(tp.Title, tp)
    trie.GetSuggestions(tp.Title)

    if i % 10000 == 0 {
      log.Printf("Inserted page #%d", i)
    }
    i++
  }

  return
}

// The remaining is a ranked autocomplete modification of https://github.com/rjohnsondev/go-trie 

const startLetter = '.'
const endLetter = 'Z'
const noLetters= int(endLetter) - int(startLetter) +1 // optimised for english... with a couple left over

type Rankable interface {
  Rank() float64
}

type branch struct {
    children []*branch
    value interface{}
    shortcut []byte
    maxRank float64 // maximum rank value below this node
}

type Trie struct {
    tree *branch
    unicodeMap map[int] int
    nextIndex int
}


func (this *Trie) GetKey(ch byte) int {
    ir := int(ch)
    index := -1
    if ir >= startLetter && ir <= endLetter {
        // we use it's key..
        index = ir - startLetter
    } else {
        mapindex, exists := this.unicodeMap[ir]
        if !exists {
            index = this.nextIndex
            this.nextIndex++
            this.unicodeMap[ir] = index
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

func (this *Trie) AddEntry(entry string, value interface{}) {
    this.AddToBranch(this.tree, []byte(entry), value)
}

func (this *Trie) AddToBranch(t *branch, remEntry []byte, value interface{}) {
    //Update ranking
    switch r := value.(type) { case Rankable: t.maxRank = math.Max(r.Rank(), t.maxRank) }

    // can we cheat?
    if t.shortcut == nil {
        t.shortcut = remEntry
        t.value = value
        t.children = nil // not needed, but it helps think things through
        return
    }

    shortcut := t.shortcut

    // are we on the right branch yet?
    if len(remEntry) == 0 && len(shortcut) == 0 {
        // we are here, set it and forget it
        t.value = value
        return
    } else {

        // find common prefix
        smallestLen := len(remEntry)
        if smallestLen > len(shortcut) {
            smallestLen = len(shortcut)
        }
        var x int
        for x = 0; x < smallestLen && shortcut[x] == remEntry[x]; x++ {

        }
        commonPrefix := shortcut[0:x]
        if x < len(shortcut) {
            // we can assign the t to a child
            ttail := shortcut[x+1:len(shortcut)]
            tkey := this.GetKey(shortcut[x])
            newTBranch := &branch {
                children: t.children,
                value: t.value,
                shortcut: ttail,
            }
            t.children = make([]*branch, noLetters, noLetters)
            t.children = this.EnsureCapacity(t.children, tkey)
            t.children[tkey] = newTBranch
            t.shortcut = commonPrefix
            t.value = nil
        } else {
            // the value of t remains
        }
        if x < len(remEntry) {
            // we can assign the v to a child
            vkey := this.GetKey(remEntry[x])
            vtail := remEntry[x+1:len(remEntry)]
            t.children = this.EnsureCapacity(t.children, vkey)
            if t.children[vkey] == nil {
                newVBranch := &branch {
                    children: nil,
                    value: nil,
                    shortcut: nil,
                }
                t.children[vkey] = newVBranch
            }
            this.AddToBranch(t.children[vkey], vtail, value)
        } else {
            // the value of v now takes up the position
            t.value = value
        }

    }
}

func (this *Trie) DumpTree() {
    fmt.Printf("\n\n")
    this.DumpBranch(this.tree, 1)
}

func (this *Trie) DumpBranch(t *branch, depth int) {
    // attempt to output a textual view of the tree
    for x := 0; x < depth; x ++ { fmt.Print("  ") }
    fmt.Printf("- cheat: %s\n", t.shortcut)
    for x := 0; x < depth; x ++ { fmt.Print("  ") }
    fmt.Printf("- value: %s\n",t.value)
    if t.children != nil {
        for x := 0; x < depth; x ++ { fmt.Print("  ") }
        fmt.Printf("- children:\n")
        for y := 0; y < len(t.children); y++ {
            if t.children[y] != nil {
                for x := 0; x < depth; x ++ { fmt.Print("  ") }
                charb := make([]byte, 1)
                charb[0] = byte(y)+startLetter
                fmt.Printf(" - %s\n", string(charb))
                this.DumpBranch(t.children[y], depth+1)
            }
        }
    }
}

func (this *Trie) GetSuggestions(entry string) ([]interface{}, bool) {
  start := make([]interface{}, 0, 10)
  return this.getSuggestions(entry, start)
}

func (this *branch) allValuesBelow(values []interface{}) ([]interface{}) {
  if this.value != nil {
    values = append(values, this.value)
  }

  for _, branch := range this.children {
    if(branch != nil) { // Dear author: why does this happen?
      values = branch.allValuesBelow(values)
    }
  }

  return values
}

func (this *Trie) getSuggestions(entry string, suggestions []interface{}) ([]interface{}, bool) {
    t := this.tree
    eb := []byte(entry)
    // it's <= here to ensure we get to the cheat comparison nil on a valid path
    for x := 0; x <= len(eb); x++ {
        // if the current branch has a cheat, make sure we match it
        s := t.shortcut
        var y int
        for y = 0; y < len(s); y++ {
            if x+y >= len(eb) {
                //return nil, true
                return t.allValuesBelow(suggestions), true
            }
            if s[y] != eb[x+y] {
                // NO WAY, RETURN NOTHING
                return suggestions, false
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
                mapindex, exists := this.unicodeMap[ir]
                if !exists {
                    //return nil, false // no mapping :/
                    return suggestions, false
                } else {
                    index = mapindex
                }
            }
            if index > len(t.children)-1 || t.children[index] == nil {
                //return nil, false
                return suggestions, false
            }
            t = t.children[index]
            eb = eb[x:]
            x = 0
        }
    }
    return t.allValuesBelow(suggestions), true
}


func (this *Trie) GetEntry(entry string) (value interface{}, validPath bool) {
    t := this.tree
    eb := []byte(entry)
    // it's <= here to ensure we get to the cheat comparison nil on a valid path
    for x := 0; x <= len(eb); x++ {
        // if the current branch has a cheat, make sure we match it
        s := t.shortcut
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
                mapindex, exists := this.unicodeMap[ir]
                if !exists {
                    return nil, false // no mapping :/
                } else {
                    index = mapindex
                }
            }
            if index > len(t.children)-1 || t.children[index] == nil {
                return nil, false
            }
            t = t.children[index]
            eb = eb[x:]
            x = 0
        }
    }
    return t.value, true
}

func NewTrie() *Trie {
    t := &Trie {
        tree: &branch {
            children: nil,
            value: nil,
            shortcut: nil,
        },
        unicodeMap: make(map[int]int),
        nextIndex: 26,
    }
    return t
}
