package db

import (
	"testing"
)

func TestConnectShouldRetryBeforeFail(t *testing.T) {

	makeDbConnectionString = func() string {
		return ""
	}

	_, err := Connect()

	if err == nil {
		t.Errorf("Failed ! Error should NOT be 'nil'. Got %v", err)
	}
}

//Need to also test a good connection, as well as inserting, but not entirely sure how to mock this out yet
