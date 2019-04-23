// Copyright 2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ethernet

import (
	"github.com/platinasystems/elib/dep"
)

type BrmAddDelHookVec struct {
	deps  dep.Deps
	hooks []BrmAddDelHook
}

func (t *BrmAddDelHookVec) Len() int {
	return t.deps.Len()
}

func (t *BrmAddDelHookVec) Get(i int) BrmAddDelHook {
	return t.hooks[t.deps.Index(i)]
}

func (t *BrmAddDelHookVec) Add(x BrmAddDelHook, ds ...*dep.Dep) {
	if len(ds) == 0 {
		t.deps.Add(&dep.Dep{})
	} else {
		t.deps.Add(ds[0])
	}
	t.hooks = append(t.hooks, x)
}
