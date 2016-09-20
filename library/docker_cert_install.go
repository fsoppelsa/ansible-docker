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

Ansible docker_install module
Requires Ansible 2.2+
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

    - name: Install Docker certificates
      docker_cert_install:
		cert_dir: ""
		cacert_path: ""
		caprivate_path: ""
		client_certpath: ""
		client_keypath: ""
      register: cert_install_result

*/

type ModuleArgs struct {
	CertDir        string
	CacertPath     string
	CaprivatePath  string
	ClientCertpath string
	ClientKeypath  string
}

// Use BootstrapCertificates

func createCertificates(options *ModuleArgs) error {
	var authOptions *auth.Options

	authOptions.CertDir = options.CertDir
	authOptions.CaCertPath = options.CacertPath
	authOptions.CaPrivateKeyPath = options.CaprivatePath
	authOptions.ClientCertPath = options.ClientCertpath
	authOptions.ClientKeyPath = options.ClientKeypath

	if err := cert.BootstrapCertificates(authOptions); err != nil {
		return err
	} else {
		return nil
	}
}

func main() {
	var response ansible.Response
	var moduleArgs *ModuleArgs

	text := ansible.ParseVariables(os.Args)

	if err := json.Unmarshal(text, &moduleArgs); err != nil {
		response.Msg = "Configuration file not valid JSON: " + os.Args[1]
		ansible.FailJson(response)
	}

	// Attempt to generate certs
	if err := createCertificates(moduleArgs); err != nil {
		response.Msg = "Certs generated ok"
		ansible.ExitJson(response)
	} else {
		response.Msg = err.Error() //"Certs generetion ERROR"
		ansible.FailJson(response)
	}

}
