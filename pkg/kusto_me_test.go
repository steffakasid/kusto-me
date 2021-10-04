package pkg_test

import (
	"os"
	"path"
	"testing"

	"github.com/steffakasid/kusto-me/internal"
	"github.com/steffakasid/kusto-me/pkg"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"sigs.k8s.io/kustomize/api/types"
)

func TestKustomizeMeSimple(t *testing.T) {
	to := pkg.KustoMe{
		ApplicationName:          "UnitTest",
		ApplicationDefaultLabels: []string{"label1:value"},
		ApplicationRootFolder:    "../test",
		ApplicationFiles:         []string{"deployment.yml"},
	}
	to.KustomizeMe(false)
	assert.FileExists(t, path.Join(to.ApplicationRootFolder, internal.KustomizationFilename))
	bt, err := os.ReadFile(path.Join(to.ApplicationRootFolder, internal.KustomizationFilename))
	assert.NoError(t, err)
	kustomization := &types.Kustomization{}
	err = yaml.Unmarshal(bt, kustomization)
	assert.NoError(t, err)
	expected := &types.Kustomization{
		MetaData:     &types.ObjectMeta{Name: to.ApplicationName},
		CommonLabels: map[string]string{"github.com.steffakasid.kusto-me/app": to.ApplicationName, "label1": "value"},
		Resources:    to.ApplicationFiles,
	}
	assert.Equal(t, expected, kustomization)

	defer func() {
		err = os.Remove(path.Join(to.ApplicationRootFolder, internal.KustomizationFilename))
		assert.NoError(t, err)
	}()
}
