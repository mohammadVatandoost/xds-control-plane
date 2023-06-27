package memory_test

import (
	. "github.com/onsi/ginkgo/v2"

	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/resources/memory"
	test_store "github.com/mohammadVatandoost/xds-conrol-plane/pkg/test/store"
)

var _ = Describe("MemoryStore template", func() {
	test_store.ExecuteStoreTests(memory.NewStore, "memory")
	test_store.ExecuteOwnerTests(memory.NewStore, "memory")
})
