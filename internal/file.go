package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	KustomizationOverlays = "overlays"
	KustomizationBase     = "base"
	KustomizationFilename = "kustomization.yaml"
)

const Permissions = 0755

func GetFiles(dir string) []string {
	files := []string{}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if ext := strings.ToLower(filepath.Ext(path)); (ext == ".yml" || ext == ".yaml") && info.Name() != KustomizationFilename {
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
	if _, err := os.Stat(targetFolder); err != nil {
		err = os.Mkdir(targetFolder, Permissions)
		if err != nil {
			panic(err)
		}
	}

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
			err := os.Mkdir(targetPath, Permissions)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetSubPath(relativePath string) string {
	if strings.Contains(relativePath, "/") {
		return relativePath[0:strings.LastIndex(relativePath, "/")]
	}
	return ""
}

func WriteYaml(y interface{}, path string) error {
	bt, err := yaml.Marshal(y)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, bt, Permissions)
	if err != nil {
		return err
	}
	return nil
}
