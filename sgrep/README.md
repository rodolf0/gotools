# sgrep, scoped grep

like grep but will show a whole section, ie: give context

example:

{
  "key1": 1,
  {
    "nested1": 2,
  }
}

sgrep'ing for nested should show {"nested1": 2,}


  config = {
    { "<html>", "</html>" },
    { "<[^>]*>", "</[^>]*>" },
  }

  --scopes=1 (how many outer levels to show)
  --open="{" --close="}" (override known open/close)
  --white-based (python like scopes)
  --drop  (drop outermost scope delimiters)
  --icase
  --color (different colors per pair)
  --pretty (ie: for python remove first indents, format json, format html)

grep

* html
* [] <> () {} pairs
* c++ comments: /* */
* tab based indents (python)

1. keep track of each opening element type, acumulate in buffer
2. slide buffer as new opening elements are found, keep at most --scopes
3. if no open/close has being provided, use all defaults
4. when a match is found, stop looking further, just found the closing elements
5. regex checks are always per line (don't cross line boundaries)


1. keep a stack of opening markers: context
2. for each line, search incrementally for any+all markers in order open/close
3. check line for regex, if match check context


a. Check if the line matches... if and only if line matches, get the context

b. how long to keep context?
** if the context is open, need to keep it
** if closed and there was no match: discard
