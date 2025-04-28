package dag

import (
	"reflect"
	"testing"
)

func TestGetPathsForWindow(t *testing.T) {
	tests := []struct {
		name           string
		algoExecPath   []string
		windowExecPath []string
		procExecPath   []string
		windowID       int
		want           []ExecutionPath
		wantErr        bool
		errorMsg       string
	}{
		{
			name:           "filter for window ID 1",
			algoExecPath:   []string{"1.2.3.4.5"},
			windowExecPath: []string{"1.1.1.3.3"},
			procExecPath:   []string{"1.1.1.2.2"},
			windowID:       1,
			want: []ExecutionPath{
				{
					AlgoPath:    "1.2.3",
					ProcessorId: "1",
				},
			},
			wantErr: false,
		},
		{
			name:           "filter for window ID 3",
			algoExecPath:   []string{"1.2.3.4.5"},
			windowExecPath: []string{"1.1.1.3.3"},
			procExecPath:   []string{"1.1.1.2.2"},
			windowID:       3,
			want: []ExecutionPath{
				{
					AlgoPath:    "4.5",
					ProcessorId: "2",
				},
			},
			wantErr: false,
		},
		{
			name:           "mismatched path lengths",
			algoExecPath:   []string{"1.2.3"},
			windowExecPath: []string{"1.1"},
			procExecPath:   []string{"1.1.1"},
			windowID:       1,
			want:           nil,
			wantErr:        true,
			errorMsg:       "path lengths do not match: algo=3, window=2, proc=3",
		},
		{
			name:           "window ID not found",
			algoExecPath:   []string{"1.2.3"},
			windowExecPath: []string{"1.1.1"},
			procExecPath:   []string{"1.1.1"},
			windowID:       4,
			want:           nil,
			wantErr:        false,
		},
		{
			name:           "single segment paths",
			algoExecPath:   []string{"1"},
			windowExecPath: []string{"1"},
			procExecPath:   []string{"1"},
			windowID:       1,
			want: []ExecutionPath{
				{
					AlgoPath:    "1",
					ProcessorId: "1",
				},
			},
			wantErr: false,
		},
		{
			name:           "multiple matches in path",
			algoExecPath:   []string{"1.2.3.1.2"},
			windowExecPath: []string{"1.1.1.1.1"},
			procExecPath:   []string{"1.1.1.1.1"},
			windowID:       1,
			want:           nil,
			wantErr:        true,
			errorMsg:       "cyclic graph discovered at position 3. aborting",
		},
		{
			name:           "empty paths",
			algoExecPath:   []string{""},
			windowExecPath: []string{""},
			procExecPath:   []string{""},
			windowID:       1,
			want:           nil,
			wantErr:        false,
		},
		{
			name:           "split processor",
			algoExecPath:   []string{"1.2.3.4.5.6"},
			windowExecPath: []string{"1.1.1.1.1.2"},
			procExecPath:   []string{"3.4.4.5.5.5"},
			windowID:       1,
			want: []ExecutionPath{
				{
					AlgoPath:    "1",
					ProcessorId: "3",
				},
				{
					AlgoPath:    "2.3",
					ProcessorId: "4",
				},
				{
					AlgoPath:    "4.5",
					ProcessorId: "5",
				},
			},

			wantErr: false,
		},
		{
			name:           "split processor with revisit",
			algoExecPath:   []string{"1.2.3.4.5.6"},
			windowExecPath: []string{"1.1.1.1.1.1"},
			procExecPath:   []string{"3.4.4.5.5.4"},
			windowID:       1,
			want: []ExecutionPath{
				{
					AlgoPath:    "1",
					ProcessorId: "3",
				},
				{
					AlgoPath:    "2.3",
					ProcessorId: "4",
				},
				{
					AlgoPath:    "4.5",
					ProcessorId: "5",
				},
				{
					AlgoPath:    "6",
					ProcessorId: "4",
				},
			},

			wantErr: false,
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

func TestPathInSubset(t *testing.T) {
	tests := []struct {
		name      string
		pathStack []string
		new       string
		want      bool
	}{
		{
			name:      "not subpath",
			pathStack: []string{"a.b.c", "d.e.f"},
			new:       "h.i",
			want:      false,
		},
		{
			name:      "is subpath",
			pathStack: []string{"a.b.c", "d.e.f"},
			new:       "b.c",
			want:      true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isSubspath := isSubsetOf(test.pathStack, test.new)
			if isSubspath != test.want {
				t.Errorf("isSubpath() = %v, want %v", isSubspath, test.want)
			}
		})
	}
}
