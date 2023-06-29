package proto_test

// import (
// 	. "github.com/onsi/ginkgo/v2"
// 	. "github.com/onsi/gomega"

// 	"github.com/mohammadVatandoost/pkg/test/matchers"
// 	envoy_metadata "github.com/mohammadVatandoost/pkg/xds/envoy/metadata/v3"
// 	util_proto "github.com/mohammadVatandoost/xds-conrol-plane/pkg/util/proto"
// )

// var _ = Describe("MarshalAnyDeterministic", func() {
// 	It("should marshal deterministically", func() {
// 		tags := map[string]string{
// 			"service": "backend",
// 			"version": "v1",
// 			"cloud":   "aws",
// 		}
// 		metadata := envoy_metadata.EndpointMetadata(tags)
// 		for i := 0; i < 100; i++ {
// 			any1, _ := util_proto.MarshalAnyDeterministic(metadata)
// 			any2, _ := util_proto.MarshalAnyDeterministic(metadata)
// 			Expect(any1).To(matchers.MatchProto(any2))
// 		}
// 	})
// })
