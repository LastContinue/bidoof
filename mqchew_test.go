package main

import "testing"

func TestJsonToRegistration(t *testing.T) {
	jsonMessage := `{"name":"Prof.Oak","email":"oak@kantolab.poke","tour":"Durice","islandType":"ice"}`

	expectedRegistration := Registration{
		Name:       "Prof.Oak",
		Email:      "oak@kantolab.poke",
		Tour:       "Durice",
		IslandType: "ice",
	}

	gotRegistration, err := jsonToRegistration(jsonMessage)

	if expectedRegistration != gotRegistration {
		t.Errorf("Failed ! Expected %v got %v", expectedRegistration, gotRegistration)
	}

	if err != nil {
		t.Errorf("Failed ! Error should be nil. Got %v", err)
	}
}

func TestJsonToRegistrationFailure(t *testing.T) {
	jsonMessage := ""

	gotRegistration, err := jsonToRegistration(jsonMessage)

	if err == nil {
		t.Errorf("Failed ! Error should be not nil. Got Registration is %v", gotRegistration)
	}
}

func TestFormatEmailBody(t *testing.T) {
	registration := Registration{
		Name:       "Prof.Oak",
		Email:      "oak@kantolab.poke",
		Tour:       "Durice",
		IslandType: "ice",
	}

	expectedEmailBody := "Thank you for registering to Illumicon, Prof.Oak!\n For your tour to Durice, please pack appropriately as it is a ice-y place!"

	gotEmailBody := formatEmailBody(&registration)

	if expectedEmailBody != gotEmailBody {
		t.Errorf("Failed with IslandType ! Expected %v got %v", expectedEmailBody, gotEmailBody)
	}

	registration.IslandType = NA_TYPE

	expectedEmailBodyNoIslandType := "Thank you for registering to Illumicon, Prof.Oak!"

	gotEmailBodyNoIslandType := formatEmailBody(&registration)

	if expectedEmailBodyNoIslandType != gotEmailBodyNoIslandType {
		t.Errorf("Failed without IslandType ! Expected %v got %v", expectedEmailBody, gotEmailBody)
	}
}
