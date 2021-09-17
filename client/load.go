package client

import (
	"fmt"
	"github.com/1xyz/hraftd-client/internal/tools"
	"log"
	"sync/atomic"
	"time"
)

type PerfCounter struct {
	nCount  int64
	nPassed int64
	nFailed int64
}

func NewPerfCounter() *PerfCounter { return &PerfCounter{nCount: 0, nPassed: 0, nFailed: 0} }
func (p *PerfCounter) String() string {
	return fmt.Sprintf("nCount = %d, nPassed = %d, nFailed = %d",
		p.nCount, p.nPassed, p.nFailed)
}
func (p *PerfCounter) Inc()       { atomic.AddInt64(&p.nCount, 1) }
func (p *PerfCounter) IncPassed() { atomic.AddInt64(&p.nPassed, 1) }
func (p *PerfCounter) IncFailed() { atomic.AddInt64(&p.nFailed, 1) }

type LoadTest struct {
	cli        *HttpClient
	nProducers int
	doneChans  []chan int
	keySize    int
	valueSize  int
	duration   time.Duration

	nPuts *PerfCounter
	nGets *PerfCounter
}

func NewLoadTest(cli *HttpClient, nProducers, keySize, valueSize int, duration time.Duration) *LoadTest {
	doneChans := make([]chan int, nProducers)
	for i := range doneChans {
		doneChans[i] = make(chan int)
	}

	return &LoadTest{
		cli:        cli,
		nProducers: nProducers,
		keySize:    keySize,
		valueSize:  valueSize,
		doneChans:  doneChans,
		duration:   duration,

		nPuts: NewPerfCounter(),
		nGets: NewPerfCounter(),
	}
}

func (l *LoadTest) Run() error {
	log.Printf("Run.. producers %d", l.nProducers)
	for i := 0; i < l.nProducers; i++ {
		log.Printf("runProducer[%d]", i)
		go func(i int) {
			log.Printf("start producer [%d]", i)
			if err := l.runProducer(i); err != nil {
				log.Printf("error = %v", err)
				return
			}
		}(i)
	}
	return nil
}

func (l *LoadTest) Stop() error {
	for i := 0; i < l.nProducers; i++ {
		go func(i int) {
			l.doneChans[i] <- 0
			close(l.doneChans[i])
		}(i)
	}
	return nil
}

func (l *LoadTest) runProducer(index int) error {
	for {
		key := tools.RandomAlphaNumeric(l.keySize)
		value := tools.RandomAlphaNumeric(l.valueSize)
		l.nPuts.Inc()
		if err := l.cli.Put(key, value); err != nil {
			log.Printf("Put err = %v", err)
			l.nPuts.IncFailed()
			return err
		}
		l.nPuts.IncPassed()

		l.nGets.Inc()
		_, err := l.cli.Get(key)
		if err != nil {
			l.nGets.IncFailed()
			log.Printf("Get err = %v", err)
			return err
		}
		l.nGets.IncPassed()

		select {
		case <-l.doneChans[index]:
			log.Printf("DoneChan signaled for index = %d", index)
			return nil
		default:
		}
	}
}

func RunLoadTest(cli *HttpClient, duration time.Duration) error {
	lt := NewLoadTest(cli, 20, 10, 10, duration)
	log.Printf("Started a new load test...")
	if err := lt.Run(); err != nil {
		log.Printf("RunLoadTest err = %v", err)
		return err
	}
	log.Printf("waiting......")
	time.Sleep(lt.duration)
	lt.Stop()
	log.Printf("Perf counter stats puts %v gets %v \n", lt.nPuts, lt.nGets)
	time.Sleep(100 * time.Millisecond)
	return nil
}
