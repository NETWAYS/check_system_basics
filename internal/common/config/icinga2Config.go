package config

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/NETWAYS/go-icingadsl"
)

var ccaFlags []icingadsl.CheckCommandArgument

func GenerateIcinga2Config(cmd *cobra.Command, commandName, executableName string, _ bool) string {
	checkCommands := make([]icingadsl.CheckCommand, 0)

	flags := cmd.Flags()

	parentDefinition := icingadsl.CheckCommand{
		Name: commandName,
		Command: icingadsl.Array{
			icingadsl.InfixExpression{
				Left:          icingadsl.Identifier("PluginContribDir"),
				InfixOperator: icingadsl.PLUS,
				Right:         icingadsl.String("/" + executableName),
			},
		},
	}

	ccaFlags = make([]icingadsl.CheckCommandArgument, 0)

	flags.VisitAll(func(flag *pflag.Flag) {
		if flag.Name != "help" && flag.Name != "debug" {
			_ = GenerateIcinga2CheckCommandArgument(flag, &ccaFlags)
		}
	})

	parentArgs := make([]icingadsl.CheckCommandArgument, len(ccaFlags))
	copy(parentArgs, ccaFlags)

	parentDefinition.Arguments = parentArgs

	checkCommands = append(checkCommands, parentDefinition)

	for _, command := range cmd.Commands() {
		// Ignore the magical cobra command "no-help"
		if command.Name() == "no-help" {
			continue
		}

		ccaFlags = make([]icingadsl.CheckCommandArgument, 0)

		// This triggers a side effect to get inherited flags
		// do not remove
		command.InheritedFlags()

		scFlags := command.Flags()
		scFlags.VisitAll(func(foo *pflag.Flag) {
			if foo.Name != "help" && foo.Name != "debug" {
				_ = GenerateIcinga2CheckCommandArgument(foo, &ccaFlags)
			}
		})

		cc := icingadsl.CheckCommand{}

		cc.Name = parentDefinition.Name + "_" + command.Name()

		cc.Command = append(cc.Command, parentDefinition.Command...)

		cc.Imports = []*icingadsl.CheckCommand{&parentDefinition}
		args := make([]icingadsl.CheckCommandArgument, len(ccaFlags)+1)
		copy(args, ccaFlags)

		subcommandArgument := icingadsl.CheckCommandArgument{
			Name:  command.Name(),
			Order: -1,
			SetIf: icingadsl.True,
		}

		args[len(ccaFlags)] = subcommandArgument
		cc.Arguments = args

		checkCommands = append(checkCommands, cc)
	}

	resultWriter := strings.Builder{}

	for _, command := range checkCommands {
		resultWriter.WriteString(command.String())
	}

	return resultWriter.String()
}

func GenerateIcinga2CheckCommandArgument(flags *pflag.Flag, returnList *[]icingadsl.CheckCommandArgument) error {
	cca := icingadsl.CheckCommandArgument{}

	cca.Name = "--" + flags.Name
	cca.Description = icingadsl.String(flags.Usage)

	switch flags.Value.Type() {
	case "bool":
		cca.SetIf = icingadsl.String(flags.Name)
		cca.SkipKey = false
	case "stringSlice":
		cca.RepeatKey = true
		cca.Value = flags.Name
	default:
		cca.Value = flags.Name
	}

	*returnList = append(*returnList, cca)

	return nil
}
