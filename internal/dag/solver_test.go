package dag

import (
	"reflect"
	"testing"
)

func TestGetPathsForWindow(t *testing.T) {
	tests := []struct {
		name           string
		algoExecPath   string
		windowExecPath string
		procExecPath   string
		windowID       string
		want           []ExecutionPath
		wantErr        bool
		errorMsg       string
	}{
		{
			name:           "filter for window ID 1",
			algoExecPath:   "1.2.3.4.5",
			windowExecPath: "1.1.1.3.3",
			procExecPath:   "1.1.1.2.2",
			windowID:       "1",
			want: []ExecutionPath{
				{
					AlgoPath:   "1.2.3",
					WindowPath: "1.1.1",
					ProcPath:   "1.1.1",
				},
			},
			wantErr: false,
		},
		{
			name:           "filter for window ID 3",
			algoExecPath:   "1.2.3.4.5",
			windowExecPath: "1.1.1.3.3",
			procExecPath:   "1.1.1.2.2",
			windowID:       "3",
			want: []ExecutionPath{
				{
					AlgoPath:   "4.5",
					WindowPath: "3.3",
					ProcPath:   "2.2",
				},
			},
			wantErr: false,
		},
		{
			name:           "mismatched path lengths",
			algoExecPath:   "1.2.3",
			windowExecPath: "1.1",
			procExecPath:   "1.1.1",
			windowID:       "1",
			want:           nil,
			wantErr:        true,
			errorMsg:       "path lengths do not match: algo=3, window=2, proc=3",
		},
		{
			name:           "window ID not found",
			algoExecPath:   "1.2.3",
			windowExecPath: "1.1.1",
			procExecPath:   "1.1.1",
			windowID:       "4",
			want:           nil,
			wantErr:        true,
			errorMsg:       "no paths found for window ID: 4",
		},
		{
			name:           "single segment paths",
			algoExecPath:   "1",
			windowExecPath: "1",
			procExecPath:   "1",
			windowID:       "1",
			want: []ExecutionPath{
				{
					AlgoPath:   "1",
					WindowPath: "1",
					ProcPath:   "1",
				},
			},
			wantErr: false,
		},
		{
			name:           "multiple matches in path",
			algoExecPath:   "1.2.3.1.2",
			windowExecPath: "1.1.1.1.1",
			procExecPath:   "1.1.1.1.1",
			windowID:       "1",
			want:           nil,
			wantErr:        true,
			errorMsg:       "cyclic graph discovered at position 3. aborting",
		},
		{
			name:           "empty paths",
			algoExecPath:   "",
			windowExecPath: "",
			procExecPath:   "",
			windowID:       "1",
			want:           nil,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPathsForWindow(
				tt.algoExecPath,
				tt.windowExecPath,
				tt.procExecPath,
				tt.windowID,
			)

			// Check error cases
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetPathsForWindow() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("GetPathsForWindow() error = %v, wantErr %v", err, tt.errorMsg)
					return
				}
				return
			}

			// Check success cases
			if err != nil {
				t.Errorf("GetPathsForWindow() unexpected error = %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPathsForWindow() = %v, want %v", got, tt.want)
			}
		})
	}
}
