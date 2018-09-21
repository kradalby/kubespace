# kubespace
Tool to create namespaces and service accounts that can safely be handed to CI or users

## Install
You need go 1.11 (other might work) and cluster-admin\* access to the kubernetes cluster.
```
go get github.com/kradalby/kubespace
```

A pre-compiled binary is also available from [kubespace.kradalby.no](https://kubespace.kradalby.no)

## Upgrade
```
go get -u github.com/kradalby/kubespace
```

## Usage

```
Usage:
  kubespace [command]

Available Commands:
  config      A brief description of your command
  create      Create a namespace with a restricted account
  delete      Delete a namespace and the restricted account
  drone       Output commands for adding secrets to drone CI
  gitlab      Output configuration for GitLab CI
  help        Help about any command

Flags:
  -h, --help               help for kubespace
  -c, --kubeconf string    Path to kubeconfig (default "/home/kradalby/.kube/config")
  -n, --namespace string   Namespace (required)

Use "kubespace [command] --help" for more information about a command.
```
