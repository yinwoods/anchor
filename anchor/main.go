// Copyright 2018 The Liman Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/yinwoods/anchor/anchor/handlers"
)

func main() {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
	flag.Lookup("v").Value.Set("5")
	glog.V(2).Info("Listening http://0.0.0.0:8090")

	handlers.ServerRun()
}
