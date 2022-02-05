package internal_test

import (
	"os"
	"path"
	"testing"

	"github.com/steffakasid/kusto-me/internal"
	"github.com/stretchr/testify/assert"
)

func TestGetFiles(t *testing.T) {
	files := internal.GetFiles("../test")
	assert.Contains(t, files, "deployment.yml")
	assert.Contains(t, files, "service.yaml")
	assert.Equal(t, 3, len(files))
	assert.NotContains(t, files, "something.txt")
}

func TestMoveFiles(t *testing.T) {
	origFile := []string{"deployment.yml"}
	baseFolder := "../test"
	tgtFolder := "subfolder"

	internal.MoveFiles(origFile, baseFolder, path.Join(baseFolder, tgtFolder))
	assert.FileExists(t, path.Join(baseFolder, tgtFolder, origFile[0]))
	defer func() {
		// Cleanup test files
		internal.MoveFiles(origFile, path.Join(baseFolder, tgtFolder), baseFolder)
		assert.FileExists(t, path.Join("../test/", baseFolder, origFile[0]))
		err := os.Remove(path.Join(baseFolder, tgtFolder))
		assert.NoError(t, err)
	}()
}

func TestCreatePath(t *testing.T) {
	baseFolder := "../test"
	subFolder := "subfolder"
	err := internal.CreatePath(subFolder, baseFolder)
	assert.NoError(t, err)
	err = os.Remove(path.Join(baseFolder, subFolder))
	assert.NoError(t, err)
}

func TestGetSubPath(t *testing.T) {
	somepath := "something/blub"
	result := internal.GetSubPath(somepath)
	assert.Equal(t, result, "something")
}

func TestGetSubPathWithoutSlash(t *testing.T) {
	somepath := "something"
	result := internal.GetSubPath(somepath)
	assert.Empty(t, result)
}

func TestWriteYaml(t *testing.T) {
	someStruct := struct {
		field  string
		field2 string
	}{
		field:  "value",
		field2: "value2",
	}

	internal.WriteYaml(someStruct, "../test/test.yml")
	assert.FileExists(t, "../test/test.yml")

	defer func() {
		err := os.Remove("../test/test.yml")
		assert.NoError(t, err)
	}()

}
