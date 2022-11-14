/**
 * Tencent is pleased to support the open source community by making Polaris available.
 *
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 *
 * Licensed under the BSD 3-Clause License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://opensource.org/licenses/BSD-3-Clause
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

package watcher

import (
	"github.com/fsnotify/fsnotify"

	"github.com/polaris-contrib/polaris-server-remote-plugin-common/log"
)

// OnModifiedEvent call this func when modified event fired.
type OnModifiedEvent func(name string)

// Watcher watch plugin binary change event, if changed, hot-reload plugin.
type Watcher struct {
	watcher      *fsnotify.Watcher
	eventChannel chan string
}

// New returns a new plugin watch.
func New(path string, callback OnModifiedEvent) *Watcher {
	var err error
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.DefaultLogger.Fatal("failed to add plugin watcher", "error", err)
	}

	pw := &Watcher{
		watcher:      watcher,
		eventChannel: make(chan string, 10),
	}

	go func(pw *Watcher) {
		for {
			select {
			case ev := <-watcher.Events:
				if ev.Op == fsnotify.Create {
					pw.eventChannel <- ev.Name
				}
				log.Info("got plugin file watch event", "event", ev)
			case err = <-watcher.Errors:
				log.Error("got error from plugin file watcher", "error_event", err)
			}
		}
	}(pw)

	go func() {
		for event := range pw.eventChannel {
			callback(event)
		}
	}()

	err = pw.watcher.Add(path)
	if err != nil {
		log.Fatal("failed to add plugin file watcher", "error", err)
	}

	return pw
}
