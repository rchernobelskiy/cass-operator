package v1alpha1

import (
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_makeImage(t *testing.T) {
	type args struct {
		repo    string
		version string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test empty inputs",
			args: args{
				repo:    "",
				version: "",
			},
			want: "datastax/dse-server:6.7.3",
		},
		{
			name: "test different public version",
			args: args{
				repo:    "",
				version: "6.8",
			},
			want: "datastax/dse-server:6.8",
		},
		{
			name: "test private repo server",
			args: args{
				repo:    "datastax.jfrog.io/secret-debug-image/dse-server",
				version: "",
			},
			want: "datastax.jfrog.io/secret-debug-image/dse-server:6.7.3",
		},
		{
			name: "test fully custom params",
			args: args{
				repo:    "jfrog.io:6789/dse-server-team/dse-server",
				version: "master.20190605.123abc",
			},
			want: "jfrog.io:6789/dse-server-team/dse-server:master.20190605.123abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeImage(tt.args.repo, tt.args.version); got != tt.want {
				t.Errorf("makeImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDseDatacenter_GetServerImage(t *testing.T) {
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       DseDatacenterSpec
		Status     DseDatacenterStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple test",
			fields: fields{
				Spec: DseDatacenterSpec{
					Repository: "jfrog.io:6789/dse-server-team/dse-server",
					Version:    "master.20190605.123abc",
				},
			},
			want: "jfrog.io:6789/dse-server-team/dse-server:master.20190605.123abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := &DseDatacenter{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if got := dc.GetServerImage(); got != tt.want {
				t.Errorf("DseDatacenter.GetServerImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDseDatacenter_GetSeedList(t *testing.T) {
	type fields struct {
		TypeMeta   metav1.TypeMeta
		ObjectMeta metav1.ObjectMeta
		Spec       DseDatacenterSpec
		Status     DseDatacenterStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "1 DC, 1 Rack, 1 Node",
			fields: fields{
				ObjectMeta: v1.ObjectMeta{
					Name:      "example-dsedatacenter",
					Namespace: "default_ns",
				},
				Spec: DseDatacenterSpec{
					Size: 1,
					Racks: []DseRack{{
						Name: "rack0",
					}},
				},
			},
			want: []string{"example-dsedatacenter-rack0-stateful-set-0.example-dsedatacenter-service.default_ns.svc.cluster.local"},
		}, {
			name: "1 DC, 2 Rack, 2 Node",
			fields: fields{
				ObjectMeta: v1.ObjectMeta{
					Name:      "example-dsedatacenter",
					Namespace: "default_ns",
				},
				Spec: DseDatacenterSpec{
					Size: 2,
					Racks: []DseRack{{
						Name: "rack0",
					}, {
						Name: "rack1",
					}},
				},
			},
			want: []string{"example-dsedatacenter-rack0-stateful-set-0.example-dsedatacenter-service.default_ns.svc.cluster.local",
				"example-dsedatacenter-rack1-stateful-set-0.example-dsedatacenter-service.default_ns.svc.cluster.local"},
		}, {
			name: "1 DC, 1 Rack, 2 Node",
			fields: fields{
				ObjectMeta: v1.ObjectMeta{
					Name:      "example-dsedatacenter",
					Namespace: "default_ns",
				},
				Spec: DseDatacenterSpec{
					Size: 2,
					Racks: []DseRack{{
						Name: "rack0",
					}},
				},
			},
			want: []string{"example-dsedatacenter-rack0-stateful-set-0.example-dsedatacenter-service.default_ns.svc.cluster.local",
				"example-dsedatacenter-rack0-stateful-set-1.example-dsedatacenter-service.default_ns.svc.cluster.local"},
		}, {
			name: "1 DC, 3 Rack, 6 Node",
			fields: fields{
				ObjectMeta: v1.ObjectMeta{
					Name:      "example-dsedatacenter",
					Namespace: "default_ns",
				},
				Spec: DseDatacenterSpec{
					Size: 6,
					Racks: []DseRack{{
						Name: "rack0",
					}, {
						Name: "rack1",
					}, {
						Name: "rack2",
					}},
				},
			},
			want: []string{"example-dsedatacenter-rack0-stateful-set-0.example-dsedatacenter-service.default_ns.svc.cluster.local",
				"example-dsedatacenter-rack1-stateful-set-0.example-dsedatacenter-service.default_ns.svc.cluster.local",
				"example-dsedatacenter-rack2-stateful-set-0.example-dsedatacenter-service.default_ns.svc.cluster.local"},
		}, {
			name: "1 DC, 0 Rack, 0 Node",
			fields: fields{
				ObjectMeta: v1.ObjectMeta{
					Name:      "example-dsedatacenter",
					Namespace: "default_ns",
				},
				Spec: DseDatacenterSpec{
					Size:  0,
					Racks: []DseRack{},
				},
			},
			want: []string{},
		}, {
			name: "1 DC, 3 Rack, 3 Node",
			fields: fields{
				ObjectMeta: v1.ObjectMeta{
					Name:      "example-dsedatacenter",
					Namespace: "default_ns",
				},
				Spec: DseDatacenterSpec{
					Size: 3,
					Racks: []DseRack{{
						Name: "rack0",
					}, {
						Name: "rack1",
					}, {
						Name: "rack2",
					}},
				},
			},
			want: []string{"example-dsedatacenter-rack0-stateful-set-0.example-dsedatacenter-service.default_ns.svc.cluster.local",
				"example-dsedatacenter-rack1-stateful-set-0.example-dsedatacenter-service.default_ns.svc.cluster.local",
				"example-dsedatacenter-rack2-stateful-set-0.example-dsedatacenter-service.default_ns.svc.cluster.local"},
		}, {
			name: "1 DC, 0 Rack, 1 Node",
			fields: fields{
				ObjectMeta: v1.ObjectMeta{
					Name:      "example-dsedatacenter",
					Namespace: "default_ns",
				},
				Spec: DseDatacenterSpec{
					Size: 1,
				},
			},
			want: []string{"example-dsedatacenter-default-stateful-set-0.example-dsedatacenter-service.default_ns.svc.cluster.local"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := &DseDatacenter{
				TypeMeta:   tt.fields.TypeMeta,
				ObjectMeta: tt.fields.ObjectMeta,
				Spec:       tt.fields.Spec,
				Status:     tt.fields.Status,
			}
			if got := in.GetSeedList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DseDatacenter.GetSeedList() = %v, want %v", got, tt.want)
			}
		})
	}
}