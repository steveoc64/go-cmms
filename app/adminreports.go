package main

import (
	"time"

	"github.com/go-humble/router"
)

func adminReports(context *router.Context) {
	print("TODO - adminReports")

	go func() {
		newTasks := 0
		rpcClient.Call("TaskRPC.Generate", time.Now(), &newTasks)
		print("generated ", newTasks, " new tasks")

	}()
}
