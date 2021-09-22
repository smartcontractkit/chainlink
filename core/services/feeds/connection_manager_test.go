package feeds

import (
	"testing"
)

func Test_connectionsManager_IsConnected(t *testing.T) {
	type fields struct {
		connections map[int64]*connection
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "inactive connection exists",
			fields: fields{
				connections: map[int64]*connection{
					1: {
						connected: false,
					},
				},
			},
			args: args{
				id: 1,
			},
			want: false,
		},
		{
			name: "active connection exists",
			fields: fields{
				connections: map[int64]*connection{
					1: {
						connected: true,
					},
				},
			},
			args: args{
				id: 1,
			},
			want: true,
		},
		{
			name: "connection does not exist",
			fields: fields{
				connections: map[int64]*connection{
					1: {
						connected: true,
					},
				},
			},
			args: args{
				id: 2,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr := &connectionsManager{
				connections: tt.fields.connections,
			}
			if got := mgr.IsConnected(tt.args.id); got != tt.want {
				t.Errorf("IsConnected() = %v, want %v", got, tt.want)
			}
		})
	}
}
