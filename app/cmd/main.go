package main

import (
	"fmt"
	"github.com/faycheng/gob/plugin"
	"github.com/faycheng/gob/task"
	"github.com/faycheng/gob/worker"
	"github.com/spf13/cobra"
	"time"
)

var (
	rate        int
	step        int
	low         int
	high        int
	concurrency int
	duration    time.Duration
	taskArgs    string
)

func validateFlags() {
	return
}

func workerOpts() []worker.WorkerOpt {
	opts := make([]worker.WorkerOpt, 0)
	opts = append(opts, worker.WithQpsWorker())
	opts = append(opts, worker.WithConstanceBucket(rate))
	return opts
}

func gob(args []string) error {
	validateFlags()
	pluginPath := args[0]
	taskKey := task.Key(args[1])
	pluginLoader := plugin.NewPluginLoader(pluginPath)
	taskSet, err := pluginLoader.Load()
	if err != nil {
		return err
	}
	taskRegistry := task.NewTaskRegistry()
	for key, task := range taskSet {
		taskRegistry.Register(key, task)
	}
	targetTask, ok := taskRegistry.Get(taskKey)
	if !ok {
		return fmt.Errorf("task(%s) not found", taskKey)
	}
	factory := worker.NewWorkerFactory()
	worker, err := factory.NewWorker(targetTask, taskArgs, workerOpts()...)
	if err != nil {
		return err
	}
	worker.Run()
	return nil
}

func main() {
	var gobCmd = &cobra.Command{
		Use: "gob [plugin_path] [task]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return gob(args)
		},
	}
	flags := gobCmd.PersistentFlags()
	flags.IntVarP(&rate, "qps", "q", 100, "qps, default 100")
	flags.IntVarP(&concurrency, "concurrency", "c", 1, "concurrency, default 1")
	flags.StringVarP(&taskArgs, "data", "d", "{}", "target args in json format, default {}")
}
