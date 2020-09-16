package main

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	authmodule "gitlab.com/balconygames/analytics/modules/auth"
	collectormodule "gitlab.com/balconygames/analytics/modules/collector"
	leaderboardmodule "gitlab.com/balconygames/analytics/modules/leaderboard"
	pixelmodule "gitlab.com/balconygames/analytics/modules/pixel"
	primarymodule "gitlab.com/balconygames/analytics/modules/primary"
	"gitlab.com/balconygames/analytics/pkg/runtime"
)

func main() {
	if err := run(); err != nil {
		log.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	action := os.Args[1]

	var s runtime.Spec
	if err := envconfig.Process("", &s); err != nil {
		return err
	}

	r := runtime.New(action, s)

	runtimes := []runtime.InitFunc{
		authmodule.New,
		primarymodule.New,
		leaderboardmodule.New,
		collectormodule.New,
		pixelmodule.New,
	}

	for _, m := range runtimes {
		if err := m(r); err != nil {
			return err
		}
	}

	return r.Run()
}
