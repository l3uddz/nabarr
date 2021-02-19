package cache

import (
	"github.com/rs/zerolog"
	"testing"
)

func TestClient_Delete(t *testing.T) {
	type fields struct {
		log zerolog.Logger
	}
	type args struct {
		bucket string
		key    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "key not present",
			fields: fields{
				log: zerolog.Logger{},
			},
			args: args{
				bucket: "delete",
				key:    "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				log: tt.fields.log,
				db:  newDb(t),
			}
			if err := c.Delete(tt.args.bucket, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
