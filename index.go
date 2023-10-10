package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"google.golang.org/api/cloudbuild/v1"
)

func main() {
	h := slog.NewJSONHandler(os.Stderr, nil)
	l := slog.New(h)
	cmd := flag.NewFlagSet("deploy", flag.ExitOnError)
	project := cmd.String("project", "", "project")
	trigger := cmd.String("trigger", "", "trigger")
	branch := cmd.String("branch", "", "branch")

	if len(os.Args) < 1 {
		fmt.Println("expected subcommands")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "deploy":
		cmd.Parse(os.Args[2:])
		ctx := context.Background()
		service, err := cloudbuild.NewService(ctx)
		if err != nil {
			panic(err)
		}

		if *project == "" || *trigger == "" || *branch == "" {
			panic("missing args")
		}

		b, err := service.Projects.Triggers.Run(
			*project,
			*trigger,
			&cloudbuild.RepoSource{BranchName: *branch},
		).Do()

		if err != nil {
			panic(err)
		}

		if b.HTTPStatusCode == 200 {
			l.Info(b.Name)
			l.Info(http.StatusText(b.HTTPStatusCode))
		} else {
			l.Error(b.Error.Message)
		}
		os.Exit(1)
	default:
		os.Exit(1)
	}
}
