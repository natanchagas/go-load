package main

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

type RampMode int

const (
	up RampMode = iota
	down
)

func (rm RampMode) String() string {

	switch rm {
	case 1:
		return "down"
	default:
		return "up"
	}
}

type Ramp struct {
	RampMode   RampMode
	RampAmount int
}

type Execution func()

type LoadConfig struct {
	Users    int
	Duration time.Duration
	Ramp     Ramp
}

func NewLoadConfig(users int, duration int, rampMode RampMode, rampAmount int) LoadConfig {

	if rampAmount > users {
		log.Fatalf("You can't ramp more users than you have available to test. Users: %d, Ramp Amout: %d\n", users, rampAmount)
	}

	return LoadConfig{
		Users:    users,
		Duration: time.Duration(duration) * time.Second,
		Ramp: Ramp{
			RampMode:   rampMode,
			RampAmount: rampAmount,
		},
	}
}

func (t LoadConfig) RampTime() int {

	rampTime := t.Users / t.Ramp.RampAmount

	if t.Users%t.Ramp.RampAmount == 0 {
		return rampTime
	} else {
		return rampTime + 1
	}
}

func (t LoadConfig) RampUpLoad(exec Execution) {
	ch := make(chan int, t.Users)
	go func(ch chan<- int) {
		for i := t.Ramp.RampAmount; ; i += t.Ramp.RampAmount {

			if i > t.Users {
				ch <- (i - t.Users)
				break
			} else {
				ch <- t.Ramp.RampAmount
			}
			time.Sleep(1 * time.Second)
		}
	}(ch)

	for {
		select {
		case v := <-ch:
			for i := 0; i < v; i++ {
				go func() {
					exec()
				}()
			}
		case <-time.After(t.Duration + time.Duration(t.RampTime())):
			log.Println("Load test has finished")
			return
		}
	}
}

func (t LoadConfig) RampUpLoadExplain(exec Execution) {
	ch := make(chan int, t.Users)
	go func(ch chan<- int) {
		for i := t.Ramp.RampAmount; ; i += t.Ramp.RampAmount {

			if i > t.Users {
				ch <- (i - t.Users)
				break
			} else {
				ch <- t.Ramp.RampAmount
			}
			time.Sleep(1 * time.Second)
		}
	}(ch)

	for {
		select {
		case v := <-ch:
			for i := 0; i < v; i++ {
				go func(i int, v int) { // Remove these parameters
					fmt.Println(i, v) // Run function here
				}(i, v) // Remove these parameters
			}
		case <-time.After(t.Duration + time.Duration(t.RampTime())):
			log.Println("Load test has finished")
			return
		}
	}
}

func (t LoadConfig) RampDownLoad(exec Execution) {
	ch := make(chan int, t.Users)
	time.Sleep(t.Duration)
	go func(ch chan<- int) {
		for i := t.Users; ; i -= t.Ramp.RampAmount {
			if i < 0 {
				i *= -1
				for j := 0; j < i; j++ {
					ch <- 1
				}
				break
			} else {
				for i := 0; i < t.Ramp.RampAmount; i++ {
					ch <- 1
				}
			}
			time.Sleep(1 * time.Second)
		}
	}(ch)

	for i := 0; ; i++ {
		if i < 3 {
			go func() {
				for {
					select {
					case <-ch:
						return
					case <-time.After(t.Duration + time.Duration(t.RampTime())):
						return
					default:
						exec()
					}
				}
			}()
		}

		if runtime.NumGoroutine() == 1 {
			log.Println("Load test has finished")
			return
		}
	}
}

func (t LoadConfig) RampDownLoadExplain(exec Execution) {
	ch := make(chan int, t.Users)
	time.Sleep(t.Duration)
	go func(ch chan<- int) {
		for i := t.Users; ; i -= t.Ramp.RampAmount {
			if i < 0 {
				i *= -1
				for j := 0; j < i; j++ {
					ch <- 1
				}
				break
			} else {
				for i := 0; i < t.Ramp.RampAmount; i++ {
					ch <- 1
				}

			}
			time.Sleep(time.Duration((int(t.Duration) / (int(t.Duration) / t.Ramp.RampAmount)) * int(time.Second)))
		}
	}(ch)

	for i := 0; ; i++ {
		if i < t.Users {
			go func(i int) { // Remove this parameters
				defer log.Println(i, "Dead") // Remove this log - Is outside of loop for educative reasons
				fmt.Println(i)               // Is outside of loop for educative reasons
				for {
					select {
					case <-ch:
						return
					case <-time.After(t.Duration + time.Duration(t.RampTime())):
						return
					default:
						continue // Run function here
					}
				}
			}(i) // Remove this parameters
		}

		if runtime.NumGoroutine() == 1 {
			log.Println("Load test has finished")
			return
		}
	}
}

func (t LoadConfig) RunLoad(exec Execution) {
	log.Println("Test will simulate", t.Users, "users over", t.Duration, "ramping", t.Ramp.RampMode.String(), t.Ramp.RampAmount, "users per second.")

	if t.Ramp.RampMode == up {
		t.RampUpLoad(exec)
	} else {
		t.RampDownLoad(exec)
	}

}
