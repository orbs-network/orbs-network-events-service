package storage

import (
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_matchMultiple(t *testing.T) {
	arg0, _ := protocol.ArgumentArrayFromNatives([]interface{}{"Raising Arizona"})
	arg1, _ := protocol.ArgumentArrayFromNatives([]interface{}{uint32(1987)})
	arg2, _ := protocol.ArgumentArrayFromNatives([]interface{}{"Nicolas Cage"})

	var singleArgArray []*protocol.ArgumentArray
	singleArgArray = append(singleArgArray, arg0)

	require.True(t, matchEvent(ARIZONA_EVENT, nil))
	require.True(t, matchEvent(VAMPIRE_EVENT, nil))
	require.True(t, matchEvent(MASHUP_EVENT, nil))

	require.True(t, matchEvent(ARIZONA_EVENT, singleArgArray))
	require.False(t, matchEvent(VAMPIRE_EVENT, singleArgArray))
	require.False(t, matchEvent(MASHUP_EVENT, singleArgArray))

	var multiArgArray []*protocol.ArgumentArray
	multiArgArray = append(multiArgArray, arg0, arg1, arg2)

	require.True(t, matchEvent(ARIZONA_EVENT, multiArgArray))
	require.False(t, matchEvent(VAMPIRE_EVENT, multiArgArray))
	require.False(t, matchEvent(MASHUP_EVENT, multiArgArray))
}
