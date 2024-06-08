# Serif Health Takehome Interview

This is my solution to the Serif Health take-home. I've moved the original README to `_README.md`.

### Setup

* Make sure [Golang](https://go.dev/doc/install) is installed.
* You'll need a [dump file](https://antm-pt-prod-dataz-nogbd-nophi-us-east1.s3.amazonaws.com/anthem/2024-06-01_anthem_index.json.gz) from Anthem stored inside of the root of the project.
* Run the source:
```
go run main.go
```

### Solution

I began trying to understand the data by printing out network file descriptions as recommended and seeing what the results were. This yielded a lot of different results where the name of the relevant
state was present, so I figured I would filter for descriptions with "NY" or "New York" in them. This produced file locations which I noticed had the substrings `_71B0_` and `_72B0_` appearing as
identifiers for each dataset. Searching for Highmark in the MRF lookup system showed datasets which also had these identifiers, which were both prefixed with the state abbreviation `_NY_`. I then looked
for other datasets which had the `_NY_` substring within them, and noticed there were some identifiers which my code was not producing any URLs for by simply filtering for descriptions with NY state in them.
I inferred that any URLs with these identifiers may also be part of the correct solution and began filtering on the file locations themselves, this time for any identifier which also had a `_NY_` prefix in
the MRF system. This produced more URLs with each identifier I had a filter for, which gave me a formula of sorts to work from. I ruled out any network file descriptions with states other than New York in
them and searched the MRF system with the remainder to collect any other identifiers with a `_NY_` prefix and built a list. From there, it seemed reasonable to me that the output I was producing contained
all the URLs which are expected.

The output URLs can be found in `urls.txt` and this solution runs in just about 4.5 minutes on my machine. It is very naive and could likely be optimized a good deal more, but I wanted to focus on
attempting to arrive at the right solution in the allotted time. It took me roughly the full 2 hours, plus some time for figuring out how to do things in Go, which I am still getting familiar with.
There are a things I'd consider working on to "productionize" this code:
* Investigate how Go's concurrency features could help speed up the program's execution.
* Parameterize the list of dataset indentifiers so that, for example, a set of file locations can be extracted for any state or multiple states.
* Wrap some tests around the the core parsing logic with more contrived inputs as a fixture. I'm not totally confident that this code isn't missing some of the correct URLs; I nearly missed a good portion of
the file locations because I forgot to start iterating over the list of in network files themselves after they were decoded. 

