// Copyright (C) 2019-2022  Nicola Murino
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package plugin

import (
	"crypto/sha256"
	"fmt"
	"os/exec"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/sftpgo/sdk/plugin/metadata"

	"github.com/drakkan/sftpgo/v2/internal/logger"
)

type metadataPlugin struct {
	config    Config
	metadater metadata.Metadater
	client    *plugin.Client
}

func newMetadaterPlugin(config Config) (*metadataPlugin, error) {
	p := &metadataPlugin{
		config: config,
	}
	if err := p.initialize(); err != nil {
		logger.Warn(logSender, "", "unable to create metadata plugin: %v, config %+v", err, config)
		return nil, err
	}
	return p, nil
}

func (p *metadataPlugin) exited() bool {
	return p.client.Exited()
}

func (p *metadataPlugin) cleanup() {
	p.client.Kill()
}

func (p *metadataPlugin) initialize() error {
	killProcess(p.config.Cmd)
	logger.Debug(logSender, "", "create new metadata plugin %#v", p.config.Cmd)
	var secureConfig *plugin.SecureConfig
	if p.config.SHA256Sum != "" {
		secureConfig.Checksum = []byte(p.config.SHA256Sum)
		secureConfig.Hash = sha256.New()
	}
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: metadata.Handshake,
		Plugins:         metadata.PluginMap,
		Cmd:             exec.Command(p.config.Cmd, p.config.Args...),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolGRPC,
		},
		Managed:      false,
		AutoMTLS:     p.config.AutoMTLS,
		SecureConfig: secureConfig,
		Logger: &logger.HCLogAdapter{
			Logger: hclog.New(&hclog.LoggerOptions{
				Name:        fmt.Sprintf("%v.%v", logSender, metadata.PluginName),
				Level:       pluginsLogLevel,
				DisableTime: true,
			}),
		},
	})
	rpcClient, err := client.Client()
	if err != nil {
		logger.Debug(logSender, "", "unable to get rpc client for plugin %#v: %v", p.config.Cmd, err)
		return err
	}
	raw, err := rpcClient.Dispense(metadata.PluginName)
	if err != nil {
		logger.Debug(logSender, "", "unable to get plugin %v from rpc client for command %#v: %v",
			metadata.PluginName, p.config.Cmd, err)
		return err
	}

	p.client = client
	p.metadater = raw.(metadata.Metadater)

	return nil
}
