// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"context"
	"fmt"
	driverapi "github.com/atomix/atomix-api/go/atomix/management/driver"
	primitiveapi "github.com/atomix/atomix-api/go/atomix/primitive"
	protocolapi "github.com/atomix/atomix-api/go/atomix/protocol"
	"github.com/atomix/atomix-go-framework/pkg/atomix/cluster"
	"github.com/atomix/atomix-go-framework/pkg/atomix/driver"
	"github.com/atomix/atomix-go-framework/pkg/atomix/driver/env"
	"github.com/atomix/atomix-go-framework/pkg/atomix/driver/proxy"
	rsmdriver "github.com/atomix/atomix-go-framework/pkg/atomix/driver/proxy/rsm"
	rsmcounterproxy "github.com/atomix/atomix-go-framework/pkg/atomix/driver/proxy/rsm/counter"
	rsmelectionproxy "github.com/atomix/atomix-go-framework/pkg/atomix/driver/proxy/rsm/election"
	rsmindexedmapproxy "github.com/atomix/atomix-go-framework/pkg/atomix/driver/proxy/rsm/indexedmap"
	rsmleaderproxy "github.com/atomix/atomix-go-framework/pkg/atomix/driver/proxy/rsm/leader"
	rsmlistproxy "github.com/atomix/atomix-go-framework/pkg/atomix/driver/proxy/rsm/list"
	rsmlockproxy "github.com/atomix/atomix-go-framework/pkg/atomix/driver/proxy/rsm/lock"
	rsmlogproxy "github.com/atomix/atomix-go-framework/pkg/atomix/driver/proxy/rsm/log"
	rsmmapproxy "github.com/atomix/atomix-go-framework/pkg/atomix/driver/proxy/rsm/map"
	rsmsetproxy "github.com/atomix/atomix-go-framework/pkg/atomix/driver/proxy/rsm/set"
	rsmvalueproxy "github.com/atomix/atomix-go-framework/pkg/atomix/driver/proxy/rsm/value"
	rsmprotocol "github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm"
	rsmcounterprotocol "github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm/counter"
	rsmelectionprotocol "github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm/election"
	rsmindexedmapprotocol "github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm/indexedmap"
	rsmleaderprotocol "github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm/leader"
	rsmlistprotocol "github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm/list"
	rsmlockprotocol "github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm/lock"
	rsmlogprotocol "github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm/log"
	rsmmapprotocol "github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm/map"
	rsmsetprotocol "github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm/set"
	rsmvalueprotocol "github.com/atomix/atomix-go-framework/pkg/atomix/storage/protocol/rsm/value"
	"github.com/atomix/atomix-go-local/pkg/atomix/local"
	"google.golang.org/grpc"
)

func NewRSMTest() *RSMTest {
	return &RSMTest{
		config: protocolapi.ProtocolConfig{
			Replicas: []protocolapi.ProtocolReplica{
				{
					ID:      "rsm-1",
					NodeID:  "node-1",
					Host:    "localhost",
					APIPort: 7001,
				},
			},
			Partitions: []protocolapi.ProtocolPartition{
				{
					PartitionID: 1,
					Replicas:    []string{"rsm-1"},
				},
			},
		},
	}
}

type RSMTest struct {
	config   protocolapi.ProtocolConfig
	protocol *rsmprotocol.Node
	drivers  []*driver.Driver
}

func (t *RSMTest) Start() error {
	t.protocol = rsmprotocol.NewNode(cluster.NewCluster(t.config, cluster.WithMemberID("rsm-1")), local.NewProtocol())
	rsmcounterprotocol.RegisterService(t.protocol)
	rsmelectionprotocol.RegisterService(t.protocol)
	rsmindexedmapprotocol.RegisterService(t.protocol)
	rsmleaderprotocol.RegisterService(t.protocol)
	rsmlistprotocol.RegisterService(t.protocol)
	rsmlockprotocol.RegisterService(t.protocol)
	rsmlogprotocol.RegisterService(t.protocol)
	rsmmapprotocol.RegisterService(t.protocol)
	rsmsetprotocol.RegisterService(t.protocol)
	rsmvalueprotocol.RegisterService(t.protocol)
	err := t.protocol.Start()
	if err != nil {
		return err
	}
	return nil
}

func (t *RSMTest) CreateProxy(primitiveID primitiveapi.PrimitiveId) (*grpc.ClientConn, error) {
	protocolFunc := func(rsmCluster cluster.Cluster, driverEnv env.DriverEnv) proxy.Protocol {
		protocol := rsmdriver.NewProtocol(rsmCluster, driverEnv)
		rsmcounterproxy.Register(protocol)
		rsmelectionproxy.Register(protocol)
		rsmindexedmapproxy.Register(protocol)
		rsmleaderproxy.Register(protocol)
		rsmlistproxy.Register(protocol)
		rsmlockproxy.Register(protocol)
		rsmlogproxy.Register(protocol)
		rsmmapproxy.Register(protocol)
		rsmsetproxy.Register(protocol)
		rsmvalueproxy.Register(protocol)
		return protocol
	}

	driverPort := 5252 + len(t.drivers)
	driver := driver.NewDriver(protocolFunc, driver.WithNamespace("test"), driver.WithDriverID("rsm"), driver.WithPort(driverPort))
	err := driver.Start()
	if err != nil {
		return nil, err
	}
	t.drivers = append(t.drivers, driver)

	driverConn, err := grpc.Dial(fmt.Sprintf("localhost:%d", driverPort), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer driverConn.Close()
	driverClient := driverapi.NewDriverClient(driverConn)

	agentPort := int32(55680 + len(t.drivers))
	agentID := driverapi.AgentId{
		Namespace: "test",
		Name:      "rsm",
	}
	agentAddress := driverapi.AgentAddress{
		Host: "localhost",
		Port: agentPort,
	}
	agentConfig := driverapi.AgentConfig{
		Protocol: t.config,
	}

	_, err = driverClient.StartAgent(context.TODO(), &driverapi.StartAgentRequest{AgentID: agentID, Address: agentAddress, Config: agentConfig})
	if err != nil {
		return nil, err
	}

	agentConn, err := grpc.Dial(fmt.Sprintf("localhost:%d", agentPort), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	agentClient := driverapi.NewAgentClient(agentConn)

	proxyOptions := driverapi.ProxyOptions{
		Read:  true,
		Write: true,
	}
	_, err = agentClient.CreateProxy(context.TODO(), &driverapi.CreateProxyRequest{ProxyID: driverapi.ProxyId{primitiveID}, Options: proxyOptions})
	if err != nil {
		return nil, err
	}
	return agentConn, nil
}

func (t *RSMTest) Stop() error {
	for _, driver := range t.drivers {
		err := driver.Stop()
		if err != nil {
			return err
		}
	}

	if t.protocol != nil {
		err := t.protocol.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}
