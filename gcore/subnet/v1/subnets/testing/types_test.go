package testing

import (
	"encoding/json"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/subnet/v1/subnets"

	"github.com/stretchr/testify/require"
)

func TestMarshallCreateStructure(t *testing.T) {
	options := subnets.CreateOpts{
		Name:                   Subnet1.Name,
		EnableDHCP:             true,
		CIDR:                   Subnet1.CIDR,
		NetworkID:              Subnet1.NetworkID,
		ConnectToNetworkRouter: true,
	}

	mp, err := options.ToSubnetCreateMap()
	require.NoError(t, err)
	s, err := json.Marshal(mp)
	require.NoError(t, err)
	require.JSONEq(t, CreateRequest, string(s))

}
