package models

import "time"

type Profile struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Attributes  *Attributes `json:"attributes,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type Attributes struct {
	Gender             string   `json:"gender,omitempty"`
	Age                int      `json:"age,omitempty,string"`
	Location           Location `json:"location,omitempty"`
	EducationLevel     string   `json:"education_level,omitempty"`
	Occupation         string   `json:"occupation,omitempty"`
	Interests          []string `json:"interests,omitempty"`
	ExpertiseLevel     string   `json:"expertise_level,omitempty"`
	TonePreference     string   `json:"tone_preference,omitempty"`
	PreferredLanguages []string `json:"preferred_languages,omitempty"`
	Intent             string   `json:"intent,omitempty"`
}

type Location struct {
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	Country    string `json:"country,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
}
