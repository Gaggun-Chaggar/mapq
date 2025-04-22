# mapq

A small `GoLang` module for querying maps, json arrays and structured json logs.

## Install

```
go get github.com/Gaggun-Chaggar/mapq
```

## What is mapq?

`mapq` is a library designed to query arrays of structured data. `mapq` was originally conceived to support testing scenarios and to make assertions on log buffers.

`mapq` currently supports slices of maps, slog strings (json objects separated by new lines) and json arrays (as strings).

## Example

```golang
...

import (
  . "https://github.com/smarty/assertions" // example uses "assertions" for all 'Should' comparisons
)

...

queryBuilder := mapq.FromSlice([]map[string]any{
  {"level": "error", "nested": map[string]any{"object": "hi"}, "all": true},
  {"level": "info", "size": "big", "all": true},
  {"level": "warn", "array": []any{1, 2}, "all": true},
  {"level": "warn", "size": "big", "all": true},
  {"level": "error", "all": true, "nested": []any{"object", "hi"}},
})

// find all maps where ("level" = "warn" AND "size" = "big") OR "level" = "error"
query := queryBuilder.Where(
  mapq.Or(
    mapq.And(
      mapq.Assert("level", ShouldEqual, "warn"),
      mapq.Assert("size", ShouldContainSubstring, "big"),
    ),
    mapq.Assert("level", ShouldEqual, "error"),
  ),
)

filteredSlice := mapq.Filter(query) // returns a []map[string]any with all matching values
hasAtLeastOne := mapq.Exists(query) // true if there is at least one match
allMapsMatch := mapq.All(query) // true if all maps in slice match
hasTwoResultsOnly := mapq.Has(2, query) // true if the number of maps in slice that match the query is 2
```
