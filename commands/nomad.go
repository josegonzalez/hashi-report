package commands

import (
	"fmt"
	"os"

	"github.com/hashicorp/nomad/api"
	"github.com/josegonzalez/cli-skeleton/command"
	"github.com/posener/complete"
	flag "github.com/spf13/pflag"
)

type NomadCommand struct {
	command.Meta
}

func (c *NomadCommand) Name() string {
	return "nomad"
}

func (c *NomadCommand) Synopsis() string {
	return "Generate a nomad report"
}

func (c *NomadCommand) Help() string {
	return command.CommandHelp(c)
}

func (c *NomadCommand) Examples() map[string]string {
	appName := os.Getenv("CLI_APP_NAME")
	return map[string]string{
		"Generate a nomad usage report": fmt.Sprintf("%s %s", appName, c.Name()),
	}
}

func (c *NomadCommand) Arguments() []command.Argument {
	args := []command.Argument{}
	return args
}

func (c *NomadCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *NomadCommand) ParsedArguments(args []string) (map[string]command.Argument, error) {
	return command.ParseArguments(args, c.Arguments())
}

func (c *NomadCommand) FlagSet() *flag.FlagSet {
	f := c.Meta.FlagSet(c.Name(), command.FlagSetClient)
	return f
}

func (c *NomadCommand) AutocompleteFlags() complete.Flags {
	return command.MergeAutocompleteFlags(
		c.Meta.AutocompleteFlags(command.FlagSetClient),
		complete.Flags{},
	)
}

func (c *NomadCommand) Run(args []string) int {
	flags := c.FlagSet()
	flags.Usage = func() { c.Ui.Output(c.Help()) }
	if err := flags.Parse(args); err != nil {
		c.Ui.Error(err.Error())
		c.Ui.Error(command.CommandErrorText(c))
		return 1
	}

	_, err := c.ParsedArguments(flags.Args())
	if err != nil {
		c.Ui.Error(err.Error())
		c.Ui.Error(command.CommandErrorText(c))
		return 1
	}

	conf := api.DefaultConfig()
	if os.Getenv("NOMAD_ADDR") != "" {
		conf.Address = os.Getenv("NOMAD_ADDR")
	}

	client, err := api.NewClient(conf)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	regions, err := client.Regions().List()
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	nodes, _, err := client.Nodes().List(nil)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	jobs, _, err := client.Jobs().List(&api.QueryOptions{
		Namespace: "*",
	})
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	allocations, _, err := client.Allocations().List(&api.QueryOptions{
		Namespace: "*",
	})
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	taskCount := 0
	for _, allocation := range allocations {
		taskCount += len(allocation.TaskStates)
	}

	c.Ui.Output(fmt.Sprintf("Nomad address: %s", conf.Address))
	c.Ui.Output(fmt.Sprintf("Regions: %v", regions))
	c.Ui.Output(fmt.Sprintf("Node count: %d", len(nodes)))
	c.Ui.Output(fmt.Sprintf("Job count: %d", len(jobs)))
	c.Ui.Output(fmt.Sprintf("Allocation count: %d", len(allocations)))
	c.Ui.Output(fmt.Sprintf("Task count: %d", taskCount))

	return 0
}
