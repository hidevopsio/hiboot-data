// Copyright 2018 John Deng (hi.devops.io@gmail.com).
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

package service

import (
	"hidevops.io/hiboot-data/examples/bolt/entity"
	"hidevops.io/hiboot-data/starter"
)

type FakeRepository struct {
	data.BaseKVRepository
}

func (r *FakeRepository) Get(params ...interface{}) error {
	if len(params) == 2 {
		key := params[0].(string)
		if key == "1" {
			u := params[1].(*entity.User)
			u.Name = "John Doe"
			u.Age = 18
		}
	}

	return nil
}
