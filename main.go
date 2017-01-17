package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

// Drone Plugin Env Ref:
//   http://readme.drone.io/usage/variables/
//   http://readme.drone.io/0.5/plugin-parameters/
//   http://readme.drone.io/0.5/usage/environment-reference/

func main() {
	app := cli.NewApp()
	app.Name = "spark plugin"
	app.Usage = "spark plugin"
	app.Action = run
	app.Version = fmt.Sprintf("1.0.%s", build)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "message",
			Usage:  "spark message",
			EnvVar: "PLUGIN_MESSAGE",
		},
		cli.StringFlag{
			Name:   "auth_token",
			Usage:  "spark auth token",
			EnvVar: "PLUGIN_AUTH_TOKEN",
		},
		cli.StringFlag{
			Name:   "roomId",
			Usage:  "spark room id",
			EnvVar: "PLUGIN_ROOMID",
		},
		cli.StringFlag{
			Name:   "roomName",
			Usage:  "spark room name",
			EnvVar: "PLUGIN_ROOMNAME",
		},
		cli.StringFlag{
			Name:   "system.link_url",
			Usage:  "drone server url",
			EnvVar: "DRONE_SERVER",
		},
		cli.StringFlag{
			Name:   "repo.owner",
			Usage:  "repository owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "repo.full_name",
			Usage:  "repository full name",
			EnvVar: "DRONE_REPO",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "commit.ref",
			Value:  "refs/heads/master",
			Usage:  "git commit ref",
			EnvVar: "DRONE_COMMIT_REF",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.author",
			Usage:  "git author name",
			EnvVar: "DRONE_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "commit.author_email",
			Usage:  "git author email",
			EnvVar: "DRONE_COMMIT_AUTHOR_EMAIL",
		},
		cli.StringFlag{
			Name:   "commit.commit_link",
			Usage:  "git commit link",
			EnvVar: "DRONE_COMMIT_LINK",
		},
		cli.StringFlag{
			Name:   "commit.message",
			Usage:  "git commit message",
			EnvVar: "DRONE_COMMIT_MESSAGE",
		},
		cli.StringFlag{
			Name:   "build.event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.IntFlag{
			Name:   "build.number",
			Usage:  "build number",
			EnvVar: "DRONE_BUILD_NUMBER",
		},
		cli.StringFlag{
			Name:   "build.status",
			Usage:  "build status",
			Value:  "success",
			EnvVar: "DRONE_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "build.link",
			Usage:  "build link",
			EnvVar: "DRONE_BUILD_LINK",
		},
		cli.Int64Flag{
			Name:   "build.started",
			Usage:  "build started",
			EnvVar: "DRONE_BUILD_STARTED",
		},
		cli.Int64Flag{
			Name:   "build.created",
			Usage:  "build created",
			EnvVar: "DRONE_BUILD_CREATED",
		},
		cli.StringFlag{
			Name:   "build.tag",
			Usage:  "build tag",
			EnvVar: "DRONE_TAG",
		},
		cli.Int64Flag{
			Name:   "job.started",
			Usage:  "job started",
			EnvVar: "DRONE_JOB_STARTED",
		},
		cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		Repo: Repo{
			Owner:    c.String("repo.owner"),
			Name:     c.String("repo.name"),
			FullName: c.String("repo.full_name"),
		},
		Build: Build{
			Tag:        c.String("build.tag"),
			Number:     c.Int("build.number"),
			Event:      c.String("build.event"),
			Status:     c.String("build.status"),
			Commit:     c.String("commit.sha"),
			Ref:        c.String("commit.ref"),
			Branch:     c.String("commit.branch"),
			Author:     c.String("commit.author"),
			Email:      c.String("commit.author_email"),
			Link:       c.String("build.link"),
			CommitLink: c.String("commit.commit_link"),
			Message:    c.String("commit.message"),
			DroneLink:  c.String("system.link_url"),
			Started:    c.Int64("build.started"),
			Created:    c.Int64("build.created"),
		},
		Job: Job{
			Started: c.Int64("job.started"),
		},
		Config: Config{
			Message:   c.String("message"),
			AuthToken: c.String("auth_token"),
			RoomID:    c.String("roomId"),
			RoomName:  c.String("roomName"),
		},
	}

	return plugin.Exec()
}
