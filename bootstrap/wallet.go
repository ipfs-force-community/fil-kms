package bootstrap

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/types"
	cli2 "github.com/filecoin-project/lotus/cli"
	"github.com/urfave/cli/v2"

	cli3 "fil-kms/app/cli"
)

var WalletImport = &cli.Command{
	Name:      "import",
	Usage:     "import keys",
	ArgsUsage: "[<path> (optional, will read from stdin if omitted)]",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "format",
			Usage: "specify input format for key",
			Value: "hex",
		},
	},
	Action: func(cctx *cli.Context) error {

		var inpdata []byte
		if !cctx.Args().Present() || cctx.Args().First() == "-" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter private key: ")
			indata, err := reader.ReadBytes('\n')
			if err != nil {
				return err
			}
			inpdata = indata

		} else {
			fdata, err := os.ReadFile(cctx.Args().First())
			if err != nil {
				return err
			}
			inpdata = fdata
		}

		var ki types.KeyInfo
		switch cctx.String("format") {
		case "hex":
			data, err := hex.DecodeString(strings.TrimSpace(string(inpdata)))
			if err != nil {
				return err
			}

			if err := json.Unmarshal(data, &ki); err != nil {
				return err
			}
		case "json":
			if err := json.Unmarshal(inpdata, &ki); err != nil {
				return err
			}
		case "gfc-json":
			var f struct {
				KeyInfo []struct {
					PrivateKey []byte
					SigType    int
				}
			}
			if err := json.Unmarshal(inpdata, &f); err != nil {
				return fmt.Errorf("failed to parse go-filecoin key: %s", err)
			}

			gk := f.KeyInfo[0]
			ki.PrivateKey = gk.PrivateKey
			switch gk.SigType {
			case 1:
				ki.Type = types.KTSecp256k1
			case 2:
				ki.Type = types.KTBLS
			default:
				return fmt.Errorf("unrecognized key type: %d", gk.SigType)
			}
		default:
			return fmt.Errorf("unrecognized format: %s", cctx.String("format"))
		}
		addr, err := clientImport(ki)
		if err != nil {
			return err
		}

		fmt.Printf("imported key %s successfully!\n", addr)
		return nil
	},
}

var WalletDelete = &cli.Command{
	Name:      "delete",
	Usage:     "Soft delete an address from the wallet - hard deletion needed for permanent removal",
	ArgsUsage: "<address> ",
	Action: func(cctx *cli.Context) error {

		if !cctx.Args().Present() || cctx.NArg() != 1 {
			return fmt.Errorf("must specify address to delete")
		}

		addr, err := address.NewFromString(cctx.Args().First())
		if err != nil {
			return err
		}

		return clientDelete(addr)
	},
}

var WalletList = &cli.Command{
	Name:  "list",
	Usage: "List wallet address",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {

		afmt := cli2.NewAppFmt(cctx.App)

		addrs, err := clientList()
		if err != nil {
			return err
		}

		for _, addr := range addrs {
			afmt.Println(addr.String())
		}

		return nil
	},
}

func clientList() ([]address.Address, error) {
	res, err := cli3.Invoke("WalletList", nil)
	if err != nil {
		return nil, err
	}

	var addrs []address.Address
	err = json.Unmarshal(res, &addrs)
	if err != nil {
		return nil, err
	}

	return addrs, nil
}

func clientImport(keyInfo types.KeyInfo) (address.Address, error) {
	params, _ := json.Marshal(keyInfo)
	retInfo, err := cli3.Invoke("WalletImport", params)
	if err != nil {
		return address.Undef, err
	}
	addr, err := address.NewFromBytes(retInfo)
	if err != nil {
		return address.Address{}, err
	}

	return addr, nil
}

func clientDelete(k address.Address) error {
	_, err := cli3.Invoke("WalletDelete", k.Bytes())
	if err != nil {
		return err
	}
	return nil
}
