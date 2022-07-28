package messagequeue

import (
	"testing"
)

func TestConnectShouldRetryBeforeFail(t *testing.T) {
	mqUrl = ""

	_, err := Connect()

	if err == nil {
		t.Errorf("Failed ! Error should NOT be 'nil'. Got %v", err)
	}
}
