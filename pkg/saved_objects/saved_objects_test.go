package saved_objects

/*
Copyright 2022 MOIA GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIgnoreFieldsOnIndexPatternProperty(t *testing.T) {
	testCases := []struct {
		desc         string
		expectFields bool
		obj          *SavedObjectOSD
		handlerFunc  http.HandlerFunc
		syncFields   bool
	}{
		{
			desc:         "should ignore fields",
			expectFields: false,
			syncFields:   false,
			obj:          &SavedObjectOSD{Type: "index-pattern", ID: "mock-pattern", SavedObjectPostPayload: SavedObjectPostPayload{Attributes: map[string]any{}, References: []Reference{}}},
			handlerFunc: func(w http.ResponseWriter, _ *http.Request) {
				obj := SavedObjectOSD{Type: "index-pattern", ID: "mock-pattern", SavedObjectPostPayload: SavedObjectPostPayload{Attributes: map[string]any{
					"fields": map[string]interface{}{
						"fieldA": "Hello world",
					},
				}, References: []Reference{}}}
				b, err := json.Marshal(&obj)
				if err != nil {
					w.WriteHeader(http.StatusOK)
				}
				w.Write(b)
			},
		},
		{
			desc:         "should not ignore fields on other types",
			expectFields: true,
			syncFields:   false,
			obj:          &SavedObjectOSD{Type: "search", ID: "mock-search", SavedObjectPostPayload: SavedObjectPostPayload{Attributes: map[string]any{}, References: []Reference{}}},
			handlerFunc: func(w http.ResponseWriter, _ *http.Request) {
				obj := SavedObjectOSD{Type: "search", ID: "mock-search", SavedObjectPostPayload: SavedObjectPostPayload{Attributes: map[string]any{
					"fields": map[string]interface{}{
						"fieldA": "Hello world",
					},
				}, References: []Reference{}}}
				b, err := json.Marshal(&obj)
				if err != nil {
					w.WriteHeader(http.StatusOK)
				}
				w.Write(b)
			},
		},
		{
			desc:         "should not ignore fields if provider is configured to sync fields",
			expectFields: true,
			syncFields:   true,
			obj:          &SavedObjectOSD{Type: "index-pattern", ID: "mock-pattern", SavedObjectPostPayload: SavedObjectPostPayload{Attributes: map[string]any{}, References: []Reference{}}},
			handlerFunc: func(w http.ResponseWriter, _ *http.Request) {
				obj := SavedObjectOSD{Type: "index-pattern", ID: "mock-pattern", SavedObjectPostPayload: SavedObjectPostPayload{Attributes: map[string]any{
					"fields": map[string]interface{}{
						"fieldA": "Hello world",
					},
				}, References: []Reference{}}}
				b, err := json.Marshal(&obj)
				if err != nil {
					w.WriteHeader(http.StatusOK)
				}
				w.Write(b)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			handler := http.NewServeMux()
			handler.HandleFunc("/_dashboards/api/saved_objects/search/mock-search", tC.handlerFunc)
			handler.HandleFunc("/_dashboards/api/saved_objects/index-pattern/mock-pattern", tC.handlerFunc)

			srv := httptest.NewServer(handler)
			defer srv.Close()

			provider := NewSavedObjectsProvider(srv.URL, http.DefaultClient, tC.syncFields)
			obj, diag := provider.GetObject(context.TODO(), tC.obj)
			if diag != nil {
				t.Error(diag)
			}

			attr := map[string]interface{}{}
			marshErr := json.Unmarshal([]byte(obj.Attributes), &attr)
			if marshErr != nil {
				t.Errorf("could not unmarshal '%v', %+v", obj.Attributes, marshErr)
			}
			_, hasFields := attr["fields"]

			if tC.expectFields != hasFields {
				t.Errorf("Expected field presence: %v but actual field presence was %v", tC.expectFields, hasFields)
			}

		})
	}
}

func TestNewSavedObjectsProvider(t *testing.T) {
	testCases := []struct {
		desc       string
		baseUrl    string
		httpClient http.Client
	}{
		{
			desc:       "must create provider",
			baseUrl:    "https://foobar.com",
			httpClient: http.Client{},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			provider := NewSavedObjectsProvider(testCase.baseUrl, &testCase.httpClient, false)
			if provider == nil {
				t.Fail()
			}
		})
	}
}

func TestGetObject(t *testing.T) {
	testCases := []struct {
		desc        string
		wantErr     bool
		obj         *SavedObjectOSD
		handlerFunc http.HandlerFunc
	}{
		{
			desc:    "must get object",
			wantErr: false,
			obj:     &SavedObjectOSD{Type: "search", ID: "mock-search", SavedObjectPostPayload: SavedObjectPostPayload{Attributes: map[string]any{}, References: []Reference{}}},
			handlerFunc: func(w http.ResponseWriter, _ *http.Request) {
				obj := SavedObjectOSD{Type: "search", ID: "mock-search", SavedObjectPostPayload: SavedObjectPostPayload{Attributes: map[string]any{}, References: []Reference{}}}
				b, err := json.Marshal(&obj)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
				w.Write(b)
			},
		},
		{
			desc:    "must not fail if object does not exist",
			wantErr: false,
			obj:     &SavedObjectOSD{Type: "search", ID: "mock-search", SavedObjectPostPayload: SavedObjectPostPayload{Attributes: map[string]any{}, References: []Reference{}}},
			handlerFunc: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
		},
		{
			desc:    "must fail with wrong payload",
			wantErr: true,
			obj:     &SavedObjectOSD{Type: "search", ID: "mock-search", SavedObjectPostPayload: SavedObjectPostPayload{Attributes: map[string]any{}, References: []Reference{}}},
			handlerFunc: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			handler := http.NewServeMux()
			handler.HandleFunc("/_dashboards/api/saved_objects/search/mock-search", tC.handlerFunc)

			srv := httptest.NewServer(handler)
			defer srv.Close()

			provider := NewSavedObjectsProvider(srv.URL, http.DefaultClient, false)
			_, diag := provider.GetObject(context.TODO(), tC.obj)
			if tC.wantErr && diag != nil {
				return
			}
			if diag != nil {
				t.Error(diag)
			}
		})
	}
}

func TestSaveObject(t *testing.T) {
	testCases := []struct {
		desc        string
		wantErr     bool
		obj         *SavedObjectOSD
		handlerFunc http.HandlerFunc
	}{
		{
			desc:    "must save object",
			wantErr: false,
			obj:     &SavedObjectOSD{Type: "search", ID: "mock-search", SavedObjectPostPayload: SavedObjectPostPayload{Attributes: map[string]any{}, References: []Reference{}}},
			handlerFunc: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			handler := http.NewServeMux()
			// we define two handle funcs here, first handler saves  the object,  second one will update the object
			handler.HandleFunc("/_dashboards/api/saved_objects/search/mock-search", tC.handlerFunc)
			handler.HandleFunc("/_dashboards/api/saved_objects/dashboard/mock-search", tC.handlerFunc)

			srv := httptest.NewServer(handler)
			defer srv.Close()

			provider := NewSavedObjectsProvider(srv.URL, http.DefaultClient, false)
			diag := provider.SaveObject(context.TODO(), tC.obj)
			if tC.wantErr && diag != nil {
				return
			}
			if diag != nil {
				t.Error(diag)
			}

			// save the same object twice to update
			tC.obj.Type = "dashboard"
			diag = provider.SaveObject(context.TODO(), tC.obj)
			if tC.wantErr && diag != nil {
				return
			}
			if diag != nil {
				t.Error(diag)
			}
		})
	}
}

func TestDeleteObject(t *testing.T) {
	testCases := []struct {
		desc        string
		wantErr     bool
		obj         *SavedObjectOSD
		handlerFunc http.HandlerFunc
	}{
		{
			desc:    "must delete object",
			wantErr: false,
			obj:     &SavedObjectOSD{Type: "search", ID: "mock-search", SavedObjectPostPayload: SavedObjectPostPayload{Attributes: map[string]any{}, References: []Reference{}}},
			handlerFunc: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		},
		{
			desc:    "must fail when HTTP status not 200 OK",
			wantErr: true,
			obj:     &SavedObjectOSD{Type: "search", ID: "mock-search", SavedObjectPostPayload: SavedObjectPostPayload{Attributes: map[string]any{}, References: []Reference{}}},
			handlerFunc: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			handler := http.NewServeMux()
			handler.HandleFunc("/_dashboards/api/saved_objects/search/mock-search", tC.handlerFunc)

			srv := httptest.NewServer(handler)
			defer srv.Close()

			provider := NewSavedObjectsProvider(srv.URL, http.DefaultClient, false)
			diag := provider.DeleteObject(context.TODO(), tC.obj)
			if tC.wantErr && diag != nil {
				return
			}
			if diag != nil {
				t.Error(diag)
			}
		})
	}
}
