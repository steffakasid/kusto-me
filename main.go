package main

import (
	"flag"
	"os"
	"path"
	"strings"

	"github.com/steffakasid/kusto-me/internal"
	"github.com/steffakasid/kusto-me/pkg"
)

type stringArrFlag []string

var (
	overlay                   bool
	name, folder              string
	overlayDir, defaultLabels stringArrFlag
)

var defaultOverlays = []string{"development", "production"}

func init() {
	flag.BoolVar(&overlay, "overlay", false, "Defines if overlay structure should be created or just a simple project")
	flag.StringVar(&name, "name", "", "Set the projectname. If not set the current directoryname is used")
	flag.StringVar(&folder, "folder", "", "Set the folder to create kustomize project. If not set current dir is used.")
	flag.Var(&overlayDir, "dir", "Define overlay directories to be used.")
	flag.Var(&defaultLabels, "label", "Add default labels to kustomization.yml. Format: <key>:<value>")
	flag.Parse()
}

func (s *stringArrFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *stringArrFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
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
