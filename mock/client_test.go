package mock_test

import (
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClient_Close(t *testing.T) {
	c := &mock.Client{}
	c.On("Close").Return(nil)

	require.NoError(t, c.Close())
}

func TestClient_Put(t *testing.T) {
	c := &mock.Client{}
	c.On("Put", uint32(0), time.Duration(0), 5*time.Second, []byte("test")).Return(1, nil)

	id, err := c.Put(0, 0, 5*time.Second, []byte("test"))

	require.Equal(t, 1, id)
	require.Nil(t, err)
}

func TestClient_Use(t *testing.T) {
	c := &mock.Client{}
	c.On("Use", "test").Return("test", nil)

	tube, err := c.Use("test")

	require.Equal(t, "test", tube)
	require.Nil(t, err)
}

func TestClient_Reserve(t *testing.T) {
	expectedJob := beanstalk.Job{ID: 1, Data: []byte("test")}

	c := &mock.Client{}
	c.On("Reserve").Return(expectedJob, nil)

	job, err := c.Reserve()

	require.Equal(t, expectedJob, job)
	require.Nil(t, err)
}

func TestClient_ReserveWithTimeout(t *testing.T) {
	expectedJob := beanstalk.Job{ID: 1, Data: []byte("test")}

	c := &mock.Client{}
	c.On("ReserveWithTimeout", 5*time.Second).Return(expectedJob, nil)

	job, err := c.ReserveWithTimeout(5 * time.Second)

	require.Equal(t, expectedJob, job)
	require.Nil(t, err)
}

func TestClient_ReserveJob(t *testing.T) {
	expectedJob := beanstalk.Job{ID: 1, Data: []byte("test")}

	c := &mock.Client{}
	c.On("ReserveJob", 1).Return(expectedJob, nil)

	job, err := c.ReserveJob(1)

	require.Equal(t, expectedJob, job)
	require.Nil(t, err)
}

func TestClient_Delete(t *testing.T) {
	c := &mock.Client{}
	c.On("Delete", 1).Return(nil)

	require.NoError(t, c.Delete(1))
}

func TestClient_Release(t *testing.T) {
	c := &mock.Client{}
	c.On("Release", 1, uint32(0), 5*time.Second).Return(nil)

	require.NoError(t, c.Release(1, 0, 5*time.Second))
}

func TestClient_Bury(t *testing.T) {
	c := &mock.Client{}
	c.On("Bury", 1, uint32(0)).Return(nil)

	require.NoError(t, c.Bury(1, 0))
}

func TestClient_Touch(t *testing.T) {
	c := &mock.Client{}
	c.On("Touch", 1).Return(nil)

	require.NoError(t, c.Touch(1))
}

func TestClient_Watch(t *testing.T) {
	c := &mock.Client{}
	c.On("Watch", "test").Return(1, nil)

	count, err := c.Watch("test")

	require.Equal(t, 1, count)
	require.Nil(t, err)
}

func TestClient_Ignore(t *testing.T) {
	c := &mock.Client{}
	c.On("Ignore", "test").Return(1, nil)

	count, err := c.Ignore("test")

	require.Equal(t, 1, count)
	require.Nil(t, err)
}

func TestClient_Peek(t *testing.T) {
	expectedJob := beanstalk.Job{ID: 1, Data: []byte("test")}

	c := &mock.Client{}
	c.On("Peek", 1).Return(expectedJob, nil)

	job, err := c.Peek(1)

	require.Equal(t, expectedJob, job)
	require.Nil(t, err)
}

func TestClient_PeekReady(t *testing.T) {
	expectedJob := beanstalk.Job{ID: 1, Data: []byte("test")}

	c := &mock.Client{}
	c.On("PeekReady").Return(expectedJob, nil)

	job, err := c.PeekReady()

	require.Equal(t, expectedJob, job)
	require.Nil(t, err)
}

func TestClient_PeekDelayed(t *testing.T) {
	expectedJob := beanstalk.Job{ID: 1, Data: []byte("test")}

	c := &mock.Client{}
	c.On("PeekDelayed").Return(expectedJob, nil)

	job, err := c.PeekDelayed()

	require.Equal(t, expectedJob, job)
	require.Nil(t, err)
}

func TestClient_PeekBuried(t *testing.T) {
	expectedJob := beanstalk.Job{ID: 1, Data: []byte("test")}

	c := &mock.Client{}
	c.On("PeekBuried").Return(expectedJob, nil)

	job, err := c.PeekBuried()

	require.Equal(t, expectedJob, job)
	require.Nil(t, err)
}

func TestClient_Kick(t *testing.T) {
	c := &mock.Client{}
	c.On("Kick", 5).Return(5, nil)

	count, err := c.Kick(5)

	require.Equal(t, 5, count)
	require.Nil(t, err)
}

func TestClient_KickJob(t *testing.T) {
	c := &mock.Client{}
	c.On("KickJob", 1).Return(nil)

	require.NoError(t, c.KickJob(1))
}

func TestClient_StatsJob(t *testing.T) {
	expectedStats := beanstalk.StatsJob{}

	c := &mock.Client{}
	c.On("StatsJob", 1).Return(expectedStats, nil)

	stats, err := c.StatsJob(1)

	require.Equal(t, expectedStats, stats)
	require.Nil(t, err)
}

func TestClient_StatsTube(t *testing.T) {
	expectedStats := beanstalk.StatsTube{}

	c := &mock.Client{}
	c.On("StatsTube", "test").Return(expectedStats, nil)

	stats, err := c.StatsTube("test")

	require.Equal(t, expectedStats, stats)
	require.Nil(t, err)
}

func TestClient_Stats(t *testing.T) {
	expectedStats := beanstalk.Stats{}

	c := &mock.Client{}
	c.On("Stats").Return(expectedStats, nil)

	stats, err := c.Stats()

	require.Equal(t, expectedStats, stats)
	require.Nil(t, err)
}

func TestClient_ListTubes(t *testing.T) {
	expectedTubes := []string{"default", "test"}

	c := &mock.Client{}
	c.On("ListTubes").Return(expectedTubes, nil)

	tubes, err := c.ListTubes()

	require.ElementsMatch(t, expectedTubes, tubes)
	require.Nil(t, err)
}

func TestClient_ListTubeUsed(t *testing.T) {
	c := &mock.Client{}
	c.On("ListTubeUsed").Return("test", nil)

	tube, err := c.ListTubeUsed()

	require.Equal(t, "test", tube)
	require.Nil(t, err)
}

func TestClient_ListTubesWatched(t *testing.T) {
	expectedTubes := []string{"default", "test"}

	c := &mock.Client{}
	c.On("ListTubesWatched").Return(expectedTubes, nil)

	tubes, err := c.ListTubesWatched()

	require.ElementsMatch(t, expectedTubes, tubes)
	require.Nil(t, err)
}

func TestClient_PauseTube(t *testing.T) {
	c := &mock.Client{}
	c.On("PauseTube", "test", 10*time.Second).Return(nil)

	require.NoError(t, c.PauseTube("test", 10*time.Second))
}
