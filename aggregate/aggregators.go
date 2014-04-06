package main

import (
	"bisect"
	"math"
	"strconv"
)

type Aggregator interface {
	Aggregate(value []byte)
	String() string
}

// Sums
type Adder float64

func (a *Adder) Aggregate(value []byte) {
	var f, _ = strconv.ParseFloat(string(value), 64)
	*a += Adder(f)
}

func (a Adder) String() string {
	return strconv.FormatFloat(float64(a), 'g', 15, 64)
}

// Counters
type Counter uint64

func (c *Counter) Aggregate(value []byte) {
	*c += Counter(1)
}

func (c Counter) String() string {
	return strconv.FormatUint(uint64(c), 10)
}

// Minimums
type Miner struct {
	min  float64
	init bool
}

func (m *Miner) Aggregate(value []byte) {
	var f, _ = strconv.ParseFloat(string(value), 64)
	if !m.init {
		m.init = true
		m.min = f
	} else if f < m.min {
		m.min = f
	}
}

func (m Miner) String() string {
	return strconv.FormatFloat(m.min, 'g', 15, 64)
}

// Maximums
type Maxer struct {
	max  float64
	init bool
}

func (m *Maxer) Aggregate(value []byte) {
	var f, _ = strconv.ParseFloat(string(value), 64)
	if !m.init {
		m.init = true
		m.max = f
	} else if f > m.max {
		m.max = f
	}
}

func (m Maxer) String() string {
	return strconv.FormatFloat(m.max, 'g', 15, 64)
}

// Averages
type Averager struct {
	sum float64
	num uint64
}

func (a *Averager) Aggregate(value []byte) {
	var f, _ = strconv.ParseFloat(string(value), 64)
	a.sum += f
	a.num++
}

func (a Averager) String() string {
	if a.num > 0 {
		return strconv.FormatFloat(a.sum/float64(a.num), 'g', 15, 64)
	}
	return "0"
}

// Firsts
type Firster []byte

func (f *Firster) Aggregate(value []byte) {
	if *f == nil {
		*f = make([]byte, len(value))
		copy(*f, value)
	}
}

func (f Firster) String() string {
	return string(f)
}

// Lasts
type Laster []byte

func (l *Laster) Aggregate(value []byte) {
	if *l == nil || len(*l) < len(value) {
		*l = make([]byte, len(value))
	}
	copy(*l, value)
}

func (l Laster) String() string {
	return string(l)
}

// Concatenation
type Concater struct {
	buffer []byte
	Delim  []byte
}

func (c *Concater) Aggregate(value []byte) {
	if c.buffer != nil {
		c.buffer = append(c.buffer, c.Delim...)
	}
	c.buffer = append(c.buffer, value...)
}

func (c Concater) String() string {
	return string(c.buffer)
}

// Standard devs
type Stdever struct {
	sum  float64
	ssum float64
	num  uint64
}

func (s *Stdever) Aggregate(value []byte) {
	var f, _ = strconv.ParseFloat(string(value), 64)
	s.sum += f
	s.ssum += f * f
	s.num++
}

func (s Stdever) String() string {
	if s.num != 0.0 {
		avg := s.sum / float64(s.num)
		return strconv.FormatFloat(math.Sqrt(s.ssum/float64(s.num)-avg*avg), 'g', 15, 64)
	}
	return "0"
}

// Medianer
type mEl float64

func (m mEl) Less(other bisect.Elem) bool {
	return m < other.(mEl)
}

type Medianer struct {
	el []bisect.Elem
}

func (m *Medianer) Aggregate(value []byte) {
	var f, _ = strconv.ParseFloat(string(value), 64)
	m.el = bisect.Insort(m.el, mEl(f))
}

func (m *Medianer) String() string {
	return strconv.FormatFloat(m.Float64(), 'g', 15, 64)
}

func (m *Medianer) Float64() float64 {
	if len(m.el) > 0 {
		var med float64
		if len(m.el)%2 == 0 {
			med = (float64(m.el[len(m.el)/2].(mEl)) + float64(m.el[len(m.el)/2-1].(mEl))) / 2.0
		} else {
			med = float64(m.el[len(m.el)/2].(mEl))
		}
		return med
	}
	return 0.0
}

// MAD
type MADer struct {
	Medianer
}

func (m *MADer) Float64() float64 {
	if len(m.el) == 0 {
		return 0.0
	}
	med := m.Medianer.Float64()
	abs := make([]bisect.Elem, 0, len(m.el))
	for _, e := range m.el {
		abs = bisect.Insort(abs, mEl(math.Abs(float64(e.(mEl)) - med)))
	}
	var mad float64
	if len(abs)%2 == 0 {
		mad = (float64(abs[len(abs)/2].(mEl)) + float64(abs[len(abs)/2-1].(mEl))) / 2.0
	} else {
		mad = float64(abs[len(abs)/2].(mEl))
	}
	return mad
}

func (m *MADer) String() string {
	return strconv.FormatFloat(m.Float64(), 'g', 15, 64)
}
