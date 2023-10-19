package taskstore

import (
	"strings"

	"github.com/gouniverse/utils"
	"github.com/mingrammer/cfmt"
)

// TaskExecuteCli - CLI tool to find a task by its alias and execute its handler
// - alias "list" is reserved. it lists all the available commands
func (store *Store) TaskExecuteCli(alias string, args []string) bool {
	argumentsMap := utils.ArgsToMap(args)
	cfmt.Infoln("Executing task: ", alias, " with arguments: ", argumentsMap)

	// Lists the available tasks
	if alias == "list" {
		for index, taskHandler := range store.TaskHandlerList() {
			cfmt.Warningln(utils.ToString(index+1) + ". Task Alias: " + taskHandler.Alias())
			cfmt.Infoln("    - Task Title: " + taskHandler.Title())
			cfmt.Infoln("    - Task Description: " + taskHandler.Description())
		}

		return true
	}

	// Finds the task and executes its handler
	for _, taskHandler := range store.TaskHandlerList() {
		if strings.EqualFold(unifyName(taskHandler.Alias()), unifyName(alias)) {
			taskHandler.SetOptions(argumentsMap)
			taskHandler.Handle()
			return true
		}
	}

	cfmt.Errorln("Unrecognized task alias: ", alias)
	return false
}

func unifyName(name string) string {
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "_", "")
	return name
}
