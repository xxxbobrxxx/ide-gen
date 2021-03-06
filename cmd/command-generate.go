package main

import (
	"github.com/spf13/cobra"
	"github.com/xxxbobrxxx/ide-gen/pkg/config"
	"github.com/xxxbobrxxx/ide-gen/pkg/idea"
	"github.com/xxxbobrxxx/ide-gen/pkg/repository"
)

type GenerateCommand struct {
	config.GlobalFlags
	repository.SourcesRootFlags
	idea.Project

	cmd *cobra.Command
}

func NewGenerateCommand() *GenerateCommand {
	command := &GenerateCommand{}

	cmd := &cobra.Command{
		Use:          "generate",
		Aliases:      []string{"gen"},
		Short:        "Clone repositories and generate IDE project",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		RunE:         command.Execute,
	}
	command.cmd = cmd

	command.Project.AddFlags(cmd.PersistentFlags())
	command.GlobalFlags.AddFlags(cmd.PersistentFlags())
	command.SourcesRootFlags.AddFlags(cmd.PersistentFlags())

	_ = command.cmd.MarkPersistentFlagRequired("config")

	return command
}

func (command *GenerateCommand) Register() *cobra.Command {
	return command.cmd
}

func (command *GenerateCommand) Execute(_ *cobra.Command, _ []string) (err error) {
	c, err := command.ReadConfig()
	if err != nil {
		return err
	}

	//: Read entries from config and flags
	projectEntries, err := c.GetProjectEntries(command.SourcesRootFlags)
	if err != nil {
		return err
	}

	//: Clone repos
	for _, projectEntry := range projectEntries {
		exists, err := projectEntry.Commander.Exists(projectEntry.Directory)
		if err != nil {
			return err
		}

		if exists {
			logger.Infof(
				"Skip clone project '%s' to '%s'", projectEntry.Name, projectEntry.Directory)
		} else {
			logger.Infof(
				"Clone project '%s' to '%s'", projectEntry.Name, projectEntry.Directory)
			err = projectEntry.Commander.Clone(projectEntry.Directory)
			if err != nil {
				return err
			}
		}
	}

	//: Idea project
	if command.Project.Root != "" {
		project := command.Project
		for _, projectEntry := range projectEntries {
			project.AddEntry(projectEntry)
		}

		logger.Infof("Writing idea project %s", project.Root)
		err = project.Write()
		if err != nil {
			return err
		}
	}

	return nil
}
