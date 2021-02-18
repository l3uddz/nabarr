package cache

import (
	"github.com/lefelys/state"
	"github.com/rs/zerolog"
	"reflect"
	"testing"
	"time"
)

func TestClient_Put(t *testing.T) {
	type fields struct {
		log zerolog.Logger
		st  state.State
	}
	type args struct {
		bucket string
		key    string
		val    []byte
		ttl    time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		sleep   time.Duration
		want    []byte
		wantErr bool
	}{
		{
			name: "with ttl",
			fields: fields{
				log: zerolog.Logger{},
				st:  state.Empty(),
			},
			args: args{
				bucket: "put",
				key:    "test",
				val:    []byte("testing"),
				ttl:    50 * time.Millisecond,
			},
			want:    []byte("testing"),
			wantErr: false,
		},
		{
			name: "ttl timed out",
			fields: fields{
				log: zerolog.Logger{},
				st:  state.Empty(),
			},
			args: args{
				bucket: "put",
				key:    "test",
				val:    []byte("testing"),
				ttl:    1 * time.Second,
			},
			sleep:   2 * time.Second,
			want:    nil,
			wantErr: true,
		},
		{
			name: "no ttl",
			fields: fields{
				log: zerolog.Logger{},
				st:  state.Empty(),
			},
			args: args{
				bucket: "put",
				key:    "test",
				val:    []byte("testing"),
				ttl:    0,
			},
			want:    []byte("testing"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				log: tt.fields.log,
				st:  tt.fields.st,
				db:  newDb(t, "nabarr_put"),
			}

			if err := c.Put(tt.args.bucket, tt.args.key, tt.args.val, tt.args.ttl); (err != nil) != tt.wantErr && tt.sleep == 0 {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}

			time.Sleep(tt.sleep)

			got, err := c.Get(tt.args.bucket, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Put() get error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Put() got = %v, want %v", got, tt.want)
			}

			if err := c.Close(); err != nil {
				t.Errorf("Close() error = %v, wantErr %v", err, nil)
			}
		})
	}
}
