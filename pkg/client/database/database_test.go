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
	"github.com/atomix/go-client/pkg/client/database/counter"
	"github.com/atomix/go-client/pkg/client/database/election"
	"github.com/atomix/go-client/pkg/client/database/indexedmap"
	"github.com/atomix/go-client/pkg/client/database/list"
	"github.com/atomix/go-client/pkg/client/database/lock"
	"github.com/atomix/go-client/pkg/client/database/log"
	"github.com/atomix/go-client/pkg/client/database/map"
	"github.com/atomix/go-client/pkg/client/database/set"
	"github.com/atomix/go-client/pkg/client/database/value"
	"github.com/atomix/go-client/pkg/client/protocol"
	"github.com/atomix/go-client/pkg/client/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDatabase(t *testing.T) {
	partitions, closers := test.StartTestPartitions(3)
	defer test.StopTestPartitions(closers)

	sessions, err := test.OpenSessions(partitions)
	assert.NoError(t, err)
	defer test.CloseSessions(sessions)

	database := &Database{
		Client:   protocol.New(nil, "default", "test", "default"),
		sessions: sessions,
	}

	primitives, err := database.GetPrimitives(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, primitives, 0)

	_, err = database.GetCounter(context.TODO(), "counter")
	assert.NoError(t, err)

	primitives, err = database.GetPrimitives(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, primitives, 1)

	primitives, err = database.GetPrimitives(context.TODO(), protocol.WithPrimitiveType(counter.Type))
	assert.NoError(t, err)
	assert.Len(t, primitives, 1)

	_, err = database.GetElection(context.TODO(), "election")
	assert.NoError(t, err)

	primitives, err = database.GetPrimitives(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, primitives, 2)

	primitives, err = database.GetPrimitives(context.TODO(), protocol.WithPrimitiveType(election.Type))
	assert.NoError(t, err)
	assert.Len(t, primitives, 1)

	_, err = database.GetIndexedMap(context.TODO(), "indexedmap")
	assert.NoError(t, err)

	primitives, err = database.GetPrimitives(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, primitives, 3)

	primitives, err = database.GetPrimitives(context.TODO(), protocol.WithPrimitiveType(indexedmap.Type))
	assert.NoError(t, err)
	assert.Len(t, primitives, 1)

	_, err = database.GetList(context.TODO(), "list")
	assert.NoError(t, err)

	primitives, err = database.GetPrimitives(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, primitives, 4)

	primitives, err = database.GetPrimitives(context.TODO(), protocol.WithPrimitiveType(list.Type))
	assert.NoError(t, err)
	assert.Len(t, primitives, 1)

	_, err = database.GetLock(context.TODO(), "lock")
	assert.NoError(t, err)

	primitives, err = database.GetPrimitives(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, primitives, 5)

	primitives, err = database.GetPrimitives(context.TODO(), protocol.WithPrimitiveType(lock.Type))
	assert.NoError(t, err)
	assert.Len(t, primitives, 1)

	_, err = database.GetLog(context.TODO(), "log")
	assert.NoError(t, err)

	primitives, err = database.GetPrimitives(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, primitives, 6)

	primitives, err = database.GetPrimitives(context.TODO(), protocol.WithPrimitiveType(log.Type))
	assert.NoError(t, err)
	assert.Len(t, primitives, 1)

	_, err = database.GetMap(context.TODO(), "map")
	assert.NoError(t, err)

	primitives, err = database.GetPrimitives(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, primitives, 7)

	primitives, err = database.GetPrimitives(context.TODO(), protocol.WithPrimitiveType(_map.Type))
	assert.NoError(t, err)
	assert.Len(t, primitives, 1)

	_, err = database.GetSet(context.TODO(), "set")
	assert.NoError(t, err)

	primitives, err = database.GetPrimitives(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, primitives, 8)

	primitives, err = database.GetPrimitives(context.TODO(), protocol.WithPrimitiveType(set.Type))
	assert.NoError(t, err)
	assert.Len(t, primitives, 1)

	_, err = database.GetValue(context.TODO(), "value")
	assert.NoError(t, err)

	primitives, err = database.GetPrimitives(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, primitives, 9)

	primitives, err = database.GetPrimitives(context.TODO(), protocol.WithPrimitiveType(value.Type))
	assert.NoError(t, err)
	assert.Len(t, primitives, 1)
}