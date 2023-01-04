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
			provider := NewSavedObjectsProvider(testCase.baseUrl, &testCase.httpClient)
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

			provider := NewSavedObjectsProvider(srv.URL, http.DefaultClient)
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

			provider := NewSavedObjectsProvider(srv.URL, http.DefaultClient)
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

			provider := NewSavedObjectsProvider(srv.URL, http.DefaultClient)
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
