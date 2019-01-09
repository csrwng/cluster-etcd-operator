package manifests

import (
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/templates/content"

	configv1 "github.com/openshift/api/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	dnsCrdFilename = "cluster-dns-01-crd.yaml"
	dnsCfgFilename = filepath.Join(manifestDir, "cluster-dns-02-config.yml")
)

// DNS generates the cluster-dns-*.yml files.
type DNS struct {
	config   *configv1.DNS
	FileList []*asset.File
}

var _ asset.WritableAsset = (*DNS)(nil)

// Name returns a human friendly name for the asset.
func (*DNS) Name() string {
	return "DNS Config"
}

// Dependencies returns all of the dependencies directly needed to generate
// the asset.
func (*DNS) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
	}
}

// Generate generates the DNS config and its CRD.
func (d *DNS) Generate(dependencies asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	dependencies.Get(installConfig)

	d.config = &configv1.DNS{
		TypeMeta: metav1.TypeMeta{
			APIVersion: configv1.SchemeGroupVersion.String(),
			Kind:       "DNS",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "cluster",
			// not namespaced
		},
		Spec: configv1.DNSSpec{
			BaseDomain: installConfig.Config.BaseDomain,
		},
	}

	configData, err := yaml.Marshal(d.config)
	if err != nil {
		return errors.Wrapf(err, "failed to create %s manifests from InstallConfig", d.Name())
	}

	crdData, err := content.GetBootkubeTemplate(dnsCrdFilename)
	if err != nil {
		return err
	}

	d.FileList = []*asset.File{
		{
			Filename: filepath.Join(manifestDir, dnsCrdFilename),
			Data:     []byte(crdData),
		},
		{
			Filename: dnsCfgFilename,
			Data:     configData,
		},
	}

	return nil
}

// Files returns the files generated by the asset.
func (d *DNS) Files() []*asset.File {
	return d.FileList
}

// Load loads the already-rendered files back from disk.
func (d *DNS) Load(f asset.FileFetcher) (bool, error) {
	crdFile, err := f.FetchByName(filepath.Join(manifestDir, dnsCrdFilename))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	cfgFile, err := f.FetchByName(dnsCfgFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	dnsConfig := &configv1.DNS{}
	if err := yaml.Unmarshal(cfgFile.Data, dnsConfig); err != nil {
		return false, errors.Wrapf(err, "failed to unmarshal %s", dnsCfgFilename)
	}

	fileList := []*asset.File{crdFile, cfgFile}

	d.FileList, d.config = fileList, dnsConfig

	return true, nil
}
