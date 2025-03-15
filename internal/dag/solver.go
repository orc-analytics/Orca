package dag

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// ExecutionPath represents a set of filtered paths for algorithms, windows and processors
type ExecutionPath struct {
	AlgoPath   string
	WindowPath string
	ProcPath   string
}

// GetPathsForWindow
// Given a set of algorithm execution paths, the processors the
// algorithms belong to, and the windows that trigger them, filter
// by window, and group by processor producing the consecutive
// algorithm executions.
//
// When processing needs to be split over processors, the order
// of execution paths in the return argument should be preserved.
//
// Example arguments
// algo_exec_path   window_exec_path  proc_exec_path
// 1	              1	                1
// 1.2.3.4.5        1.1.1.3.3	        1.1.1.2.2
// 1.2.5.6.7.8      1.1.1.24.3.3	    1.1.1.1.2.2
// 1.2.4.5.6        1.1.1.1.24        1.1.1.1.1
//
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
func GetPathsForWindow(
	algoExecPath string,
	windowExecPath string,
	procExecPath string,
	windowID int,
) ([]ExecutionPath, error) {
	windowIDInt := strconv.Itoa(windowID)

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
	inRun := false
	var startIdx int
	var visitedAlgos []string
	var currentProcessor string
	for i, windowSegmentId := range windowSegments {
		if slices.Contains(visitedAlgos, algoSegments[i]) {
			return nil, fmt.Errorf("cyclic graph discovered at position %v. aborting", i)
		}

		visitedAlgos = append(visitedAlgos, algoSegments[i])
		if windowSegmentId == windowIDInt && !inRun {
			inRun = true
			startIdx = i
			currentProcessor = procSegments[i]
		}

		// handle case where the processor changes in a run
		if inRun && (procSegments[i] != currentProcessor || windowSegmentId != windowIDInt) {
			currentProcessor = procSegments[i]
			result := ExecutionPath{
				AlgoPath:   strings.Join(algoSegments[startIdx:i], "."),
				WindowPath: strings.Join(windowSegments[startIdx:i], "."),
				ProcPath:   strings.Join(procSegments[startIdx:i], "."),
			}
			startIdx = i
			results = append(results, result)
		}
	}

	// catch the edge case where the window runs to the end
	if inRun && windowSegments[len(windowSegments)-1] == windowIDInt {
		result := ExecutionPath{
			AlgoPath:   strings.Join(algoSegments[startIdx:], "."),
			WindowPath: strings.Join(windowSegments[startIdx:], "."),
			ProcPath:   strings.Join(procSegments[startIdx:], "."),
		}
		results = append(results, result)
	}

	// If no paths were found, return nil
	if len(results) == 0 {
		return nil, nil
	}

	return results, nil
}
