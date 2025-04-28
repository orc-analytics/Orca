package dag

import (
	"fmt"
	"log/slog"
	"slices"
	"strconv"
	"strings"
)

// ExecutionPath represents a set of filtered paths for algorithms, windows and processors
type ExecutionPath struct {
	AlgoPath    string
	ProcessorId int
}

// ExecutionStage represents a set of paths that can be executed in parallel
type ExecutionStage struct {
	Paths []ProcessorPath
	Level int // Depth in the execution graph
}

// ProcessorPath represents algorithms to be executed on a specific processor
type ProcessorPath struct {
	AlgoPath    string
	ProcessorId int
}

// ExecutionPlan represents the complete execution strategy
type ExecutionPlan struct {
	Stages []ExecutionStage // Ordered stages of execution
}

// isSubsetOf accepts a list of execution paths and if if the new execution path is a
// subset of existing, returns true
//
// Example:
//
// isSubsetOf("e.f.g.h", "f.g")
// > true
//
// isSubsetOf("e.f.g.h", "i.k")
// > false
func isSubsetOf(existing string, new string) bool {
	if strings.Contains(existing, new) {
		return true
	}
	return false
}

// appendResults will extend the results stack, taking in to account:
// - if the algo segment is a subset of what is already present, then do nothing
// - if it extends what is already present then replace it
// TODO: capture potential bug where the algoritm IDs are large. e.g.:
//
//	45.6.74
//	6.7 <- would be matched in the above, but this is not correct.
//	need to actually match like so:
//
// .45.6.74.
// .6.7.
// this would be robust
func appendResults(results []ExecutionPath, result ExecutionPath) []ExecutionPath {
	for ii, subResult := range results {
		subAlgoPath := fmt.Sprintf(".%v.", subResult.AlgoPath)
		newAlgoPath := fmt.Sprintf(".%v.", result.AlgoPath)

		// is new result subset of what already exists?
		if isSubsetOf(subAlgoPath, newAlgoPath) {
			return results // then do nothing
		} else if isSubsetOf(newAlgoPath, subAlgoPath) { // does result extend what already exists?
			// then replace the results
			results[ii] = result
			return results
		}
	}
	// else, new results
	results = append(results, result)
	return results
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
// Examples:
//
// algo_exec_path   window_exec_path  proc_exec_path
// 1	              1	                1
// 1.2.3.4.5        1.1.1.3.3	        1.1.1.2.2
// 1.2.5.6.7.8      1.1.1.24.3.3	    1.1.1.1.2.2
// 1.2.4.5.6        1.1.1.1.24        1.1.1.1.1
//
// when filtering for window id of '1'
//
// algo_exec_path   window_exec_path  proc_exec_path
// 1	              1	                1
// 1.2.3	          1.1.1	            1.1.1
// 1.2.5	          1.1.1             1.1.1
// 1.2.4.5          1.1.1.1           1.1.1.1
//
// when filtering for window id of '3'
//
// algo_exec_path   window_exec_path  proc_exec_path
// 4.5	            3.3               2.2
// 7.8  	          3.3               2.2
func GetPathsForWindow(
	algoExecPaths []string,
	windowExecPaths []string,
	procExecPaths []string,
	windowID int,
) ([]ExecutionPath, error) {
	windowIdStr := strconv.Itoa(windowID)
	if len(algoExecPaths) != len(windowExecPaths) || len(windowExecPaths) != len(procExecPaths) {
		return nil, fmt.Errorf(
			"number of graph paths do not match: algo=%d, window=%d, proc=%d",
			len(algoExecPaths),
			len(windowExecPaths),
			len(procExecPaths),
		)
	}

	var results []ExecutionPath

	for ii := range algoExecPaths {
		algoExecPath := algoExecPaths[ii]
		windowExecPath := windowExecPaths[ii]
		procExecPath := procExecPaths[ii]

		algoSegments := strings.Split(algoExecPath, ".")
		windowSegments := strings.Split(windowExecPath, ".")
		procSegments := strings.Split(procExecPath, ".")

		if len(algoSegments) != len(windowSegments) || len(windowSegments) != len(procSegments) {
			return nil, fmt.Errorf("path lengths do not match: algo=%d, window=%d, proc=%d",
				len(algoSegments), len(windowSegments), len(procSegments))
		}

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
			if windowSegmentId == windowIdStr && !inRun {
				inRun = true
				startIdx = i
				currentProcessor = procSegments[i]
			}

			// handle case where the window ends
			if inRun && (windowSegmentId != windowIdStr) {
				inRun = false
				procId, err := strconv.Atoi(procSegments[i-1])
				if err != nil {
					slog.Error("could not convert processor id to int", procSegments[i-1])
					return nil, err
				}
				result := ExecutionPath{
					AlgoPath:    strings.Join(algoSegments[startIdx:i], "."),
					ProcessorId: procId,
				}
				startIdx = i
				results = appendResults(results, result)
				continue
			}
			// handle case where the processor changes in a run
			if inRun && (procSegments[i] != currentProcessor) {
				currentProcessor = procSegments[i]

				procId, err := strconv.Atoi(procSegments[i-1])
				if err != nil {
					slog.Error("could not convert processor id to int", procSegments[i-1])
					return nil, err
				}
				result := ExecutionPath{
					AlgoPath:    strings.Join(algoSegments[startIdx:i], "."),
					ProcessorId: procId,
				}
				startIdx = i
				results = appendResults(results, result)
			}

		}

		// catch the edge case where the window runs to the end
		if inRun && windowSegments[len(windowSegments)-1] == windowIdStr {
			procId, err := strconv.Atoi(currentProcessor)
			if err != nil {
				slog.Error("could not convert processor id to int", currentProcessor)
				return nil, err
			}
			result := ExecutionPath{
				AlgoPath:    strings.Join(algoSegments[startIdx:], "."),
				ProcessorId: procId,
			}
			results = append(results, result)
		}

	}
	return results, nil
}
