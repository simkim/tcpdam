package tcpdam_test

import "github.com/simkim/tcpdam"
import "testing"

// disclamer : I don't now how to test code in got yet ...

var (
	la = ":9999"
	ra = "example.com:80"
)

func TestCreateDam(t *testing.T) {
	dam := tcpdam.NewDam(la, ra, 1, 1)
	if dam == nil {
		t.Errorf("Dam should not be nil")
	}
}

func TestOpeningDam(t *testing.T) {
	dam := tcpdam.NewDam(la, ra, 1, 1)
	if dam == nil {
		t.Errorf("Dam should not be nil")
	}
	err := dam.Open()
	if err != nil {
		t.Errorf("Can't open the dam : %s", err.Error)
	}
}
