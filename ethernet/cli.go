// Copyright 2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ethernet

import (
	"github.com/platinasystems/elib/cli"
	"github.com/platinasystems/vnet"
	"github.com/platinasystems/vnet/ip"

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
			ns := fi.Name(im)

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
			intf := n.Si.Name(v)
			lladdr := n.Ethernet.String()

			ai := ip.AdjNil
			ln := 0
			if ai, as, ok = im.GetReachable(&prefix, n.Si); ok {
				for i := range as {
					adj_lines = as[i].String(im)
				}
				if ln == 0 {
					fmt.Fprintf(w, "%6v%20v dev %10v lladdr %v      adjacency %v:%v\n", ns, ipAddr, intf, lladdr, ai, adj_lines)
				} else {
					fmt.Fprintf(w, "%6v%20v dev %10v lladdr %v      adjacency %v:%v\n", "", "unexpected extras", "", "", ai, adj_lines)
				}
				ln++
			} else {
				fmt.Fprintf(w, "%6v%20v dev %10v lladdr %v      adjacency %v:%v\n", ns, ipAddr, intf, lladdr, ai, "not found")
			}

			if cf.detail {
				//no additional details for now
			}
		}
	}
	return
}

func (m *Main) fdbBridgeShow(c cli.Commander, w cli.Writer, in *cli.Input) (err error) {
	var brmPerPort map[int32]uint32

	brmPerPort = make(map[int32]uint32)

	for stag, br := range bridgeByStag {
		fmt.Fprintf(w, "bridgeByStag[%v] %s\n", stag, br)
	}
	fmt.Fprintf(w, "fdbBrmToBri\n")
	for brm, bri := range fdbBrmToBri {
		if count, ok := brmPerPort[bri.port]; ok {
			brmPerPort[bri.port] = count + 1
		} else {
			brmPerPort[bri.port] = 1
		}
		fmt.Fprintf(w, "%+v %+v\n", brm, bri)
	}
	fmt.Fprintf(w, "brmPerPort\n")
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
			Action:    m.fdbBridgeShow,
		},
	}
	for i := range cmds {
		v.CliAdd(&cmds[i])
	}
}
