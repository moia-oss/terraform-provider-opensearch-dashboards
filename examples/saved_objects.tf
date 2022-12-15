resource "opensearch_saved_object" "test-query" {
  obj_id = "test-query"
  type   = "query"
  attributes = jsonencode({
    "title" : "Test Query",
    "description" : "",
    "query" : {
      "query" : "\"firstSearchTerm\" AND \"SecondSearchTerm\"",
      "language" : "kuery"
    },
    "filters" : [
      {
        "meta" : {
          "index" : "1234",
          "negate" : false,
          "disabled" : false,
          "type" : "phrase",
          "key" : "foo.key",
          "params" : {
            "query" : "value-to-search-for"
          }
        },
        "query" : {
          "match_phrase" : {
            "foo.key" : "value-to-search-for"
          }
        },
        "$state" : {
          "store" : "fooState"
        }
      }
    ]
  })
}
