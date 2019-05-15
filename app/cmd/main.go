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

func gob(pluginPath, taskKey string) error {
	validateFlags()
	pluginLoader := plugin.NewPluginLoader(pluginPath)
	plugin, err := pluginLoader.Load()
	if err != nil {
		return err
	}
	taskRegistry := task.NewTaskRegistry()
	err = plugin.Run()
	if err != nil {
		return err
	}
	taskSet, err := plugin.Tasks()
	if err != nil {
		return err
	}
	for key, task := range taskSet {
		taskRegistry.Register(key, task)
	}
	targetTask, ok := taskRegistry.Get(string(taskKey))
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
			if len(args) != 2 {
				return fmt.Errorf("bad arguments(%+v)", args)
			}
			return gob(args[0], args[1])
		},
	}
	flags := gobCmd.PersistentFlags()
	flags.IntVarP(&rate, "qps", "q", 100, "qps")
	flags.IntVarP(&concurrency, "concurrency", "c", 1, "concurrency")
	flags.StringVarP(&taskArgs, "data", "d", "{}", "target args in json format")
	if err := gobCmd.Execute(); err != nil {
		panic(err)
	}
}
