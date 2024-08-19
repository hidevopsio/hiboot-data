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
	"context"
	"encoding/json"
	"errors"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/idgen"
	"hiboot-data/examples/sqlx/entity"
	"hiboot-data/starter/redis"
	"hiboot-data/starter/sqlx"
	"strconv"
)

type UserService struct {
	at.Scope `value:"request"`
	db       *sqlx.DB
	cache    *redis.Client
}

func init() {
	// register UserService
	app.Register(newUserService)
}

// will inject sqlx.Repository that configured in hiboot-data/starter/gorm
func newUserService(db *sqlx.DB, redisClient *redis.Client) *UserService {

	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,  -- id: Unsigned 64-bit integer with auto-increment
		name VARCHAR(255) NOT NULL,                     -- name: Non-nullable string
		username VARCHAR(255) NOT NULL UNIQUE,          -- username: Non-nullable string, unique
		password VARCHAR(255) NOT NULL,                 -- password: Non-nullable string
		email VARCHAR(255) NOT NULL UNIQUE,             -- email: Non-nullable string, unique
		age TINYINT UNSIGNED CHECK (age >= 0 AND age <= 130),  -- age: Unsigned integer, limited between 0 and 130
		gender TINYINT UNSIGNED CHECK (gender >= 0 AND gender <= 2)  -- gender: Unsigned integer, limited between 0 and 2
);

`
	// Execute the schema
	_, err := db.Exec(schema)
	if err != nil {
		log.Fatalf("failed to create table for schema: %v ", schema)
	}

	return &UserService{
		db:    db,
		cache: redisClient,
	}
}

func (s *UserService) AddUser(user *entity.User) error {
	if user == nil {
		return errors.New("user is not allowed to be nil")
	}
	if user.Id == 0 {
		user.Id, _ = idgen.Next()
	}

	// SQLX insert query
	query := `INSERT INTO users (id, name, username, password, email, age, gender) 
			  VALUES (:id, :name, :username, :password, :email, :age, :gender)`

	// NamedExec binds named parameters from a struct to the query
	_, err := s.db.NamedExec(query, user)
	if err != nil {
		return err
	}

	// Cache the user
	return s.cacheUser(user)
}

func (s *UserService) cacheUser(user *entity.User) (err error) {
	var userJSON []byte
	userJSON, err = json.Marshal(user)
	if err != nil {
		return
	}
	err = s.cache.Set(context.Background(), strconv.FormatUint(user.Id, 16), string(userJSON), 0).Err()
	return
}

func (s *UserService) GetUser(id uint64) (*entity.User, error) {
	user := &entity.User{}

	// Attempt to get the user from the cache
	res, err := s.cache.Get(context.Background(), strconv.FormatUint(id, 16)).Result()
	if err == nil {
		err = json.Unmarshal([]byte(res), user)
		if err == nil {
			return user, nil
		}
	}

	// If the cache lookup failed or unmarshalling failed, query the database
	query := `SELECT * FROM users WHERE id = ?`
	err = s.db.Get(user, query, id)
	if err != nil {
		return nil, err
	}

	// Cache the user
	if err := s.cacheUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetAll() ([]entity.User, error) {
	users := []entity.User{}

	// SQLX select query to fetch all users
	query := `SELECT * FROM users`

	// Execute the query and map the results to the users slice
	err := s.db.Select(&users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) DeleteUser(id uint64) error {
	// SQLX delete query
	query := `DELETE FROM users WHERE id = ?`

	// Execute the delete query
	_, err := s.db.Exec(query, id)
	return err
}
