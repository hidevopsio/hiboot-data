# hiboot-data

<p align="center">
  <img src="https://raw.githubusercontent.com/hidevopsio/hiboot-data/logo/hiboot-data.png?raw=true" alt="hiboot">
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
  <a href="https://goreportcard.com/report/hidevops.io/hiboot-data">
      <img src="https://goreportcard.com/badge/hidevops.io/hiboot-data" />
  </a>
  <a href="https://godoc.org/hidevops.io/hiboot-data">
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
calling the app.Register() function inside the init() function of your configuration pkg.

For more details, see https://godoc.org/hidevops.io/hiboot/pkg/starter

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
	"hidevops.io/hiboot-data/starter"
	"sync"
)

type Repository interface {
	data.KVRepository
}

type repository struct {
	data.BaseKVRepository
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
		return nil, nil, nil, data.InvalidDataSourceError
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
	return data.InvalidDataSourceError
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
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/hiboot/pkg/log"
)

type boltConfiguration struct {
	app.Configuration
	// the properties member name must be Bolt if the mapstructure is bolt,
	// so that the reference can be parsed
	BoltProperties properties `mapstructure:"bolt"`
}

func init() {
	app.Register(new(boltConfiguration))
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

Below is the example, for more details, please see [example](https://hidevops.io/hiboot-data/tree/master/examples/bolt)

```go


package service

import (
	"hidevops.io/hiboot-data/examples/bolt/entity"
	"hidevops.io/hiboot-data/starter/bolt"
	"hidevops.io/hiboot/pkg/app"
)

type UserService struct {
	repository bolt.Repository
}

func init() {
	app.Register(newUserService)
}

// will inject BoltRepository that configured in hidevops.io/hiboot-data/starter/bolt
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


## Development Guide

### Git workflow
Below, we outline one of the more common Git workflows that core developers use. Other Git workflows are also valid.

### Fork the main repository

* Go to https://github.com/hidevopsio/hiboot-data
* Click the "Fork" button (at the top right)

### Clone your fork

The commands below require that you have $GOPATH set ($GOPATH docs). We highly recommend you put Istio's code into your GOPATH. Note: the commands below will not work if there is more than one directory in your $GOPATH.

```bash
export GITHUB_USER=your-github-username
mkdir -p $GOPATH/src/github.com/hidevopsio
cd $GOPATH/src/github.com/hidevopsio
git clone https://github.com/$GITHUB_USER/hiboot-data
cd hiboot-data
git remote add upstream 'https://github.com/hidevopsio/hiboot-data'
git config --global --add http.followRedirects 1
```

### Create a branch and make changes

```bash
git checkout -b my-feature
# Then make your code changes
```

### Keeping your fork in sync

```bash
git fetch upstream
git rebase upstream/master
```

Note: If you have write access to the main repositories (e.g. github.com/hidevopsio/hiboot-data), you should modify your Git configuration so that you can't accidentally push to upstream:
