/*
Copyright (c) 2016, Fabrizio Soppelsa
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of ansible-swarm nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

=============================================================================

Ansible docker_cert_install module
Requires Ansible 2.2+

GOOS=linux GOARCH=amd64 \
	go build -o library/docker_cert_install library/docker_cert_install.go
*/

package main

import (
	"encoding/json"
	"github.com/docker/machine/libmachine/auth"
	"github.com/docker/machine/libmachine/cert"
	"github.com/fsoppelsa/ansible"
	"os"
)

/*
  Sample:

    - name: Generate certificates
      docker_cert_install:
        cert_dir: "/etc/docker/"
        cacert_path: "/etc/docker/ca.pem"
        caprivate_path: "/etc/docker/ca-key.pem"
        servercert_path: "/etc/docker/server.pem"
        servercertkey_path: "/etc/docker/server-key.pem"
*/

type ModuleArgs struct {
	Cert_dir           string
	Cacert_path        string
	Caprivate_path     string
	Servercert_path    string
	Servercertkey_path string
}

func createCertificates(options *ModuleArgs) error {
	var authOptions auth.Options

	authOptions.CertDir = options.Cert_dir
	authOptions.CaCertPath = options.Cacert_path
	authOptions.CaPrivateKeyPath = options.Caprivate_path
	authOptions.ServerCertPath = options.Servercert_path
	authOptions.ServerKeyPath = options.Servercertkey_path

	err := cert.GenerateCACertificate(
		authOptions.CaCertPath,
		authOptions.CaPrivateKeyPath,
		"org",
		2048)

	if err != nil {
		return err
	}

	err = cert.GenerateCert(&cert.Options{
		Hosts:       []string{"127.0.0.1"},
		CertFile:    authOptions.ServerCertPath,
		KeyFile:     authOptions.ServerKeyPath,
		CAFile:      authOptions.CaCertPath,
		CAKeyFile:   authOptions.CaPrivateKeyPath,
		Org:         "org",
		Bits:        2048,
		SwarmMaster: false,
	})

	return err
}

func main() {
	var response ansible.Response
	var moduleArgs ModuleArgs

	text := ansible.ParseVariables(os.Args)

	if err := json.Unmarshal(text, &moduleArgs); err != nil {
		response.Msg = "Configuration file not valid JSON: " + os.Args[1]
		ansible.FailJson(response)
	}

	// Create certs
	if err := createCertificates(&moduleArgs); err != nil {
		response.Msg = "Certs generated ok"
		ansible.ExitJson(response)
	} else {
		response.Msg = "Certs generation ERROR"
		ansible.ExitJson(response)
	}

}
