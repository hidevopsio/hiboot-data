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

package gorm

import (
	"hidevops.io/hiboot/pkg/at"
)

type Config struct {
	Decrypt    bool   `json:"decrypt" default:"true"`
	DecryptKey string `json:"decrypt_key"`
}

type Properties struct {
	// annotation ConfigurationProperties
	at.ConfigurationProperties `value:"gorm"`

	Type            string `json:"type" default:"mysql"` // mysql, postgres, sqlite3, mssql,
	Host            string `json:"host" default:"mysql-dev"`
	Port            string `json:"port" default:"3306"`
	Database        string `json:"database"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Charset         string `json:"charset" default:"utf8"`
	ParseTime       bool   `json:"parse_time" default:"true"`
	Loc             string `json:"loc" default:"Asia/Shanghai"`
	Config          Config `json:"config"`
	ConnMaxLifetime string `json:"connMaxLifetime" default:"60s"`
	MaxIdleConns    int    `json:"maxIdle_conns" default:"20"`
	MaxOpenConns    int    `json:"maxOpen_conns" default:"200"`
	AutoReconnect   bool   `json:"auto_reconnect" default:"true"`
	RetryTimes      int    `json:"retry_times" default:"-1"`
	Interval        string `json:"interval" default:"3s"`
	NowRetryTimes   int    `json:"now_retry_times" default:"0"`
}
