package v1alpha1

import (
	common_api "github.com/mohammadVatandoost/xds-conrol-plane/api/common/v1alpha1"
	"github.com/mohammadVatandoost/xds-conrol-plane/pkg/util/pointer"
)

func (x *To) GetDefault() interface{} {
	if len(x.Rules) == 0 {
		return Rule{
			Default: RuleConf{
				BackendRefs: []common_api.BackendRef{{
					TargetRef: x.TargetRef,
					Weight:    pointer.To(uint(1)),
				}},
			},
		}
	}

	return x.Rules[0]
}
