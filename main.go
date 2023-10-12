package main

import (
	"bufio"
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

func command(name string) cmd {
	f := flag.NewFlagSet(name, flag.ExitOnError)
	return cmd{
		flag:    f,
		project: f.String("project", "", "project"),
		trigger: f.String("trigger", "", "trigger"),
		branch:  f.String("branch", "", "branch"),
	}
}

func getProject() (string, error) {
	v, err := os.ReadFile(".config")
	if err != nil {
		return "", err
	}

	fmt.Println("Project:", string(v))

	return string(v), nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	deploy := command("deploy")
	status := command("status")

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

		if *deploy.trigger == "" {
			panic("trigger is required")
		}

		if *deploy.project == "" {
			p, err := getProject()

			if err != nil {
				panic(err)
			}

			*deploy.project = p
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
		ctx := context.Background()
		service, err := cloudbuild.NewService(ctx)
		if err != nil {
			panic(err)
		}

		if *status.trigger == "" {
			panic("trigger is required")
		}

		if *status.project == "" {
			p, err := getProject()

			if err != nil {
				panic(err)
			}

			*status.project = p
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
	case "config":
		fmt.Print("Default project: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		err := scanner.Err()

		if err != nil {
			panic(err)
		}

		text := scanner.Text()

		f, err := os.OpenFile(".config", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

		if err != nil {
			panic(err)
		}

		if _, err = f.WriteString(text); err != nil {
			panic(err)
		}

		defer f.Close()
	default:
		os.Exit(1)
	}
}
