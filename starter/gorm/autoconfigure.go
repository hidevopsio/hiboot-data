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
	"database/sql"
	"fmt"
	"github.com/hidevopsio/hiboot-data/utils"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
)

const Profile = "gorm"

type DB struct {
	at.Scope `value:"prototype"`
	*gorm.DB
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
	var sqlDB *sql.DB
	// return cached if it is healthy
	if c.db != nil {
		sqlDB, err = c.db.DB.DB()
		if err == nil && sqlDB.Ping() == nil {
			// If the connection is still alive, return the existing connection
			db = c.db
			return
		}
		log.Warnf("lost connection to database, attempting to reconnect...")
		c.db = nil
	}

	// create new connection if it is unhealthy
	var report string
	report = fmt.Sprintf("database %v@%v:%v", c.prop.Username, c.prop.Host, c.prop.Port)
	log.Infof("create new connection to %v", report)
	db = new(DB)
	password := c.prop.Password
	if c.prop.Config.Decrypt {
		password = utils.Decrypt(password, c.prop.Config.DecryptKey)
	}
	loc := strings.Replace(c.prop.Loc, "/", "%2F", -1)
	databaseName := strings.Replace(c.prop.Database, "-", "_", -1)

	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=%v&loc=%v",
		c.prop.Username,
		password,
		c.prop.Host,
		c.prop.Port,
		databaseName,
		c.prop.Charset,
		c.prop.ParseTime,
		loc,
	)
	db.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Errorf("failed to connect %v, err: %v", report, err)
		return
	}

	if sqlDB.Ping() == nil {
		// If the connection is alive, assign the new connection
		c.db = db
		log.Infof("%v is connected", report)
	}
	return
}
