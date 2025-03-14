package dag

import (
	"fmt"
	"strings"
)

// GetPathsForWindow
// Given a set of algorithm execution paths, the processors the
// algorithms belong to, and the windows that trigger them, filter
// by window, and group by processor producing the consecutive
// algorithm executions
//
// Example arguments

// algo_exec_path   window_exec_path  proc_exec_path
// 1	              1	                1
// 1.2.3.4.5        1.1.1.3.3	        1.1.1.2.2
// 1.2.5.6.7.8      1.1.1.24.3.3	    1.1.1.1.2.2
// 1.2.4.5.6        1.1.1.1.24        1.1.1.1.1

// when filtering for '1'
//
// algo_exec_path   window_exec_path  proc_exec_path
// 1	              1	                1
// 1.2.3	          1.1.1	            1.1.1
// 1.2.5	          1.1.1             1.1.1
// 1.2.4.5          1.1.1.1           1.1.1.1
//
// when filtering for '3'
//
// algo_exec_path   window_exec_path  proc_exec_path
// 4.5	            3.3               2.2
// 7.8  	          3.3               2.2

// ExecutionPath represents a set of filtered paths for algorithms, windows and processors
type ExecutionPath struct {
	AlgoPath   string
	WindowPath string
	ProcPath   string
}

func GetPathsForWindow(
	algoExecPath string,
	windowExecPath string,
	procExecPath string,
	windowID string,
) ([]ExecutionPath, error) {
	// Split paths into segments
	algoSegments := strings.Split(algoExecPath, ".")
	windowSegments := strings.Split(windowExecPath, ".")
	procSegments := strings.Split(procExecPath, ".")

	// Validate input lengths match
	if len(algoSegments) != len(windowSegments) || len(windowSegments) != len(procSegments) {
		return nil, fmt.Errorf("path lengths do not match: algo=%d, window=%d, proc=%d",
			len(algoSegments), len(windowSegments), len(procSegments))
	}

	var results []ExecutionPath

	// Find the indices where windowID appears in window path
	for i, windowSegmentId := range windowSegments {
		if windowSegmentId == windowID {
			// For each match, take all segments up to that point
			result := ExecutionPath{
				AlgoPath:   strings.Join(algoSegments[:i+1], "."),
				WindowPath: strings.Join(windowSegments[:i+1], "."),
				ProcPath:   strings.Join(procSegments[:i+1], "."),
			}
			results = append(results, result)
		}
	}

	// If no paths were found, return an error
	if len(results) == 0 {
		return nil, fmt.Errorf("no paths found for window ID: %s", windowID)
	}

	return results, nil
}
