package sheets

import (
	"context"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/PoolHealth/PoolHealthServer/common"
	"github.com/PoolHealth/PoolHealthServer/pkg/log"
)

func TestSheetClient_HasSheet(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Test HasSheet",
			args: args{
				id: "17AB97B62cfI6EqnC3RPQqoMLFF7qgsF3_939o0ikhds",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(log.NewLogger(logrus.New()), ".google/")
			require.NoError(t, s.Start(context.Background()))
			got, err := s.HasSheet(context.Background(), tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("HasSheet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HasSheet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSheetClient_GetPools(t *testing.T) {
	type args struct {
		sheetID string
	}
	tests := []struct {
		name    string
		args    args
		want    []common.PoolData
		wantErr bool
	}{
		{
			name: "Test GetPools",
			args: args{
				sheetID: "17AB97B62cfI6EqnC3RPQqoMLFF7qgsF3_939o0ikhds",
			},
			want: []common.PoolData{
				{
					Name: "Pool 1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s := New(log.NewLogger(logrus.New()), ".google/")
			require.NoError(t, s.Start(ctx))
			got, err := s.GetPools(ctx, tt.args.sheetID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPools() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPools() got = %v, want %v", got, tt.want)
			}
		})
	}
}
