package main

import (
	"reflect"
	"testing"
	"time"
)

func TestNewLoadConfig(t *testing.T) {
	assert := func(t testing.TB, got LoadConfig, want LoadConfig) {
		t.Helper()

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Got: %v,\nWant: %v", got, want)
		}
	}

	t.Run("Ramp up load test", func(t *testing.T) {
		loadConfig := NewLoadConfig(10, 10, up, 2)

		desired := LoadConfig{
			Users:    10,
			Duration: time.Duration(10 * time.Second),
			Ramp: Ramp{
				RampMode:   up,
				RampAmount: 2,
			},
		}

		assert(t, loadConfig, desired)
	})

	t.Run("Ramp down load test", func(t *testing.T) {
		loadConfig := NewLoadConfig(20, 20, down, 4)

		desired := LoadConfig{
			Users:    20,
			Duration: time.Duration(20 * time.Second),
			Ramp: Ramp{
				RampMode:   down,
				RampAmount: 4,
			},
		}

		assert(t, loadConfig, desired)
	})

}

func TestRampTime(t *testing.T) {
	assert := func(t testing.TB, got int, want int) {
		t.Helper()

		if got != want {
			t.Errorf("Got: %v,\nWant: %v", got, want)
		}
	}

	t.Run("Exact division", func(t *testing.T) {
		loadConfig := NewLoadConfig(10, 10, up, 2)
		rampTime := loadConfig.RampTime()

		desired := 5

		assert(t, rampTime, desired)
	})

	t.Run("Inexact division", func(t *testing.T) {
		loadConfig := NewLoadConfig(10, 10, up, 3)
		rampTime := loadConfig.RampTime()

		desired := 4

		assert(t, rampTime, desired)
	})
}
