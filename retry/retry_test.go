package retry

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type MockSleeperInterface struct {
	CallCount    int
	LastDuration time.Duration
	Durations    []time.Duration
}

func (m *MockSleeperInterface) Sleep(d time.Duration) {
	m.CallCount++
	m.LastDuration = d
	m.Durations = append(m.Durations, d)
}

func TestRetry(t *testing.T) {
	t.Log("it does not retry if no error")
	{
		retryCnt := 0

		err := Times(2).Try(func(attempt uint) error {
			retryCnt++
			return nil
		})

		require.NoError(t, err)
		require.Equal(t, 1, retryCnt)
	}

	t.Log("it does retry if error")
	{
		attemptCnt := 0
		err := Times(2).Try(func(attempt uint) error {
			attemptCnt++
			return errors.New("error")
		})

		require.Error(t, err)
		require.Equal(t, "error", err.Error())
		require.Equal(t, 3, attemptCnt)
	}

	t.Log("it does not retry if Times=0")
	{
		attemptCnt := 0

		err := Times(0).Try(func(attempt uint) error {
			attemptCnt++
			return errors.New("error")
		})

		require.Error(t, err)
		require.Equal(t, "error", err.Error())
		require.Equal(t, 1, attemptCnt)
	}

	t.Log("it does a total attempt of 2 if Times=1")
	{
		attemptCnt := 0

		err := Times(1).Try(func(attempt uint) error {
			attemptCnt++
			return errors.New("error")
		})

		require.Error(t, err)
		require.Equal(t, "error", err.Error())
		require.Equal(t, 2, attemptCnt)
	}

	t.Log("it does a total attempt of 5 if Times=4")
	{
		attemptCnt := 0

		err := Times(4).Try(func(attempt uint) error {
			attemptCnt++
			return errors.New("error")
		})

		require.Error(t, err)
		require.Equal(t, "error", err.Error())
		require.Equal(t, 5, attemptCnt)
	}

	t.Log("it does not wait before first execution")
	{
		mockSleeper := &MockSleeperInterface{}
		attemptCnt := 0

		err := NewWithSleeper(1, 3*time.Second, mockSleeper).Try(func(attempt uint) error {
			attemptCnt++
			return errors.New("error")
		})

		require.Error(t, err)
		require.Equal(t, "error", err.Error())
		require.Equal(t, 2, attemptCnt)
		require.Equal(t, 1, mockSleeper.CallCount)
		require.Equal(t, 3*time.Second, mockSleeper.LastDuration)
	}

	t.Log("it waits before second execution with correct duration")
	{
		mockSleeper := &MockSleeperInterface{}
		attemptCnt := 0

		err := NewWithSleeper(1, 4*time.Second, mockSleeper).Try(func(attempt uint) error {
			attemptCnt++
			return errors.New("error")
		})

		require.Error(t, err)
		require.Equal(t, "error", err.Error())
		require.Equal(t, 2, attemptCnt)
		require.Equal(t, 1, mockSleeper.CallCount)
		require.Equal(t, 4*time.Second, mockSleeper.LastDuration)
	}

	t.Log("it stops retrying when abort indicates it")
	{
		type pair struct {
			error error
			abort bool
		}

		pairs := []pair{
			{errors.New("error-1"), false},
			{errors.New("error-2"), false},
			{errors.New("error-3"), true},
			{errors.New("error-4"), false},
		}
		attemptCnt := -1

		err := Times(3).TryWithAbort(func(attempt uint) (error, bool) {
			attemptCnt++
			pair := pairs[attempt]

			return pair.error, pair.abort
		})

		require.Error(t, err)
		require.Equal(t, "error-3", err.Error())
		require.Equal(t, 2, attemptCnt)
	}
}

func TestWait(t *testing.T) {
	t.Log("it creates retry model with wait time")
	{
		helper := Wait(3 * time.Second)
		require.Equal(t, 3*time.Second, helper.waitTime)
	}

	t.Log("it creates retry model with wait time")
	{
		helper := Wait(3 * time.Second)
		helper.Wait(5 * time.Second)
		require.Equal(t, 5*time.Second, helper.waitTime)
	}
}

func TestTimes(t *testing.T) {
	t.Log("it creates retry model with retry times")
	{
		helper := Times(3)
		require.Equal(t, uint(3), helper.retry)
	}

	t.Log("it sets retry times")
	{
		helper := Times(3)
		helper.Times(5)
		require.Equal(t, uint(5), helper.retry)
	}
}

func TestMockSleeperInterface(t *testing.T) {
	t.Log("it calls sleeper with correct duration and count")
	{
		mockSleeper := &MockSleeperInterface{}

		err := NewWithSleeper(2, 100*time.Millisecond, mockSleeper).Try(func(attempt uint) error {
			return errors.New("error")
		})

		require.Error(t, err)
		require.Equal(t, 2, mockSleeper.CallCount)
		require.Equal(t, 100*time.Millisecond, mockSleeper.LastDuration)
	}

	t.Log("it does not call sleeper if no error")
	{
		mockSleeper := &MockSleeperInterface{}

		err := NewWithSleeper(2, 100*time.Millisecond, mockSleeper).Try(func(attempt uint) error {
			return nil
		})

		require.NoError(t, err)
		require.Equal(t, 0, mockSleeper.CallCount)
	}

	t.Log("it does not call sleeper if wait time is zero")
	{
		mockSleeper := &MockSleeperInterface{}

		err := NewWithSleeper(2, 0, mockSleeper).Try(func(attempt uint) error {
			return errors.New("error")
		})

		require.Error(t, err)
		require.Equal(t, 0, mockSleeper.CallCount)
	}

	t.Log("it respects abort and does not call sleeper after abort")
	{
		mockSleeper := &MockSleeperInterface{}
		attemptCnt := 0

		err := NewWithSleeper(5, 100*time.Millisecond, mockSleeper).TryWithAbort(func(attempt uint) (error, bool) {
			attemptCnt++
			if attempt == 1 {
				return errors.New("abort-error"), true
			}
			return errors.New("error"), false
		})

		require.Error(t, err)
		require.Equal(t, "abort-error", err.Error())
		require.Equal(t, 2, attemptCnt)
		require.Equal(t, 1, mockSleeper.CallCount)
	}
}
