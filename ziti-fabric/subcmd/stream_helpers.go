/*
	Copyright NetFoundry, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package subcmd

import (
	"github.com/openziti/foundation/channel2"
	"github.com/michaelquigley/pfxlog"
	"sync"
)

func waitForChannelClose(ch channel2.Channel) {
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)

	ch.AddCloseHandler(&closeWatcher{waitGroup})

	waitGroup.Wait()
}

type closeWatcher struct {
	waitGroup *sync.WaitGroup
}

func (watcher *closeWatcher) HandleClose(ch channel2.Channel) {
	pfxlog.Logger().Info("Management channel to controller closed. Shutting down.")
	watcher.waitGroup.Done()
}
