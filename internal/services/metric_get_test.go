package services

import (
	"context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestMetricgetService_Get_Table(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		getter *MockMetricGetGetter
	}

	type args struct {
		ctx context.Context
		id  types.MetricID
	}

	type want struct {
		metric *types.Metrics
		err    bool
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
		setup  func(f fields, args args)
	}{
		{
			name: "successfully get metric",
			fields: fields{
				getter: NewMockMetricGetGetter(ctrl),
			},
			args: args{
				ctx: context.Background(),
				id:  types.MetricID{ID: "metric1", Type: types.Counter},
			},
			want: want{
				metric: &types.Metrics{ID: "metric1", Type: types.Counter, Delta: func(i int64) *int64 { return &i }(10)},
				err:    false,
			},
			setup: func(f fields, args args) {
				f.getter.EXPECT().
					Get(args.ctx, args.id).
					Return(&types.Metrics{ID: "metric1", Type: types.Counter, Delta: func(i int64) *int64 { return &i }(10)}, nil)
			},
		},
		{
			name: "getter returns error",
			fields: fields{
				getter: NewMockMetricGetGetter(ctrl),
			},
			args: args{
				ctx: context.Background(),
				id:  types.MetricID{ID: "metric2", Type: types.Gauge},
			},
			want: want{
				metric: nil,
				err:    true,
			},
			setup: func(f fields, args args) {
				f.getter.EXPECT().
					Get(args.ctx, args.id).
					Return(nil, errors.New("getter error"))
			},
		},
		{
			name: "metric not found returns nil",
			fields: fields{
				getter: NewMockMetricGetGetter(ctrl),
			},
			args: args{
				ctx: context.Background(),
				id:  types.MetricID{ID: "metric3", Type: types.Counter},
			},
			want: want{
				metric: nil,
				err:    false,
			},
			setup: func(f fields, args args) {
				f.getter.EXPECT().
					Get(args.ctx, args.id).
					Return(nil, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt.fields, tt.args)
			}

			service := NewMetricGetService(tt.fields.getter)
			gotMetric, err := service.Get(tt.args.ctx, tt.args.id)

			if tt.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want.metric, gotMetric)
		})
	}
}
