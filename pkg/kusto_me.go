package pkg

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/steffakasid/kusto-me/internal"
	"gopkg.in/yaml.v2"
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
	ApplicationOverlays      []string
}

func (k KustoMe) KustomizeMe(overlay bool) {
	targetPath := path.Join(k.ApplicationRootFolder, internal.KustomizationFilename)
	kustomization := k.init(targetPath)

	meta := kustomization.MetaData
	if meta == nil {
		meta = &types.ObjectMeta{}
	}
	meta.Name = k.ApplicationName

	commonLabels := map[string]string{labelKeyApp: k.ApplicationName}
	if len(k.ApplicationDefaultLabels) > 0 {
		for _, l := range k.ApplicationDefaultLabels {
			label := strings.Split(l, ":")
			commonLabels[label[0]] = label[1]
		}
	}

	kustomization.CommonLabels = mergeMaps(commonLabels, kustomization.CommonLabels)
	kustomization.Resources = mergeArrays(k.ApplicationFiles, kustomization.Resources)

	if overlay {
		targetPath = path.Join(k.ApplicationBaseFolder, internal.KustomizationFilename)
		k.WriteOverlays()
	}

	err := internal.WriteYaml(kustomization, targetPath)
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
		}
		err = yaml.UnmarshalStrict(bt, kustomizationYaml)
		if err != nil {
			fmt.Println(err)
		}
	}
	return kustomizationYaml
}

func (k KustoMe) WriteOverlays() {
	// TODO: Check if overlay exists
	k.initOverlayStructure()

	internal.MoveFiles(k.ApplicationFiles, k.ApplicationRootFolder, k.ApplicationBaseFolder)

	for _, o := range k.ApplicationOverlays {
		if err := internal.CreatePath(o, k.ApplicationOverlayFolder); err != nil {
			panic(err)
		}
		kustomization := k.CreateOverlay(o)
		if err := internal.WriteYaml(kustomization, path.Join(k.ApplicationOverlayFolder, o, internal.KustomizationFilename)); err != nil {
			panic(err)
		}
	}
}

func (k KustoMe) initOverlayStructure() {
	if err := os.Mkdir(k.ApplicationBaseFolder, internal.Permissions); err != nil {
		panic(err)
	}
	if err := os.Mkdir(k.ApplicationOverlayFolder, internal.Permissions); err != nil {
		panic(err)
	}
}

func (k KustoMe) CreateOverlay(name string) types.Kustomization {
	return types.Kustomization{
		NamePrefix: name[0:3] + "-",
		CommonLabels: map[string]string{
			labelKeyVariant: name,
		},
		Bases: []string{
			path.Join("../../", internal.KustomizationBase),
		},
	}
}
