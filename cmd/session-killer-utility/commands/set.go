package commands

import (
	"log"
	"os"
	"os/exec"

	"github.com/Galdoba/hometools/cmd/session-killer-utility/commands/skucfg"
	"github.com/Galdoba/hometools/pkg/config"
	"github.com/urfave/cli/v2"
)

func Set() *cli.Command {
	return &cli.Command{
		Name:   "set",                                //
		Usage:  "set opens config file with $EDITOR", //
		Action: setConfiguration,                     //
		//Subcommands: []*cli.Command{{Name: "health", Usage: "edts config health", Action: ReadConfig}}, //

	}
}

func setConfiguration(c *cli.Context) error {
	cf, err := config.New(c.App.Name, config.JSON)
	if err != nil {
		return err
	}
	cfg := skucfg.Config{DailyAllowance: 90}

	if err := cf.Write(&cfg); err != nil {
		return err
	}

	cmd := exec.Command(os.Getenv("EDITOR"), cf.Path())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
