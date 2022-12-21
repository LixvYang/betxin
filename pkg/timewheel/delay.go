package timewheel

import (
	"math/rand"
	"strconv"
	"time"
)

var (
	tw       = New(time.Second, 3600)
	everyMap = make(map[string]bool)
)

func init() {
	tw.Start()
}

// Delay executes job after waiting the given duration
func Delay(duration time.Duration, key string, job func()) {
	tw.AddJob(duration, key, job)
}

// At executes job at given time
func At(at time.Time, key string, job func()) {
	// tw.AddJob(at.Sub(time.Now()), key, job)
	tw.AddJob(time.Until(at), key, job)
}

// Every execfutes job at every time
func Every(every time.Duration, job func()) {
	randomInt := strconv.Itoa(int(rand.Float64()))
	// until randomInt is unique
	for everyMap[randomInt] {
		randomInt = strconv.Itoa(int(rand.Float64()))
	}

	everyMap[randomInt] = true

	t := time.NewTicker(every)
	defer t.Stop()
	for range t.C {
		At(time.Now().Add(every), randomInt, job)
	}
}

// Cancel stops a pending job
func Cancel(key string) {
	tw.RemoveJob(key)
}
