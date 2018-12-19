package kubespace

import (
	"bytes"
	b64 "encoding/base64"
	"errors"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	kubernetesErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"strings"
	"text/template"
)

var clusterRoleName string = "kubespace-clusterrole"

type Client struct {
	client *kubernetes.Clientset
	config *restclient.Config
}

func NewClient(kubeconf string) (*Client, error) {

	config, err := clientcmd.BuildConfigFromFlags("", kubeconf)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	client := &Client{
		client: clientset,
		config: config,
	}

	return client, nil
}

func (c *Client) ListNamespaces() (*corev1.NamespaceList, error) {
	list, err := c.client.CoreV1().Namespaces().List(metav1.ListOptions{
		LabelSelector: "kubespace=true",
	})
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (c *Client) CreateIfNotExistClusterRole() error {
	err := c.CreateClusterRole()

	if kubernetesErrors.IsAlreadyExists(err) {
		return nil
	}

	return err
}

func (c *Client) CreateClusterRole() error {
	role := rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name:   clusterRoleName,
			Labels: getLabels(),
		},
		Rules: []rbacv1.PolicyRule{
			rbacv1.PolicyRule{
				APIGroups: []string{""},
				Resources: []string{"nodes", "persistentvolumes", "namespaces"},
				Verbs:     []string{"list"},
			},
			rbacv1.PolicyRule{
				APIGroups: []string{"storage.k8s.io"},
				Resources: []string{"storageclasses"},
				Verbs:     []string{"get", "list", "watch"},
			},
		},
	}

	_, err := c.client.RbacV1().ClusterRoles().Create(&role)
	return err
}

func (c *Client) DeleteClusterRole() error {
	err := c.client.RbacV1().ClusterRoles().Delete(clusterRoleName, &metav1.DeleteOptions{})
	return err
}

func (c *Client) CreateNamespace(namespace string) error {
	_, err := c.client.CoreV1().Namespaces().Create(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   namespace,
			Labels: getLabels(),
		},
	})
	return err
}

func (c *Client) DeleteNamespace(namespace string) error {
	err := c.client.CoreV1().Namespaces().Delete(namespace, &metav1.DeleteOptions{})
	return err
}

func (c *Client) CreateServiceAccount(namespace string) error {
	serviceAccountName := getServiceAccountName(namespace)

	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccountName,
			Namespace: namespace,
			Labels:    getLabels(),
		},
	}

	serviceAccount, err := c.client.CoreV1().ServiceAccounts(namespace).Create(serviceAccount)

	return err
}

func (c *Client) DeleteServiceAccount(namespace string) error {
	err := c.client.CoreV1().ServiceAccounts(namespace).Delete(getServiceAccountName(namespace), &metav1.DeleteOptions{})
	return err
}

func (c *Client) CreateRole(namespace string) error {
	roleName := getRoleName(namespace)

	role := rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      roleName,
			Namespace: namespace,
			Labels:    getLabels(),
		},
		Rules: []rbacv1.PolicyRule{
			rbacv1.PolicyRule{
				APIGroups: []string{"", "extensions", "apps"},
				Resources: []string{"*"},
				Verbs:     []string{"*"},
			},
			rbacv1.PolicyRule{
				APIGroups: []string{"batch"},
				Resources: []string{"jobs", "cronjobs"},
				Verbs:     []string{"*"},
			},
		},
	}

	_, err := c.client.RbacV1().Roles(namespace).Create(&role)
	return err
}

func (c *Client) DeleteRole(namespace string) error {
	err := c.client.RbacV1().Roles(namespace).Delete(getRoleName(namespace), &metav1.DeleteOptions{})
	return err
}

func (c *Client) CreateIfNotExistServiceAccountClusterRoleBinding(namespace string) error {
	err := c.CreateServiceAccountClusterRoleBinding(namespace)

	if kubernetesErrors.IsAlreadyExists(err) {
		return nil
	}

	return err
}

func (c *Client) CreateServiceAccountClusterRoleBinding(namespace string) error {
	serviceAccountName := getServiceAccountName(namespace)
	roleBindingName := getClusterRoleBindingName(namespace)

	roleBinding := rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:   roleBindingName,
			Labels: getLabels(),
		},
		Subjects: []rbacv1.Subject{{
			Name:      serviceAccountName,
			Kind:      "ServiceAccount",
			Namespace: namespace,
		}},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     clusterRoleName,
			APIGroup: "rbac.authorization.k8s.io",
		}}

	_, err := c.client.RbacV1().ClusterRoleBindings().Create(&roleBinding)

	return err
}

func (c *Client) DeleteServiceAccountClusterRoleBinding(namespace string) error {
	err := c.client.RbacV1().ClusterRoleBindings().Delete(getClusterRoleBindingName(namespace), &metav1.DeleteOptions{})
	return err
}

func (c *Client) CreateServiceAccountRoleBinding(namespace string) error {
	serviceAccountName := getServiceAccountName(namespace)
	roleBindingName := getRoleBindingName(namespace)
	roleName := getRoleName(namespace)

	roleBinding := rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      roleBindingName,
			Namespace: namespace,
			Labels:    getLabels(),
		},
		Subjects: []rbacv1.Subject{{
			Name:      serviceAccountName,
			Kind:      "ServiceAccount",
			Namespace: namespace,
		}},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     roleName,
			APIGroup: "rbac.authorization.k8s.io",
		}}

	_, err := c.client.RbacV1().RoleBindings(namespace).Create(&roleBinding)

	return err
}

func (c *Client) DeleteServiceAccountRoleBinding(namespace string) error {
	err := c.client.RbacV1().RoleBindings(namespace).Delete(getRoleBindingName(namespace), &metav1.DeleteOptions{})
	return err
}

func (c *Client) getServiceAccount(namespace string) (*corev1.ServiceAccount, error) {
	serviceAccountName := getServiceAccountName(namespace)
	serviceAccount, err := c.client.CoreV1().ServiceAccounts(namespace).Get(serviceAccountName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return serviceAccount, nil
}

func (c *Client) getSecretName(namespace string) (string, error) {
	sa, err := c.getServiceAccount(namespace)
	if err != nil {
		return "", nil
	}

	// This should probably be changed
	for _, secret := range sa.Secrets {
		if strings.Contains(secret.Name, "token") {
			return secret.Name, nil
		}

	}
	return "", errors.New("Could not find secret name")
}

func (c *Client) getSecret(namespace string) (*corev1.Secret, error) {
	secretName, err := c.getSecretName(namespace)
	if err != nil {
		return nil, err
	}

	secret, err := c.client.CoreV1().Secrets(namespace).Get(secretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (c *Client) GetCertificate(namespace string) (string, error) {
	secret, err := c.getSecret(namespace)
	if err != nil {
		return "", err
	}

	return string(secret.Data["ca.crt"]), nil
}

func (c *Client) GetCertificateB64(namespace string) (string, error) {
	cert, err := c.GetCertificate(namespace)
	if err != nil {
		return "", err
	}

	certB64 := b64.StdEncoding.EncodeToString([]byte(cert))

	return certB64, nil
}

func (c *Client) GetToken(namespace string) (string, error) {
	secret, err := c.getSecret(namespace)
	if err != nil {
		return "", err
	}

	return string(secret.Data["token"]), nil
}

func (c *Client) GetEndpoint() string {
	return c.config.Host
}

func (c *Client) CreateConfiguration(namespace string) (string, error) {
	type Config struct {
		Namespace   string
		Endpoint    string
		Certificate string
		Token       string
	}

	config := `
---
apiVersion: v1
kind: Config

clusters:
- cluster:
    certificate-authority-data: {{.Certificate}}
    server: {{.Endpoint}}
  name: cluster

users:
- name: {{.Namespace}}-user
  user:
    client-key-data: {{.Certificate}}
    token: {{.Token}}

contexts:
- context:
    cluster: cluster
    namespace: {{.Namespace}}
    user: {{.Namespace}}-user
  name: {{.Namespace}}

current-context: {{.Namespace}}`

	endpoint := c.GetEndpoint()

	certificate, err := c.GetCertificateB64(namespace)
	if err != nil {
		return "", err
	}

	token, err := c.GetToken(namespace)
	if err != nil {
		return "", err
	}

	data := Config{
		Namespace:   namespace,
		Endpoint:    endpoint,
		Certificate: certificate,
		Token:       token,
	}

	tmpl := template.Must(template.New("config").Parse(config))

	var generatedTemplate bytes.Buffer
	if err := tmpl.Execute(&generatedTemplate, data); err != nil {
		return "", err
	}

	return generatedTemplate.String(), nil
}

func getServiceAccountName(namespace string) string {
	return namespace + "-user"
}

func getRoleBindingName(namespace string) string {
	return namespace + "-user-view"
}

func getClusterRoleBindingName(namespace string) string {
	return namespace + "-user-clusterrole-binding"
}

func getRoleName(namespace string) string {
	return namespace + "-user-full-access"
}

func getLabels() map[string]string {
	return map[string]string{
		"kubespace": "true",
	}
}
