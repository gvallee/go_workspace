// Copyright (c) 2021, NVIDIA CORPORATION. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package workspace

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gvallee/go_util/pkg/util"
)

func TestCreate(t *testing.T) {
	var err error
	newWorkspace := new(Config)
	// By default the workspace configuration directory and file are created in $HOME
	// which we do not want while testing
	newWorkspace.ConfDir, err = ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("unable to create test directory: %s", err)
	}
	defer os.RemoveAll(newWorkspace.ConfDir)
	newWorkspace.Basedir = newWorkspace.ConfDir
	newWorkspace.Name = "test_workspace"

	// First time we call Load without a configuration file, it should fail and the configuration
	// file be created
	err = newWorkspace.Load()
	if err == nil {
		t.Fatalf("workspace creation without a pre-existing configuration file succeeded")
	}

	configFile := filepath.Join(newWorkspace.ConfDir, "."+newWorkspace.Name, "workspace.conf")
	if !util.PathExists(configFile) {
		t.Fatalf("workspace configuration file was not properly created: %s", err)
	}

	// If we load the workspace again, it should succeed
	err = newWorkspace.Load()
	if err != nil {
		t.Fatalf("loading the workspace failed: %s", err)
	}
}
