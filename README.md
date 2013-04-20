# grushtools

An exercise in Go taking aggregate/pivot tools from the [crush-tools](https://code.google.com/p/crush-tools/) project.

Currently only reads from **stdin** and writes to **stdout**.

### examples
* List beers launched per year
`cat beers.in | ./grushtools -d $'\x09' -k Year -t Beer-Title -b ', '`
* Check how many beers were launched that year
`cat beers.in | ./grushtools -d $'\t' -k Year -c Beer-Title`
* Total beers launched
`cat beers.in | ./grushtools -d $'\t' -c Beer-Title`

### building
1. $(GOPATH/env.sh) # Will set GOPATH to find our modules
2. go build
3. ./grushtools -h
