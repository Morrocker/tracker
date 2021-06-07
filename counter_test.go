package tracker

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCounter(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want Counter
	}{
		{
			name: "new object",
			args: args{
				"object A",
			},
			want: &counter{
				name:    "object A",
				current: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCounter(tt.args.name)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_counter_SetCurrent(t *testing.T) {
	type fields struct {
		name    string
		current uint64
	}
	type args struct {
		n uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "basic set current test",
			fields: fields{
				name:    "simple counter",
				current: 0,
			},
			args: args{
				n: 15,
			},
		},
		{
			name: "basic set current test",
			fields: fields{
				name:    "simple counter",
				current: 40,
			},
			args: args{
				n: 15,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &counter{
				name:    tt.fields.name,
				current: tt.fields.current,
			}
			g.SetCurrent(tt.args.n)

			assert.Equal(t, tt.args.n, g.current, "Expected and actual do not match. Expected: %d Actual: %d", tt.args.n, g.current)
		})
	}
}

func Test_counter_Current(t *testing.T) {
	type fields struct {
		name    string
		current uint64
	}
	type args struct {
		n int64
	}
	type want struct {
		err  bool
		want uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "basic Current() test. 0 + 15",
			fields: fields{
				name:    "simple counter",
				current: 0,
			},
			args: args{
				n: 15,
			},
			want: want{
				err:  false,
				want: 15,
			},
		},
		{
			name: "basic Current() test. 15 + 15",
			fields: fields{
				name:    "simple counter",
				current: 15,
			},
			args: args{
				n: 15,
			},
			want: want{
				err:  false,
				want: 30,
			},
		},
		{
			name: "basic Current() test. 15 + 15",
			fields: fields{
				name:    "simple counter",
				current: 15,
			},
			args: args{
				n: -25,
			},
			want: want{
				err:  true,
				want: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &counter{
				name:    tt.fields.name,
				current: tt.fields.current,
			}
			got, err := g.Current(tt.args.n)

			if tt.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, got, tt.want.want)
		})
	}
}

func Test_counter_RawValue(t *testing.T) {
	type fields struct {
		name    string
		current uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   uint64
	}{
		{
			name: "basic RawValue() test",
			fields: fields{
				name:    "simple counter",
				current: 55,
			},
			want: 55,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &counter{
				name:    tt.fields.name,
				current: tt.fields.current,
			}
			if got := g.RawValue(); got != tt.want {
				assert.Equal(t, got, tt.want)
			}
		})
	}
}

func Test_counter_Value(t *testing.T) {
	type fields struct {
		name      string
		current   uint64
		unitsFunc func(uint64) string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "basic Value() test",
			fields: fields{
				name:    "simple counter",
				current: 55,
				unitsFunc: func(n uint64) string {
					return fmt.Sprintf("%d", n)
				},
			},
			want: "55",
		},
		{
			name: "basic Value() test",
			fields: fields{
				name:    "simple counter",
				current: 0,
				unitsFunc: func(n uint64) string {
					return fmt.Sprintf("%d", n)
				},
			},
			want: "0",
		},
		{
			name: "basic Value() test",
			fields: fields{
				name:    "simple counter",
				current: 0,
			},
			want: "unitsFunction not set",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &counter{
				name:      tt.fields.name,
				current:   tt.fields.current,
				unitsFunc: tt.fields.unitsFunc,
			}
			if got := g.Value(); got != tt.want {
				assert.Equal(t, got, tt.want)
			}
		})
	}
}

func Test_counter_UnitsFunc(t *testing.T) {
	type fields struct {
		name      string
		current   uint64
		unitsFunc func(uint64) string
	}
	type args struct {
		f func(uint64) string
	}
	type want struct {
		err bool
		f   func(uint64) string
		ret string
	}
	testFunc := func(n uint64) string {
		return fmt.Sprintf("%d", n)
	}
	testFunc2 := func(n uint64) string {
		return fmt.Sprintf("bad%d", n)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "basic UnitsFunc() test",
			fields: fields{
				name:    "good function",
				current: 55,
			},
			args: args{
				f: testFunc,
			},
			want: want{
				err: false,
				f:   testFunc,
				ret: "55",
			},
		},
		{
			name: "bad UnitsFunc() test",
			fields: fields{
				name:    "bad function match",
				current: 55,
			},
			args: args{
				f: testFunc,
			},
			want: want{
				err: true,
				f:   testFunc2,
				ret: "bad55",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &counter{
				name:      tt.fields.name,
				current:   tt.fields.current,
				unitsFunc: tt.fields.unitsFunc,
			}
			g.UnitsFunc(tt.args.f)
			ptr1 := reflect.ValueOf(g.unitsFunc).Pointer()
			ptr2 := reflect.ValueOf(tt.want.f).Pointer()
			val := g.Value()
			if tt.want.err {
				assert.NotEqual(t, ptr1, ptr2)
				assert.NotEqual(t, tt.want.ret, val)
			} else {
				assert.Equal(t, ptr1, ptr2)
				assert.Equal(t, tt.want.ret, val)
			}
		})
	}
}
