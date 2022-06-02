package k8portforwarder

// Based on https://github.com/justinbarrick/go-k8s-portforward

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/httpstream"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

type K8PortForwarderConfig struct {
	Context   string `mapstructure:"context"`
	Namespace string `mapstructure:"namespace"`
	Name      string `mapstructure:"name"`
	Port      int    `mapstructure:"port"`
	LocalPort int    `mapstructure:"localPort"`
}

type K8PortForwarder struct {
	Key                 string
	Config              *rest.Config
	Clientset           kubernetes.Interface
	PortForwarderConfig K8PortForwarderConfig
	stopChan            chan struct{}
	readyChan           chan struct{}
}

func NewK8PortForwarder(key string, config K8PortForwarderConfig) (*K8PortForwarder, error) {

	if ok := validatePortForwarderConfig(config); !ok {
		return nil, errors.New("Invalid config for " + key)
	}

	kpf := &K8PortForwarder{
		Key:                 key,
		PortForwarderConfig: config,
	}

	var err error
	kpf.Config, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{
			CurrentContext: kpf.PortForwarderConfig.Context,
		},
	).ClientConfig()

	if err != nil {
		return kpf, errors.Wrap(err, "Could not load kubernetes configuration file")
	}

	kpf.Clientset, err = kubernetes.NewForConfig(kpf.Config)
	if err != nil {
		return kpf, errors.Wrap(err, "Could not create kubernetes client")
	}

	return kpf, nil
}

func validatePortForwarderConfig(config K8PortForwarderConfig) bool {
	//TODO
	if config.Context == "" {
		return false
	}
	if config.Name == "" {
		return false
	}
	return true
}

func (kpf *K8PortForwarder) Test() {
	fmt.Printf("TEST: %#v\n", kpf.PortForwarderConfig)
}

func (kpf *K8PortForwarder) Start(ctx context.Context) error {
	kpf.stopChan = make(chan struct{}, 1)
	readyChan := make(chan struct{}, 1)
	errChan := make(chan error, 1)

	dialer, err := kpf.dialer(ctx)
	if err != nil {
		return errors.Wrap(err, "Could not create a dialer")
	}
	ports := []string{
		fmt.Sprintf("%d:%d", kpf.PortForwarderConfig.LocalPort, kpf.PortForwarderConfig.Port),
	}

	discard := ioutil.Discard
	pf, err := portforward.New(dialer, ports, kpf.stopChan, readyChan, discard, discard)
	if err != nil {
		return errors.Wrap(err, "Could not port forward")
	}

	go func() {
		errChan <- pf.ForwardPorts()
	}()

	select {
	case err = <-errChan:
		return errors.Wrap(err, "Could not create port forward")
	case <-readyChan:
		return nil
	}

	return nil
}

func (kpf *K8PortForwarder) Stop() {
	kpf.stopChan <- struct{}{}
}

func (kpf *K8PortForwarder) dialer(ctx context.Context) (httpstream.Dialer, error) {

	//TODO - services!

	url := kpf.Clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Namespace(kpf.PortForwarderConfig.Namespace).
		Name(kpf.PortForwarderConfig.Name).
		SubResource("portforward").URL()

	transport, upgrader, err := spdy.RoundTripperFor(kpf.Config)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create round tripper")
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", url)
	return dialer, nil
}
