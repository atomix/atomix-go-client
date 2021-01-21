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

package value

import (
	api "github.com/atomix/api/go/atomix/primitive/value"
	"github.com/atomix/go-framework/pkg/atomix/meta"
	"github.com/google/uuid"
)

// Option is a value option
type Option interface {
	apply(options *options)
}

// options is value options
type options struct {
	clientID string
}

func applyOptions(opts ...Option) options {
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	options := &options{
		clientID: id.String(),
	}
	for _, opt := range opts {
		opt.apply(options)
	}
	return *options
}

// WithClientID sets the client identifier
func WithClientID(id string) Option {
	return &clientIDOption{
		clientID: id,
	}
}

type clientIDOption struct {
	clientID string
}

func (o *clientIDOption) apply(options *options) {
	options.clientID = o.clientID
}

// SetOption is an option for Set calls
type SetOption interface {
	beforeSet(request *api.SetRequest)
	afterSet(response *api.SetResponse)
}

// IfMatch updates the value if the version matches the given version
func IfMatch(object meta.Object) SetOption {
	return matchOption{object}
}

type matchOption struct {
	object meta.Object
}

func (o matchOption) beforeSet(request *api.SetRequest) {
	request.Value.ObjectMeta = o.object.Meta().Proto()
}

func (o matchOption) afterSet(response *api.SetResponse) {

}
