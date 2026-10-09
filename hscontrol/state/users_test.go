package state

import (
	"testing"

	"github.com/juanfont/headscale/hscontrol/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenameUserUpdatesNodeStore(t *testing.T) {
	s := newRetagTestState(t)

	user := s.CreateUserForTest("old-name")
	otherUser := s.CreateUserForTest("other")

	node := s.CreateRegisteredNodeForTest(user, "user-node")
	node.User = user
	s.PutNodeInStoreForTest(*node)

	// A tagged node keeps the user that created it.
	tagged := s.CreateRegisteredNodeForTest(user, "tagged-node")
	tagged.User = user
	tagged.Tags = []string{"tag:server"}
	s.PutNodeInStoreForTest(*tagged)

	other := s.CreateRegisteredNodeForTest(otherUser, "other-node")
	other.User = otherUser
	s.PutNodeInStoreForTest(*other)

	_, c, err := s.RenameUser(types.UserID(user.ID), "new-name")
	require.NoError(t, err)
	assert.ElementsMatch(t, []types.NodeID{node.ID, tagged.ID}, c.PeersChanged)

	for _, id := range []types.NodeID{node.ID, tagged.ID} {
		nv, ok := s.GetNodeByID(id)
		require.True(t, ok)
		assert.Equal(t, "new-name", nv.User().Name())
	}

	nv, ok := s.GetNodeByID(other.ID)
	require.True(t, ok)
	assert.Equal(t, "other", nv.User().Name())

	_, c, err = s.RenameUser(types.UserID(user.ID), "new-name")
	require.NoError(t, err)
	assert.Empty(t, c.PeersChanged)
}
