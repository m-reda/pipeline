package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSchedulerStart(t *testing.T) {
	schedulerStart()

	require.Equal(t, len(mainPipeline.Schedule), len(scheduler[mainPipeline.ID]))

}

func TestSchedulerSet(t *testing.T) {
	id := "-id-"
	schedule := []string{"0 0 0 10 *", "* * * * 7", "0 0 * * 0"}
	schedulerSet(id, schedule)

	require.Equal(t, 3, len(scheduler[id]))
}
