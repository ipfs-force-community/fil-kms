package bootstrap

import (
	"fil-kms/app/http/router"
	"fil-kms/app/service/sign_service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/go-homedir"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"fil-kms/app/global/variables"
	"github.com/urfave/cli/v2"
)

const FlagFilKMSRepo = "storage"

var runCmd = &cli.Command{
	Name:  "run",
	Usage: "start daemon process",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  FlagFilKMSRepo,
			Value: "~/.filkms",
		},
		&cli.StringFlag{
			Name:     "aks",
			EnvVars:  []string{"FIL_KMS_AKS"},
			Required: true,
		},
	},
	Action: run,
	Before: initBeforeRun,
}

func run(ctx *cli.Context) error {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	log.Info("closed")
	return nil
}

func initBeforeRun(cctx *cli.Context) error {
	variables.AKS = cctx.String("aks")

	dir, err := homedir.Expand(cctx.String(FlagFilKMSRepo))
	if err != nil {
		log.Errorw("could not expand repo location", "error", err)
		return err
	} else {
		log.Infof("fil-kms repo: %s", dir)
	}

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		log.Errorf("create repo dir %v failed,err: %v", dir, err)
		return err
	}

	err = sign_service.Init(filepath.Join(dir, "storage"))
	if err != nil {
		log.Errorf("create repo dir %v failed,err: %v", dir, err)
		return err
	}

	err = initServer()
	if err != nil {
		log.Errorf("init server failed,listen: %v,err: %v", cctx.String("listen"), err)
		return err
	}
	return nil
}

func initServer() error {
	ip := variables.ServerIP
	port := variables.ServerPort

	if variables.IsDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := router.InitRouter()
	var serverLocal *http.Server
	var serverRpc *http.Server

	if ip != "127.0.0.1" && ip != "0.0.0.0" {
		serverLocal = &http.Server{
			Addr:         fmt.Sprintf("%s:%d", "127.0.0.1", port),
			ReadTimeout:  time.Second * 3,
			WriteTimeout: time.Second * 60,
			Handler:      engine,
		}
	}
	serverRpc = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", ip, port),
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 3,
		Handler:      engine,
	}

	if serverLocal != nil {
		go func() {
			err := serverLocal.ListenAndServe()
			if err != nil {
				panic(err)
			}
		}()
	}

	if serverRpc != nil {
		go func() {
			err := serverRpc.ListenAndServe()
			if err != nil {
				panic(err)
			}
		}()
	}

	log.Infof("listening %s:%d", ip, port)

	return nil
}
