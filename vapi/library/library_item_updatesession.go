/*
Copyright (c) 2018 VMware, Inc. All Rights Reserved.

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

package library

import (
	"context"
	"net/http"
	"time"

	"github.com/vmware/govmomi/vapi/internal"
)

// UpdateSession is used to create an initial update session
type UpdateSession struct {
	ID                        string `json:"id,omitempty"`
	LibraryItemID             string `json:"library_item_id,omitempty"`
	LibraryItemContentVersion string `json:"library_item_content_version,omitempty"`
	// ErrorMessage              struct {
	//	ID             string   `json:"id,omitempty"`
	//	DefaultMessage string   `json:"default_message,omitempty"`
	//	Args           []string `json:"args,omitempty"`
	// } `json:"error_message,omitempty"`
	ClientProgress int64      `json:"client_progress,omitempty"`
	State          string     `json:"state,omitempty"`
	ExpirationTime *time.Time `json:"expiration_time,omitempty"`
}

// CreateLibraryItemUpdateSession creates a new library item
func (c *Manager) CreateLibraryItemUpdateSession(ctx context.Context, session UpdateSession) (string, error) {
	url := internal.URL(c, internal.LibraryItemUpdateSession)
	spec := struct {
		CreateSpec UpdateSession `json:"create_spec"`
	}{session}
	var res string
	return res, c.Do(ctx, url.Request(http.MethodPost, spec), &res)
}

// GetLibraryItemUpdateSession gets the update session information with status
func (c *Manager) GetLibraryItemUpdateSession(ctx context.Context, id string) (*UpdateSession, error) {
	url := internal.URL(c, internal.LibraryItemUpdateSession).WithID(id)
	var res UpdateSession
	return &res, c.Do(ctx, url.Request(http.MethodGet), &res)
}

// ListLibraryItemUpdateSession gets the list of update sessions
func (c *Manager) ListLibraryItemUpdateSession(ctx context.Context) (*[]string, error) {
	url := internal.URL(c, internal.LibraryItemUpdateSession)
	var res []string
	return &res, c.Do(ctx, url.Request(http.MethodGet), &res)
}

// CancelLibraryItemUpdateSession cancels an update session
func (c *Manager) CancelLibraryItemUpdateSession(ctx context.Context, id string) error {
	url := internal.URL(c, internal.LibraryItemUpdateSession).WithID(id).WithAction("cancel")
	return c.Do(ctx, url.Request(http.MethodPost), nil)
}

// CompleteLibraryItemUpdateSession completes an update session
func (c *Manager) CompleteLibraryItemUpdateSession(ctx context.Context, id string) error {
	url := internal.URL(c, internal.LibraryItemUpdateSession).WithID(id).WithAction("complete")
	return c.Do(ctx, url.Request(http.MethodPost), nil)
}

// DeleteLibraryItemUpdateSession deletes an update session
func (c *Manager) DeleteLibraryItemUpdateSession(ctx context.Context, id string) error {
	url := internal.URL(c, internal.LibraryItemUpdateSession).WithID(id)
	return c.Do(ctx, url.Request(http.MethodDelete), nil)
}

// FailLibraryItemUpdateSession fails an update session
func (c *Manager) FailLibraryItemUpdateSession(ctx context.Context, id string) error {
	url := internal.URL(c, internal.LibraryItemUpdateSession).WithID(id).WithAction("fail")
	return c.Do(ctx, url.Request(http.MethodPost), nil)
}

// KeepAliveLibraryItemUpdateSession keeps an inactive update session alive.
func (c *Manager) KeepAliveLibraryItemUpdateSession(ctx context.Context, id string) error {
	url := internal.URL(c, internal.LibraryItemUpdateSession).WithID(id).WithAction("keep-alive")
	return c.Do(ctx, url.Request(http.MethodPost), nil)
}

// WaitOnLibraryItemUpdateSession blocks until the update session is no longer
// in the ACTIVE state.
func (c *Manager) WaitOnLibraryItemUpdateSession(
	ctx context.Context, sessionID string,
	interval time.Duration, intervalCallback func()) error {

	// Wait until the upload operation is complete to return.
	for {
		session, err := c.GetLibraryItemUpdateSession(ctx, sessionID)
		if err != nil {
			return err
		}
		if session.State != "ACTIVE" {
			return nil
		}
		time.Sleep(interval)
		if intervalCallback != nil {
			intervalCallback()
		}
	}
}
