package kubetools

import (
	"path/filepath"

	"github.com/sirupsen/logrus"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type KubeTools struct {
	config    *rest.Config
	clientset *kubernetes.Clientset
	context   string
}

type ConnectionType string

func NewKubetools(kubecontext string) (*KubeTools, error) {
	k := &KubeTools{
		context: kubecontext,
	}
	if err := k.setConfig(); err != nil {
		return nil, err
	}

	if err := k.setClientset(); err != nil {
		return nil, err
	}

	return k, nil
}

func (k *KubeTools) setConfig() error {
	homedir := homedir.HomeDir()
	kubeconfig := filepath.Join(homedir, ".kube", "config")
	var err error
	k.config, err = clientcmd.BuildConfigFromFlags(k.context, kubeconfig)
	if err != nil {
		logrus.Panic(err)
	}
	return err
}

func (k *KubeTools) setClientset() error {
	var err error
	k.clientset, err = kubernetes.NewForConfig(k.config)
	if err != nil {
		logrus.Panic(err)
	}
	return err
}

func (k *KubeTools) TestForward() {

}

// func (k *KubeTools) Clientset() kubernetes.Clientset {
// 	if k.clientset == nil {
// 		err := k.setClientset()
// 		if err != nil {
// 			logrus.Panic(err)
// 		}
// 	}
// 	return *k.clientset
// }

// func (k *KubeTools) IsProcessActive(namespace string, podName string, containerName string, pid int32) bool {

// 	clientset := k.Clientset()
// 	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
// 	if k8errors.IsNotFound(err) {
// 		logrus.Debugf("Pod %s not found in default namespace", podName)
// 		return false
// 	} else if statusError, isStatus := err.(*k8errors.StatusError); isStatus {
// 		logrus.Fatalf("Error getting pod %s %v", podName, statusError.ErrStatus.Message)
// 	} else if err != nil {
// 		logrus.Panic(err.Error())
// 	} else if pod.Status.Phase != "Running" {
// 		logrus.Debug("Pod % is not running. Status: %s", pod.Status.Phase)
// 		return false
// 	} else {
// 		logrus.Debugf("Pod %s in %s namespace is Running", podName, namespace)
// 		cmd := []string{
// 			"sh",
// 			"-c",
// 			fmt.Sprintf("[ -d \"/proc/%d\" ]", pid),
// 		}
// 		req := clientset.CoreV1().RESTClient().Post().Resource("pods").Name(podName).Namespace(namespace).SubResource("exec")
// 		option := &v1.PodExecOptions{
// 			Container: containerName,
// 			Command:   cmd,
// 			Stdin:     true,
// 			Stdout:    true,
// 			Stderr:    true,
// 			TTY:       true,
// 		}

// 		// FIXME - set stdin/out/err to dummy
// 		stdin := os.Stdin
// 		stdout := os.Stdout
// 		stderr := os.Stderr

// 		if stdin == nil {
// 			option.Stdin = false
// 		}
// 		req.VersionedParams(
// 			option,
// 			scheme.ParameterCodec,
// 		)
// 		exec, err := remotecommand.NewSPDYExecutor(k.config, "POST", req.URL())
// 		if err != nil {
// 			logrus.Error(err)
// 		}
// 		err = exec.Stream(remotecommand.StreamOptions{
// 			Stdin:  stdin,
// 			Stdout: stdout,
// 			Stderr: stderr,
// 		})
// 		if err != nil {
// 			errString := fmt.Sprintf("%s", err)
// 			{
// 				matched, err := regexp.MatchString(`.*connection timed out`, errString)
// 				if err != nil {
// 					logrus.Error(err)
// 				}
// 				if matched {
// 					logrus.Debugf("Time out: %d", pid)
// 					return true
// 				}
// 			}
// 			{
// 				matched, err := regexp.MatchString(`.*No such file or directory`, errString)
// 				if err != nil {
// 					logrus.Error(err)
// 				}
// 				if matched {
// 					logrus.Debugf("The process %d does not seem to be working", pid)
// 					return false
// 				}
// 			}
// 			return false
// 		}
// 		return true

// 		/*
// 			some examples
// 			https://miminar.fedorapeople.org/_preview/openshift-enterprise/registry-redeploy/go_client/executing_remote_processes.html
// 			https://stackoverflow.com/questions/43314689/example-of-exec-in-k8ss-pod-by-using-go-client
// 			https://sourcegraph.com/github.com/kubernetes/kubernetes@6900f8856f8cd9a6c94a156b9e4a9fee0c16f807/-/blob/pkg/kubectl/cmd/exec.go
// 			https://gitlab.com/gitlab-org/gitlab-runner/-/blob/master/executors/kubernetes/exec.go
// 		*/
// 	}
// 	return true
// }

// func (k *KubeTools) GetCurrentReplicas(namespace string, deploymentName string) int32 {
// 	clientset := k.Clientset()
// 	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	replicas := deployment.Spec.Replicas
// 	return *replicas
// }

// func (k *KubeTools) SetReplicas(namespace string, deploymentName string, replicas int32) {
// 	clientset := k.Clientset()
// 	scale := &autoscalingv1.Scale{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      deploymentName,
// 			Namespace: namespace,
// 		},
// 		Spec: autoscalingv1.ScaleSpec{
// 			Replicas: replicas,
// 		},
// 	}
// 	d, err := clientset.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deploymentName, scale, metav1.UpdateOptions{})
// 	if err != nil {
// 		logrus.Errorf("Scaling error: %s", err)
// 		return
// 	}
// 	logrus.Infof("Replicas after scaling: %d", d.Spec.Replicas)
// }

// func (k *KubeTools) DeletePodsBySelector(namespace string, selector map[string]interface{}, numberOfPodsLimit int) {

// 	matchLabels := make(map[string]string)
// 	for k, v := range selector {
// 		matchLabels[k] = v.(string)
// 	}

// 	labelSelector := metav1.LabelSelector{
// 		MatchLabels: matchLabels,
// 	}

// 	clientset := k.Clientset()
// 	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
// 		LabelSelector: labels.Set(labelSelector.MatchLabels).String(),
// 		Limit:         int64(numberOfPodsLimit),
// 	})
// 	if err != nil {
// 		logrus.Errorf("Get pods error: %s", err)
// 		return
// 	}
// 	for i := range pods.Items {
// 		pod := pods.Items[i]
// 		logrus.Infof("Delete pod: %s", pod.Name)
// 		clientset.CoreV1().Pods(namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
// 	}

// }
