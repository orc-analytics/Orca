package datalayer_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/reflect/protoreflect"

	datalayer "github.com/predixus/pdb_framework/internal/datalayer"
)

type MockProtoMessage struct {
	protoreflect.ProtoMessage
}

func TestRemoveIndexWorks(t *testing.T) {
	testCases := []struct {
		input  []int
		output []int
		r_idx  int
	}{
		{[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, []int{0, 1, 2, 3, 4, 6, 7, 8, 9}, 5},
		{[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, 0},
		{[]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, []int{0, 1, 2, 3, 4, 5, 6, 7, 9}, 8},
	}

	for _, tc := range testCases {

		input_ptr := make([]*int, len(tc.input))
		for i, v := range tc.input {
			input_ptr[i] = &v
		}
		output_ptr := make([]*int, len(tc.output))
		for i, v := range tc.output {
			output_ptr[i] = &v
		}

		res, err := datalayer.RemoveIndex(input_ptr, tc.r_idx)
		assert.True(t, len(res) == len(output_ptr), "Correct length of result")
		assert.Nil(t, err)

		for i, v := range res {
			assert.True(t, *v == *output_ptr[i], "Assignment is correct")
		}
	}
}

func TestRemoveIndexOutOfRangeGivesError(t *testing.T) {
	test_arr := []int{0, 1, 2, 3}
	input_ptr := make([]*int, len(test_arr))

	for i, v := range test_arr {
		input_ptr[i] = &v
	}

	res, err := datalayer.RemoveIndex(input_ptr, 4)
	assert.NotNil(t, err)
	assert.Nil(t, res)

	res, err = datalayer.RemoveIndex(input_ptr, -1)
	assert.NotNil(t, err)
	assert.Nil(t, res)
}
