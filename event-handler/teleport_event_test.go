/*
Copyright 2015-2021 Gravitational, Inc.

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

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	auditlogpb "github.com/gravitational/teleport/api/gen/proto/go/teleport/auditlog/v1"
	"github.com/gravitational/teleport/api/types/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/gravitational/teleport-plugins/event-handler/lib"
)

func TestNew(t *testing.T) {
	e := &events.SessionPrint{
		Metadata: events.Metadata{
			ID:   "test",
			Type: "mock",
		},
	}

	event, err := NewTeleportEvent(eventToJSON(t, events.AuditEvent(e)), "cursor")
	require.NoError(t, err)
	assert.Equal(t, "test", event.ID)
	assert.Equal(t, "mock", event.Type)
	assert.Equal(t, "cursor", event.Cursor)
}

func TestGenID(t *testing.T) {
	e := &events.SessionPrint{}

	event, err := NewTeleportEvent(eventToJSON(t, events.AuditEvent(e)), "cursor")
	require.NoError(t, err)
	assert.NotEmpty(t, event.ID)
}

func TestSessionEnd(t *testing.T) {
	e := &events.SessionUpload{
		Metadata: events.Metadata{
			Type: "session.upload",
		},
		SessionMetadata: events.SessionMetadata{
			SessionID: "session",
		},
	}

	event, err := NewTeleportEvent(eventToJSON(t, events.AuditEvent(e)), "cursor")
	require.NoError(t, err)
	assert.NotEmpty(t, event.ID)
	assert.NotEmpty(t, event.SessionID)
	assert.True(t, event.IsSessionEnd)
}

func TestFailedLogin(t *testing.T) {
	e := &events.UserLogin{
		Metadata: events.Metadata{
			Type: "user.login",
		},
		Status: events.Status{
			Success: false,
		},
	}

	event, err := NewTeleportEvent(eventToJSON(t, events.AuditEvent(e)), "cursor")
	require.NoError(t, err)
	assert.NotEmpty(t, event.ID)
	assert.True(t, event.IsFailedLogin)
}

func TestSuccessLogin(t *testing.T) {
	e := &events.UserLogin{
		Metadata: events.Metadata{
			Type: "user.login",
		},
		Status: events.Status{
			Success: true,
		},
	}

	event, err := NewTeleportEvent(eventToJSON(t, events.AuditEvent(e)), "cursor")
	require.NoError(t, err)
	assert.NotEmpty(t, event.ID)
	assert.False(t, event.IsFailedLogin)
}

func eventToJSON(t *testing.T, e events.AuditEvent) *auditlogpb.EventUnstructured {
	data, err := lib.FastMarshal(e)
	require.NoError(t, err)
	str := &structpb.Struct{}
	err = str.UnmarshalJSON(data)
	require.NoError(t, err)
	id := e.GetID()
	if id == "" {
		hash := sha256.Sum256(data)
		id = hex.EncodeToString(hash[:])
	}
	return &auditlogpb.EventUnstructured{
		Type:         e.GetType(),
		Unstructured: str,
		Id:           id,
		Index:        e.GetIndex(),
		Time:         timestamppb.New(e.GetTime()),
	}
}
