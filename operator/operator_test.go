package operator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResult(t *testing.T) {
	type fields struct {
		op Operator
	}

	type args struct {
		left, right int
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			"2+2",
			fields{&Add{}},
			args{
				left:  2,
				right: 2,
			},
			4,
		},
		{
			"12+7",
			fields{&Add{}},
			args{
				left:  12,
				right: 7,
			},
			19,
		},
		{
			"2-2",
			fields{&Subtract{}},
			args{
				left:  2,
				right: 2,
			},
			0,
		},
		{
			"12-7",
			fields{&Subtract{}},
			args{
				left:  12,
				right: 7,
			},
			5,
		},
		{
			"12-9",
			fields{&Subtract{}},
			args{
				left:  12,
				right: 9,
			},
			3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := tt.fields.op
			assert.Equalf(t, tt.want, op.Result(tt.args.left, tt.args.right), "result for %s", op)
			// if got := op.Result(tt.args.left, tt.args.right); got != tt.want {
			// 	t.Errorf("Add.Result() = %v, want %v", got, tt.want)
			// }
		})
	}
}
