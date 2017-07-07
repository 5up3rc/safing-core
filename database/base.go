// Copyright Safing ICS Technologies GmbH. Use of this source code is governed by the AGPL license that can be found in the LICENSE file.

package database

import (
	"errors"
	"safing/database/dbutils"
	"strings"

	"github.com/ipfs/go-datastore"
	uuid "github.com/satori/go.uuid"
)

type Base struct {
	dbKey *datastore.Key
	meta  *dbutils.Meta
}

func (m *Base) SetKey(key *datastore.Key) {
	m.dbKey = key
}

func (m *Base) GetKey() *datastore.Key {
	return m.dbKey
}

func (m *Base) FmtKey() string {
	return m.dbKey.String()
}

func (m *Base) Meta() *dbutils.Meta {
	return m.meta
}

func (m *Base) CreateObject(namespace *datastore.Key, name string, model Model) error {
	var newKey datastore.Key
	if name == "" {
		newKey = namespace.ChildString(getTypeName(model)).Instance(strings.Replace(uuid.NewV4().String(), "-", "", -1))
	} else {
		newKey = namespace.ChildString(getTypeName(model)).Instance(name)
	}
	m.dbKey = &newKey
	return Create(*m.dbKey, model)
}

func (m *Base) SaveObject(model Model) error {
	if m.dbKey == nil {
		return errors.New("cannot save new object, use Create() instead")
	}
	return Update(*m.dbKey, model)
}

func (m *Base) Delete() error {
	if m.dbKey == nil {
		return errors.New("cannot delete object unsaved object")
	}
	return Delete(*m.dbKey)
}
