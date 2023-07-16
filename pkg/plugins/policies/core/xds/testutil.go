package xds

import (
	_ "embed"

	envoy_cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	. "github.com/onsi/gomega"

	core_xds "github.com/mohammadVatandoost/xds-conrol-plane/pkg/core/xds"
	util_proto "github.com/mohammadVatandoost/xds-conrol-plane/pkg/util/proto"
	clusters_builder "github.com/mohammadVatandoost/xds-conrol-plane/pkg/xds/envoy/clusters"
)

func ResourceArrayShouldEqual(resources core_xds.ResourceList, expected []string) {
	Expect(resources).To(HaveLen(len(expected)))

	for i, r := range resources {
		actual, err := util_proto.ToYAML(r.Resource)
		Expect(err).ToNot(HaveOccurred())

		Expect(actual).To(MatchYAML(expected[i]))
	}
	Expect(resources).To(HaveLen(len(expected)))
}

type NameConfigurer struct {
	Name string
}

func (n *NameConfigurer) Configure(c *envoy_cluster.Cluster) error {
	c.Name = n.Name
	return nil
}

func WithName(name string) clusters_builder.ClusterBuilderOpt {
	return clusters_builder.ClusterBuilderOptFunc(func(builder *clusters_builder.ClusterBuilder) {
		builder.AddConfigurer(&NameConfigurer{Name: name})
	})
}
