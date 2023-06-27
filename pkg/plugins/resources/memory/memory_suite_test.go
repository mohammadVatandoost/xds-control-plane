package memory_test

import (
	"testing"

	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/test"
)

func TestClient(t *testing.T) {
	test.RunSpecs(t, "In-memory ResourceStore Suite")
}
