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

package pubsub

import (
	"github.com/emitter-io/emitter/internal/event"
	"github.com/emitter-io/emitter/internal/message"
	"github.com/emitter-io/emitter/internal/provider/storage"
	"github.com/emitter-io/emitter/internal/service"
)

// Notifier represents a cluster-wide notifier.
type notifier interface {
	NotifySubscribe(message.Subscriber, *event.Subscription)
	NotifyUnsubscribe(message.Subscriber, *event.Subscription)
}

// Service represents a publish service.
type Service struct {
	auth     service.Authorizer         // The authorizer to use.
	store    storage.Storage            // The storage provider to use.
	notifier notifier                   // The notifier to use.
	trie     *message.Trie              // The subscription matching trie.
	handlers map[uint32]service.Handler // The emitter request handlers.
}

// New creates a new publisher service.
func New(auth service.Authorizer, store storage.Storage, notifier notifier, trie *message.Trie, handlers []service.Handler) *Service {
	s := &Service{
		auth:     auth,
		store:    store,
		notifier: notifier,
		trie:     trie,
		handlers: make(map[uint32]service.Handler, len(handlers)),
	}

	for _, h := range handlers {
		s.handlers[h.Type()] = h
	}
	return s
}