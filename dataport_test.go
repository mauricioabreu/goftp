package main

import "testing"

func TestDataportAddress(t *testing.T) {
	dataport, err := parse("10,2,0,2,4,15")
	if err != nil {
		t.Fatal(err)
	}
	if dataport.ip != "10.2.0.2" {
		t.Fatalf("Wrong IP value. Got %v want %v", dataport.ip, "10.2.0.2")
	}
	if dataport.port != 1039 {
		t.Fatalf("Wrong port value. Got %v want %v", dataport.ip, 1039)
	}
}
