package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/opencontainers/runc/libcontainer"
	"github.com/urfave/cli"
)

var startCommand = cli.Command{
	Name:  "start",
	Usage: "executes the user defined process in a created container",
	ArgsUsage: `<container-id>

Where "<container-id>" is your name for the instance of the container that you
are starting. The name you provide for the container instance must be unique on
your host.`,
	Description: `The start command executes the user defined process in a created container.`,
	Action: func(context *cli.Context) error {
		fmt.Println("RUNC: start command: checkArgs")
		if err := checkArgs(context, 1, exactArgs); err != nil {
			return err
		}

		fmt.Println("RUNC: start command: get Container")

		container, err := getContainer(context)
		if err != nil {
			return err
		}

		fmt.Println("RUNC: start command: get Container")
		status, err := container.Status()
		if err != nil {
			return err
		}
		fmt.Printf("RUNC: start command: Container status = %d\n", status)
		switch status {
		case libcontainer.Created:
			fmt.Printf("RUNC: start command: the container is created\n")
			notifySocket, err := notifySocketStart(context, os.Getenv("NOTIFY_SOCKET"), container.ID())

			fmt.Printf("RUNC: notify socket start\n")
			if notifySocket != nil {
				fmt.Printf("RUNC: notify socket %s\n", notifySocket.host)
			}

			if err != nil {
				return err
			}

			fmt.Printf("RUNC: exec container\n")

			if err := container.Exec(); err != nil {
				return err
			}

			fmt.Printf("RUNC: exec container returns\n")
			if notifySocket != nil {
				fmt.Printf("RUNC: Wait for container\n")
				return notifySocket.waitForContainer(container)
			}
			fmt.Printf("RUNC: start container returns\n")
			return nil
		case libcontainer.Stopped:
			fmt.Printf("RUNC: start command: the container is stopped\n")
			return errors.New("cannot start a container that has stopped")
		case libcontainer.Running:
			return errors.New("cannot start an already running container")
		default:
			return fmt.Errorf("cannot start a container in the %s state", status)
		}
	},
}
