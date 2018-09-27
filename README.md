# hiboot-data

<p align="center">
  <img src="https://github.com/hidevopsio/hiboot-data/blob/master/hiboot-data.png?raw=true" alt="hiboot">
</p>

<p align="center">
  <a href="https://travis-ci.org/hidevopsio/hiboot-data?branch=master">
    <img src="https://travis-ci.org/hidevopsio/hiboot-data.svg?branch=master" alt="Build Status"/>
  </a>
  <a href="https://codecov.io/gh/hidevopsio/hiboot-data">
    <img src="https://codecov.io/gh/hidevopsio/hiboot-data/branch/master/graph/badge.svg" />
  </a>
  <a href="https://opensource.org/licenses/Apache-2.0">
      <img src="https://img.shields.io/badge/License-Apache%202.0-green.svg" />
  </a>
  <a href="https://goreportcard.com/report/github.com/hidevopsio/hiboot-data">
      <img src="https://goreportcard.com/badge/github.com/hidevopsio/hiboot-data" />
  </a>
  <a href="https://godoc.org/github.com/hidevopsio/hiboot-data">
      <img src="https://godoc.org/github.com/golang/gddo?status.svg" />
  </a>
</p>

hiboot-data is the collection of Hiboot data starter, include bolt, etcd, gorm

* bolt - bolt database starter
* etcd - etcd key value store starter
* gorm - gorm database starter orm starter, support mysql, postgres, mssql, sqlite

## Auto-configured Starter

Hiboot auto-configuration attempts to automatically configure your Hiboot application based on the pkg dependencies that you have added.
For example, if bolt is imported in you main.go, and you have not manually configured any database connection,
then Hiboot auto-configures an database bolt for any service to inject.

You need to opt-in to auto-configuration by embedding app.Configuration in your configuration and
calling the app.AutoConfiguration() function inside the init() function of your configuration pkg.

For more details, see https://godoc.org/github.com/hidevopsio/hiboot/pkg/starter

## Creating Your Own Starter

A full Hiboot starter for a library may contain the following structs:
	autoconfigure - object that handle the auto-configuration code.
	properties - object that contains properties which will be injected configurable default values or user specified values
If you work in a company that develops shared go packages, or if you work on an open-source or commercial project,
you might want to develop your own auto-configured starter. starter can be implemented in external packages and
can be imported by any go applications.

## Understanding Auto-configured Starter

Under the hood, auto-configuration is implemented with standard struct. Additional embedded field app.Configuration.
AutoConfiguration used to constrain when the auto-configuration should apply. Usually, auto-configuration struct use
`after:"fooConfiguration"` or `missing:"fooConfiguration"` tags. This ensures that auto-configuration applies only
when relevant configuration are found and when you have not declared your own configuration.

## Example

For example, if you want to make bolt starter,

### Properties

First, define Properties for injecting external configurations

```go

// properties
type properties struct {
	Database string      `json:"database" default:"hiboot.db"`
	Mode     os.FileMode `json:"mode" default:"0600"`
	Timeout  int64       `json:"timeout" default:"2"`
}

```

### Repository

Implementing BoltRepository to wrap up bolt client APIs

```go


package bolt

import (
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"github.com/hidevopsio/hiboot-data/starter"
	"sync"
)

type Repository interface {
	starter.KVRepository
}

type repository struct {
	starter.BaseKVRepository
	db         *bolt.DB
	dataSource DataSource
}

var repo *repository
var once sync.Once
var InvalidPropertiesError = errors.New("properties must not be nil")

func GetRepository() *repository {
	once.Do(func() {
		repo = &repository{}
	})
	return repo
}

func (r *repository) parse(params ...interface{}) ([]byte, []byte, interface{}, error) {
	if r.db == nil {
		return nil, nil, nil, starter.InvalidDataSourceError
	}
	return r.Parse(params...)
}

// Open bolt database
func (r *repository) SetDataSource(d interface{}) {
	if d != nil {
		r.dataSource = d.(DataSource)
		r.db = r.dataSource.DB()
	}
}

func (r *repository) DataSource() interface{} {
	return r.dataSource
}

func (r *repository) CloseDataSource() error {
	if r.dataSource != nil {
		return r.dataSource.Close()
	}
	return starter.InvalidDataSourceError
}

// Put inserts a key:value pair into the database
func (r *repository) Put(params ...interface{}) error {
	bucketName, key, value, err := r.parse(params...)
	if err != nil {
		return err
	}

	err = r.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		// marshal data to bytes
		b, err := json.Marshal(value)
		if err != nil {
			return err
		}

		err = bucket.Put(key, b)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

// Get retrieves a key:value pair from the database
func (r *repository) Get(params ...interface{}) error {
	bucketName, key, value, err := r.parse(params...)
	if err != nil {
		return err
	}
	var result []byte
	err = r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b != nil {
			v := b.Get(key)
			if v != nil {
				result = make([]byte, len(v))
				copy(result, v)
			}
		} else {
			result = []byte("")
		}
		return nil
	})
	// TODO: if result len is 0, return errors.New("no record found")
	if err == nil {
		err = json.Unmarshal(result, value)
	}
	return err
}

// Delete removes a key:value pair from the database
func (r *repository) Delete(params ...interface{}) error {
	bucketName, key, _, err := r.parse(params...)
	if err != nil {
		return err
	}
	err = r.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		err = bucket.Delete(key)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

```

### Auto Configuration

Auto Configuration is responsible for creating injectable instances that can be injected in your services directory through constructor function.

```go
package bolt

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

import (
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type boltConfiguration struct {
	app.Configuration
	// the properties member name must be Bolt if the mapstructure is bolt,
	// so that the reference can be parsed
	BoltProperties properties `mapstructure:"bolt"`
}

func init() {
	app.AutoConfiguration(new(boltConfiguration))
}

func (c *boltConfiguration) dataSource() DataSource {
	dataSource := GetDataSource()
	if !dataSource.IsOpened() {
		err := dataSource.Open(&c.BoltProperties)
		if err != nil {
			log.Error(err.Error())
		}
	}
	return dataSource
}

func (c *boltConfiguration) BoltRepository() Repository {
	repository := GetRepository()
	repository.SetDataSource(c.dataSource())
	return repository
}

```

### Usage

After bolt starter is built, you can inject it directly in your application.

Below is the example, for more details, please see [example](https://github.com/hidevopsio/hiboot-data/tree/master/examples/bolt)

```go


package service

import (
	"github.com/hidevopsio/hiboot-data/examples/bolt/entity"
	"github.com/hidevopsio/hiboot-data/starter/bolt"
	"github.com/hidevopsio/hiboot/pkg/app"
)

type UserService struct {
	repository bolt.Repository
}

func init() {
	app.Component(newUserService)
}

// will inject BoltRepository that configured in github.com/hidevopsio/hiboot-data/starter/bolt
func newUserService(repository bolt.Repository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) AddUser(user *entity.User) error {
	return s.repository.Put(user)
}

func (s *UserService) GetUser(id string) (*entity.User, error) {
	var user entity.User
	err := s.repository.Get(id, &user)
	return &user, err
}

func (s *UserService) DeleteUser(id string) error {
	return s.repository.Delete(id, &entity.User{})
}

```