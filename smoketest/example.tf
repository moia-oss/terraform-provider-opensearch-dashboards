resource "opensearch_saved_object" "ref_terraform_provider_test_search" {
  obj_id = "1d2131a0-8fe9-11ed-a0ec-dd170dad7eb3"
  type   = "search"
  attributes = jsonencode(
    {
      "columns": [
        "level",
        "tripId",
        "kubernetes.namespace_name",
        "component"
      ],
      "description": "",
      "hits": 0,
      "kibanaSavedObjectMeta": {
        "searchSourceJSON": "{\"highlightAll\":true,\"version\":true,\"query\":{\"query\":\"level:ERROR\",\"language\":\"kuery\"},\"filter\":[{\"meta\":{\"alias\":null,\"negate\":false,\"disabled\":false,\"type\":\"phrase\",\"key\":\"component\",\"params\":{\"query\":\"rts\"},\"indexRefName\":\"kibanaSavedObjectMeta.searchSourceJSON.filter[0].meta.index\"},\"query\":{\"match_phrase\":{\"component\":\"rts\"}},\"$state\":{\"store\":\"appState\"}}],\"indexRefName\":\"kibanaSavedObjectMeta.searchSourceJSON.index\"}"
      },
      "sort": [],
      "title": "terraform-provider-test-search",
      "version": 1
    }
  )

  references {
    id   = "application-index-pattern"
    name = "kibanaSavedObjectMeta.searchSourceJSON.index"
    type = "index-pattern"
  }

  references {
    id   = "application-index-pattern"
    name = "kibanaSavedObjectMeta.searchSourceJSON.filter[0].meta.index"
    type = "index-pattern"
  }

}

resource "opensearch_saved_object" "ref_terraform_provider_test_visualization" {
  obj_id = "7fa60370-8ff1-11ed-a0ec-dd170dad7eb3"
  type   = "visualization"
  attributes = jsonencode(
    {
      "description": "",
      "kibanaSavedObjectMeta": {
        "searchSourceJSON": "{\"query\":{\"query\":\"level:ERROR\",\"language\":\"kuery\"},\"filter\":[{\"$state\":{\"store\":\"appState\"},\"meta\":{\"alias\":null,\"disabled\":false,\"key\":\"component\",\"negate\":false,\"params\":{\"query\":\"rts\"},\"type\":\"phrase\",\"indexRefName\":\"kibanaSavedObjectMeta.searchSourceJSON.filter[0].meta.index\"},\"query\":{\"match_phrase\":{\"component\":\"rts\"}}}]}"
      },
      "savedSearchRefName": "search_0",
      "title": "terraform-provider-test-visualization",
      "uiStateJSON": "{\"vis\":{\"params\":{\"sort\":{\"columnIndex\":null,\"direction\":null}}}}",
      "version": 1,
      "visState": "{\"title\":\"terraform-provider-test-visualization\",\"type\":\"table\",\"aggs\":[{\"id\":\"1\",\"enabled\":true,\"type\":\"count\",\"params\":{},\"schema\":\"metric\"},{\"id\":\"2\",\"enabled\":true,\"type\":\"terms\",\"params\":{\"field\":\"level\",\"orderBy\":\"1\",\"order\":\"desc\",\"size\":5,\"otherBucket\":false,\"otherBucketLabel\":\"Other\",\"missingBucket\":false,\"missingBucketLabel\":\"Missing\"},\"schema\":\"bucket\"}],\"params\":{\"perPage\":10,\"showPartialRows\":false,\"showMetricsAtAllLevels\":false,\"sort\":{\"columnIndex\":null,\"direction\":null},\"showTotal\":false,\"totalFunc\":\"sum\",\"percentageCol\":\"\"}}"
    }
  )

  references {
    id   = "1d2131a0-8fe9-11ed-a0ec-dd170dad7eb3"
    name = "search_0"
    type = "search"
  }

  references {
    id   = "application-index-pattern"
    name = "kibanaSavedObjectMeta.searchSourceJSON.filter[0].meta.index"
    type = "index-pattern"
  }

}

resource "opensearch_saved_object" "ref_terraform_provider_test_dashboard" {
  obj_id = "0e625c00-8ff5-11ed-93ec-ebe1e8735d27"
  type   = "dashboard"
  attributes = jsonencode(
    {
      "description": "",
      "hits": 0,
      "kibanaSavedObjectMeta": {
        "searchSourceJSON": "{\"query\":{\"query\":\"level:ERROR\",\"language\":\"kuery\"},\"filter\":[{\"meta\":{\"alias\":null,\"negate\":false,\"disabled\":false,\"type\":\"phrase\",\"key\":\"component\",\"params\":{\"query\":\"rts\"},\"indexRefName\":\"kibanaSavedObjectMeta.searchSourceJSON.filter[0].meta.index\"},\"query\":{\"match_phrase\":{\"component\":\"rts\"}},\"$state\":{\"store\":\"appState\"}}]}"
      },
      "optionsJSON": "{\"hidePanelTitles\":false,\"useMargins\":true}",
      "panelsJSON": "[{\"embeddableConfig\":{},\"gridData\":{\"h\":15,\"i\":\"252f4483-f5f2-4e67-a161-0511c25088d1\",\"w\":24,\"x\":0,\"y\":0},\"panelIndex\":\"252f4483-f5f2-4e67-a161-0511c25088d1\",\"version\":\"1.3.2\",\"panelRefName\":\"panel_0\"}]",
      "refreshInterval": {
        "pause": true,
        "value": 0
      },
      "timeFrom": "now-1h",
      "timeRestore": true,
      "timeTo": "now",
      "title": "terraform-provider-test-dashboard",
      "version": 1
    }
  )

  references {
    id   = "application-index-pattern"
    name = "kibanaSavedObjectMeta.searchSourceJSON.filter[0].meta.index"
    type = "index-pattern"
  }

  references {
    id   = "7fa60370-8ff1-11ed-a0ec-dd170dad7eb3"
    name = "panel_0"
    type = "visualization"
  }

}

resource "opensearch_saved_object" "ref_terraform_provider_test_query" {
  obj_id = "terraform-provider-test-query"
  type   = "query"
  attributes = jsonencode(
    {
      "description": "",
      "filters": [
        {
          "$state": {
            "store": "appState"
          },
          "meta": {
            "alias": null,
            "disabled": false,
            "index": "application-index-pattern",
            "key": "component",
            "negate": false,
            "params": {
              "query": "rts"
            },
            "type": "phrase"
          },
          "query": {
            "match_phrase": {
              "component": "rts"
            }
          }
        }
      ],
      "query": {
        "language": "kuery",
        "query": "level: WARN"
      },
      "timefilter": {
        "from": "now-1h",
        "refreshInterval": {
          "pause": true,
          "value": 0
        },
        "to": "now"
      },
      "title": "terraform-provider-test-query"
    }
  )

}

resource "opensearch_saved_object" "ref_terraform_provider_test_url" {
  obj_id = "terraform_provider_test_url"
  type   = "url"
  attributes = jsonencode(
    {
      "accessCount": 0,
      "accessDate": 1599746912897,
      "createDate": 1599746912897,
      "url": "/app/kibana#/dashboard/%5Bdispatching%5D-Dashboard?_g=(filters:!(),refreshInterval:(display:Off,pause:!f,value:0),time:(from:now-4h,to:now))\u0026_a=(description:%27%27,filters:!((%27$state%27:(store:appState),meta:(alias:!n,disabled:!f,index:%2757e52300-e235-11ea-a48a-77ea726c762e%27,key:kubernetes.labels.app.keyword,negate:!t,params:(query:dispatching2),type:phrase),query:(match_phrase:(kubernetes.labels.app.keyword:dispatching2))),(%27$state%27:(store:appState),meta:(alias:!n,disabled:!f,index:%2757e52300-e235-11ea-a48a-77ea726c762e%27,key:level.keyword,negate:!t,params:(query:ERROR),type:phrase),query:(match_phrase:(level.keyword:ERROR))),(%27$state%27:(store:appState),meta:(alias:!n,disabled:!f,index:%2757e52300-e235-11ea-a48a-77ea726c762e%27,key:level.keyword,negate:!f,params:(query:WARN),type:phrase),query:(match_phrase:(level.keyword:WARN)))),fullScreenMode:!f,options:(darkTheme:!t),panels:!((embeddableConfig:(columns:!(level,logger,traceId,message),sort:!(%27@timestamp%27,desc)),gridData:(h:35,i:%275%27,w:48,x:0,y:35),id:%5Bdispatching%5D-All,panelIndex:%275%27,type:search,version:%277.7.0%27),(embeddableConfig:(spy:(mode:(fill:!f,name:!n)),vis:(colors:(DEBUG:%23CCA300,ERROR:%2399440A),legendOpen:!f)),gridData:(h:10,i:%276%27,w:48,x:0,y:0),id:%5Bdispatching%5D-Histogram,panelIndex:%276%27,type:visualization,version:%277.7.0%27),(embeddableConfig:(spy:(mode:(fill:!f,name:!n)),vis:(params:(sort:(columnIndex:!n,direction:!n)))),gridData:(h:20,i:%277%27,w:16,x:16,y:10),id:%5Bdispatching%5D-Aggregation-by-Level,panelIndex:%277%27,type:visualization,version:%277.7.0%27),(embeddableConfig:(vis:(params:(sort:(columnIndex:2,direction:desc)))),gridData:(h:25,i:%278%27,w:16,x:32,y:10),id:%5Bdispatching%5D-Aggregation-by-Logger,panelIndex:%278%27,type:visualization,version:%277.7.0%27),(embeddableConfig:(vis:(params:(sort:(columnIndex:!n,direction:!n)))),gridData:(h:20,i:%2710%27,w:16,x:0,y:10),id:%5Bdispatching%5D-Aggregation-by-Application,panelIndex:%2710%27,type:visualization,version:%277.7.0%27)),query:(language:lucene,query:%27!!message:(%22VSO%22,%20%22io.moia.triphandling.Error$NoTripWithThisId%22,%20%22io.moia.triphandling.Error$DbVersionConflict%22,%20%22io.moia.triphandling.Error$TripInWrongState%22)%27),timeRestore:!t,title:%27%5Bdispatching%5D%20Dashboard%27,viewMode:view)"
    }
  )

}

