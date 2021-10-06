/*
   Copyright 2021 The DbunderFS Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package mysql

import (
	"database/sql"
	"github.com/kos-v/dbunderfs/src/db"
	"strings"
	"sync"
)

type Instance struct {
	DSN DSN

	sync.RWMutex
	pool *sql.DB
}

func (inst *Instance) Close() error {
	if inst.pool == nil {
		return nil
	}
	return inst.pool.Close()
}

func (inst *Instance) Connect() (*sql.DB, error) {
	if inst.HasConnection() {
		return inst.pool, nil
	}

	pool, err := inst.generatePool(inst.DSN.ToString())
	if err != nil {
		return nil, err
	}

	inst.Lock()
	inst.pool = pool
	inst.Unlock()

	return inst.pool, nil
}

func (inst *Instance) GetDSN() db.DSN {
	return &inst.DSN
}

func (inst *Instance) Exec(query string, args ...interface{}) (sql.Result, error) {
	inst.Lock()
	defer inst.Unlock()

	return inst.pool.Exec(inst.prepareQuery(query), args...)
}

func (inst *Instance) GetDriverName() string {
	return "mysql"
}

func (inst *Instance) GetPool() *sql.DB {
	return inst.pool
}

func (inst *Instance) HasConnection() bool {
	if inst.pool != nil {
		if err := inst.pool.Ping(); err == nil {
			return true
		}
	}

	return false
}

func (inst *Instance) Reconnect() (*sql.DB, error) {
	pool, err := inst.generatePool(inst.DSN.ToString())
	if err != nil {
		return nil, err
	}

	inst.Lock()
	defer inst.Unlock()

	if inst.HasConnection() {
		if err = inst.Close(); err != nil {
			return nil, err
		}
	}
	inst.pool = pool

	return inst.pool, nil
}

func (inst *Instance) Query(query string, args ...interface{}) (*sql.Rows, error) {
	inst.Lock()
	defer inst.Unlock()

	return inst.pool.Query(inst.prepareQuery(query), args...)
}

func (inst *Instance) QueryRow(query string, args ...interface{}) *sql.Row {
	inst.Lock()
	defer inst.Unlock()

	return inst.pool.QueryRow(inst.prepareQuery(query), args...)
}

func (inst *Instance) generatePool(dsn string) (*sql.DB, error) {
	pool, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(); err != nil {
		return nil, err
	}

	pool.Query("SET max_allowed_packet = ?", 1024*1024*64)
	return pool, err
}

func (inst *Instance) prepareQuery(query string) string {
	// TODO: Move to a replace function
	return strings.Replace(query, "{%t_prefix%}", inst.DSN.GetTablePrefix(), -1)
}
