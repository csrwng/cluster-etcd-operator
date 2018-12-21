package etcdconfig

import (
	"fmt"

	kubeclient "k8s.io/client-go/kubernetes"

	"github.com/openshift/library-go/pkg/operator/events"
)

func setupEventRecorder(kubeClient kubeclient.Interface) events.Recorder {
	controllerRef, err := events.GetControllerReferenceForCurrentPod(kubeClient, operatorNamespace, nil)
	if err != nil {
		return &logRecorder{}
	}

	eventsClient := kubeClient.CoreV1().Events(controllerRef.Namespace)
	return events.NewRecorder(eventsClient, "openshift-master-dns-operator", controllerRef)
}

type logRecorder struct{}

func (logRecorder) Event(reason, message string) {
	log.WithValues("reason", reason).Info(message)
}

func (logRecorder) Eventf(reason, messageFmt string, args ...interface{}) {
	log.WithValues("reason", reason).Info(fmt.Sprintf(messageFmt, args...))
}

func (logRecorder) Warning(reason, message string) {
	log.WithValues("reason", reason).Info(message)
}

func (logRecorder) Warningf(reason, messageFmt string, args ...interface{}) {
	log.WithValues("reason", reason).Info(fmt.Sprintf(messageFmt, args...))
}
