package tracker

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGauge(t *testing.T) {
	type args struct {
		name    string
		current uint64
		total   uint64
	}
	tests := []struct {
		name string
		args args
		want Gauge
	}{
		{
			name: "new gauge",
			args: args{
				name:    "object A",
				current: 0,
				total:   100,
			},
			want: &gauge{
				name:    "object A",
				current: 0,
				total:   100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewGauge(tt.args.name, tt.args.total)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_gauge_SetCurrent(t *testing.T) {
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
		want   uint64
	}{
		{
			name: "basic set current test",
			fields: fields{
				name:    "simple gauge",
				current: 0,
			},
			args: args{
				n: 15,
			},
			want: 15,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &gauge{
				name:    tt.fields.name,
				current: tt.fields.current,
			}
			g.SetCurrent(tt.args.n)

			assert.Equal(t, tt.want, g.current)
		})
	}
}

func Test_gauge_Current(t *testing.T) {
	type fields struct {
		name    string
		current uint64
	}
	type args struct {
		n int64
	}
	type want struct {
		err  bool
		curr uint64
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
				name:    "simple gauge",
				current: 0,
			},
			args: args{
				n: 15,
			},
			want: want{
				err:  false,
				curr: 15,
			},
		},
		{
			name: "basic Current() test. 15 + 15",
			fields: fields{
				name:    "simple gauge",
				current: 15,
			},
			args: args{
				n: 15,
			},
			want: want{
				err:  false,
				curr: 30,
			},
		},
		{
			name: "basic Current() test. 15 - 20",
			fields: fields{
				name:    "simple gauge",
				current: 15,
			},
			args: args{
				n: -20,
			},
			want: want{
				err:  true,
				curr: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &gauge{
				name:    tt.fields.name,
				current: tt.fields.current,
			}
			got, err := g.Current(tt.args.n)
			if tt.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, got, tt.want.curr)
		})
	}
}

func Test_gauge_SetTotal(t *testing.T) {
	type fields struct {
		name  string
		total uint64
	}
	type args struct {
		n uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint64
	}{
		{
			name: "basic set current test",
			fields: fields{
				name:  "simple gauge",
				total: 0,
			},
			args: args{
				n: 15,
			},
			want: 15,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &gauge{
				name:  tt.fields.name,
				total: tt.fields.total,
			}
			g.SetTotal(tt.args.n)

			assert.Equal(t, tt.want, g.total)
		})
	}
}

func Test_gauge_Total(t *testing.T) {
	type fields struct {
		name  string
		total uint64
	}
	type args struct {
		n int64
	}
	type want struct {
		err   bool
		total uint64
		ret   uint64
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
				name:  "simple gauge",
				total: 0,
			},
			args: args{
				n: 15,
			},
			want: want{
				err:   false,
				total: 15,
				ret:   15,
			},
		},
		{
			name: "basic Current() test. 15 + 15",
			fields: fields{
				name:  "simple gauge",
				total: 15,
			},
			args: args{
				n: -25,
			},
			want: want{
				err:   true,
				total: 15,
				ret:   0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &gauge{
				name:  tt.fields.name,
				total: tt.fields.total,
			}
			got, err := g.Total(tt.args.n)
			if err != nil {
				assert.Error(t, err)
				assert.Equal(t, g.total, tt.want.total)
				assert.Equal(t, got, tt.want.ret)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, g.total, tt.want.total)
				assert.Equal(t, got, tt.want.ret)
			}
		})
	}
}

func Test_gauge_RawValues(t *testing.T) {
	type fields struct {
		name    string
		current uint64
		total   uint64
	}
	type want struct {
		current uint64
		total   uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "basic RawValue() test",
			fields: fields{
				name:    "simple gauge",
				current: 55,
				total:   100,
			},
			want: want{
				current: 55,
				total:   100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &gauge{
				name:    tt.fields.name,
				current: tt.fields.current,
				total:   tt.fields.total,
			}
			curr, tot := g.RawValues()
			assert.Equal(t, curr, tt.want.current)
			assert.Equal(t, tot, tt.want.total)
		})
	}
}

func Test_gauge_Value(t *testing.T) {
	type fields struct {
		name      string
		current   uint64
		total     uint64
		unitsFunc func(uint64) string
	}
	type want struct {
		current string
		total   string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "basic Value() test",
			fields: fields{
				name:    "simple gauge",
				current: 55,
				total:   100,
				unitsFunc: func(n uint64) string {
					return fmt.Sprintf("%d", n)
				},
			},
			want: want{
				current: "55",
				total:   "100",
			},
		},
		{
			name: "basic Value() test",
			fields: fields{
				name:    "simple gauge",
				current: 55,
				total:   100,
			},
			want: want{
				current: "unitsFunction not set",
				total:   "unitsFunction not set",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &gauge{
				name:      tt.fields.name,
				current:   tt.fields.current,
				total:     tt.fields.total,
				unitsFunc: tt.fields.unitsFunc,
			}
			curr, tot := g.Values()
			assert.Equal(t, tt.want.current, curr)
			assert.Equal(t, tt.want.total, tot)
		})
	}
}

func Test_gauge_UnitsFunc(t *testing.T) {
	type fields struct {
		name    string
		current uint64
		total   uint64
	}
	type args struct {
		f func(uint64) string
	}
	type want struct {
		err   bool
		f     func(uint64) string
		curr  string
		total string
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
				total:   100,
			},
			args: args{
				f: testFunc,
			},
			want: want{
				err:   false,
				f:     testFunc,
				curr:  "55",
				total: "100",
			},
		},
		{
			name: "bad UnitsFunc() test",
			fields: fields{
				name:    "bad function match",
				current: 55,
				total:   100,
			},
			args: args{
				f: testFunc,
			},
			want: want{
				err:   true,
				f:     testFunc2,
				curr:  "55",
				total: "100",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &gauge{
				name:    tt.fields.name,
				total:   tt.fields.total,
				current: tt.fields.current,
			}
			g.UnitsFunc(tt.args.f)
			ptr1 := reflect.ValueOf(g.unitsFunc).Pointer()
			ptr2 := reflect.ValueOf(tt.want.f).Pointer()
			curr, tot := g.Values()
			if tt.want.err {
				assert.NotEqual(t, ptr1, ptr2)
			} else {
				assert.Equal(t, ptr1, ptr2)
			}
			assert.Equal(t, tt.want.curr, curr)
			assert.Equal(t, tt.want.total, tot)
		})
	}
}
