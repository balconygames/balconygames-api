package geo

import "testing"

func TestLookupCountryCode(t *testing.T) {
	db := New()

	country, err := db.Resolve("8.8.8.8")
	if err != nil {
		t.Error(err)
	}
	if country != "us" {
		t.Errorf("wrong look up country is '%s'", country)
	}
}
