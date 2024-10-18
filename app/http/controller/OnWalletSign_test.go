package controller

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fil-kms/app/utils"
	"fmt"
	"net/http"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/builtin/v10/market"
	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/assert"
)

func TestOnWalletSign(t *testing.T) {
	aks := "xxxxx"
	client, err := address.NewFromString("f3tginto55sulyzzjoy3byof7u7t5vicwa372z4hn4zdhjkdwdkhalsn6mytk53paoc63uf575rfdvqetyvpka")
	assert.NoError(t, err)
	provider, err := address.NewFromString("t01000")
	assert.NoError(t, err)
	provider2, err := address.NewFromString("t01001")
	assert.NoError(t, err)
	label, err := market.NewLabelFromString("test")
	assert.NoError(t, err)
	pieceCID, err := cid.Decode("baga6ea4seaqihyevzzuoyacvl5umwtksq5tgq3f2er4m4kqg4dndevymj3vdufi")
	assert.NoError(t, err)
	assert.NoError(t, err)
	props := []*market.DealProposal{
		{
			PieceCID:             pieceCID,
			PieceSize:            100000,
			Client:               client,
			Provider:             provider,
			Label:                label,
			StoragePricePerEpoch: abi.NewTokenAmount(0),
			ProviderCollateral:   abi.NewTokenAmount(0),
			ClientCollateral:     abi.NewTokenAmount(0),
		},
		{
			PieceCID:             pieceCID,
			PieceSize:            1024 * 1024 * 1024 * 1024 * 1024 * 2,
			Client:               client,
			Provider:             provider,
			Label:                label,
			StoragePricePerEpoch: abi.NewTokenAmount(0),
			ProviderCollateral:   abi.NewTokenAmount(0),
			ClientCollateral:     abi.NewTokenAmount(0),
		},
		{
			PieceCID:             pieceCID,
			PieceSize:            100000,
			Client:               client,
			Provider:             provider2,
			Label:                label,
			StoragePricePerEpoch: abi.NewTokenAmount(0),
			ProviderCollateral:   abi.NewTokenAmount(0),
			ClientCollateral:     abi.NewTokenAmount(0),
		},
	}

	for i, p := range props {
		buf := new(bytes.Buffer)
		assert.NoError(t, p.MarshalCBOR(buf))
		propData := buf.Bytes()

		ow := onWalletSign{
			Addr: client,
			Msg:  propData,
		}
		data, err := json.Marshal(ow)
		assert.NoError(t, err)
		fmt.Println(string(data))

		req, err := http.NewRequest(http.MethodPost, "http://192.168.200.18:10025/sign", bytes.NewReader(data))
		assert.NoError(t, err)

		reqSign := utils.Sign(data, []byte(aks))
		req.Header.Set("Authorization", utils.Authorize(aks, hex.EncodeToString(reqSign)))
		fmt.Println("Authorization: ", utils.Authorize(aks, hex.EncodeToString(reqSign)))

		resp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		if i != 0 {
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			continue
		}

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		buf = new(bytes.Buffer)
		_, err = buf.ReadFrom(resp.Body)
		assert.NoError(t, err)
		fmt.Println(buf.String())
	}
}
