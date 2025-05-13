// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && cgo) || (darwin && cgo) || (freebsd && cgo)

package plugin

/*
#cgo linux LDFLAGS: -ldl
#include <dlfcn.h>
#include <limits.h>
#include <stdlib.h>
#include <stdint.h>

#include <stdio.h>

static uintptr_t pluginOpen(const char* path, char** err) {
	void* h = dlopen(path, RTLD_NOW|RTLD_GLOBAL);
	if (h == NULL) {
		*err = (char*)dlerror();
	}
	return (uintptr_t)h;
}

static void* pluginLookup(uintptr_t h, const char* name, char** err) {
	void* r = dlsym((void*)h, name);
	if (r == NULL) {
		*err = (char*)dlerror();
	}
	return r;
}
*/
import "C"

import (
	"errors"
	"sync"
)

// func open(name string) (*Plugin, error) {
// 	return nil, nil
// }

func open(name string) (*Plugin, error) {
	return nil, nil
}

func lookup(p *Plugin, symName string) (Symbol, error) {
	if s := p.syms[symName]; s != nil {
		return s, nil
	}
	return nil, errors.New("plugin: symbol " + symName + " not found in plugin " + p.pluginpath)
}

var (
	pluginsMu  sync.Mutex
	sysPlugins map[string]*Plugin
)

// // lastmoduleinit is defined in package runtime.
// func lastmoduleinit() (pluginpath string, syms map[string]any, inittasks []*initTask, errstr string)

// lastmoduleinit is defined in package runtime.
<<<<<<< HEAD
func lastmoduleinit() (pluginpath string, syms map[string]any, inittasks []*initTask, errstr string) {
	return "", nil, nil, ""
}
=======
//func lastmoduleinit() (pluginpath string, syms map[string]any, inittasks []*initTask, errstr string)
>>>>>>> 936cc6b8fb8001f6c7cefd7fc906fb618ab39656

// doInit is defined in package runtime.
//
//go:linkname doInit runtime.doInit
func doInit(t []*initTask)

type initTask struct {
	// fields defined in runtime.initTask. We only handle pointers to an initTask
	// in this package, so the contents are irrelevant.
}
