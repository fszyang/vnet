// Copyright 2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ethernet

import (
	"github.com/platinasystems/elib/dep"
)

type BridgeAddDelHookVec struct {
	deps  dep.Deps
	hooks []BridgeAddDelHook
}

func (t *BridgeAddDelHookVec) Len() int {
	return t.deps.Len()
}

func (t *BridgeAddDelHookVec) Get(i int) BridgeAddDelHook {
	return t.hooks[t.deps.Index(i)]
}

func (t *BridgeAddDelHookVec) Add(x BridgeAddDelHook, ds ...*dep.Dep) {
	if len(ds) == 0 {
		t.deps.Add(&dep.Dep{})
	} else {
		t.deps.Add(ds[0])
	}
	t.hooks = append(t.hooks, x)
}
