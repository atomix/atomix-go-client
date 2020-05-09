// Copyright 2019-present Open Networking Foundation.
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

package database

import (
	"context"
	"fmt"
	databaseapi "github.com/atomix/api/proto/atomix/database"
	"github.com/atomix/go-client/pkg/client/database/counter"
	"github.com/atomix/go-client/pkg/client/database/election"
	"github.com/atomix/go-client/pkg/client/database/indexedmap"
	"github.com/atomix/go-client/pkg/client/database/leader"
	"github.com/atomix/go-client/pkg/client/database/list"
	"github.com/atomix/go-client/pkg/client/database/lock"
	"github.com/atomix/go-client/pkg/client/database/log"
	"github.com/atomix/go-client/pkg/client/database/map"
	"github.com/atomix/go-client/pkg/client/database/partition"
	"github.com/atomix/go-client/pkg/client/database/set"
	"github.com/atomix/go-client/pkg/client/database/value"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/atomix/go-client/pkg/client/protocol"
	"github.com/atomix/go-client/pkg/client/util/net"
	"sort"
)

func New(ctx context.Context, protocol *protocol.Client, databaseProto *databaseapi.Database, opts ...Option) (*Database, error) {
	options := applyOptions(opts...)

	// Ensure the partitions are sorted in case the controller sent them out of order.
	partitionProtos := databaseProto.Partitions
	sort.Slice(partitionProtos, func(i, j int) bool {
		return partitionProtos[i].PartitionID < partitionProtos[j].PartitionID
	})

	// Iterate through the partitions and create gRPC client connections for each partition.
	partitions := make([]partition.Partition, len(databaseProto.Partitions))
	for i, partitionProto := range partitionProtos {
		ep := partitionProto.Endpoints[0]
		partitions[i] = partition.Partition{
			ID:      int(partitionProto.PartitionID),
			Address: net.Address(fmt.Sprintf("%s:%d", ep.Host, ep.Port)),
		}
	}

	// Iterate through partitions and open sessions
	sessions := make([]*partition.Session, len(partitions))
	for i, p := range partitions {
		session, err := partition.NewSession(ctx, p, partition.WithSessionTimeout(options.sessionTimeout))
		if err != nil {
			return nil, err
		}
		sessions[i] = session
	}

	return &Database{
		Client:   protocol,
		sessions: sessions,
	}, nil
}

// Database manages the primitives in a set of partitions
type Database struct {
	*protocol.Client
	sessions []*partition.Session
}

// GetCounter gets or creates a Counter with the given name
func (d *Database) GetCounter(ctx context.Context, name string) (counter.Counter, error) {
	return counter.New(ctx, primitive.NewName(d.Namespace, d.Name, d.Scope, name), d.sessions)
}

// GetElection gets or creates an Election with the given name
func (d *Database) GetElection(ctx context.Context, name string, opts ...election.Option) (election.Election, error) {
	return election.New(ctx, primitive.NewName(d.Namespace, d.Name, d.Scope, name), d.sessions)
}

// GetIndexedMap gets or creates a Map with the given name
func (d *Database) GetIndexedMap(ctx context.Context, name string) (indexedmap.IndexedMap, error) {
	return indexedmap.New(ctx, primitive.NewName(d.Namespace, d.Name, d.Scope, name), d.sessions)
}

// GetLeaderLatch gets or creates a LeaderLatch with the given name
func (d *Database) GetLeaderLatch(ctx context.Context, name string, opts ...leader.Option) (leader.Latch, error) {
	return leader.New(ctx, primitive.NewName(d.Namespace, d.Name, d.Scope, name), d.sessions)
}

// GetList gets or creates a List with the given name
func (d *Database) GetList(ctx context.Context, name string) (list.List, error) {
	return list.New(ctx, primitive.NewName(d.Namespace, d.Name, d.Scope, name), d.sessions)
}

// GetLock gets or creates a Lock with the given name
func (d *Database) GetLock(ctx context.Context, name string) (lock.Lock, error) {
	return lock.New(ctx, primitive.NewName(d.Namespace, d.Name, d.Scope, name), d.sessions)
}

// GetLog gets or creates a Log with the given name
func (d *Database) GetLog(ctx context.Context, name string) (log.Log, error) {
	return log.New(ctx, primitive.NewName(d.Namespace, d.Name, d.Scope, name), d.sessions)
}

// GetMap gets or creates a Map with the given name
func (d *Database) GetMap(ctx context.Context, name string, opts ..._map.Option) (_map.Map, error) {
	return _map.New(ctx, primitive.NewName(d.Namespace, d.Name, d.Scope, name), d.sessions, opts...)
}

// GetSet gets or creates a Set with the given name
func (d *Database) GetSet(ctx context.Context, name string) (set.Set, error) {
	return set.New(ctx, primitive.NewName(d.Namespace, d.Name, d.Scope, name), d.sessions)
}

// GetValue gets or creates a Value with the given name
func (d *Database) GetValue(ctx context.Context, name string) (value.Value, error) {
	return value.New(ctx, primitive.NewName(d.Namespace, d.Name, d.Scope, name), d.sessions)
}

// Close closes the database
func (d *Database) Close(ctx context.Context) error {
	var returnErr error
	for _, session := range d.sessions {
		err := session.Close()
		if err != nil {
			returnErr = err
		}
	}
	return returnErr
}

var _ counter.Client = &Database{}
var _ election.Client = &Database{}
var _ indexedmap.Client = &Database{}
var _ leader.Client = &Database{}
var _ list.Client = &Database{}
var _ lock.Client = &Database{}
var _ log.Client = &Database{}
var _ _map.Client = &Database{}
var _ set.Client = &Database{}
var _ value.Client = &Database{}