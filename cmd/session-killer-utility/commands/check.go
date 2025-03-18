package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/Galdoba/hometools/cmd/session-killer-utility/commands/skucfg"
	"github.com/Galdoba/hometools/pkg/config"
	"github.com/urfave/cli/v2"
)

var flag_SOFT = "soft"
var allowance_file = "/tmp/session-allowance"
var config_dir = ""

func Check() *cli.Command {
	return &cli.Command{
		Name:   "check",                                          //
		Usage:  "read session-allowence file and reduce it by 1", //
		Action: checkAllowance,                                   //
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    flag_SOFT,
				Usage:   "do not reduce allowance",
				Value:   false,
				Aliases: []string{"s"},
			},
		},
	}
}

func checkAllowance(c *cli.Context) error {
	cfg, err := readCfg(c.App.Name)
	if err != nil {
		return err
	}
	////////////////
	bt, err := os.ReadFile(allowance_file)
	if err != nil {
		switch errors.Is(err, os.ErrNotExist) {
		case true:
			//create max
			return createNewAllowanceFile(cfg)
		case false:
			return fmt.Errorf("unexpected error reading file: %v", err)
		}
	}
	al, err := strconv.Atoi(string(bt))
	if err != nil {
		return fmt.Errorf("failed to get allowance from '%v'", err)
	}
	if !c.Bool(flag_SOFT) {
		al--
	}
	fmt.Fprintf(os.Stdout, "%v", al)

	switch al {
	case 3, 2, 1:
		countdown(al)
	case 0:
		shutdown()
	}
	return nil
}

func run(comm string, args ...string) error {
	cmd := exec.Cmd{}
	cmd.Path = comm
	cmd.Args = args
	return cmd.Run()
}

func countdown(i int) {
	cmd := exec.Cmd{}
	cmd.Path = "mplayer"
	switch i {
	case 3:
		cmd.Args = []string{fmt.Sprintf("%v3_minutes.mp3", config_dir)}
	case 2:
		cmd.Args = []string{fmt.Sprintf("%v2_minutes.mp3", config_dir)}
	case 1:
		cmd.Args = []string{fmt.Sprintf("%v1_minute.mp3", config_dir)}
	}
	cmd.Run()
}

func shutdown() {
	cmd := exec.Cmd{}
	cmd.Path = "xfce4-session-logout"
	cmd.Args = []string{"--logout"}
	cmd.Run()
}

func createNewAllowanceFile(cfg *skucfg.Config) error {
	f, err := os.Create(allowance_file)
	if err != nil {
		return fmt.Errorf("allowance file creation failed: %v", err)
	}
	_, err = f.WriteString(fmt.Sprintf("%v", cfg.DailyAllowance))
	if err != nil {
		return fmt.Errorf("failed to set allowance file: %v", err)
	}
	return nil
}

func readCfg(name string) (*skucfg.Config, error) {
	cf, err := config.New(name, config.JSON)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
	config_dir = cf.Dir()
	btCfg, err := cf.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read config struct: %v", err)
	}
	cfg := &skucfg.Config{}
	if err := cfg.FromBytes(btCfg); err != nil {
		return nil, fmt.Errorf("failed to convert bytes to config struct: %v", err)
	}
	return cfg, nil
}
