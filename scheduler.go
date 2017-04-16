package main

import (
	"github.com/robfig/cron"
	"sync"
)

var (
	scheduler       = make(map[string][]*cron.Cron)
	schedulerLocker = sync.Mutex{}
)

func schedulerStart() {

	pipelines := loadAllPipelines()
	for _, p := range pipelines {
		schedulerSet(p.ID, p.Schedule)
	}
}

func schedulerSet(id string, expressions []string) {
	schedulerLocker.Lock()
	defer schedulerLocker.Unlock()

	// if the pipeline exist in the scheduler stop it
	if pCron, ok := scheduler[id]; ok {
		for _, c := range pCron {
			c.Stop()
		}
	}

	scheduler[id] = []*cron.Cron{}

	for _, exp := range expressions {
		c := cron.New()
		c.AddFunc(exp, func() {
			pipelineRun(id, nil)
		})
		c.Start()
		scheduler[id] = append(scheduler[id], c)
	}
}
