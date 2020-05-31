/**********************************************************************************
* Copyright (c) 2009-2020 Misakai Ltd.
* This program is free software: you can redistribute it and/or modify it under the
* terms of the GNU Affero General Public License as published by the  Free Software
* Foundation, either version 3 of the License, or(at your option) any later version.
*
* This program is distributed  in the hope that it  will be useful, but WITHOUT ANY
* WARRANTY;  without even  the implied warranty of MERCHANTABILITY or FITNESS FOR A
* PARTICULAR PURPOSE.  See the GNU Affero General Public License  for  more details.
*
* You should have  received a copy  of the  GNU Affero General Public License along
* with this program. If not, see<http://www.gnu.org/licenses/>.
************************************************************************************/

package presence

import (
	"encoding/json"
	"time"

	"github.com/emitter-io/emitter/internal/event"
	"github.com/emitter-io/emitter/internal/message"
	"github.com/emitter-io/emitter/internal/provider/logging"
)

// Request represents a presence request
type Request struct {
	Key     string `json:"key"`     // The channel key for this request.
	Channel string `json:"channel"` // The target channel for this request.
	Status  bool   `json:"status"`  // Specifies that a status response should be sent.
	Changes *bool  `json:"changes"` // Specifies that the changes should be notified.
}

// EventType represents a presence event type
type EventType string

// Various event types
const (
	EventTypeStatus      = EventType("status")
	EventTypeSubscribe   = EventType("subscribe")
	EventTypeUnsubscribe = EventType("unsubscribe")
)

// ------------------------------------------------------------------------------------

// Response represents a state notification.
type Response struct {
	Request uint16    `json:"req,omitempty"` // The corresponding request ID.
	Time    int64     `json:"time"`          // The UNIX timestamp.
	Event   EventType `json:"event"`         // The event, must be "status", "subscribe" or "unsubscribe".
	Channel string    `json:"channel"`       // The target channel for the notification.
	Who     []Info    `json:"who"`           // The subscriber ids.
}

// ForRequest sets the request ID in the response for matching
func (r *Response) ForRequest(id uint16) {
	r.Request = id
}

// ------------------------------------------------------------------------------------

// Info represents a presence info for a single connection.
type Info struct {
	ID       string `json:"id"`                 // The subscriber ID.
	Username string `json:"username,omitempty"` // The subscriber username set by client ID.
}

// ------------------------------------------------------------------------------------

// Notification represents a state notification.
type Notification struct {
	Time    int64        `json:"time"`    // The UNIX timestamp.
	Event   EventType    `json:"event"`   // The event, must be "status", "subscribe" or "unsubscribe".
	Channel string       `json:"channel"` // The target channel for the notification.
	Who     Info         `json:"who"`     // The subscriber id.
	Ssid    message.Ssid `json:"-"`       // The ssid to dispatch the notification on.
}

// NewNotification creates a new notification payload.
func NewNotification(event EventType, ev *event.Subscription) *Notification {
	return &Notification{
		Ssid:    message.NewSsidForPresence(ev.Ssid),
		Time:    time.Now().UTC().Unix(),
		Event:   event,
		Channel: string(ev.Channel),
		Who: Info{
			ID:       ev.ConnID(),
			Username: string(ev.User),
		},
	}
}

// Encode encodes the presence notifications and returns a payload to send.
func (e *Notification) Encode() ([]byte, bool) {
	encoded, err := json.Marshal(e)
	if err != nil {
		logging.LogError("presence", "encoding presence notification", err)
		return nil, false
	}

	return encoded, true
}