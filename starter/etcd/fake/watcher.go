package fake

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
)

type Watcher struct {
	mock.Mock
}

// Watch watches on a key or prefix. The watched events will be returned
// through the returned channel. If revisions waiting to be sent over the
// watch are compacted, then the watch will be canceled by the server, the
// client will post a compacted error watch response, and the channel will close.
func (w *Watcher) Watch(ctx context.Context, key string, opts ...clientv3.OpOption) clientv3.WatchChan {
	args := w.Called(nil, key)
	return args[0].(clientv3.WatchChan)
}

// Close closes the watcher and cancels all watch requests.
func (w *Watcher) Close() error {
	return nil
}
