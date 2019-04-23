// Copyright 2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ethernet

import (
	"github.com/platinasystems/elib/cli"
	"github.com/platinasystems/vnet"
	"github.com/platinasystems/vnet/ip"
	"github.com/platinasystems/xeth"

	"fmt"
	"net"
)

type showNeighborConfig struct {
	ip4       bool
	ip6       bool
	detail    bool
	showTable string
}

func (m *Main) showIpNeighbor(c cli.Commander, w cli.Writer, in *cli.Input) (err error) {
	cf := showNeighborConfig{}
	v := m.ipNeighborMain.v
	for !in.End() {
		switch {
		case in.Parse("ip4"):
			cf.ip4 = true
		case in.Parse("ip6"):
			cf.ip6 = true
		case in.Parse("d%*etail"):
			cf.detail = true
		case in.Parse("t%*able %s", &cf.showTable):
		default:
			err = cli.ParseError
			return
		}
	}
	//if not explicity specified, show both
	if !cf.ip4 && !cf.ip6 {
		cf.ip4 = true
		cf.ip6 = true
	}

	em := GetMain(v)

	for ipFamily, nf := range em.ipNeighborFamilies {
		im := nf.m
		if ip.Family(ipFamily) == ip.Ip4 && !cf.ip4 {
			continue
		}
		if ip.Family(ipFamily) == ip.Ip6 && !cf.ip6 {
			continue
		}
		for _, i := range nf.indexByAddress {
			n := &nf.pool.neighbors[i]
			fi := im.FibIndexForSi(n.Si)
			ns := im.FibNameForIndex(fi)

			if cf.showTable != "" && ns != cf.showTable {
				continue
			}

			var (
				ok        bool
				as        []ip.Adjacency
				adj_lines []string
				prefix    net.IPNet
			)

			prefix.IP = n.Ip
			prefix.Mask = net.CIDRMask(32, 32)
			if ip.Family(ipFamily) == ip.Ip6 {
				prefix.Mask = net.CIDRMask(128, 128)
			}

			ipAddr := n.Ip.String()
			//mac := n.Ethernet.String()
			intf := fmt.Sprint(vnet.SiName{V: v, Si: n.Si})
			lladdr := n.Ethernet.String()

			ai := ip.AdjNil
			ln := 0
			if ai, as, ok = im.GetReachable(&prefix, n.Si); ok {
				for i := range as {
					adj_lines = as[i].AdjLines(im)
				}
				if ln == 0 {
					fmt.Fprintf(w, "%10v%20v dev %10v lladdr %v      adjacency %v:%v\n", ns, ipAddr, intf, lladdr, ai, adj_lines)
				} else {
					fmt.Fprintf(w, "%10v%20v dev %10v lladdr %v      adjacency %v:%v\n", "", "unexpected extras", "", "", ai, adj_lines)
				}
				ln++
			} else {
				fmt.Fprintf(w, "%10v%20v dev %10v lladdr %v      adjacency %v:%v\n", ns, ipAddr, intf, lladdr, ai, "not found")
			}

			if cf.detail {
				//no additional details for now
			}
		}
	}
	return
}

func (m *Main) showMac(c cli.Commander, w cli.Writer, in *cli.Input) (err error) {
	fmt.Fprintf(w, "%20v%20v%15v%25v%20v%10v\n", "namespace", "bridgeName", "stag", "macAddr", "dev", "local")
	v := m.Vnet
	for _, br := range m.bridges {
		swIf := v.SwIf(br.si)
		stag := swIf.Id(v)
		fmt.Fprintf(w, "%20v%20v%15v%25v%20v%10v\n",
			xeth.Netns(br.netns), swIf.Name, stag, br.address, swIf.Name, "yes")
		for si, _ := range br.members {
			fmt.Fprintf(w, "%55v%25v%20v%10v\n", "", si.GetAddress(v), vnet.SiName{V: v, Si: si}, "yes")
		}
		fmt.Fprintf(w, "\n")
		for mac, dev := range br.macs {
			fmt.Fprintf(w, "%55v%25v%20v%10v\n", "", mac, dev.devName, "no")
		}
	}
	return
}

func (m *Main) showBridge(c cli.Commander, w cli.Writer, in *cli.Input) (err error) {
	fmt.Fprintf(w, "%20v%20v%15v   %-20v\n", "namespace", "bridgeName", "stag", "interfaces")
	v := m.Vnet
	for _, br := range m.bridges {
		line := 0
		swIf := v.SwIf(br.si)
		stag := swIf.Id(v)
		fmt.Fprintf(w, "%20v%20v%15v",
			xeth.Netns(br.netns), swIf.Name, stag)
		for si, _ := range br.members {
			if line > 0 {
				fmt.Fprintf(w, "\n%55v", "")
			}
			fmt.Fprintf(w, "   %-20v", vnet.SiName{V: v, Si: si})
			line++
		}
		fmt.Fprintf(w, "\n")
	}
	return
}

func (m *Main) fdbBridgeShow(c cli.Commander, w cli.Writer, in *cli.Input) (err error) {
	var brmPerPort map[int32]uint32

	brmPerPort = make(map[int32]uint32)

	fmt.Fprintf(w, "bridgeByStag\n")
	for _, br := range bridgeByStag {
		fmt.Fprintln(w, br)
	}
	fmt.Fprintf(w, "\nfdbBrmToBri\n")
	for brm, bri := range fdbBrmToBri {
		if count, ok := brmPerPort[bri.portIfindex]; ok {
			brmPerPort[bri.portIfindex] = count + 1
		} else {
			brmPerPort[bri.portIfindex] = 1
		}
		m := vnet.PortsByIndex[bri.memberIfindex]
		fmt.Fprintf(w, "%v %+v %+v\n", m.Ifname, brm, bri)
	}
	fmt.Fprintf(w, "\nbrmPerPort\n")
	for port, count := range brmPerPort {
		fmt.Fprintf(w, "port %v, count %v\n", port, count)
	}

	return
}

func (m *Main) cliInit(v *vnet.Vnet) {
	cmds := [...]cli.Command{
		cli.Command{
			Name:      "show neighbor",
			ShortHelp: "show neighbors",
			Action:    m.showIpNeighbor,
		},
		cli.Command{
			Name:      "show bridge",
			ShortHelp: "help",
			Action:    m.showBridge,
		},
		cli.Command{
			Name:      "show br",
			ShortHelp: "help",
			Action:    m.fdbBridgeShow,
		},
		cli.Command{
			Name:      "show mac",
			ShortHelp: "help",
			Action:    m.showMac,
		},
	}
	for i := range cmds {
		v.CliAdd(&cmds[i])
	}
}
