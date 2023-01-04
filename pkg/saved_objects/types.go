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

type SavedObjectOSD struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
	SavedObjectPostPayload
}

type SavedObjectPostPayload struct {
	Attributes map[string]any `json:"attributes,omitempty"`
	References []Reference    `json:"references,omitempty"`
}

type SavedObjectTF struct {
	Type string `json:"type,omitempty"`
	ID   string `json:"id,omitempty"`
	// attributes is a string encoded json
	Attributes string      `json:"attributes,omitempty"`
	References []Reference `json:"references,omitempty"`
}

type Reference struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
	Type string `json:"type"`
}
