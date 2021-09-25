package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
	"sigs.k8s.io/kustomize/api/types"
)

type stringArrFlag []string

var (
	overlay                   bool
	name, folder              string
	overlayDir, defaultLabels stringArrFlag
)

var defaultOverlays = []string{"development", "production"}

const (
	filename = "kustomization.yaml"
	overlays = "overlays"
	base     = "base"
)

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
		pwd = folder
	} else {
		pwd, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}
	baseFolder := pwd

	if overlay {
		baseFolder = path.Join(pwd, base)
		os.Mkdir(baseFolder, 0755)
	}
	baseKustomize := path.Join(baseFolder, filename)

	files := GetFiles(pwd)
	if overlay {
		overlayFolder := path.Join(pwd, overlays)
		if err := os.Mkdir(overlayFolder, 0755); err != nil {
			panic(err)
		}

		MoveFiles(files, pwd, baseFolder)
		if len(overlayDir) > 0 {
			for _, o := range overlayDir {
				err = CreatePath(o, overlayFolder)
				if err != nil {
					panic(err)
				}
				err = CreateOverlay(o, overlayFolder)
				if err != nil {
					panic(err)
				}
			}
		} else {
			for _, o := range defaultOverlays {
				err = CreatePath(o, overlayFolder)
				if err != nil {
					panic(err)
				}
				err = CreateOverlay(o, overlayFolder)
				if err != nil {
					panic(err)
				}
			}
		}
	}
	kustomization := BaseKustomization(path.Base(pwd), files)
	err = WriteYaml(kustomization, baseKustomize)
	if err != nil {
		panic(err)
	}
}

func BaseKustomization(name string, resources []string) types.Kustomization {
	meta := &types.ObjectMeta{Name: name}
	// TODO: Check if kustomization.yml already exists and read it in (maybe controlled by flag)
	commonLabels := map[string]string{"app": name}
	if len(defaultLabels) > 0 {
		for _, l := range defaultLabels {
			label := strings.Split(l, ":")
			commonLabels[label[0]] = label[1]
		}
	}

	kustomization := types.Kustomization{
		MetaData:     meta,
		CommonLabels: commonLabels,
		Resources:    resources,
	}
	return kustomization
}

func CreateOverlay(name string, overlyFolder string) error {
	kustomization := types.Kustomization{
		NamePrefix: name[0:3] + "-",
		CommonLabels: map[string]string{
			"variant": name,
		},
		Bases: []string{
			path.Join("../../", base),
		},
	}
	err := WriteYaml(kustomization, path.Join(overlyFolder, name, filename))
	return err

}

func GetFiles(dir string) []string {
	files := []string{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if ext := strings.ToLower(filepath.Ext(path)); (ext == ".yml" || ext == ".yaml") && info.Name() != filename {
				subpath := strings.Replace(strings.Replace(path, dir+"/", "", 1), info.Name(), "", 1)
				files = append(files, fmt.Sprintf("%s%s", subpath, info.Name()))
			}
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return files
}

func MoveFiles(files []string, pwd, targetFolder string) {
	for _, f := range files {
		subPath := GetSubPath(f)

		err := CreatePath(subPath, targetFolder)
		if err != nil {
			panic(err)
		}
		err = os.Rename(fmt.Sprintf("%s/%s", pwd, f), fmt.Sprintf("%s/%s", targetFolder, f))
		if err != nil {
			panic(err)
		}
	}
}

func CreatePath(subPath, targetPath string) error {
	pathElem := strings.Split(subPath, "/")
	for _, p := range pathElem {
		targetPath = fmt.Sprintf("%s/%s", targetPath, p)
		if _, err := os.Stat(targetPath); err != nil {
			err := os.Mkdir(targetPath, 0755)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetSubPath(relativePath string) string {
	return relativePath[0:strings.LastIndex(relativePath, "/")]
}

func WriteYaml(y interface{}, path string) error {
	bt, err := yaml.Marshal(y)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, bt, 0777)
	if err != nil {
		return err
	}
	return nil
}
