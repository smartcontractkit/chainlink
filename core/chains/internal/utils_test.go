package internal

import (
	"reflect"
	"testing"
)

func TestNewPageToken(t *testing.T) {
	type args struct {
		t *PageToken
	}
	tests := []struct {
		name    string
		args    args
		want    *PageToken
		wantErr bool
	}{
		{
			name: "empty",
			args: args{t: &PageToken{}},
			want: &PageToken{Page: 0, Size: defaultSize},
		},
		{
			name: "page set, size unset",
			args: args{t: &PageToken{Page: 1}},
			want: &PageToken{Page: 1, Size: defaultSize},
		},
		{
			name: "page set, size set",
			args: args{t: &PageToken{Page: 3, Size: 10}},
			want: &PageToken{Page: 3, Size: 10},
		},
		{
			name: "page unset, size set",
			args: args{t: &PageToken{Size: 17}},
			want: &PageToken{Page: 0, Size: 17},
		},
	}
	for _, tt := range tests {
		enc := tt.args.t.Encode()
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPageToken(enc)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPageToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPageToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
