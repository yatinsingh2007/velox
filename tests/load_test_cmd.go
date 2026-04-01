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
	Language          string
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
	requestsPerLanguage = 20 // Number of concurrent requests per language
	pollInterval = 500 * time.Millisecond
)

// Language configurations with appropriate test programs
var languageConfigs = []struct {
	Language   string
	SourceCode string
	TestCases  []TestCase
}{
	{
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
		},
	},
	{
		Language: "c",
		SourceCode: `#include <stdio.h>
#include <stdlib.h>

long long fibonacci(int n) {
    if (n <= 1) return n;
    long long *fib = malloc((n + 1) * sizeof(long long));
    fib[0] = 0;
    fib[1] = 1;
    for (int i = 2; i <= n; i++) {
        fib[i] = fib[i-1] + fib[i-2];
    }
    long long result = fib[n];
    free(fib);
    return result;
}

int main() {
    int n;
    while(scanf("%d", &n) != EOF) {
        printf("%lld\n", fibonacci(n));
    }
    return 0;
}`,
		TestCases: []TestCase{
			{TestCaseID: 1, Input: "10", ExpectedOutput: "55"},
			{TestCaseID: 2, Input: "20", ExpectedOutput: "6765"},
			{TestCaseID: 3, Input: "15", ExpectedOutput: "610"},
		},
	},
	{
		Language: "java",
		SourceCode: `import java.util.Scanner;

public class Main {
    public static long fibonacci(int n) {
        if (n <= 1) return n;
        long[] fib = new long[n + 1];
        fib[0] = 0;
        fib[1] = 1;
        for (int i = 2; i <= n; i++) {
            fib[i] = fib[i-1] + fib[i-2];
        }
        return fib[n];
    }
    
    public static void main(String[] args) {
        Scanner sc = new Scanner(System.in);
        while(sc.hasNextInt()) {
            int n = sc.nextInt();
            System.out.println(fibonacci(n));
        }
        sc.close();
    }
}`,
		TestCases: []TestCase{
			{TestCaseID: 1, Input: "10", ExpectedOutput: "55"},
			{TestCaseID: 2, Input: "20", ExpectedOutput: "6765"},
			{TestCaseID: 3, Input: "15", ExpectedOutput: "610"},
		},
	},
	{
		Language: "python",
		SourceCode: `def fibonacci(n):
    if n <= 1:
        return n
    fib = [0] * (n + 1)
    fib[1] = 1
    for i in range(2, n + 1):
        fib[i] = fib[i-1] + fib[i-2]
    return fib[n]

import sys
for line in sys.stdin:
    n = int(line.strip())
    print(fibonacci(n))`,
		TestCases: []TestCase{
			{TestCaseID: 1, Input: "10", ExpectedOutput: "55"},
			{TestCaseID: 2, Input: "20", ExpectedOutput: "6765"},
			{TestCaseID: 3, Input: "15", ExpectedOutput: "610"},
		},
	},
	{
		Language: "node",
		SourceCode: `const readline = require('readline');

function fibonacci(n) {
    if (n <= 1) return n;
    const fib = new Array(n + 1);
    fib[0] = 0;
    fib[1] = 1;
    for (let i = 2; i <= n; i++) {
        fib[i] = fib[i-1] + fib[i-2];
    }
    return fib[n];
}

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout,
    terminal: false
});

rl.on('line', (line) => {
    const n = parseInt(line);
    console.log(fibonacci(n));
});`,
		TestCases: []TestCase{
			{TestCaseID: 1, Input: "10", ExpectedOutput: "55"},
			{TestCaseID: 2, Input: "20", ExpectedOutput: "6765"},
			{TestCaseID: 3, Input: "15", ExpectedOutput: "610"},
		},
	},
	{
		Language: "csharp",
		SourceCode: `using System;

class Program
{
    static long Fibonacci(int n)
    {
        if (n <= 1) return n;
        long[] fib = new long[n + 1];
        fib[0] = 0;
        fib[1] = 1;
        for (int i = 2; i <= n; i++)
        {
            fib[i] = fib[i-1] + fib[i-2];
        }
        return fib[n];
    }
    
    static void Main()
    {
        string line;
        while ((line = Console.ReadLine()) != null)
        {
            int n = int.Parse(line);
            Console.WriteLine(Fibonacci(n));
        }
    }
}`,
		TestCases: []TestCase{
			{TestCaseID: 1, Input: "10", ExpectedOutput: "55"},
			{TestCaseID: 2, Input: "20", ExpectedOutput: "6765"},
			{TestCaseID: 3, Input: "15", ExpectedOutput: "610"},
		},
	},
	{
		Language: "ts",
		SourceCode: `import * as readline from 'readline';

function fibonacci(n: number): number {
    if (n <= 1) return n;
    const fib: number[] = new Array(n + 1);
    fib[0] = 0;
    fib[1] = 1;
    for (let i = 2; i <= n; i++) {
        fib[i] = fib[i-1] + fib[i-2];
    }
    return fib[n];
}

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout,
    terminal: false
});

rl.on('line', (line: string) => {
    const n = parseInt(line);
    console.log(fibonacci(n));
});`,
		TestCases: []TestCase{
			{TestCaseID: 1, Input: "10", ExpectedOutput: "55"},
			{TestCaseID: 2, Input: "20", ExpectedOutput: "6765"},
			{TestCaseID: 3, Input: "15", ExpectedOutput: "610"},
		},
	},
}

func main() {
	totalRequests := len(languageConfigs) * requestsPerLanguage

	fmt.Printf("╔════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║         MULTI-LANGUAGE CONCURRENT LOAD TEST                    ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════════════╝\n\n")
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Languages:          %d (%s)\n", len(languageConfigs), getLanguageList())
	fmt.Printf("  Requests/Language:  %d\n", requestsPerLanguage)
	fmt.Printf("  Total Requests:     %d\n", totalRequests)
	fmt.Printf("  Poll Interval:      %v\n", pollInterval)
	fmt.Printf("  Test Cases/Request: 3 (Fibonacci)\n\n")

	var wg sync.WaitGroup
	var successCount int64
	var failureCount int64

	metrics := make([]RequestMetrics, 0, totalRequests)
	var metricsLock sync.Mutex

	tr := &http.Transport{
		MaxIdleConns:        totalRequests,
		MaxIdleConnsPerHost: totalRequests,
		IdleConnTimeout:     90 * time.Second,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   60 * time.Second,
	}

	startTime := time.Now()
	fmt.Printf("Starting load test at %s...\n\n", startTime.Format("15:04:05"))

	requestID := 0
	// Launch requests for each language concurrently
	for _, config := range languageConfigs {
		for i := 0; i < requestsPerLanguage; i++ {
			wg.Add(1)
			requestID++
			go func(id int, lang string, code string, testCases []TestCase) {
				defer wg.Done()

				payload := SubmissionRequest{
					Language:      lang,
					SourceCode:    code,
					TestCases:     testCases,
					TimeLimitMs:   5000,
					MemoryLimitKb: 262144,
				}

				metric := executeRequest(client, id, lang, payload)

				if metric.Success {
					atomic.AddInt64(&successCount, 1)
				} else {
					atomic.AddInt64(&failureCount, 1)
				}

				metricsLock.Lock()
				metrics = append(metrics, metric)
				metricsLock.Unlock()
			}(requestID, config.Language, config.SourceCode, config.TestCases)
		}
	}

	wg.Wait()
	totalWallTime := time.Since(startTime)

	// Calculate and print comprehensive statistics
	printDetailedResults(metrics, successCount, failureCount, totalWallTime)
}

func executeRequest(client *http.Client, id int, language string, payload SubmissionRequest) RequestMetrics {
	metric := RequestMetrics{
		ID:       id,
		Language: language,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		metric.Success = false
		metric.ErrorMessage = "JSON marshal error: " + err.Error()
		return metric
	}

	// 1. Submit
	submitStart := time.Now()
	resp, err := client.Post(baseURL+"/submit", "application/json", bytes.NewBuffer(jsonData))
	metric.SubmitDuration = time.Since(submitStart)

	if err != nil {
		metric.Success = false
		metric.ErrorMessage = "Submit error: " + err.Error()
		return metric
	}

	// Accept both 200 OK and 202 Accepted for async submissions
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		metric.Success = false
		metric.StatusCode = resp.StatusCode
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		metric.ErrorMessage = fmt.Sprintf("Submit failed with status %d: %s", resp.StatusCode, string(body))
		return metric
	}

	var subResp struct {
		SubmissionID string `json:"submission_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&subResp); err != nil {
		metric.Success = false
		metric.ErrorMessage = "Failed to decode submit response: " + err.Error()
		resp.Body.Close()
		return metric
	}
	resp.Body.Close()
	metric.SubmissionID = subResp.SubmissionID

	// 2. Poll until completion
	maxPolls := 120 // Safety limit: 120 polls * 500ms = 60 seconds max
	for metric.PollCount < maxPolls {
		time.Sleep(pollInterval)
		metric.PollCount++

		pollReqStart := time.Now()
		statusResp, err := client.Get(fmt.Sprintf("%s/status?submission_id=%s", baseURL, metric.SubmissionID))
		pollReqDuration := time.Since(pollReqStart)
		metric.TotalPollTime += pollReqDuration

		if err != nil {
			metric.Success = false
			metric.ErrorMessage = "Poll error: " + err.Error()
			return metric
		}

		if statusResp.StatusCode != http.StatusOK {
			metric.Success = false
			metric.StatusCode = statusResp.StatusCode
			body, _ := io.ReadAll(statusResp.Body)
			statusResp.Body.Close()
			metric.ErrorMessage = fmt.Sprintf("Poll failed with status %d: %s", statusResp.StatusCode, string(body))
			return metric
		}

		var statusData SubmissionResponse
		if err := json.NewDecoder(statusResp.Body).Decode(&statusData); err != nil {
			metric.Success = false
			metric.ErrorMessage = "Poll decode error: " + err.Error()
			statusResp.Body.Close()
			return metric
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

		return metric
	}

	// If we exit the loop, we hit the max poll limit
	metric.Success = false
	metric.ErrorMessage = fmt.Sprintf("Exceeded max polls (%d), submission may still be pending", maxPolls)
	return metric
}

func getLanguageList() string {
	langs := make([]string, len(languageConfigs))
	for i, cfg := range languageConfigs {
		langs[i] = cfg.Language
	}
	result := ""
	for i, lang := range langs {
		if i > 0 {
			result += ", "
		}
		result += lang
	}
	return result
}

func printDetailedResults(metrics []RequestMetrics, successCount, failureCount int64, totalWallTime time.Duration) {
	totalCount := len(metrics)

	// Separate metrics by language
	langMetrics := make(map[string][]RequestMetrics)
	for _, m := range metrics {
		langMetrics[m.Language] = append(langMetrics[m.Language], m)
	}

	fmt.Printf("╔════════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║                       TEST RESULTS                             ║\n")
	fmt.Printf("╚════════════════════════════════════════════════════════════════╝\n\n")

	// Overall Statistics
	fmt.Printf("┌─ Overall Statistics ────────────────────────────────────────┐\n")
	fmt.Printf("│ Total Requests:         %-6d                               │\n", totalCount)
	fmt.Printf("│ Successful:             %-6d (%.2f%%)                     │\n",
		successCount, float64(successCount)/float64(totalCount)*100)
	fmt.Printf("│ Failed:                 %-6d (%.2f%%)                     │\n",
		failureCount, float64(failureCount)/float64(totalCount)*100)
	fmt.Printf("│ Total Wall Time:        %-40s │\n", totalWallTime)
	fmt.Printf("│ Throughput:             %.2f req/sec                        │\n",
		float64(totalCount)/totalWallTime.Seconds())
	fmt.Printf("└──────────────────────────────────────────────────────────────┘\n\n")

	// Per-Language Breakdown
	fmt.Printf("┌─ Per-Language Results ──────────────────────────────────────┐\n")
	
	languages := []string{"cpp", "c", "java", "python", "node", "csharp", "ts"}
	for _, lang := range languages {
		langMetricsSlice := langMetrics[lang]
		if len(langMetricsSlice) == 0 {
			continue
		}

		successLang := int64(0)
		failedLang := int64(0)
		var totalWaitLang []time.Duration

		for _, m := range langMetricsSlice {
			if m.Success {
				successLang++
				totalWaitLang = append(totalWaitLang, m.TotalWaitDuration)
			} else {
				failedLang++
			}
		}

		avgTime := time.Duration(0)
		if len(totalWaitLang) > 0 {
			sort.Slice(totalWaitLang, func(i, j int) bool { return totalWaitLang[i] < totalWaitLang[j] })
			avgTime = calculateAvg(totalWaitLang)
		}

		fmt.Printf("│  %-10s  Success: %3d/%3d  Avg Time: %-16s│\n",
			lang, successLang, len(langMetricsSlice), avgTime)
	}
	fmt.Printf("└──────────────────────────────────────────────────────────────┘\n\n")

	// Collect successful metrics for detailed analysis
	var (
		submitDurations    []time.Duration
		totalWaitDurations []time.Duration
		internalExecTimes  []time.Duration
		overheadDurations  []time.Duration
		pollCounts         []int
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
		}
	}

	if len(submitDurations) == 0 {
		fmt.Printf("⚠️  No successful requests to analyze.\n\n")
		printErrorAnalysis(metrics, failureCount)
		return
	}

	sort.Slice(submitDurations, func(i, j int) bool { return submitDurations[i] < submitDurations[j] })
	sort.Slice(totalWaitDurations, func(i, j int) bool { return totalWaitDurations[i] < totalWaitDurations[j] })
	sort.Slice(internalExecTimes, func(i, j int) bool { return internalExecTimes[i] < internalExecTimes[j] })
	sort.Slice(overheadDurations, func(i, j int) bool { return overheadDurations[i] < overheadDurations[j] })
	sort.Slice(pollCounts, func(i, j int) bool { return pollCounts[i] < pollCounts[j] })
	sort.Slice(totalPollTimes, func(i, j int) bool { return totalPollTimes[i] < totalPollTimes[j] })

	// Submit Phase Statistics
	fmt.Printf("┌─ Submit Phase (Initial POST /submit) ───────────────────────┐\n")
	printTimingStats("│ Submit Time", submitDurations)
	fmt.Printf("└──────────────────────────────────────────────────────────────┘\n\n")

	// Total Wait Time
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
	fmt.Printf("└──────────────────────────────────────────────────────────────┘\n\n")

	// Error Breakdown (if any)
	if failureCount > 0 {
		printErrorAnalysis(metrics, failureCount)
	}

	// Performance Summary
	fmt.Printf("┌─ Performance Summary ────────────────────────────────────────┐\n")
	fmt.Printf("│  Average End-to-End:    %-37s │\n", calculateAvg(totalWaitDurations))
	fmt.Printf("│  Average Execution:     %-37s │\n", calculateAvg(internalExecTimes))
	fmt.Printf("│  Average Overhead:      %-37s │\n", calculateAvg(overheadDurations))
	fmt.Printf("│  Average Polls:         %-37.2f │\n", calculateAvgInt(pollCounts))
	fmt.Printf("│                                                              │\n")
	fmt.Printf("│  Success Rate:          %.2f%%                               │\n",
		float64(successCount)/float64(totalCount)*100)
	fmt.Printf("└──────────────────────────────────────────────────────────────┘\n")
}

func printErrorAnalysis(metrics []RequestMetrics, failureCount int64) {
	fmt.Printf("┌─ Error Analysis ─────────────────────────────────────────────┐\n")
	
	errorsByLang := make(map[string]int)
	errorTypes := make(map[string]int)
	var errorSamples []string

	for _, m := range metrics {
		if !m.Success {
			errorsByLang[m.Language]++
			
			errKey := m.ErrorMessage
			if len(errKey) > 50 {
				errKey = errKey[:50]
			}
			errorTypes[errKey]++

			if len(errorSamples) < 3 {
				errorSamples = append(errorSamples, fmt.Sprintf("[%s] %s", m.Language, m.ErrorMessage))
			}
		}
	}

	fmt.Printf("│ Errors by Language:                                          │\n")
	for lang, count := range errorsByLang {
		fmt.Printf("│  • %-10s: %4d errors                                  │\n", lang, count)
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
	fmt.Printf("└──────────────────────────────────────────────────────────────┘\n\n")
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