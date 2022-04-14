package counter

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

// cutTime Time for which need to find requests
const cutTime = 60 * time.Second

type Counter struct {
	f     *os.File
	dates []time.Time
	mu    *sync.Mutex
}

func MustCounter(ctx context.Context, path string) *Counter {
	f, err := os.OpenFile(path, os.O_WRONLY, 0755)
	if err != nil {
		log.Fatalln(err)
	}

	go func(ctx context.Context, f *os.File) {
		select {
		case <-ctx.Done():
			err := f.Close()
			if err != nil {
				log.Println(err)
			}
		}
	}(ctx, f)

	return &Counter{
		f:     f,
		dates: fillDates(path),
		mu:    &sync.Mutex{},
	}
}

// CountRequests Count correct requests and delete expired
func (c Counter) CountRequests(requestTime time.Time) (int, error) {
	count := len(c.dates)

	if len(c.dates) == 0 {
		c.mu.Lock()
		c.dates = append(c.dates, requestTime)
		c.mu.Unlock()
	}

	for i := len(c.dates) - 1; i >= 0; i-- {
		if requestTime.Sub(c.dates[i]) > cutTime {
			c.mu.Lock()

			c.dates = c.dates[i:]
			count -= i + 1
			c.dates = append(c.dates, requestTime)

			c.mu.Unlock()
			break
		}
	}

	err := c.updateFile()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// updateFile Truncate file and write new struct
func (c Counter) updateFile() error {
	err := c.f.Truncate(0)
	if err != nil {
		return fmt.Errorf("failed to Truncate: %v", err)
	}

	data, err := json.Marshal(&c.dates)
	if err != nil {
		return fmt.Errorf("failed to Marshal: %v", err)
	}

	_, err = c.f.WriteAt(data, 0)
	if err != nil {
		return fmt.Errorf("failed to WriteAt: %v", err)
	}

	return nil
}

// fillDates Read dates from file and fill local struct
func fillDates(path string) []time.Time {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	dates := make([]time.Time, 0, 1000)
	err = json.Unmarshal(data, &dates)
	if err != nil {
		log.Println(err)
		return dates
	}

	return dates
}
