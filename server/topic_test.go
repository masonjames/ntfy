package server

import (
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"heckel.io/ntfy/v2/model"
)

func TestTopic_CancelSubscribersExceptUser(t *testing.T) {
	subFn := func(v *visitor, msg *model.Message) error {
		return nil
	}
	canceled1 := atomic.Bool{}
	cancelFn1 := func() {
		canceled1.Store(true)
	}
	canceled2 := atomic.Bool{}
	cancelFn2 := func() {
		canceled2.Store(true)
	}
	to := newTopic("mytopic")
	to.Subscribe(subFn, "", cancelFn1)
	to.Subscribe(subFn, "u_phil", cancelFn2)

	to.CancelSubscribersExceptUser("u_phil")
	require.True(t, canceled1.Load())
	require.False(t, canceled2.Load())
}

func TestTopic_CancelSubscribersUser(t *testing.T) {
	t.Parallel()

	subFn := func(v *visitor, msg *model.Message) error {
		return nil
	}
	canceled1 := atomic.Bool{}
	cancelFn1 := func() {
		canceled1.Store(true)
	}
	canceled2 := atomic.Bool{}
	cancelFn2 := func() {
		canceled2.Store(true)
	}
	canceled3 := atomic.Bool{}
	cancelFn3 := func() { canceled3.Store(true) }
	to := newTopic("mytopic")
	to.Subscribe(subFn, "u_another", cancelFn1)
	to.Subscribe(subFn, "u_phil", cancelFn2)
	to.Subscribe(subFn, "u_phil", cancelFn3)

	to.CancelSubscriberUser("u_phil")
	require.False(t, canceled1.Load())
	require.True(t, canceled2.Load())
	require.True(t, canceled3.Load())
}

func TestTopic_CancelSubscribersOutsideLock(t *testing.T) {
	for _, mode := range []string{"user", "except"} {
		t.Run(mode, func(t *testing.T) {
			to := newTopic("mytopic")
			canceled := false
			to.Subscribe(func(v *visitor, msg *model.Message) error { return nil }, "u_phil", func() {
				// Context logging and Unsubscribe may reenter the topic. TryLock
				// detects the old lock-held callback without hanging the test.
				unlocked := to.mu.TryLock()
				if unlocked {
					to.mu.Unlock()
				}
				require.True(t, unlocked)
				canceled = true
			})
			if mode == "user" {
				to.CancelSubscriberUser("u_phil")
			} else {
				to.CancelSubscribersExceptUser("u_other")
			}
			require.True(t, canceled)
		})
	}
}

func TestTopic_CancelSubscriberEmptyUser(t *testing.T) {
	to := newTopic("mytopic")
	canceled := false
	to.Subscribe(func(v *visitor, msg *model.Message) error { return nil }, "", func() { canceled = true })
	to.CancelSubscriberUser("")
	require.False(t, canceled)
}

func TestTopic_Keepalive(t *testing.T) {
	t.Parallel()

	to := newTopic("mytopic")
	to.lastAccess = time.Now().Add(-1 * time.Hour)
	to.Keepalive()
	require.True(t, to.LastAccess().Unix() >= time.Now().Unix()-2)
	require.True(t, to.LastAccess().Unix() <= time.Now().Unix()+2)
}

func TestTopic_Subscribe_DuplicateID(t *testing.T) {
	t.Parallel()
	to := newTopic("mytopic")

	//lint:ignore SA1019 Fix random seed to force same number generation
	rand.Seed(1)
	a := rand.Int()
	to.subscribers[a] = &topicSubscriber{
		userID:     "a",
		subscriber: nil,
		cancel:     func() {},
	}

	subFn := func(v *visitor, msg *model.Message) error {
		return nil
	}

	//lint:ignore SA1019 Force rand.Int to generate the same id once more
	rand.Seed(1)
	id := to.Subscribe(subFn, "b", func() {})
	res := to.subscribers[id]

	require.NotEqual(t, id, a)
	require.Equal(t, "b", res.userID, "b")
}
