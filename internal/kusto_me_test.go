package internal_test

import (
	"os"
	"path"
	"testing"

	"github.com/steffakasid/kusto-me/internal"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"sigs.k8s.io/kustomize/api/types"
)

func TestKustomizeMeSimple(t *testing.T) {
	to := internal.KustoMe{
		ApplicationName:          "UnitTest",
		ApplicationDefaultLabels: []string{"label1:value"},
		ApplicationRootFolder:    "../test",
		ApplicationFiles:         []string{"deployment.yml", "crd.yml"},
	}
	to.KustomizeMe(false)
	assert.FileExists(t, path.Join(to.ApplicationRootFolder, internal.KustomizationFilename))
	bt, err := os.ReadFile(path.Join(to.ApplicationRootFolder, internal.KustomizationFilename))
	assert.NoError(t, err)
	kustomization := &types.Kustomization{}
	err = yaml.Unmarshal(bt, kustomization)
	assert.NoError(t, err)

	assert.Equal(t, to.ApplicationName, kustomization.MetaData.Name)
	expectedLabels := map[string]string{"github.com.steffakasid.kusto-me/app": to.ApplicationName, "label1": "value"}
	assert.Equal(t, expectedLabels, kustomization.CommonLabels)
	assert.Contains(t, kustomization.Crds, "crd.yml")
	assert.NotContains(t, kustomization.Crds, "deployment.yml")
	assert.Contains(t, kustomization.Resources, "deployment.yml")
	assert.NotContains(t, kustomization.Resources, "someotherfile.txt")
	assert.NotContains(t, kustomization.Resources, "crd.yml")

	defer func() {
		err = os.Remove(path.Join(to.ApplicationRootFolder, internal.KustomizationFilename))
		assert.NoError(t, err)
	}()
}
