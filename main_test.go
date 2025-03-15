package benchtune

import (
	"context"
	"fmt"
	"testing"
)

func TestSession(t *testing.T) {
	/*	param := WorkloadParameter{
			"request rate",
			"--rate %d",
			[]int{100, 1000, 10000},
		}

		metric := WorkloadMetric{
			"response rate",
			"Requests/sec:   %f",
			true, // higher better
		}

		workload := Workload{
			"HTTP load (wrk2)",
			"wrk -t16 -c64 -d20s -s wrk-urls.lua http://192.168.1.2/ --latency",
			[]WorkloadParameter{param},
			[]WorkloadMetric{metric},
		} */

	workload := Workload{
		"Test sleep",
		"sleep 0.3",
		[]WorkloadParameter{},
		[]WorkloadMetric{},
	}

	ctx, cancel := context.WithCancel(context.Background())
	// FIXME call cancel() when detecting ctrl-c

	err, stdin, stdout := NewShell(ctx)
	if err != nil {
		t.Fatal(err)
	}

	tuneables := []*Tuneable{
		{
			Name:    "TCP congestion control",
			Default: 0,
			Len:     3,
			Cmd: func(i int) string {
				items := []string{
					"cubic",
					"reno",
					"bbr",
				}

				//				return fmt.Sprintf("echo %s >/proc/sys/net/ipv4/tcp_congestion_control", items[i])
				return fmt.Sprintf("echo %s >/tmp/tcp_congestion_control", items[i])
			},
		},
	}

	session := NewSession(&workload, tuneables, stdin, stdout, stdin, stdout)

	err = session.SweepIndependent()
	if err != nil {
		t.Fatal(err)
	}

	session.StatsPrint()
	cancel()
}
