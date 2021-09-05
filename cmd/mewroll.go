package cmd

import (
	"github.com/penguin-statistics/mewroll/internal/cmd"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"time"
)

func RunCliApp() {
	app := cli.NewApp()

	app.Name = "mewroll"
	app.Usage = "mewroll provides a way for the Mew community to randomly randomdrawer an arbitrary number of comments from a constraint scope of comments, in a willing to come out a winner in a classical present drawing activity"
	app.Flags = []cli.Flag {
		&cli.StringFlag{
			Name:        "deduplication",
			Aliases:     []string{"d"},
			Usage:       "Apply deduplication on data source retrieved. Possible values are: 'disabled' to disable deduplication, 'single' to deduplicate a user's possibility to win down to 1/total, and 'eliminate' to eliminate a user's possibility to win if they have more than 1 entity existed in the data source",
			Value:       "single",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name: "comment",
			Usage: "randomly draw some winners using comment as data source",
			Flags: []cli.Flag {
				&cli.StringFlag{
					Name:        "id",
					Usage:       "The thought id to extract comments from",
					Required: true,
				},
				&cli.IntFlag{
					Name:        "count",
					Usage:       "How many comments should be extracted",
					Value:		 1,
				},
			},
			Action:  func(c *cli.Context) error {
				return cmd.NewCommentDrawer(&cmd.CommentDrawerConfig{
					ThoughtID: c.String("id"),
					Deduplication: c.String("deduplication"),
					Count: c.Int("count"),
				}).Draw()
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		cli.ShowAppHelpAndExit(c, -1)
		return nil
	}

	app.Compiled = time.Now()
	app.EnableBashCompletion = true

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
