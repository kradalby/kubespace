# kubespace
Tool to create namespaces and service accounts that can safely be handed to CI or users

## Install
kubespace requires Python 3.7 or newer. The other requirement is kubectl in your path and cluster-admin\* access to the kubernetes cluster.

```
pip install git+https://github.com/kradalby/kubespace.git
```

\* The ability to create namespaces

## Upgrade
```
pip install -U git+https://github.com/kradalby/kubespace.git
```

## Usage

```
usage:
kubespace <command> [<args>]
The most commonly used sandwich commands are:
   create     Create a namespace and a restrictive role
   delete     Delete a namespace and a restrictive role
   config     Print the kubeconfig for the restrictive role
   -n/--namespace name of the namespace to manage
```
