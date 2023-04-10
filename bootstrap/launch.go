package bootstrap

import (
	log2 "log"
	"net/netip"
	"os"

	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli/v2"

	"fil-kms/app/global/variables"
)

var log = logging.Logger("client")

func Launch() {
	local := []*cli.Command{runCmd, WalletImport, WalletDelete, WalletList}
	app := &cli.App{
		Name:     "fil-kms",
		Usage:    "manage filecoin key",
		Commands: local,

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "ip",
				EnvVars: []string{"FIL_KMS_IP"},
				Value:   "127.0.0.1",
			},
			&cli.IntFlag{
				Name:    "port",
				EnvVars: []string{"FIL_KMS_PORT"},
				Value:   10025,
			},
			&cli.BoolFlag{
				Name:   "debug",
				Hidden: true,
			},
		},
		Before: beforeAnything,
	}

	err := app.Run(os.Args)
	if err != nil {
		log2.Println(err)
	}

}

func beforeAnything(cctx *cli.Context) error {
	_ = logging.SetLogLevel("*", "INFO")
	isDebug := cctx.Bool("debug")
	if isDebug {
		_ = logging.SetLogLevel("*", "DEBUG")
	}

	_ip := cctx.String("ip")
	_, err := netip.ParseAddr(_ip)
	if err != nil {
		return err
	}

	variables.ServerIP = _ip
	variables.ServerPort = cctx.Int("port")

	return nil
}
