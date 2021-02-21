package service

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/vps2/accounttesttask/pkg/log"
)

//StatisticsSvc представляет сборщик статистики. Методы типа могу вызываться из разных горутин.
type StatisticsSvc struct {
	readOps        int64
	readOpsPerSec  int64
	writeOps       int64
	writeOpsPerSec int64
}

func NewStatisticsSvc(ctx context.Context, pollInterval time.Duration) *StatisticsSvc {
	statistics := &StatisticsSvc{
		readOps:        0,
		writeOps:       0,
		readOpsPerSec:  0,
		writeOpsPerSec: 0,
	}

	go func() {
		ticker := time.NewTicker(pollInterval)
	loop:
		for {
			readOpsAtBeginning := statistics.TotalReadOperations()
			writeOpsAtBeginning := statistics.TotalWriteOperations()

			select {
			case <-ticker.C:
				totalReadOps := statistics.TotalReadOperations()
				totalWriteOps := statistics.TotalWriteOperations()

				readOps := totalReadOps - readOpsAtBeginning
				writeOps := totalWriteOps - writeOpsAtBeginning

				if readOps < 0 || writeOps < 0 {
					continue
				}

				readOpsPerSec := readOps / int64(pollInterval.Seconds())
				writeOpsPerSec := writeOps / int64(pollInterval.Seconds())

				atomic.StoreInt64(&statistics.readOpsPerSec, readOpsPerSec)
				atomic.StoreInt64(&statistics.writeOpsPerSec, writeOpsPerSec)

				log.Infof("read operations per sec: %d, total read operations: %d, write operations per sec: %d, total write operations: %d\n",
					readOpsPerSec,
					totalReadOps,
					writeOpsPerSec,
					totalWriteOps)
			case <-ctx.Done():
				break loop
			}
		}

		ticker.Stop()
	}()

	return statistics
}

func (svc *StatisticsSvc) IncReadOperations() {
	atomic.AddInt64(&svc.readOps, 1)
}

func (svc *StatisticsSvc) IncWriteOperations() {
	atomic.AddInt64(&svc.writeOps, 1)
}

func (svc *StatisticsSvc) Reset() {
	atomic.StoreInt64(&svc.readOps, 0)
	atomic.StoreInt64(&svc.writeOps, 0)
}

func (svc *StatisticsSvc) TotalReadOperations() int64 {
	return atomic.LoadInt64(&svc.readOps)
}

func (svc *StatisticsSvc) TotalWriteOperations() int64 {
	return atomic.LoadInt64(&svc.writeOps)
}

func (svc *StatisticsSvc) ReadOperationsPerSecond() int64 {
	return atomic.LoadInt64(&svc.readOpsPerSec)
}

func (svc *StatisticsSvc) WriteOperationsPerSecond() int64 {
	return atomic.LoadInt64(&svc.writeOpsPerSec)
}
