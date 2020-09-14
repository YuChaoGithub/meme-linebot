package models

import (
	"reflect"
	"testing"
)

func TestGetAll(t *testing.T) {
	// Stub and driver.
	db, teardown := newTestDB(t)
	defer teardown()

	m := MemeModel{db}

	// When.
	entries, err := m.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	// Want.
	wantMemes := []MemeEntry{
		{Name: "adios", Link: "6UegMI2.png"},
		{Name: "bonjour", Link: "qg8sB6f.png"},
		{Name: "honest work", Link: "BPCZHUi.png"},
		{Name: "it ain't much, but it's honest work", Link: "BPCZHUi.png"},
		{Name: "我就爛", Link: "t9WaxTw.png"},
	}
	for i := range wantMemes {
		wantMemes[i].Name += nameSuffix
		wantMemes[i].Link = imgurBaseLink + wantMemes[i].Link
	}

	if !reflect.DeepEqual(entries, wantMemes) {
		t.Errorf("want:\n%v\ngot:\n%v", wantMemes, entries)
	}
}

func TestGet(t *testing.T) {
	// Testcases.
	tests := []struct {
		testName string
		memeName string
		wantURL  string
	}{
		{"Chinese", "我就爛", "t9WaxTw.png"},
		{"With single quote", "it ain't much, but it's honest work", "BPCZHUi.png"},
		{"Doesn't exist", "ahhhhhhh", ""},
	}

	// Perform tests.
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// Stub and driver.
			db, teardown := newTestDB(t)
			defer teardown()

			m := MemeModel{db}

			// When.
			url, err := m.Get(tc.memeName)

			// Want.
			wantURL := tc.wantURL
			if err == nil {
				wantURL = imgurBaseLink + wantURL
			}

			if url != wantURL {
				t.Errorf("want %v; got %v", wantURL, url)
			}
		})
	}
}

func TestGetFuzzy(t *testing.T) {
	// Testcases.
	tests := []struct {
		testName string
		memeName string
		wantURL  string
	}{
		{"Chinese", "就爛", "t9WaxTw.png"},
		{"almost match", "bonjer", "qg8sB6f.png"},
		{"exact match", "adios", "6UegMI2.png"},
	}

	// Perform tests.
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// Stub and driver.
			db, teardown := newTestDB(t)
			defer teardown()

			m := MemeModel{db}

			// When.
			url, err := m.GetFuzzy(tc.memeName)

			// Want.
			wantURL := tc.wantURL
			if err == nil {
				wantURL = imgurBaseLink + wantURL
			}

			if url != wantURL {
				t.Errorf("want %v; got %v", wantURL, url)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	// Testcases.
	tests := []struct {
		testName string
		memeName string
		memeURL  string
		hasErr   bool
	}{
		{"Existing Entry", "我就爛", "t9WaxTw.png", true},
		{"Existing URL but different name", "爛", "t9WaxTw.png", false},
		{"Doesn't exist", "ah", "txt.png", false},
	}

	// Perform tests.
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// Stub and driver.
			db, teardown := newTestDB(t)
			defer teardown()

			m := MemeModel{db}

			// When.
			err := m.Insert(tc.memeName, tc.memeURL)
			hasErr := err != nil

			// Want.
			if tc.hasErr != hasErr {
				t.Errorf("want %v; got %v", tc.hasErr, hasErr)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	// Testcases.
	tests := []struct {
		testName string
		memeName string
		hasErr   bool
	}{
		{"Existing Entry", "我就爛", false},
		{"Doesn't exist", "ah", false},
	}

	// Perform tests.
	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// Stub and driver.
			db, teardown := newTestDB(t)
			defer teardown()

			m := MemeModel{db}

			// When.
			err := m.Delete(tc.memeName)
			hasErr := err != nil

			// Want.
			if tc.hasErr != hasErr {
				t.Errorf("want %v; got %v", tc.hasErr, hasErr)
			}
		})
	}
}
