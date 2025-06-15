package services

import (
	"context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/sbilibin2017/yandex-go-advanced/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestMetricUpdateService_Update_Table(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ptrInt64 := func(i int64) *int64 {
		return &i
	}

	ptrFloat64 := func(f float64) *float64 {
		return &f
	}

	type fields struct {
		saver  *MockMetricUpdateSaver
		getter *MockMetricUpdateGetter
	}
	type args struct {
		metrics []types.Metrics
	}
	type want struct {
		err bool
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
		setup  func(f fields, args args)
	}{
		{
			name: "counter metric with existing value",
			fields: fields{
				saver:  NewMockMetricUpdateSaver(ctrl),
				getter: NewMockMetricUpdateGetter(ctrl),
			},
			args: args{
				metrics: []types.Metrics{
					{
						ID:    "metric1",
						MType: types.Counter,
						Delta: ptrInt64(10),
					},
				},
			},
			want: want{err: false},
			setup: func(f fields, args args) {
				existing := &types.Metrics{
					ID:    "metric1",
					MType: types.Counter,
					Delta: ptrInt64(5),
				}
				f.getter.EXPECT().
					Get(gomock.Any(), types.MetricID{ID: "metric1", MType: types.Counter}).
					Return(existing, nil)
				f.saver.EXPECT().
					Save(gomock.Any(), types.Metrics{
						ID:    "metric1",
						MType: types.Counter,
						Delta: ptrInt64(15),
					}).
					Return(nil)
			},
		},
		{
			name: "gauge metric saves as is",
			fields: fields{
				saver:  NewMockMetricUpdateSaver(ctrl),
				getter: NewMockMetricUpdateGetter(ctrl),
			},
			args: args{
				metrics: []types.Metrics{
					{
						ID:    "metric2",
						MType: types.Gauge,
						Value: ptrFloat64(3.14),
					},
				},
			},
			want: want{err: false},
			setup: func(f fields, args args) {
				f.saver.EXPECT().
					Save(gomock.Any(), args.metrics[0]).
					Return(nil)
			},
		},
		{
			name: "getter returns error",
			fields: fields{
				saver:  NewMockMetricUpdateSaver(ctrl),
				getter: NewMockMetricUpdateGetter(ctrl),
			},
			args: args{
				metrics: []types.Metrics{
					{
						ID:    "metric3",
						MType: types.Counter,
						Delta: ptrInt64(1),
					},
				},
			},
			want: want{err: true},
			setup: func(f fields, args args) {
				f.getter.EXPECT().
					Get(gomock.Any(), types.MetricID{ID: "metric3", MType: types.Counter}).
					Return(nil, errors.New("getter error"))
			},
		},
		{
			name: "saver returns error",
			fields: fields{
				saver:  NewMockMetricUpdateSaver(ctrl),
				getter: NewMockMetricUpdateGetter(ctrl),
			},
			args: args{
				metrics: []types.Metrics{
					{
						ID:    "metric4",
						MType: types.Gauge,
						Value: ptrFloat64(2.71),
					},
				},
			},
			want: want{err: true},
			setup: func(f fields, args args) {
				f.saver.EXPECT().
					Save(gomock.Any(), args.metrics[0]).
					Return(errors.New("save error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt.fields, tt.args)
			}
			svc := NewMetricUpdateService(tt.fields.saver, tt.fields.getter)

			err := svc.Update(context.Background(), tt.args.metrics)
			if tt.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
