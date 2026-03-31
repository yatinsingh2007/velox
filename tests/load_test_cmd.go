package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type TestCase struct {
	TestCaseID     int    `json:"test_case_id"`
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
}

type TestCaseResult struct {
	TestCaseID int    `json:"test_case_id"`
	Status     string `json:"status"`
	TimeMs     int64  `json:"time_ms"`
	MemoryKb   int64  `json:"memory_kb"`
}

type SubmissionRequest struct {
	Language      string     `json:"language"`
	SourceCode    string     `json:"source_code"`
	TestCases     []TestCase `json:"test_cases"`
	TimeLimitMs   int        `json:"time_limit_ms"`
	MemoryLimitKb int        `json:"memory_limit_kb"`
}

type SubmissionResponse struct {
	SubmissionID string           `json:"submission_id"`
	Status       string           `json:"status,omitempty"`
	OverallState string           `json:"overall_state,omitempty"`
	Results      []TestCaseResult `json:"results,omitempty"`
}

type RequestMetrics struct {
	ID                int
	SubmitDuration    time.Duration
	TotalWaitDuration time.Duration
	InternalExecTime  time.Duration
	OverheadTime      time.Duration
	PollCount         int
	TotalPollTime     time.Duration
	StatusCode        int
	Success           bool
	ErrorMessage      string
	SubmissionID      string
}

const (
	baseURL      = "http://localhost:8080"
	count        = 100
	pollInterval = 500 * time.Millisecond
)

func main() {
	payload := SubmissionRequest{
		Language: "cpp",
		SourceCode: `#include <iostream>
#include <vector>
using namespace std;

long long fibonacci(int n) {
    if (n <= 1) return n;
    vector<long long> fib(n + 1);
    fib[0] = 0;
    fib[1] = 1;
    for (int i = 2; i <= n; i++) {
        fib[i] = fib[i-1] + fib[i-2];
    }
    return fib[n];
}

int main() {
    int n;
    while(cin >> n) {
        cout << fibonacci(n) << endl;
    }
    return 0;
}`,
		TestCases: []TestCase{
			{TestCaseID: 1, Input: "10", ExpectedOutput: "55"},
			{TestCaseID: 2, Input: "20", ExpectedOutput: "6765"},
			{TestCaseID: 3, Input: "15", ExpectedOutput: "610"},
			{TestCaseID: 4, Input: "25", ExpectedOutput: "75025"},
			{TestCaseID: 5, Input: "30", ExpectedOutput: "832040"},
			{TestCaseID: 6, Input: "5", ExpectedOutput: "5"},
			{TestCaseID: 7, Input: "35", ExpectedOutput: "9227465"},
			{TestCaseID: 8, Input: "12", ExpectedOutput: "144"},
			{TestCaseID: 9, Input: "18", ExpectedOutput: "2584"},
			{TestCaseID: 10, Input: "22", ExpectedOutput: "17711"},
		},
		TimeLimitMs:   5000,
		MemoryLimitKb: 262144,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	fmt.Printf("╔════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║           CONCURRENT LOAD TEST WITH POLLING                    ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════════════╝\n\n")
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Total Requests:     %d\n", count)
	fmt.Printf("  Test Cases:         %d\n", len(payload.TestCases))
	fmt.Printf("  Poll Interval:      %v\n", pollInterval)
	fmt.Printf("  Code Complexity:    Fibonacci (Dynamic Programming)\n")
	fmt.Printf("  Time Limit:         %dms\n", payload.TimeLimitMs)
	fmt.Printf("  Memory Limit:       %d KB\n\n", payload.MemoryLimitKb)

	var wg sync.WaitGroup
	var successCount int64
	var failureCount int64

	metrics := make([]RequestMetrics, count)
	var metricsLock sync.Mutex

	tr := &http.Transport{
		MaxIdleConns:        count,
		MaxIdleConnsPerHost: count,
		IdleConnTimeout:     90 * time.Second,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   60 * time.Second,
	}

	startTime := time.Now()
	fmt.Printf("Starting load test at %s...\n\n", startTime.Format("15:04:05"))

	for i := range count {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			metric := RequestMetrics{ID: id + 1}
			var subID string

			// 1. Submit
			submitStart := time.Now()
			resp, err := client.Post(baseURL+"/submit", "application/json", bytes.NewBuffer(jsonData))
			metric.SubmitDuration = time.Since(submitStart)

			if err != nil {
				atomic.AddInt64(&failureCount, 1)
				metric.Success = false
				metric.ErrorMessage = "Submit error: " + err.Error()
				metricsLock.Lock()
				metrics[id] = metric
				metricsLock.Unlock()
				return
			}

			// Accept both 200 OK and 202 Accepted for async submissions
			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
				atomic.AddInt64(&failureCount, 1)
				metric.Success = false
				metric.StatusCode = resp.StatusCode
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				metric.ErrorMessage = fmt.Sprintf("Submit failed with status %d: %s", resp.StatusCode, string(body))
				metricsLock.Lock()
				metrics[id] = metric
				metricsLock.Unlock()
				return
			}

			var subResp struct {
				SubmissionID string `json:"submission_id"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&subResp); err != nil {
				atomic.AddInt64(&failureCount, 1)
				metric.Success = false
				metric.ErrorMessage = "Failed to decode submit response: " + err.Error()
				resp.Body.Close()
				metricsLock.Lock()
				metrics[id] = metric
				metricsLock.Unlock()
				return
			}
			resp.Body.Close()
			subID = subResp.SubmissionID
			metric.SubmissionID = subID

			// 2. Poll until completion
			maxPolls := 120 // Safety limit: 120 polls * 500ms = 60 seconds max
			for metric.PollCount < maxPolls {
				time.Sleep(pollInterval)
				metric.PollCount++

				pollReqStart := time.Now()
				statusResp, err := client.Get(fmt.Sprintf("%s/status?submission_id=%s", baseURL, subID))
				pollReqDuration := time.Since(pollReqStart)
				metric.TotalPollTime += pollReqDuration

				if err != nil {
					atomic.AddInt64(&failureCount, 1)
					metric.Success = false
					metric.ErrorMessage = "Poll error: " + err.Error()
					metricsLock.Lock()
					metrics[id] = metric
					metricsLock.Unlock()
					return
				}

				if statusResp.StatusCode != http.StatusOK {
					atomic.AddInt64(&failureCount, 1)
					metric.Success = false
					metric.StatusCode = statusResp.StatusCode
					body, _ := io.ReadAll(statusResp.Body)
					statusResp.Body.Close()
					metric.ErrorMessage = fmt.Sprintf("Poll failed with status %d: %s", statusResp.StatusCode, string(body))
					metricsLock.Lock()
					metrics[id] = metric
					metricsLock.Unlock()
					return
				}

				var statusData SubmissionResponse
				if err := json.NewDecoder(statusResp.Body).Decode(&statusData); err != nil {
					atomic.AddInt64(&failureCount, 1)
					metric.Success = false
					metric.ErrorMessage = "Poll decode error: " + err.Error()
					statusResp.Body.Close()
					metricsLock.Lock()
					metrics[id] = metric
					metricsLock.Unlock()
					return
				}
				statusResp.Body.Close()

				// Check if still pending
				if statusData.Status == "pending" {
					continue
				}

				// Finished successfully!
				metric.TotalWaitDuration = time.Since(submitStart)
				metric.Success = true
				metric.StatusCode = 200

				var totalExecMs int64
				for _, res := range statusData.Results {
					totalExecMs += res.TimeMs
				}
				metric.InternalExecTime = time.Duration(totalExecMs) * time.Millisecond
				metric.OverheadTime = metric.TotalWaitDuration - metric.InternalExecTime

				atomic.AddInt64(&successCount, 1)
				metricsLock.Lock()
				metrics[id] = metric
				metricsLock.Unlock()
				return
			}

			// If we exit the loop, we hit the max poll limit
			atomic.AddInt64(&failureCount, 1)
			metric.Success = false
			metric.ErrorMessage = fmt.Sprintf("Exceeded max polls (%d), submission may still be pending", maxPolls)
			metricsLock.Lock()
			metrics[id] = metric
			metricsLock.Unlock()
		}(i)
	}

	wg.Wait()
	totalWallTime := time.Since(startTime)

	// Calculate comprehensive statistics
	printDetailedResults(metrics, successCount, failureCount, totalWallTime)
}

func printDetailedResults(metrics []RequestMetrics, successCount, failureCount int64, totalWallTime time.Duration) {
	var (
		submitDurations    []time.Duration
		totalWaitDurations []time.Duration
		internalExecTimes  []time.Duration
		overheadDurations  []time.Duration
		pollCounts         []int
		avgPollTimes       []time.Duration
		totalPollTimes     []time.Duration
	)

	for _, m := range metrics {
		if m.Success {
			submitDurations = append(submitDurations, m.SubmitDuration)
			totalWaitDurations = append(totalWaitDurations, m.TotalWaitDuration)
			internalExecTimes = append(internalExecTimes, m.InternalExecTime)
			overheadDurations = append(overheadDurations, m.OverheadTime)
			pollCounts = append(pollCounts, m.PollCount)
			totalPollTimes = append(totalPollTimes, m.TotalPollTime)
			if m.PollCount > 0 {
				avgPollTimes = append(avgPollTimes, m.TotalPollTime/time.Duration(m.PollCount))
			}
		}
	}

	sort.Slice(submitDurations, func(i, j int) bool { return submitDurations[i] < submitDurations[j] })
	sort.Slice(totalWaitDurations, func(i, j int) bool { return totalWaitDurations[i] < totalWaitDurations[j] })
	sort.Slice(internalExecTimes, func(i, j int) bool { return internalExecTimes[i] < internalExecTimes[j] })
	sort.Slice(overheadDurations, func(i, j int) bool { return overheadDurations[i] < overheadDurations[j] })
	sort.Slice(pollCounts, func(i, j int) bool { return pollCounts[i] < pollCounts[j] })
	sort.Slice(totalPollTimes, func(i, j int) bool { return totalPollTimes[i] < totalPollTimes[j] })
	sort.Slice(avgPollTimes, func(i, j int) bool { return avgPollTimes[i] < avgPollTimes[j] })

	fmt.Printf("╔════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                       TEST RESULTS                             ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════════════╝\n\n")

	// Overall Statistics
	fmt.Printf("┌─ Overall Statistics ────────────────────────────────────────┐\n")
	fmt.Printf("│ Total Requests:         %-6d                               │\n", count)
	fmt.Printf("│ Successful:             %-6d (%.2f%%)                     │\n",
		successCount, float64(successCount)/float64(count)*100)
	fmt.Printf("│ Failed:                 %-6d (%.2f%%)                     │\n",
		failureCount, float64(failureCount)/float64(count)*100)
	fmt.Printf("│ Total Wall Time:        %-40s │\n", totalWallTime)
	fmt.Printf("│ Throughput:             %.2f req/sec                        │\n",
		float64(count)/totalWallTime.Seconds())
	fmt.Printf("└──────────────────────────────────────────────────────────────┘\n\n")

	if len(submitDurations) == 0 {
		fmt.Printf("⚠️  No successful requests to analyze.\n\n")

		// Show error breakdown
		if failureCount > 0 {
			fmt.Printf("┌─ Error Analysis ─────────────────────────────────────────────┐\n")
			errorTypes := make(map[string]int)
			var errorSamples []string

			for _, m := range metrics {
				if !m.Success && m.ErrorMessage != "" {
					// Categorize errors
					errKey := m.ErrorMessage
					if len(errKey) > 50 {
						errKey = errKey[:50]
					}
					errorTypes[errKey]++

					if len(errorSamples) < 3 {
						errorSamples = append(errorSamples, m.ErrorMessage)
					}
				}
			}

			fmt.Printf("│ Error Types:                                                 │\n")
			for errType, count := range errorTypes {
				if len(errType) > 54 {
					errType = errType[:51] + "..."
				}
				fmt.Printf("│  • %-54s: %2d │\n", errType, count)
			}

			if len(errorSamples) > 0 {
				fmt.Printf("│                                                              │\n")
				fmt.Printf("│ Sample Errors:                                               │\n")
				for i, msg := range errorSamples {
					if len(msg) > 56 {
						msg = msg[:53] + "..."
					}
					fmt.Printf("│  %d. %-58s │\n", i+1, msg)
				}
			}
			fmt.Printf("└──────────────────────────────────────────────────────────────┘\n")
		}
		return
	}

	// Submit Phase Statistics
	fmt.Printf("┌─ Submit Phase (Initial POST /submit) ───────────────────────┐\n")
	printTimingStats("│ Submit Time", submitDurations)
	fmt.Printf("└──────────────────────────────────────────────────────────────┘\n\n")

	// Total Wait Time (Submit + Poll until completion)
	fmt.Printf("┌─ Total Wait Time (Submit → Final Status) ───────────────────┐\n")
	printTimingStats("│ Total Wait", totalWaitDurations)
	fmt.Printf("└──────────────────────────────────────────────────────────────┘\n\n")

	// Internal Execution Time
	fmt.Printf("┌─ Internal Execution Time (Actual code execution) ───────────┐\n")
	printTimingStats("│ Exec Time", internalExecTimes)
	fmt.Printf("└──────────────────────────────────────────────────────────────┘\n\n")

	// Overhead Time
	fmt.Printf("┌─ System Overhead (Total Wait - Execution Time) ─────────────┐\n")
	printTimingStats("│ Overhead", overheadDurations)
	avgOverhead := calculateAvg(overheadDurations)
	avgTotal := calculateAvg(totalWaitDurations)
	if avgTotal > 0 {
		overheadPercent := float64(avgOverhead.Nanoseconds()) / float64(avgTotal.Nanoseconds()) * 100
		fmt.Printf("│ Overhead %% of Total:   %.2f%%                              │\n", overheadPercent)
	}
	fmt.Printf("└──────────────────────────────────────────────────────────────┘\n\n")

	// Polling Statistics
	fmt.Printf("┌─ Polling Statistics ─────────────────────────────────────────┐\n")
	printIntStats("│ Poll Count", pollCounts)
	fmt.Printf("│                                                              │\n")
	printTimingStats("│ Total Poll Time", totalPollTimes)
	if len(avgPollTimes) > 0 {
		fmt.Printf("│                                                              │\n")
		printTimingStats("│ Avg Time/Poll", avgPollTimes)
	}
	fmt.Printf("└──────────────────────────────────────────────────────────────┘\n\n")

	// Top Slowest Requests
	fmt.Printf("┌─ Top 10 Slowest Requests (Total Wait Time) ─────────────────┐\n")
	slowestCount := min(10, len(totalWaitDurations))
	for i := len(totalWaitDurations) - 1; i >= len(totalWaitDurations)-slowestCount; i-- {
		rank := len(totalWaitDurations) - i
		fmt.Printf("│  #%-2d  %52s │\n", rank, totalWaitDurations[i])
	}
	fmt.Printf("└──────────────────────────────────────────────────────────────┘\n\n")

	// Top Fastest Requests
	fmt.Printf("┌─ Top 10 Fastest Requests (Total Wait Time) ─────────────────┐\n")
	fastestCount := min(10, len(totalWaitDurations))
	for i := 0; i < fastestCount; i++ {
		fmt.Printf("│  #%-2d  %52s │\n", i+1, totalWaitDurations[i])
	}
	fmt.Printf("└──────────────────────────────────────────────────────────────┘\n\n")

	// Error Breakdown
	if failureCount > 0 {
		fmt.Printf("┌─ Error Breakdown ────────────────────────────────────────────┐\n")
		errorTypes := make(map[string]int)
		var errorExamples []string

		for _, m := range metrics {
			if !m.Success {
				errType := "Unknown Error"
				if m.ErrorMessage != "" {
					if len(m.ErrorMessage) > 30 {
						errType = m.ErrorMessage[:30]
					} else {
						errType = m.ErrorMessage
					}
				} else if m.StatusCode > 0 {
					errType = fmt.Sprintf("HTTP %d", m.StatusCode)
				}
				errorTypes[errType]++
				if len(errorExamples) < 5 {
					errorExamples = append(errorExamples, m.ErrorMessage)
				}
			}
		}

		for errType, eCount := range errorTypes {
			if len(errType) > 50 {
				errType = errType[:47] + "..."
			}
			fmt.Printf("│  %-50s: %4d │\n", errType, eCount)
		}

		if len(errorExamples) > 0 {
			fmt.Printf("│                                                              │\n")
			fmt.Printf("│  Sample Error Messages:                                      │\n")
			for i, msg := range errorExamples {
				if len(msg) > 56 {
					msg = msg[:53] + "..."
				}
				fmt.Printf("│  %d. %-58s │\n", i+1, msg)
			}
		}
		fmt.Printf("└──────────────────────────────────────────────────────────────┘\n\n")
	}

	// Performance Summary
	fmt.Printf("┌─ Performance Summary ────────────────────────────────────────┐\n")
	fmt.Printf("│  Average End-to-End:    %-37s │\n", calculateAvg(totalWaitDurations))
	fmt.Printf("│  Average Execution:     %-37s │\n", calculateAvg(internalExecTimes))
	fmt.Printf("│  Average Overhead:      %-37s │\n", calculateAvg(overheadDurations))
	fmt.Printf("│  Average Polls:         %-37.2f │\n", calculateAvgInt(pollCounts))
	fmt.Printf("│                                                              │\n")

	if len(totalWaitDurations) > 0 {
		successRate := float64(len(totalWaitDurations)) / float64(count) * 100
		fmt.Printf("│  Success Rate:          %.2f%%                               │\n", successRate)
	}

	fmt.Printf("└──────────────────────────────────────────────────────────────┘\n")
}

func printTimingStats(label string, durations []time.Duration) {
	if len(durations) == 0 {
		return
	}

	min := durations[0]
	max := durations[len(durations)-1]
	avg := calculateAvg(durations)
	p50 := durations[len(durations)*50/100]
	p90 := durations[len(durations)*90/100]
	p95 := durations[len(durations)*95/100]
	p99 := durations[len(durations)*99/100]

	fmt.Printf("│  %-15s Min: %-10s  Avg: %-10s  Max: %-10s│\n",
		label+":", min, avg, max)
	fmt.Printf("│  %-15s P50: %-10s  P90: %-10s  P95: %-10s│\n",
		"", p50, p90, p95)
	fmt.Printf("│  %-15s P99: %-10s                              │\n",
		"", p99)
}

func printIntStats(label string, values []int) {
	if len(values) == 0 {
		return
	}

	min := values[0]
	max := values[len(values)-1]
	avg := calculateAvgInt(values)
	p50 := values[len(values)*50/100]
	p90 := values[len(values)*90/100]
	p95 := values[len(values)*95/100]
	p99 := values[len(values)*99/100]

	fmt.Printf("│  %-15s Min: %-4d  Avg: %-6.2f  Max: %-4d  Median: %-4d │\n",
		label+":", min, avg, max, p50)
	fmt.Printf("│  %-15s P90: %-4d  P95: %-4d      P99: %-4d            │\n",
		"", p90, p95, p99)
}

func calculateAvg(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	var total int64
	for _, d := range durations {
		total += d.Nanoseconds()
	}
	return time.Duration(total / int64(len(durations)))
}

func calculateAvgInt(values []int) float64 {
	if len(values) == 0 {
		return 0
	}
	var total int
	for _, v := range values {
		total += v
	}
	return float64(total) / float64(len(values))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}