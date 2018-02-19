package netem_parser

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name    string
		args    args
		want    *NetemData
		wantErr bool
	}{
		{
			"Some packet drop",
			args{`qdisc netem 8013: root refcnt 2 limit 1000 loss 3%
 Sent 100520 bytes 1036 pkt (dropped 34, overlimits 0 requeues 0)
 backlog 0b 0p requeues 0`},
			&NetemData{Total: 1036, Dropped: 34, Reordered: 0, Bytes: 100520},
			false,
		},
		{
			"Some packet drop and reordering",
			args{`qdisc netem 8012: root refcnt 2 limit 1000 delay 100.0ms loss 10% duplicate 10% reorder 25% gap 1
 Sent 98126 bytes 1011 pkt (dropped 112, overlimits 0 requeues 242)
 backlog 0b 0p requeues 242`},
			&NetemData{Total: 1011, Dropped: 112, Reordered: 242, Bytes: 98126},
			false,
		},
		{
			"Not a netem stat",
			args{`qdisc noqueue 0: root refcnt 2
 Sent 1024 bytes 1 pkt (dropped 0, overlimits 0 requeues 0)
 backlog 0b 0p requeues 0`},
			nil,
			true,
		},
		{
			"Empty",
			args{``},
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
