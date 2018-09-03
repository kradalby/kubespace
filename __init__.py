#!/usr/bin/env python3

import base64
from subprocess import run



# run('kubectl create -f ', shell=True, check=True)


def get_secret_name(namespace):
    sa = '{}-user'.format(namespace) 
    cli = 'kubectl -n {namespace} get serviceaccount {service_account}'.format(namespace=namespace, service_account=sa) + ' -o jsonpath="{.secrets[0].name}"'

    r = run(cli, shell=True, check=True, capture_output=True)
    return r.stdout.decode('utf-8')


def get_token_b64(namespace, secret_name):
    cli = 'kubectl -n {namespace} get secret {secret_name}'.format(namespace=namespace, secret_name=secret_name) + ' -o jsonpath="{.data.token}"'
    r = run(cli, shell=True, check=True, capture_output=True)
    return r.stdout.decode('utf-8')
    
    
def get_cert_b64(namespace, secret_name):
    cli = 'kubectl -n {namespace} get secret {secret_name}'.format(namespace=namespace, secret_name=secret_name) + '  -o jsonpath="{.data.ca\.crt}"'
    r = run(cli, shell=True, check=True, capture_output=True)
    return r.stdout.decode('utf-8')


def get_endpoint():
    cli = 'kubectl config view -o jsonpath="{.clusters[0].cluster.server}"' 
    r = run(cli, shell=True, check=True, capture_output=True)
    return r.stdout.decode('utf-8')



def create_yaml(yml):
    cli = ''' cat <<EOF | kubectl create -f -
{yml}
EOF'''.format(yml=yml)
    print(cli)
    r = run(cli, shell=True, check=True)
    print(r)


def create_namespace(namespace):
    yaml = '''
---
apiVersion: v1
kind: Namespace
metadata:
  name: {namespace}
'''.format(namespace=namespace)
    create_yaml(yaml)


def create_service_account(namespace):
    yaml = '''
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {namespace}-user
  namespace: {namespace}
'''.format(namespace=namespace)
    create_yaml(yaml)


def create_role(namespace):
    yaml = '''
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: {namespace}-user-full-access
  namespace: {namespace}
rules:
- apiGroups: ["", "extensions", "apps"]
  resources: ["*"]
  verbs: ["*"]
- apiGroups: ["batch"]
  resources:
  - jobs
  - cronjobs
  verbs: ["*"]
'''.format(namespace=namespace)
    create_yaml(yaml)


def create_role_binding(namespace):
    yaml = '''
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1beta1
metadata:
  name: {namespace}-user-view
  namespace: {namespace}
subjects:
- kind: ServiceAccount
  name: {namespace}-user
  namespace: {namespace}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: {namespace}-user-full-access
'''.format(namespace=namespace)
    create_yaml(yaml)


def create_config(namespace, certificate, token, endpoint):
    yaml = '''
---
apiVersion: v1
kind: Config

clusters:
- cluster:
    certificate-authority-data: {certificate}
    # You'll need the API endpoint of your Cluster here:
    server: {endpoint}
  name: esa-cluster

users:
- name: {namespace}-user
  user:
    client-key-data: {certificate}
    token: {token}

contexts:
- context:
    cluster: esa-cluster
    namespace: {namespace}
    user: {namespace}-user
  name: {namespace}

current-context: {namespace}

'''.format(
        namespace=namespace,
        certificate=certificate,
        token=token,
        endpoint=endpoint
        )
    return yaml


if __name__ == '__main__':
    # create_namespace("derp")
    # create_service_account("derp")
    # create_role("derp")
    # create_role_binding("derp")
    # create_config("derp")
    secret = get_secret_name("derp")
    print(get_cert_b64("derp", secret))
    print(base64.b64decode(get_token_b64("derp", secret)))
    print(get_endpoint())
    print(create_config(
        "derp",
        get_cert_b64("derp", secret),
        base64.b64decode(get_token_b64("derp", secret)).decode('utf-8'),
        get_endpoint()
    ))
