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

package sqlx

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // MySQL 驱动
	"github.com/hidevopsio/hiboot-data/utils"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // SQLite 驱动
)

const Profile = "sqlx"

type DB struct {
	at.Scope `value:"prototype"`
	*sqlx.DB
}

type configuration struct {
	at.AutoConfiguration

	prop *properties
	db   *DB
}

func newConfiguration(prop *properties) *configuration {
	return &configuration{prop: prop}
}

func init() {
	app.Register(newConfiguration, new(properties))
}

func (c *configuration) DB() (db *DB, err error) {
	var sqlDB *sqlx.DB
	// return cached if it is healthy
	if c.db != nil {
		sqlDB = c.db.DB
		if err == nil && sqlDB.Ping() == nil {
			// If the connection is still alive, return the existing connection
			db = c.db
			return
		}
		log.Warnf("lost connection to database, attempting to reconnect...")
		c.db = nil
	}

	db = new(DB)

	var dsn string
	var report string
	if c.prop.Type == "sqlite3" {
		dsn = c.prop.Database
		report = fmt.Sprintf("database %v %v", c.prop.Type, c.prop.Database)
	} else {
		password := c.prop.Password
		if c.prop.Config.Decrypt {
			password = utils.Decrypt(password, c.prop.Config.DecryptKey)
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			c.prop.Username, password, c.prop.Host, c.prop.Port, c.prop.Database)
		report = fmt.Sprintf("database %v %v %v@%v:%v", c.prop.Type, c.prop.Database, c.prop.Username, c.prop.Host, c.prop.Port)
	}

	// create new connection if it is unhealthy
	log.Infof("create a new database connection to %v", report)

	db.DB, err = sqlx.Connect(c.prop.Type, dsn)
	if err != nil {
		log.Errorf("failed to connect db: %v", err)
		return
	}

	if sqlDB.Ping() == nil {
		// If the connection is alive, assign the new connection
		c.db = db
		log.Infof("%v is connected", report)
	}

	return
}
