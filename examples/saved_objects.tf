resource "opensearch_saved_object" "test-query" {
  obj_id = "test-query"
  type   = "query"
  attributes = jsonencode({
    "title" : "Test Query",
    "description" : "",
    "query" : {
      "query" : "\"Processstate\" AND \"systemtest\"",
      "language" : "kuery"
    },
    "filters" : [
      {
        "meta" : {
          "index" : "1234",
          "negate" : false,
          "disabled" : false,
          "type" : "phrase",
          "key" : "kubernetes.labels.app",
          "params" : {
            "query" : "execution-coordinator"
          }
        },
        "query" : {
          "match_phrase" : {
            "kubernetes.labels.app" : "execution-coordinator"
          }
        },
        "$state" : {
          "store" : "appState"
        }
      }
    ]
  })
}
