package policies

import (
	_ "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshaccesslog"
	_ "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshcircuitbreaker"
	_ "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshfaultinjection"
	_ "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshhealthcheck"
	_ "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshhttproute"
	_ "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshloadbalancingstrategy"
	_ "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshproxypatch"
	_ "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshratelimit"
	_ "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshretry"
	_ "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshtcproute"
	_ "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshtimeout"
	_ "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshtrace"
	_ "github.com/mohammadVatandoost/xds-conrol-plane/pkg/plugins/policies/meshtrafficpermission"
)
