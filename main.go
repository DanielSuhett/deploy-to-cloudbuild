package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/olekukonko/tablewriter"
	"google.golang.org/api/cloudbuild/v1"
)

type cmd struct {
	flag    *flag.FlagSet
	project *string
	trigger *string
	branch  *string
}

func Command(name string) cmd {
	f := flag.NewFlagSet(name, flag.ExitOnError)
	return cmd{
		flag:    f,
		project: f.String("project", "", "project"),
		trigger: f.String("trigger", "", "trigger"),
		branch:  f.String("branch", "", "branch"),
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	deploy := Command("deploy")
	status := Command("status")

	if len(os.Args) < 1 {
		fmt.Println("expected subcommands")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "deploy":
		deploy.flag.Parse(os.Args[2:])
		ctx := context.Background()
		service, err := cloudbuild.NewService(ctx)
		if err != nil {
			panic(err)
		}

		if *deploy.project == "" || *deploy.trigger == "" || *deploy.branch == "" {
			panic("missing args")
		}

		b, err := service.Projects.Triggers.Run(
			*deploy.project,
			*deploy.trigger,
			&cloudbuild.RepoSource{BranchName: *deploy.branch},
		).Do()

		if err != nil {
			panic(err)
		}

		if b.HTTPStatusCode == 200 {
			logger.Info(b.Name)
			logger.Info(http.StatusText(b.HTTPStatusCode))
		} else {
			logger.Error(b.Error.Message)
		}
		os.Exit(1)
	case "status":
		status.flag.Parse(os.Args[2:])
		fmt.Println(os.Args[2:])
		ctx := context.Background()
		service, err := cloudbuild.NewService(ctx)
		if err != nil {
			panic(err)
		}

		if *status.project == "" || *status.trigger == "" {
			panic("missing args")
		}

		filter := fmt.Sprintf("substitutions.TRIGGER_NAME=\"%s\"", *status.trigger)

		b, err := service.Projects.Builds.List(
			*status.project,
		).Filter(filter).PageSize(5).Do()

		if err != nil {
			panic(err)
		}

		table := tablewriter.NewWriter(os.Stdout)

		table.SetHeader([]string{"Build", "Status", "Branch"})

		for _, build := range b.Builds {
			row := []string{build.Id, build.Status, build.Substitutions["BRANCH_NAME"]}
			table.Append(row)
		}

		table.SetBorder(false)
		table.SetCenterSeparator("|")
		table.SetColumnSeparator("|")
		table.SetRowSeparator("-")
		table.Render()

		os.Exit(1)
	default:
		os.Exit(1)
	}
}
