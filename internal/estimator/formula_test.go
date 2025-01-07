package estimator

import (
	"math"
	"testing"

	"github.com/guregu/null/v5"

	"github.com/PoolHealth/PoolHealthServer/common"
)

func TestCalculateChlorine(t *testing.T) {
	type args struct {
		volume          float64
		lastMeasurement common.Measurement
		lastAdditives   common.Additives
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "Test 1",
			args: args{
				volume: 300000,
				lastMeasurement: common.Measurement{
					Chlorine: null.FloatFrom(4.76),
				},
				lastAdditives: common.Additives{
					Products: map[common.ChemicalProduct]float64{
						common.CalciumHypochlorite65Percent: 1,
					},
				},
			},
			want: 2.274 + 4.39,
		},
		{
			name: "Test 2",
			args: args{
				volume: 300000,
				lastMeasurement: common.Measurement{
					Chlorine: null.FloatFrom(2.8),
				},
				lastAdditives: common.Additives{
					Products: map[common.ChemicalProduct]float64{
						common.TCCA90PercentTablets: 2,
					},
				},
			},
			want: 4.012,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateChlorine(tt.args.volume, tt.args.lastMeasurement, tt.args.lastAdditives)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateChlorine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			roundedGot := math.Floor(got*1000) / 1000
			if roundedGot != tt.want {
				t.Errorf("CalculateChlorine() got = %v, want %v", roundedGot, tt.want)
			}
		})
	}
}
