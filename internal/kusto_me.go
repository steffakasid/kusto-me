package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/ghodss/yaml"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/kustomize/api/types"
)

const (
	labelKeyApp     = "github.com.steffakasid.kusto-me/app"
	labelKeyVariant = "variant"
)

type KustoMe struct {
	ApplicationName          string
	ApplicationDefaultLabels []string
	ApplicationRootFolder    string
	ApplicationBaseFolder    string
	ApplicationOverlayFolder string
	ApplicationFiles         []string
	CRDs                     []string
	ApplicationOverlays      []string
}

func (k KustoMe) KustomizeMe(overlay bool) {
	targetPath := path.Join(k.ApplicationRootFolder, KustomizationFilename)
	kustomization := k.init(targetPath)

	meta := kustomization.MetaData
	if meta == nil {
		meta = &types.ObjectMeta{}
	}
	meta.Name = k.ApplicationName
	kustomization.MetaData = meta

	commonLabels := map[string]string{labelKeyApp: k.ApplicationName}
	if len(k.ApplicationDefaultLabels) > 0 {
		for _, l := range k.ApplicationDefaultLabels {
			label := strings.Split(l, ":")
			commonLabels[label[0]] = label[1]
		}
	}

	k.identifyCRDs()

	kustomization.CommonLabels = mergeMaps(commonLabels, kustomization.CommonLabels)
	kustomization.Resources = mergeArrays(k.ApplicationFiles, kustomization.Resources)
	kustomization.Crds = mergeArrays(k.CRDs, kustomization.Crds)

	if overlay {
		targetPath = path.Join(k.ApplicationBaseFolder, KustomizationFilename)
		k.WriteOverlays(*kustomization)
	}

	err := WriteYaml(kustomization, targetPath)
	if err != nil {
		panic(err)
	}
}

func (k KustoMe) init(filePath string) *types.Kustomization {
	kustomizationYaml := &types.Kustomization{}
	if _, err := os.Stat(filePath); err == nil {
		bt, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println(err)
		} else {
			err = yaml.Unmarshal(bt, kustomizationYaml)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return kustomizationYaml
}

func (k *KustoMe) identifyCRDs() {
	for _, file := range k.ApplicationFiles {
		bt, err := os.ReadFile(path.Join(k.ApplicationRootFolder, file))
		if err != nil {
			fmt.Println(err)
		} else {
			crds := &apiextensionsv1.CustomResourceDefinition{}
			jsonbt, err := yaml.YAMLToJSON(bt)
			if err != nil {
				fmt.Println(err)
			}
			err = json.Unmarshal(jsonbt, crds)
			if err != nil {
				fmt.Println(err)
			}
			if crds.Kind == "CustomResourceDefinition" {
				k.CRDs = appendStrings(k.CRDs, file)
				k.ApplicationFiles = removeFromArray(k.ApplicationFiles, file)
			}
		}
	}
}

func (k KustoMe) WriteOverlays(base types.Kustomization) {
	// TODO: Check if overlay exists
	k.initOverlayStructure()

	MoveFiles(k.ApplicationFiles, k.ApplicationRootFolder, k.ApplicationBaseFolder)

	for _, o := range k.ApplicationOverlays {
		if err := CreatePath(o, k.ApplicationOverlayFolder); err != nil {
			panic(err)
		}
		kustomization := k.CreateOverlay(o, base)
		if err := WriteYaml(kustomization, path.Join(k.ApplicationOverlayFolder, o, KustomizationFilename)); err != nil {
			panic(err)
		}
	}
}

func (k KustoMe) initOverlayStructure() {
	if err := os.Mkdir(k.ApplicationBaseFolder, Permissions); err != nil {
		panic(err)
	}
	if err := os.Mkdir(k.ApplicationOverlayFolder, Permissions); err != nil {
		panic(err)
	}
}

func (k KustoMe) CreateOverlay(name string, base types.Kustomization) types.Kustomization {
	// TODO: Needs tests
	overlay := types.Kustomization{
		NamePrefix: name[0:3] + "-",
		CommonLabels: map[string]string{
			labelKeyVariant: name,
		},
		Resources: []string{},
		Crds:      []string{},
	}

	for _, res := range base.Resources {
		overlay.Resources = append(overlay.Resources, path.Join("../../", KustomizationBase, res))
	}
	for _, crd := range base.Crds {
		overlay.Crds = append(overlay.Crds, path.Join("../../", KustomizationBase, crd))
	}

	return overlay
}
