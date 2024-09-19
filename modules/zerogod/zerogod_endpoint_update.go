package zerogod

import (
	"fmt"
	"strings"

	"github.com/bettercap/bettercap/v2/modules/syn_scan"
	"github.com/bettercap/bettercap/v2/network"
	"github.com/evilsocket/islazy/str"
	"github.com/grandcat/zeroconf"
)

func (mod *ZeroGod) updateEndpointMeta(address string, endpoint *network.Endpoint, svc *zeroconf.ServiceEntry) {
	mod.Debug("found endpoint %s for address %s", endpoint.HwAddress, address)

	// TODO: this is shit and needs to be refactored

	// update mdns metadata
	meta := make(map[string]string)

	svcType := svc.Service

	meta[fmt.Sprintf("mdns:%s:name", svcType)] = svc.ServiceName()
	meta[fmt.Sprintf("mdns:%s:hostname", svcType)] = svc.HostName

	// TODO: include all
	if len(svc.AddrIPv4) > 0 {
		meta[fmt.Sprintf("mdns:%s:ipv4", svcType)] = svc.AddrIPv4[0].String()
	}

	if len(svc.AddrIPv6) > 0 {
		meta[fmt.Sprintf("mdns:%s:ipv6", svcType)] = svc.AddrIPv6[0].String()
	}

	meta[fmt.Sprintf("mdns:%s:port", svcType)] = fmt.Sprintf("%d", svc.Port)

	for _, field := range svc.Text {
		field = str.Trim(field)
		if len(field) == 0 {
			continue
		}

		key := ""
		value := ""

		if strings.Contains(field, "=") {
			parts := strings.SplitN(field, "=", 2)
			key = parts[0]
			value = parts[1]
		} else {
			key = field
		}

		meta[fmt.Sprintf("mdns:%s:info:%s", svcType, key)] = value
	}

	mod.Debug("meta for %s: %v", address, meta)

	endpoint.OnMeta(meta)

	// update ports
	ports := endpoint.Meta.GetOr("ports", map[int]*syn_scan.OpenPort{}).(map[int]*syn_scan.OpenPort)
	if _, found := ports[svc.Port]; !found {
		ports[svc.Port] = &syn_scan.OpenPort{
			Proto:   "tcp",
			Port:    svc.Port,
			Service: network.GetServiceByPort(svc.Port, "tcp"),
		}
	}

	endpoint.Meta.Set("ports", ports)
}
