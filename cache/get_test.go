package cache

import (
	"github.com/lefelys/state"
	"github.com/rs/zerolog"
	"reflect"
	"testing"
	"time"
)

func TestClient_Get(t *testing.T) {
	type fields struct {
		log zerolog.Logger
		st  state.State
	}
	type args struct {
		bucket string
		key    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		sleep   time.Duration
		put     bool
		ttl     time.Duration
		want    []byte
		wantErr bool
	}{
		{
			name: "no value",
			fields: fields{
				log: zerolog.Logger{},
				st:  state.Empty(),
			},
			args: args{
				bucket: "get",
				key:    "test",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with value",
			fields: fields{
				log: zerolog.Logger{},
				st:  state.Empty(),
			},
			args: args{
				bucket: "get",
				key:    "test",
			},
			sleep:   1 * time.Second,
			put:     true,
			ttl:     2 * time.Second,
			want:    []byte("test"),
			wantErr: false,
		},
		{
			name: "no value post ttl",
			fields: fields{
				log: zerolog.Logger{},
				st:  state.Empty(),
			},
			args: args{
				bucket: "get",
				key:    "test",
			},
			sleep:   1 * time.Second,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				log: tt.fields.log,
				st:  tt.fields.st,
				db:  newDb(t, "get"),
			}

			if tt.put {
				if err := c.Put(tt.args.bucket, tt.args.key, tt.want, tt.ttl); (err != nil) != tt.wantErr && tt.sleep == 0 {
					t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			time.Sleep(tt.sleep)

			got, err := c.Get(tt.args.bucket, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}

			if err := c.Close(); err != nil {
				t.Errorf("Close() error = %v, wantErr %v", err, nil)
			}
		})
	}
}
