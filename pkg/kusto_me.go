package pkg

import (
	"os"
	"path"
	"strings"

	"github.com/steffakasid/kusto-me/internal"
	"sigs.k8s.io/kustomize/api/types"
)

const (
	labelKeyApp     = "app"
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
	meta := &types.ObjectMeta{Name: k.ApplicationName}
	// TODO: Check if kustomization.yml already exists and read it in (maybe controlled by flag)
	commonLabels := map[string]string{labelKeyApp: k.ApplicationName}
	if len(k.ApplicationDefaultLabels) > 0 {
		for _, l := range k.ApplicationDefaultLabels {
			label := strings.Split(l, ":")
			commonLabels[label[0]] = label[1]
		}
	}

	kustomization := types.Kustomization{
		MetaData:     meta,
		CommonLabels: commonLabels,
		Resources:    k.ApplicationFiles,
	}

	targetPath := path.Join(k.ApplicationRootFolder, internal.KustomizationFilename)
	if overlay {
		targetPath = path.Join(k.ApplicationBaseFolder, internal.KustomizationFilename)
		k.WriteOverlays()
	}

	err := internal.WriteYaml(kustomization, targetPath)
	if err != nil {
		panic(err)
	}
}

func (k KustoMe) WriteOverlays() {
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
