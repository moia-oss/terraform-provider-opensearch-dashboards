package saved_objects

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
