package main

import (
	"os"
	"path"

	flag "github.com/spf13/pflag"
	"github.com/steffakasid/kusto-me/internal"
	"github.com/steffakasid/kusto-me/pkg"
)

var (
	overlay                   bool
	name, folder              string
	overlayDir, defaultLabels []string
)

var defaultOverlays = []string{"development", "production"}

func init() {
	flag.BoolVarP(&overlay, "overlay", "o", false, "Defines if overlay structure should be created or just a simple project")
	flag.StringVarP(&name, "name", "n", "", "Set the projectname. If not set the current directoryname is used")
	flag.StringVarP(&folder, "folder", "f", "", "Set the folder to create kustomize project. If not set current dir is used.")
	flag.StringArrayVarP(&overlayDir, "dir", "d", []string{}, "Define overlay directories to be used.")
	flag.StringArrayVarP(&defaultLabels, "label", "l", []string{}, "Add default labels to kustomization.yml. Format: <key>:<value>")
	flag.Parse()
}

func main() {
	var (
		pwd string
		err error
	)

	if folder != "" {
		if _, err := os.Stat(folder); err != nil {
			panic(err)
		}
		pwd = folder
	} else {
		pwd, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}
	appName := path.Base(pwd)
	if name != "" {
		appName = name
	}

	kustoMe := pkg.KustoMe{
		ApplicationName:          appName,
		ApplicationRootFolder:    pwd,
		ApplicationBaseFolder:    path.Join(pwd, internal.KustomizationBase),
		ApplicationOverlayFolder: path.Join(pwd, internal.KustomizationOverlays),
		ApplicationFiles:         internal.GetFiles(pwd),
	}

	if len(overlayDir) > 0 {
		kustoMe.ApplicationOverlays = overlayDir
	} else {
		kustoMe.ApplicationOverlays = defaultOverlays
	}

	kustoMe.KustomizeMe(overlay)
}
