// Copyright 2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ethernet

import (
	"github.com/platinasystems/elib/dep"
)

type MacAddDelHookVec struct {
	deps  dep.Deps
	hooks []MacAddDelHook
}

func (t *MacAddDelHookVec) Len() int {
	return t.deps.Len()
}

func (t *MacAddDelHookVec) Get(i int) MacAddDelHook {
	return t.hooks[t.deps.Index(i)]
}

func (t *MacAddDelHookVec) Add(x MacAddDelHook, ds ...*dep.Dep) {
	if len(ds) == 0 {
		t.deps.Add(&dep.Dep{})
	} else {
		t.deps.Add(ds[0])
	}
	t.hooks = append(t.hooks, x)
}
