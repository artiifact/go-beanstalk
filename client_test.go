package beanstalk_test

import (
	"bytes"
	"fmt"
	"github.com/IvanLutokhin/go-beanstalk"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestClient_Put(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"put 0 0 0 13\r\ntest draining\r\n",
				"put 1 1 1 16\r\ntest job too big\r\n",
				"put 0 30 90 18\r\ntest expected CRLF\r\n",
				"put 100 0 1800 11\r\ntest buried\r\n",
				"put 1 5 600 4\r\ntest\r\n",
			},
			[]string{
				"DRAINING\r\n",
				"JOB_TOO_BIG\r\n",
				"EXPECTED_CRLF\r\n",
				"BURIED 1\r\n",
				"INSERTED 1\r\n",
			},
		),
	)

	if _, err := c.Put(0, 0, 0, []byte("test draining")); err != beanstalk.ErrDraining {
		t.Fatal("expected draining error, but got", err)
	}

	if _, err := c.Put(1, 1*time.Second, 1*time.Second, []byte("test job too big")); err != beanstalk.ErrJobTooBig {
		t.Fatal("expected job too big error, but got", err)
	}

	if _, err := c.Put(0, 30*time.Second, 90*time.Second, []byte("test expected CRLF")); err != beanstalk.ErrExpectedCRLF {
		t.Fatal("expected expected CRLF error, but got", err)
	}

	if _, err := c.Put(100, 0, 30*time.Minute, []byte("test buried")); err != beanstalk.ErrBuried {
		t.Fatal("expected buried error, but got", err)
	}

	id, err := c.Put(1, 5*time.Second, 10*time.Minute, []byte("test"))

	if err != nil {
		t.Error(err)
	}

	if id != 1 {
		t.Fatal("expected 1, but got", id)
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_Use(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"use test\r\n",
			},
			[]string{
				"USING test\r\n",
			},
		),
	)

	tube, err := c.Use("test")

	if err != nil {
		t.Error(err)
	}

	if !strings.EqualFold("test", tube) {
		t.Fatal(fmt.Sprintf("expected 'test', but got '%s'", tube))
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_Reserve(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"reserve\r\n",
				"reserve\r\n",
				"reserve\r\n",
			},
			[]string{
				"TIMED_OUT\r\n",
				"DEADLINE_SOON\r\n",
				"RESERVED 1 4\r\ntest\r\n",
			},
		),
	)

	if _, err := c.Reserve(); err != beanstalk.ErrTimedOut {
		t.Fatal("expected timed out error, but got", err)
	}

	if _, err := c.Reserve(); err != beanstalk.ErrDeadlineSoon {
		t.Fatal("expected deadline soon error, but got", err)
	}

	job, err := c.Reserve()

	if err != nil {
		t.Error(err)
	}

	if job.ID != 1 {
		t.Fatal("expected 1, but got", job.ID)
	}

	if bytes.Compare([]byte("test"), job.Data) != 0 {
		t.Fatal(fmt.Sprintf("expected 'test', but got '%s'", string(job.Data)))
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_ReserveWithTimeout(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"reserve-with-timeout 5\r\n",
				"reserve-with-timeout 60\r\n",
				"reserve-with-timeout 600\r\n",
			},
			[]string{
				"TIMED_OUT\r\n",
				"DEADLINE_SOON\r\n",
				"RESERVED 1 4\r\ntest\r\n",
			},
		),
	)

	if _, err := c.ReserveWithTimeout(5 * time.Second); err != beanstalk.ErrTimedOut {
		t.Fatal("expected timed out error, but got", err)
	}

	if _, err := c.ReserveWithTimeout(60 * time.Second); err != beanstalk.ErrDeadlineSoon {
		t.Fatal("expected deadline soon error, but got", err)
	}

	job, err := c.ReserveWithTimeout(10 * time.Minute)

	if err != nil {
		t.Error(err)
	}

	if job.ID != 1 {
		t.Fatal("expected 1, but got", job.ID)
	}

	if bytes.Compare([]byte("test"), job.Data) != 0 {
		t.Fatal(fmt.Sprintf("expected 'test', but got '%s'", string(job.Data)))
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_ReserveJob(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"reserve-job 1\r\n",
				"reserve-job 1\r\n",
			},
			[]string{
				"NOT_FOUND\r\n",
				"RESERVED 1 4\r\ntest\r\n",
			},
		),
	)

	if _, err := c.ReserveJob(1); err != beanstalk.ErrNotFound {
		t.Fatal("expected not found error, but got", err)
	}

	job, err := c.ReserveJob(1)

	if err != nil {
		t.Error(err)
	}

	if job.ID != 1 {
		t.Fatal("expected 1, but got", job.ID)
	}

	if bytes.Compare([]byte("test"), job.Data) != 0 {
		t.Fatal(fmt.Sprintf("expected 'test', but got '%s'", string(job.Data)))
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_Delete(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"delete 1\r\n",
				"delete 1\r\n",
			},
			[]string{
				"NOT_FOUND\r\n",
				"DELETED\r\n",
			},
		),
	)

	if err := c.Delete(1); err != beanstalk.ErrNotFound {
		t.Fatal("expected not found error, but got", err)
	}

	err := c.Delete(1)

	if err != nil {
		t.Error(err)
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_Release(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"release 1 0 5\r\n",
				"release 1 10 60\r\n",
				"release 1 999 600\r\n",
			},
			[]string{
				"NOT_FOUND\r\n",
				"BURIED\r\n",
				"RELEASED\r\n",
			},
		),
	)

	if err := c.Release(1, 0, 5*time.Second); err != beanstalk.ErrNotFound {
		t.Fatal("expected not found error, but got", err)
	}

	if err := c.Release(1, 10, 60*time.Second); err != beanstalk.ErrBuried {
		t.Fatal("expected buried error, but got", err)
	}

	err := c.Release(1, 999, 10*time.Minute)

	if err != nil {
		t.Error(err)
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_Bury(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"bury 999 100\r\n",
				"bury 1 10\r\n",
			},
			[]string{
				"NOT_FOUND\r\n",
				"BURIED\n",
			},
		),
	)

	if err := c.Bury(999, 100); err != beanstalk.ErrNotFound {
		t.Fatal("expected not found error, but got", err)
	}

	if err := c.Bury(1, 10); err != nil {
		t.Fatal(err)
	}

	if err := c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_Touch(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"touch 1\r\n",
				"touch 1\r\n",
			},
			[]string{
				"NOT_FOUND\r\n",
				"TOUCHED\r\n",
			},
		),
	)

	if err := c.Touch(1); err != beanstalk.ErrNotFound {
		t.Fatal("expected not found error, but got", err)
	}

	err := c.Touch(1)

	if err != nil {
		t.Error(err)
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_Watch(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"watch test\r\n",
			},
			[]string{
				"WATCHING 1\r\n",
			},
		),
	)

	count, err := c.Watch("test")

	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Fatal("expected 1, but got", count)
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_Ignore(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"ignore test\r\n",
				"ignore test\r\n",
			},
			[]string{
				"NOT_IGNORED\r\n",
				"WATCHING 1\n",
			},
		),
	)

	if _, err := c.Ignore("test"); err != beanstalk.ErrNotIgnored {
		t.Fatal("expected not ignored error, but got", err)
	}

	count, err := c.Ignore("test")

	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Fatal("expected 1, but got", count)
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_Peek(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"peek 1\r\n",
				"peek 1\r\n",
			},
			[]string{
				"NOT_FOUND\r\n",
				"FOUND 1 4\r\ntest\r\n",
			},
		),
	)

	if _, err := c.Peek(1); err != beanstalk.ErrNotFound {
		t.Fatal("expected not found error, but got", err)
	}

	job, err := c.Peek(1)

	if err != nil {
		t.Error(err)
	}

	if job.ID != 1 {
		t.Fatal("expected 1, but got", job.ID)
	}

	if bytes.Compare([]byte("test"), job.Data) != 0 {
		t.Fatal(fmt.Sprintf("expected 'test', but got '%s'", string(job.Data)))
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_PeekReady(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"peek-ready\r\n",
				"peek-ready\r\n",
			},
			[]string{
				"NOT_FOUND\r\n",
				"FOUND 1 4\r\ntest\r\n",
			},
		),
	)

	if _, err := c.PeekReady(); err != beanstalk.ErrNotFound {
		t.Fatal("expected not found error, but got", err)
	}

	job, err := c.PeekReady()

	if err != nil {
		t.Error(err)
	}

	if job.ID != 1 {
		t.Fatal("expected 1, but got", job.ID)
	}

	if bytes.Compare([]byte("test"), job.Data) != 0 {
		t.Fatal(fmt.Sprintf("expected 'test', but got '%s'", string(job.Data)))
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_PeekDelayed(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"peek-delayed\r\n",
				"peek-delayed\r\n",
			},
			[]string{
				"NOT_FOUND\r\n",
				"FOUND 1 4\r\ntest\r\n",
			},
		),
	)

	if _, err := c.PeekDelayed(); err != beanstalk.ErrNotFound {
		t.Fatal("expected not found error, but got", err)
	}

	job, err := c.PeekDelayed()

	if err != nil {
		t.Error(err)
	}

	if job.ID != 1 {
		t.Fatal("expected 1, but got", job.ID)
	}

	if bytes.Compare([]byte("test"), job.Data) != 0 {
		t.Fatal(fmt.Sprintf("expected 'test', but got '%s'", string(job.Data)))
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_PeekBuried(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"peek-buried\r\n",
				"peek-buried\r\n",
			},
			[]string{
				"NOT_FOUND\r\n",
				"FOUND 1 4\r\ntest\r\n",
			},
		),
	)

	if _, err := c.PeekBuried(); err != beanstalk.ErrNotFound {
		t.Fatal("expected not found error, but got", err)
	}

	job, err := c.PeekBuried()

	if err != nil {
		t.Error(err)
	}

	if job.ID != 1 {
		t.Fatal("expected 1, but got", job.ID)
	}

	if bytes.Compare([]byte("test"), job.Data) != 0 {
		t.Fatal(fmt.Sprintf("expected 'test', but got '%s'", string(job.Data)))
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_Kick(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"kick 3\r\n",
			},
			[]string{
				"KICKED 3\r\n",
			},
		),
	)

	count, err := c.Kick(3)

	if err != nil {
		t.Error(err)
	}

	if count != 3 {
		t.Fatal("expected 3, but got", count)
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_KickJob(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"kick-job 1\r\n",
				"kick-job 1\r\n",
			},
			[]string{
				"NOT_FOUND\r\n",
				"KICKED\r\n",
			},
		),
	)

	if err := c.KickJob(1); err != beanstalk.ErrNotFound {
		t.Fatal("expected not found error, but got", err)
	}

	err := c.KickJob(1)

	if err != nil {
		t.Error(err)
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_StatsJob(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"stats-job 1\r\n",
				"stats-job 1\r\n",
			},
			[]string{
				"NOT_FOUND\r\n",
				"OK 144\r\n" +
					"---\n" +
					"id: 1\n" +
					"tube: default\n" +
					"state: ready\n" +
					"pri: 0\n" +
					"age: 12\n" +
					"delay: 0\n" +
					"ttr: 1\n" +
					"time-left: 0\n" +
					"file: 0\n" +
					"reserves: 0\n" +
					"timeouts: 0\n" +
					"releases: 0\n" +
					"buries: 0\n" +
					"kicks: 0\n" +
					"\r\n",
			},
		),
	)

	if _, err := c.StatsJob(1); err != beanstalk.ErrNotFound {
		t.Fatal("expected not found error, but got", err)
	}

	stats, err := c.StatsJob(1)

	if err != nil {
		t.Error(err)
	}

	if stats.ID != 1 {
		t.Fatal("expected 1, but got", stats.ID)
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_StatsTube(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"stats-tube test\r\n",
				"stats-tube default\r\n",
			},
			[]string{
				"NOT_FOUND\r\n",
				"OK 265\r\n" +
					"---\n" +
					"name: default\n" +
					"current-jobs-urgent: 0\n" +
					"current-jobs-ready: 0\n" +
					"current-jobs-reserved: 0\n" +
					"current-jobs-delayed: 0\n" +
					"current-jobs-buried: 0\n" +
					"total-jobs: 0\n" +
					"current-using: 3\n" +
					"current-watching: 3\n" +
					"current-waiting: 0\n" +
					"cmd-delete: 0\n" +
					"cmd-pause-tube: 0\n" +
					"pause: 0\n" +
					"pause-time-left: 0\n" +
					"\r\n",
			},
		),
	)

	if _, err := c.StatsTube("test"); err != beanstalk.ErrNotFound {
		t.Fatal("expected not found error, but got", err)
	}

	stats, err := c.StatsTube("default")

	if err != nil {
		t.Error(err)
	}

	if !strings.EqualFold("default", stats.Name) {
		t.Fatal(fmt.Sprintf("expected 'default', but got '%s'", stats.Name))
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_Stats(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"stats\r\n",
			},
			[]string{
				"OK 900\r\n" +
					"---\n" +
					"current-jobs-urgent: 0\n" +
					"current-jobs-ready: 0\n" +
					"current-jobs-reserved: 0\n" +
					"current-jobs-delayed: 0\n" +
					"current-jobs-buried: 0\n" +
					"cmd-put: 0\n" +
					"cmd-peek: 0\n" +
					"cmd-peek-ready: 0\n" +
					"cmd-peek-delayed: 0\n" +
					"cmd-peek-buried: 0\n" +
					"cmd-reserve: 0\n" +
					"cmd-reserve-with-timeout: 0\n" +
					"cmd-delete: 0\n" +
					"cmd-release: 0\n" +
					"cmd-use: 0\n" +
					"cmd-watch: 0\n" +
					"cmd-ignore: 0\n" +
					"cmd-bury: 0\n" +
					"cmd-kick: 0\n" +
					"cmd-touch: 0\n" +
					"cmd-stats: 1\n" +
					"cmd-stats-job: 0\n" +
					"cmd-stats-tube: 0\n" +
					"cmd-list-tubes: 0\n" +
					"cmd-list-tube-used: 0\n" +
					"cmd-list-tubes-watched: 0\n" +
					"cmd-pause-tube: 0\n" +
					"job-timeouts: 0\n" +
					"total-jobs: 0\n" +
					"max-job-size: 65535\n" +
					"current-tubes: 1\n" +
					"current-connections: 3\n" +
					"current-producers: 0\n" +
					"current-workers: 0\n" +
					"current-waiting: 0\n" +
					"total-connections: 3\n" +
					"pid: 1\n" +
					"version: 1.10\n" +
					"rusage-utime: 0.148125\n" +
					"rusage-stime: 0.014812\n" +
					"uptime: 1864\n" +
					"binlog-oldest-index: 0\n" +
					"binlog-current-index: 0\n" +
					"binlog-records-migrated: 0\n" +
					"binlog-records-written: 0\n" +
					"binlog-max-size: 10485760\n" +
					"id: f40521014b63360d\n" +
					"hostname: 671db3de0474\n" +
					"\r\n",
			},
		),
	)

	stats, err := c.Stats()

	if err != nil {
		t.Error(err)
	}

	if stats.PID != 1 {
		t.Fatal("expected 1, but got", stats.PID)
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_ListTubes(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"list-tubes\r\n",
			},
			[]string{
				"OK 21\r\n---\n- default\n- test\n\r\n",
			},
		),
	)

	tubes, err := c.ListTubes()

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual([]string{"default", "test"}, tubes) {
		t.Fatal(fmt.Sprintf("expected 'default, test', but got '%s'", strings.Join(tubes, ", ")))
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_ListTubeUsed(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"list-tube-used\r\n",
			},
			[]string{
				"USING test\r\n",
			},
		),
	)

	tube, err := c.ListTubeUsed()

	if err != nil {
		t.Error(err)
	}

	if !strings.EqualFold("test", tube) {
		t.Fatal(fmt.Sprintf("expected 'test', but got '%s'", tube))
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_ListTubesWatched(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"list-tubes-watched\r\n",
			},
			[]string{
				"OK 14\r\n---\n- default\n\r\n",
			},
		),
	)

	tubes, err := c.ListTubesWatched()

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual([]string{"default"}, tubes) {
		t.Fatal(fmt.Sprintf("expected 'default', but got '%s'", strings.Join(tubes, ", ")))
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}

func TestClient_PauseTube(t *testing.T) {
	c := beanstalk.NewClient(
		beanstalk.NewMockConn(
			[]string{
				"pause-tube test 60\r\n",
				"pause-tube test 0\r\n",
			},
			[]string{
				"NOT_FOUND\r\n",
				"PAUSED\r\n",
			},
		),
	)

	if err := c.PauseTube("test", 1*time.Minute); err != beanstalk.ErrNotFound {
		t.Fatal("expected not found error, but got", err)
	}

	err := c.PauseTube("test", 0)

	if err != nil {
		t.Error(err)
	}

	if err = c.Close(); err != nil {
		t.Error(err)
	}
}
