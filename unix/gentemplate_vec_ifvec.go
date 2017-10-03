// autogenerated: do not edit!
// generated from gentemplate [gentemplate -d Package=unix -id ifVec -d VecType=interfaceVec -d Type=*tuntap_interface github.com/platinasystems/go/elib/vec.tmpl]

// Copyright 2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unix

import (
	"github.com/platinasystems/go/elib"
)

type interfaceVec []*tuntap_interface

func (p *interfaceVec) Resize(n uint) {
	old_cap := uint(cap(*p))
	new_len := uint(len(*p)) + n
	if new_len > old_cap {
		new_cap := elib.NextResizeCap(new_len)
		q := make([]*tuntap_interface, new_len, new_cap)
		copy(q, *p)
		*p = q
	}
	*p = (*p)[:new_len]
}

func (p *interfaceVec) validate(new_len uint, zero *tuntap_interface) **tuntap_interface {
	old_cap := uint(cap(*p))
	old_len := uint(len(*p))
	if new_len <= old_cap {
		// Need to reslice to larger length?
		if new_len > old_len {
			*p = (*p)[:new_len]
			for i := old_len; i < new_len; i++ {
				(*p)[i] = zero
			}
		}
		return &(*p)[new_len-1]
	}
	return p.validateSlowPath(zero, old_cap, new_len, old_len)
}

func (p *interfaceVec) validateSlowPath(zero *tuntap_interface, old_cap, new_len, old_len uint) **tuntap_interface {
	if new_len > old_cap {
		new_cap := elib.NextResizeCap(new_len)
		q := make([]*tuntap_interface, new_cap, new_cap)
		copy(q, *p)
		for i := old_len; i < new_cap; i++ {
			q[i] = zero
		}
		*p = q[:new_len]
	}
	if new_len > old_len {
		*p = (*p)[:new_len]
	}
	return &(*p)[new_len-1]
}

func (p *interfaceVec) Validate(i uint) **tuntap_interface {
	var zero *tuntap_interface
	return p.validate(i+1, zero)
}

func (p *interfaceVec) ValidateInit(i uint, zero *tuntap_interface) **tuntap_interface {
	return p.validate(i+1, zero)
}

func (p *interfaceVec) ValidateLen(l uint) (v **tuntap_interface) {
	if l > 0 {
		var zero *tuntap_interface
		v = p.validate(l, zero)
	}
	return
}

func (p *interfaceVec) ValidateLenInit(l uint, zero *tuntap_interface) (v **tuntap_interface) {
	if l > 0 {
		v = p.validate(l, zero)
	}
	return
}

func (p *interfaceVec) ResetLen() {
	if *p != nil {
		*p = (*p)[:0]
	}
}

func (p interfaceVec) Len() uint { return uint(len(p)) }
