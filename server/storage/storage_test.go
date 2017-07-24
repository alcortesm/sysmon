package storage_test

import (
	"math"
	"testing"

	"github.com/alcortesm/sysmon/server/storage"
	"github.com/alcortesm/sysmon/stat"
)

func mustNew(t *testing.T, capacity int) storage.Storage {
	s, err := storage.NewOld(capacity)
	if err != nil {
		t.Fatalf("capacity %d: %s", capacity, err)
	}
	return s
}

func TestInitiallyEmpty(t *testing.T) {
	s := mustNew(t, 10)
	usage := s.CPUUsage()
	if len(usage) != 0 {
		t.Errorf("CPUUsage was not empty: %v", usage)
	}
}

func TestLengthOfCPUUsage(t *testing.T) {
	for _, test := range [...]struct {
		nStats    int
		nCPUUsage int
	}{
		{0, 0},
		{1, 0},
		{2, 1},
		{3, 2},
		{4, 3},
	} {
		s := mustNew(t, 10)
		st, err := stat.New()
		if err != nil {
			t.Fatalf("cannot read stats: %s", err)
		}
		for i := 0; i < test.nStats; i++ {
			s.Insert(st)
		}
		usage := s.CPUUsage()
		if len(usage) != test.nCPUUsage {
			t.Errorf("wrong length of CPUUsage, expected %d, got %d",
				test.nCPUUsage, len(usage))
		}
	}
}

type mockCPUer struct {
	total    int
	cpuTotal int
}

func (s *mockCPUer) Total() int    { return s.total }
func (s *mockCPUer) TotalCPU() int { return s.cpuTotal }

func TestCorrectCPUUsage(t *testing.T) {
	s := mustNew(t, 10)
	s.Insert(&mockCPUer{0, 0})
	s.Insert(&mockCPUer{100, 20})
	s.Insert(&mockCPUer{200, 60})
	obtained := s.CPUUsage()
	expected := []float64{40, 20}
	if !approx(obtained, expected, 0.001) {
		t.Errorf("wrong CPU usage:\nexpected: %v\nobtained: %v\n",
			expected, obtained)
	}
}

func approx(a, b []float64, tolerance float64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !float64Equal(a[i], b[i], tolerance) {
			return false
		}
	}
	return true
}

func TestItForgetsOutdated(t *testing.T) {
	s := mustNew(t, 2)
	s.Insert(&mockCPUer{0, 0})
	s.Insert(&mockCPUer{100, 20})
	s.Insert(&mockCPUer{200, 60})
	obtained := s.CPUUsage()
	expected := []float64{40}
	if !approx(obtained, expected, 0.001) {
		t.Errorf("wrong CPU usage:\nexpected: %v\nobtained: %v\n",
			expected, obtained)
	}
}

func float64Equal(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

func TestLongRun(t *testing.T) {
	s := mustNew(t, 4)
	for _, test := range [...]struct {
		input    stat.CPUer
		expected []float64
	}{
		{&mockCPUer{0, 0}, []float64{}},
		{&mockCPUer{100, 20}, []float64{20}},
		{&mockCPUer{200, 60}, []float64{40, 20}},
		{&mockCPUer{300, 70}, []float64{10, 40, 20}},
		{&mockCPUer{400, 100}, []float64{30, 10, 40}},
		{&mockCPUer{500, 200}, []float64{100, 30, 10}},
		{&mockCPUer{600, 200}, []float64{0, 100, 30}},
		{&mockCPUer{700, 220}, []float64{20, 0, 100}},
		{&mockCPUer{800, 240}, []float64{20, 20, 0}},
		{&mockCPUer{900, 250}, []float64{10, 20, 20}},
	} {
		s.Insert(test.input)
		obtained := s.CPUUsage()
		if !approx(obtained, test.expected, 0.001) {
			t.Errorf("wrong CPU usage:\nexpected: %v\nobtained: %v\n",
				test.expected, obtained)
		}
	}
}
