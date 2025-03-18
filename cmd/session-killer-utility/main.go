package main

/*

 */
import (
	"fmt"
	"os"

	"github.com/Galdoba/hometools/cmd/session-killer-utility/commands"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "session-killer-utility"
	app.Version = "0.0.1"
	app.Usage = "monitor and control user's session time"
	app.Flags = []cli.Flag{}

	app.Commands = []*cli.Command{
		commands.Set(),
		commands.Check(),
	}

	//ПО ОКОНЧАНИЮ ДЕЙСТВИЯ
	app.After = func(c *cli.Context) error {
		return nil
	}
	args := os.Args
	if err := app.Run(args); err != nil {
		errOut := fmt.Sprintf("error: %v: %v", app.Name, err.Error())
		println(errOut)

		os.Exit(1)
	}

}
