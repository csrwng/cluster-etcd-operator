package externaldns

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "k8s.io/client-go/kubernetes"

	configv1 "github.com/openshift/api/config/v1"
)

const (
	awsCredsNamespace = "kube-system"
	awsCredsName      = "aws-creds"
)

func newAWSProvider(client kubeclient.Interface, namespace string, dnsConfig *configv1.DNS, clusterVersionConfig *configv1.ClusterVersion) provider {
	return &AWSProvider{
		namespace:            namespace,
		dnsConfig:            dnsConfig,
		clusterVersionConfig: clusterVersionConfig,
		client:               client,
	}
}

type AWSProvider struct {
	dnsConfig            *configv1.DNS
	clusterVersionConfig *configv1.ClusterVersion
	client               kubeclient.Interface
	namespace            string
}

func (p *AWSProvider) Creds() (*corev1.Secret, error) {

	systemCreds, err := p.client.CoreV1().Secrets(awsCredsNamespace).Get(awsCredsName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("Failed to obtain AWS system credentials: %v", err)
	}

	creds := &corev1.Secret{}
	creds.Name = "external-dns-aws-creds"
	creds.Namespace = p.namespace
	creds.Data = map[string][]byte{
		"aws_access_key_id":     systemCreds.Data["aws_access_key_id"],
		"aws_secret_access_key": systemCreds.Data["aws_secret_access_key"],
	}
	return creds, nil
}

func (p *AWSProvider) Env() []corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name: "AWS_SECRET_ACCESS_KEY",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "external-dns-aws-creds",
					},
					Key: "aws_secret_access_key",
				},
			},
		},
		{
			Name: "AWS_ACCESS_KEY_ID",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "external-dns-aws-creds",
					},
					Key: "aws_access_key_id",
				},
			},
		},
	}
}

func (p *AWSProvider) Args() []string {
	return []string{
		"--provider=aws",
		fmt.Sprintf("--domain-filter=%s", p.dnsConfig.Spec.BaseDomain),
		"--aws-zone-type=private",
		fmt.Sprintf("--aws-zone-tags=openshiftClusterID=%s", p.clusterVersionConfig.Spec.ClusterID),
	}
}
