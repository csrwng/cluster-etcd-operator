package bootkube

import (
	"os"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/templates/content"
)

const (
	hostEtcdServiceKubeSystemFileName = "host-etcd-service.yaml"
)

var _ asset.WritableAsset = (*HostEtcdServiceKubeSystem)(nil)

// HostEtcdServiceKubeSystem is the constant to represent contents of etcd-service.yaml file
type HostEtcdServiceKubeSystem struct {
	fileName string
	FileList []*asset.File
}

// Dependencies returns all of the dependencies directly needed by the asset
func (t *HostEtcdServiceKubeSystem) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Name returns the human-friendly name of the asset.
func (t *HostEtcdServiceKubeSystem) Name() string {
	return "HostEtcdServiceKubeSystem"
}

// Generate generates the actual files by this asset
func (t *HostEtcdServiceKubeSystem) Generate(parents asset.Parents) error {
	t.fileName = hostEtcdServiceKubeSystemFileName
	data, err := content.GetBootkubeTemplate(t.fileName)
	if err != nil {
		return err
	}
	t.FileList = []*asset.File{
		{
			Filename: filepath.Join(content.TemplateDir, t.fileName),
			Data:     []byte(data),
		},
	}
	return nil
}

// Files returns the files generated by the asset.
func (t *HostEtcdServiceKubeSystem) Files() []*asset.File {
	return t.FileList
}

// Load returns the asset from disk.
func (t *HostEtcdServiceKubeSystem) Load(f asset.FileFetcher) (bool, error) {
	file, err := f.FetchByName(filepath.Join(content.TemplateDir, hostEtcdServiceKubeSystemFileName))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	t.FileList = []*asset.File{file}
	return true, nil
}
