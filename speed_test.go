package tracker

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/morrocker/benchmark"
	"github.com/stretchr/testify/assert"
)

func TestNewSpeed(t *testing.T) {
	type args struct {
		g *int64
		n uint
	}
	var n int64 = 10
	tests := []struct {
		name string
		args args
		want Speed
	}{
		{
			name: "simple newSpeed",
			args: args{
				g: &n,
				n: 5,
			},
			want: &speed{
				target: &n,
				rate:   benchmark.NewSingleRate(5),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSpeed(tt.args.g, tt.args.n)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_speed_SampleSize(t *testing.T) {
	type fields struct {
		rate benchmark.SingleRate
	}
	type args struct {
		n uint
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint
	}{
		{
			name: "simple test",
			fields: fields{
				rate: benchmark.NewSingleRate(10),
			},
			args: args{
				n: 10,
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &speed{
				rate: tt.fields.rate,
			}
			got := s.SampleSize(tt.args.n)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_speed_StartMeasure(t *testing.T) {
	type fields struct {
		target *int64
		rate   benchmark.SingleRate
	}
	type args struct {
		time   time.Duration
		length int
	}
	ss := uint(5)
	var tgt int64
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "simple measure",
			fields: fields{
				target: &tgt,
				rate:   benchmark.NewSingleRate(ss),
			},
			args: args{
				time:   1 * time.Second,
				length: 3,
			},
		},
	}
	sum := int64(0)
	for _, tt := range tests {
		tgt = 0
		r := rand.Int63n(5000)
		t.Run(tt.name, func(t *testing.T) {
			s := &speed{
				target: tt.fields.target,
				rate:   tt.fields.rate,
			}
			for x := 0; x < tt.args.length; x++ {
				end := s.StartMeasure()
				time.Sleep(tt.args.time)
				tgt += r / int64(tt.args.time)
				end()
			}
			ssz, tot, ln := s.rate.Values()
			assert.Equal(t, sum, tot)
			assert.Equal(t, tt.args.length, ln)
			assert.Equal(t, ssz, ss)
		})
	}
}

func Test_speed_StartStopAutoMeasure(t *testing.T) {
	type fields struct {
		target *int64
		rate   benchmark.SingleRate
	}
	type args struct {
		time   time.Duration
		length int
	}
	type want struct {
		total  int64
		length int
	}
	ss := uint(5)
	var tgt int64
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "simple measure",
			fields: fields{
				target: &tgt,
				rate:   benchmark.NewSingleRate(ss),
			},
			args: args{
				time:   1000 * time.Millisecond,
				length: 10,
			},
			want: want{
				total:  0,
				length: int(ss),
			},
		},
	}
	for _, tt := range tests {
		tgt = 0
		t.Run(tt.name, func(t *testing.T) {
			s := &speed{
				target: tt.fields.target,
				rate:   tt.fields.rate,
			}
			s.StartAutoMeasure(tt.args.time)
			time.Sleep(tt.args.time + tt.args.time/2)

			arr := []int64{}
			for x := 0; x < tt.args.length; x++ {
				r := rand.Int63n(5000)
				tgt += r
				if x < int(ss) {
					arr = append(arr, r)
				} else {
					arr = append(arr[1:], r)
				}
				time.Sleep(tt.args.time)
			}
			s.StopAutoMeasure()
			ssz, tot, ln := s.rate.Values()
			for x := 0; x < len(arr); x++ {
				tt.want.total += arr[x]
			}
			assert.Equal(t, tt.want.total, tot)
			assert.Equal(t, tt.want.length, ln)
			assert.Equal(t, ssz, ss)
		})
	}
}

func Test_speed_AutoMeasureReset(t *testing.T) {
	type fields struct {
		target *int64
		rate   benchmark.SingleRate
	}
	type args struct {
		time   time.Duration
		length int
	}
	type want struct {
		total  int64
		length int
	}
	ss := uint(5)
	var tgt int64
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "simple measure",
			fields: fields{
				target: &tgt,
				rate:   benchmark.NewSingleRate(ss),
			},
			args: args{
				time:   1000 * time.Millisecond,
				length: 10,
			},
			want: want{
				total:  0,
				length: int(ss),
			},
		},
	}
	for _, tt := range tests {
		tgt = 0
		t.Run(tt.name, func(t *testing.T) {
			s := &speed{
				target: tt.fields.target,
				rate:   tt.fields.rate,
			}
			s.StartAutoMeasure(tt.args.time)
			for x := 0; x < 3; x++ {
				time.Sleep(tt.args.time / 3)
				s.StartAutoMeasure(tt.args.time)
			}
			time.Sleep(tt.args.time + tt.args.time/2)

			arr := []int64{}
			for x := 0; x < tt.args.length; x++ {
				r := rand.Int63n(5000)
				tgt += r
				if x < int(ss) {
					arr = append(arr, r)
				} else {
					arr = append(arr[1:], r)
				}
				time.Sleep(tt.args.time)
			}
			s.StopAutoMeasure()
			ssz, tot, ln := s.rate.Values()
			for x := 0; x < len(arr); x++ {
				tt.want.total += arr[x]
			}
			assert.Equal(t, tt.want.total, tot)
			assert.Equal(t, tt.want.length, ln)
			assert.Equal(t, ssz, ss)
		})
	}
}

func Test_speed_RawRate(t *testing.T) {
	type fields struct {
		target *int64
		rate   benchmark.SingleRate
	}
	type args struct {
		seconds int64
		length  int
	}
	ss := uint(5)
	var tgt int64
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "simple measure",
			fields: fields{
				target: &tgt,
				rate:   benchmark.NewSingleRate(ss),
			},
			args: args{
				seconds: 1,
				length:  3,
			},
		},
	}
	for _, tt := range tests {
		tgt = 0
		t.Run(tt.name, func(t *testing.T) {
			var r, total int64
			s := &speed{
				target: tt.fields.target,
				rate:   tt.fields.rate,
			}
			for x := 0; x < tt.args.length; x++ {
				r = rand.Int63n(5000)
				end := s.StartMeasure()
				time.Sleep(time.Duration(tt.args.seconds) * time.Second)
				tgt += r / int64(tt.args.seconds)
				total += r / int64(tt.args.seconds)
				end()
			}
			got := s.RawRate()
			assert.Equal(t, total/int64(tt.args.length), got)
		})
	}
}

func Test_speed_UnitsFunc(t *testing.T) {
	type fields struct {
		target *int64
		rate   benchmark.SingleRate
	}
	type args struct {
		fn func(int64) string
	}
	var tgt int64
	var ss uint = 5
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "simple function test",
			fields: fields{
				target: &tgt,
				rate:   benchmark.NewSingleRate(ss),
			},
			args: args{
				fn: func(x int64) string {
					return fmt.Sprintf("%d", x)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &speed{
				target: tt.fields.target,
				rate:   tt.fields.rate,
			}
			s.UnitsFunc(tt.args.fn)
			var total int64
			for x := 0; x < 3; x++ {
				end := s.StartMeasure()
				r := rand.Int63n(5000)
				tgt += r
				total += r
				time.Sleep(1 * time.Second)
				end()
				got := s.Rate()
				assert.Equal(t, strconv.Itoa(int(total/int64(x+1))), got)
			}
		})
	}
}

func Test_speed_Reset(t *testing.T) {
	type fields struct {
		target *int64
		rate   benchmark.SingleRate
	}
	type args struct {
		fn func(int64) string
	}
	type want struct {
		sampleSize uint
		total      int64
		listLength int
	}
	var tgt int64
	var ss uint = 5
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "simple function test",
			fields: fields{
				target: &tgt,
				rate:   benchmark.NewSingleRate(ss),
			},
			args: args{
				fn: func(x int64) string {
					return fmt.Sprintf("%d", x)
				},
			},
			want: want{
				sampleSize: ss,
				total:      0,
				listLength: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &speed{
				target: tt.fields.target,
				rate:   tt.fields.rate,
			}
			s.UnitsFunc(tt.args.fn)
			var total int64
			end := s.StartMeasure()
			r := rand.Int63n(5000)
			tgt += r
			total += r
			time.Sleep(1 * time.Second)
			end()
			s.Reset()
			ss, tot, ll := s.rate.Values()
			assert.Equal(t, tt.want.sampleSize, ss)
			assert.Equal(t, tt.want.total, tot)
			assert.Equal(t, tt.want.listLength, ll)
		})
	}
}
