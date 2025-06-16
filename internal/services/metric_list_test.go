package services

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestMetricListService_List_Table(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		lister *MockMetricListLister
	}

	type args struct {
		ctx context.Context
	}

	type want struct {
		metrics []types.Metrics
		err     bool
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
		setup  func(f fields, args args)
	}{
		{
			name: "successfully list metrics",
			fields: fields{
				lister: NewMockMetricListLister(ctrl),
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				metrics: []types.Metrics{
					{ID: "metric1", Type: types.Counter, Delta: func(i int64) *int64 { return &i }(10)},
					{ID: "metric2", Type: types.Gauge, Value: func(f float64) *float64 { return &f }(3.14)},
				},
				err: false,
			},
			setup: func(f fields, args args) {
				f.lister.EXPECT().
					List(args.ctx).
					Return([]types.Metrics{
						{ID: "metric1", Type: types.Counter, Delta: func(i int64) *int64 { return &i }(10)},
						{ID: "metric2", Type: types.Gauge, Value: func(f float64) *float64 { return &f }(3.14)},
					}, nil)
			},
		},
		{
			name: "lister returns error",
			fields: fields{
				lister: NewMockMetricListLister(ctrl),
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				metrics: nil,
				err:     true,
			},
			setup: func(f fields, args args) {
				f.lister.EXPECT().
					List(args.ctx).
					Return(nil, errors.New("list error"))
			},
		},
		{
			name: "empty list returns no error",
			fields: fields{
				lister: NewMockMetricListLister(ctrl),
			},
			args: args{
				ctx: context.Background(),
			},
			want: want{
				metrics: []types.Metrics{},
				err:     false,
			},
			setup: func(f fields, args args) {
				f.lister.EXPECT().
					List(args.ctx).
					Return([]types.Metrics{}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt.fields, tt.args)
			}

			service := NewMetricListService(tt.fields.lister)
			gotMetrics, err := service.List(tt.args.ctx)

			if tt.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want.metrics, gotMetrics)
		})
	}
}
