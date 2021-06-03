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

	"github.com/gvallee/go_software_build/pkg/builder"
	"github.com/gvallee/go_util/pkg/util"
	"github.com/gvallee/kv/pkg/kv"
)

const (
	defaultWPMode  = 0755
	configFileName = "workspace.conf"
)

type Config struct {
	// Name is the name of the workspace
	Name string

	// WorkspaceConfigFile is the path the configuration file for the user's workspace
	ConfigFile string

	// ConfDir is where the configuration directory and file are. The configuration directory is
	// separate from the base directory of the workspace (where all the data is), so we can have
	// the configuration file is a pre-defined location (HOME by default) and the data stored somewhere
	// with plenty of storage space
	ConfDir string

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

	// RunDir is the drectory from which all the jobs/experiments are executed
	RunDir string

	// MpiDir is the directory where MPI is installed that should be used in the context of the workspace
	MpiDir string

	// MpirunArgs is the list of arguments that the users wants to be passed in when running mpirun commands
	MpirunArgs string
}

func (w *Config) Init() error {
	if !util.IsDir(w.Basedir) {
		return fmt.Errorf("the workspace's base directory %s does not exist", w.Basedir)
	}

	w.ScratchDir = filepath.Join(w.Basedir, "scratch")
	w.DownloadDir = filepath.Join(w.Basedir, "download")
	w.SrcDir = filepath.Join(w.Basedir, "src")
	w.BuildDir = filepath.Join(w.Basedir, "build")
	w.InstallDir = filepath.Join(w.Basedir, "install")
	w.RunDir = filepath.Join(w.Basedir, "run")

	if !util.PathExists(w.DownloadDir) {
		// We use mkdirall for the first one so that is the basedirectory does not exist, it creates it
		err := os.MkdirAll(w.DownloadDir, defaultWPMode)
		if err != nil {
			return fmt.Errorf("unable to create the workspace's download directory %s: %s", w.DownloadDir, err)
		}
	}

	if !util.PathExists(w.ScratchDir) {
		err := os.Mkdir(w.ScratchDir, defaultWPMode)
		if err != nil {
			return fmt.Errorf("unable to create the workspace's scratch directory %s: %s", w.ScratchDir, err)
		}
	}

	if !util.PathExists(w.InstallDir) {
		err := os.Mkdir(w.InstallDir, defaultWPMode)
		if err != nil {
			return fmt.Errorf("unable to create the workspace's install directory %s: %s", w.InstallDir, err)
		}
	}

	if !util.PathExists(w.BuildDir) {
		err := os.Mkdir(w.BuildDir, defaultWPMode)
		if err != nil {
			return fmt.Errorf("unable to create the workspace's build directory %s: %s", w.BuildDir, err)
		}
	}

	if !util.PathExists(w.SrcDir) {
		err := os.Mkdir(w.SrcDir, defaultWPMode)
		if err != nil {
			return err
		}
	}

	if !util.PathExists(w.RunDir) {
		err := os.Mkdir(w.RunDir, defaultWPMode)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Config) getPathToConfigDir() string {
	return filepath.Join(w.ConfDir, "."+w.Name)
}

func (w *Config) getConfigFilePath() string {
	if w.ConfDir == "" {
		w.ConfDir = os.Getenv("HOME")
	}
	return filepath.Join(w.getPathToConfigDir(), configFileName)
}

func (w *Config) createDefaultConfigFile() error {
	// Sanity checks
	if w.ConfDir == "" {
		return fmt.Errorf("configuration directory is undefined")
	}
	if w.ConfigFile == "" {
		return fmt.Errorf("configuration file is undefined")
	}

	configDir := w.getPathToConfigDir()
	fmt.Printf("DBG - checking whether %s exists", configDir)
	if !util.PathExists(configDir) {
		fmt.Printf("DBG - creating directory %s\n", configDir)
		err := os.MkdirAll(configDir, defaultWPMode)
		if err != nil {
			return err
		}
	}

	// If a base directory was not specified up front, use HOME by default
	if w.Basedir == "" {
		w.Basedir = os.Getenv("HOME")
	}
	w.Basedir = filepath.Join(w.Basedir, w.Name+"_ws")
	if !util.PathExists(w.Basedir) {
		err := os.MkdirAll(w.Basedir, defaultWPMode)
		if err != nil {
			return err
		}
	}
	content := "dir=" + w.Basedir + "\n"
	fmt.Printf("DBG - creating %s\n", w.ConfigFile)
	f, err := os.Create(w.ConfigFile)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	// close is deferred and we need to make sure the content is written to the file asap
	err = f.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (w *Config) ParseCfg() error {
	if w.ConfigFile == "" {
		return fmt.Errorf("configuration file is undefined")
	}

	kvs, err := kv.LoadKeyValueConfig(w.ConfigFile)
	if err != nil {
		return err
	}

	for _, keyvalue := range kvs {
		if keyvalue.Key == "dir" {
			w.Basedir = keyvalue.Value
		} else {
			return fmt.Errorf("invalid key (%s)", keyvalue.Key)
		}
	}

	return nil
}

func (w *Config) Load() error {
	// Check if the configuration dir/file exists
	if w.ConfigFile == "" {
		w.ConfigFile = w.getConfigFilePath()
	}
	if !util.FileExists(w.ConfigFile) {
		err := w.createDefaultConfigFile()
		if err != nil {
			return err
		}
		err = w.Init()
		if err != nil {
			return err
		}
		fmt.Printf("warning! new configuration created (%s), please review and customize before re-running the same command", w.ConfigFile)
		return fmt.Errorf("new configuration file created, it needs review")
	}

	// If we get here, we can parse the content of the configuration file
	err := w.ParseCfg()
	if err != nil {
		return err
	}

	return nil
}

func (w *Config) InstallSoftware(softwareName string, softwareURL string, configArgsFn func() []string) error {
	b := new(builder.Builder)
	b.Env.ScratchDir = w.ScratchDir
	b.Env.InstallDir = w.InstallDir
	b.Env.BuildDir = filepath.Join(w.BuildDir, softwareName)
	b.GetConfigureExtraArgs = configArgsFn
	// fixme: it should fail when not specified but it does not other than when trying to untar
	b.Env.SrcDir = b.Env.BuildDir
	// fixme: builder should take care of this
	if !util.PathExists(filepath.Join(b.Env.BuildDir, softwareName)) {
		os.MkdirAll(filepath.Join(b.Env.BuildDir, softwareName), 0777)
	}
	b.App.Name = softwareName
	b.App.URL = softwareURL
	err := b.Load(true)
	if err != nil {
		return err
	}
	res := b.Install()
	if res.Err != nil {
		return res.Err
	}

	return nil
}
