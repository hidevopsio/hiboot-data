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

package redis

import (
	"github.com/hidevopsio/hiboot/pkg/at"
)

type Config struct {
	Decrypt    bool   `json:"decrypt" default:"true"`
	DecryptKey string `json:"decrypt_key"`
}

type properties struct {
	at.ConfigurationProperties `value:"redis"`

	Host     string `json:"host" default:"redis-master"`
	Port     string `json:"port" default:"6379"`
	Password string `json:"password"`
	DB       int    `json:"client" default:"0"`
	Config   Config `json:"config"`
	Timeout  int    `json:"timeout" default:"2"`
}
