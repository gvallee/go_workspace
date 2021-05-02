//
// Copyright (c) 2021, NVIDIA CORPORATION. All rights reserved.
//
// See LICENSE.txt for license information
//

package workspace

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gvallee/go_util/pkg/util"
)

const (
	defaultWPMode = 0755
)

type Config struct {
	// WorkspaceConfigFile is the path the configuration file for the user's workspace
	ConfigFile string

	// Basedir is the workspace's base directory
	Basedir string

	// WorkspaceDownloadDir is the directory where all the downloads are saved
	DownloadDir string

	// ScratchDir is the directory used as scratch in the context of the workspace
	ScratchDir string

	// InstallDir is the directory where all software is installed for the workspace
	InstallDir string

	// BuildDir is the directory where software related to the workspace is configured and installed
	BuildDir string

	// SrcDir is the directory where source code is saved in the context of the workspace
	SrcDir string

	// MpiDir is the directory where MPI is installed that should be used in the context of the workspace
	MpiDir string

	// MpirunArgs is the list of arguments that the users wants to be passed in when running mpirun commands
	MpirunArgs string
}

func (wp *Config) Init() error {
	if !util.IsDir(wp.Basedir) {
		return fmt.Errorf("the workspace's base directory %s does not exist", wp.Basedir)
	}

	wp.ScratchDir = filepath.Join(wp.Basedir, "scratch")
	wp.DownloadDir = filepath.Join(wp.Basedir, "download")
	wp.SrcDir = filepath.Join(wp.Basedir, "src")
	wp.BuildDir = filepath.Join(wp.Basedir, "build")
	wp.InstallDir = filepath.Join(wp.Basedir, "install")

	if !util.PathExists(wp.DownloadDir) {
		err := os.Mkdir(wp.DownloadDir, defaultWPMode)
		if err != nil {
			return fmt.Errorf("unable to create the workspace's download directory %s: %s", wp.DownloadDir, err)
		}
	}

	if !util.PathExists(wp.ScratchDir) {
		err := os.Mkdir(wp.ScratchDir, defaultWPMode)
		if err != nil {
			return fmt.Errorf("unable to create the workspace's scratch directory %s: %s", wp.ScratchDir, err)
		}
	}

	if !util.PathExists(wp.InstallDir) {
		err := os.Mkdir(wp.InstallDir, defaultWPMode)
		if err != nil {
			return fmt.Errorf("unable to create the workspace's install directory %s: %s", wp.InstallDir, err)
		}
	}

	if !util.PathExists(wp.BuildDir) {
		err := os.Mkdir(wp.BuildDir, defaultWPMode)
		if err != nil {
			return fmt.Errorf("unable to create the workspace's build directory %s: %s", wp.BuildDir, err)
		}
	}

	if !util.PathExists(wp.SrcDir) {
		err := os.Mkdir(wp.SrcDir, defaultWPMode)
		if err != nil {
			return err
		}
	}

	return nil
}
