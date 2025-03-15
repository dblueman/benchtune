package benchtune

import (
	"errors"
	"fmt"
	"io"
)

// SystemMetric
// cpuload system time, user time, temperature, fan

type WorkloadParameter struct {
	Name   string
	Format string
	Values []int
}

type WorkloadMetric struct {
	Name     string
	Format   string
	Positive bool
}

// predefine RuntimeMetric and ExitcodeMetric

type Workload struct {
	Name       string
	Command    string
	Parameters []WorkloadParameter
	Metrics    []WorkloadMetric
}

type Tuneable struct {
	Name    string
	Cmd     func(int) string // generate command to apply setting
	Default int              // index into entries
	Len     int              // number of settings, high if linear
}

type Session struct {
	workload       *Workload
	tuneables      []*Tuneable
	tuneableStdin  io.WriteCloser
	tuneableStdout io.ReadCloser
	workloadStdin  io.WriteCloser
	workloadStdout io.ReadCloser
}

func NewSession(workload *Workload, tuneables []*Tuneable, tuneableStdin io.WriteCloser, tuneableStdout io.ReadCloser, workloadStdin io.WriteCloser, workloadStdout io.ReadCloser) *Session {
	return &Session{
		workload:       workload,
		tuneables:      tuneables,
		tuneableStdin:  tuneableStdin,
		tuneableStdout: tuneableStdout,
		workloadStdin:  workloadStdin,
		workloadStdout: workloadStdout,
	}
}

func (s *Session) Apply(args string) error {
	fmt.Printf("[%s]", args)

	_, err := s.tuneableStdin.Write([]byte(args + "\n"))
	if err != nil {
		return fmt.Errorf("Apply: %w", err)
	}

	buf := make([]byte, 1024)
	n, err := s.tuneableStdout.Read(buf)
	if err != nil {
		return fmt.Errorf("Apply: %w", err)
	}

	if string(buf[:n]) != "$ " {
		return fmt.Errorf("Apply: unexpected output [%s]", buf[:n])
	}
	fmt.Println()

	return nil
}

func (s *Session) Benchmark(args string) error {
	fmt.Printf("[%s]", args)

	_, err := s.workloadStdin.Write([]byte(args + "\n"))
	if err != nil {
		return fmt.Errorf("Apply: %w", err)
	}

	buf := make([]byte, 1024)
	n, err := s.workloadStdout.Read(buf)
	if err != nil {
		return fmt.Errorf("Apply: %w", err)
	}

	if string(buf[:n]) != "$ " {
		return fmt.Errorf("Apply: unexpected output [%s]", buf[:n])
	}
	fmt.Println()

	return nil
}

func (s *Session) SweepIndependent() error {
	if len(s.tuneables) == 0 {
		return errors.New("no tunables to sweep")
	}

	count := 5

	for _, tuneable := range s.tuneables {
		for i := range count {
			cmd := tuneable.Cmd(i * tuneable.Len / count)
			err := s.Apply(cmd)
			if err != nil {
				return fmt.Errorf("SweepIndependent: %w", err)
			}

			err = s.Benchmark(s.workload.Command)
			if err != nil {
				return fmt.Errorf("SweepIndependent: %w", err)
			}
		}
	}

	return nil
}

func (s *Session) StatsPrint() {
	// FIXME implement
}
