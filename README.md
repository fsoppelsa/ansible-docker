# ansible-docker
Plays for managing Docker and Docker infrastructures.

Requirements:
* Ansible 2.2+
* Golang 1.7.1

In **playbooks/**:
* `docker_install_ubuntu.yml`: Installs the last stable (tried on Ubuntu 16.04)
* `docker_install_113.yml`: Installs the latest edge version from binaries (untested yet)

In **library/**:
* `docker_cert_install.go`: Module for creating certificates. Imports from libmachine

Cross-compile to the target arch:

```
cd library
go get .
GOOS=linux GOARCH=amd64 go build docker_cert_install.go
```

```
ansible-playbook -M library -i inventory playbooks/docker_install_ubuntu.yml
```

Work in progress.
