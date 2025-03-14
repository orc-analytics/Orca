package dag

// GetPathsForWindow
// Given a set of algorithm execution paths, the processors the
// algorithms belong to, and the windows that trigger them, filter
// by window, consecutive algorithm executions.
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

func GetPathsForWindow(
	algo_exec_path string,
	window_exec_path string,
	proc_exec_path string,
) error {
	return nil
}
