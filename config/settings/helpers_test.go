package settings

import (
	"testing"

	"github.com/spf13/viper"
)

func Test_bindEnvs(t *testing.T) {
	type SingleStruct struct {
		StringField string
		IntField    int
		BoolField   bool
	}

	type NestedStruct struct {
		StringField string
		IntField    int
		BoolField   bool
		Nested      SingleStruct
	}

	type NestedPointerStruct struct {
		StringField string
		IntField    int
		BoolField   bool
		Nested      *SingleStruct
	}

	type NestedEmbededStruct struct {
		SingleStruct
	}

	type UnexpectedStruct struct {
		ExportedField   string
		unexportedField string
	}

	type args struct {
		i     interface{}
		parts []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "single struct",
			args: args{
				i:     SingleStruct{},
				parts: []string{},
			},
			want: []string{
				"stringfield",
				"intfield",
				"boolfield",
			},
		},

		{
			name: "single pointer struct",
			args: args{
				i:     &SingleStruct{},
				parts: []string{},
			},
			want: []string{
				"stringfield",
				"intfield",
				"boolfield",
			},
		},
		{
			name: "nested struct",
			args: args{
				i:     NestedStruct{},
				parts: []string{},
			},
			want: []string{
				"stringfield",
				"intfield",
				"boolfield",
				"nested.stringfield",
				"nested.intfield",
				"nested.boolfield",
			},
		},
		{
			name: "nested pointer struct",
			args: args{
				i: &NestedPointerStruct{
					Nested: &SingleStruct{},
				},
				parts: []string{},
			},
			want: []string{
				"stringfield",
				"intfield",
				"boolfield",
				"nested.stringfield",
				"nested.intfield",
				"nested.boolfield",
			},
		},
		{
			name: "nested embeded struct",
			args: args{
				i:     NestedEmbededStruct{},
				parts: []string{},
			},
			want: []string{
				"singlestruct_stringfield",
				"singlestruct_intfield",
				"singlestruct_boolfield",
			},
		},
		{
			name: "unexpected struct",
			args: args{
				i:     UnexpectedStruct{},
				parts: []string{},
			},
			want: []string{
				"exportedfield",
			},
		},
		{
			name: "unitialized pointer struct",
			args: args{
				i:     &NestedPointerStruct{},
				parts: []string{},
			},
			want: []string{
				"stringfield",
				"intfield",
				"boolfield",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bindEnvs(tt.args.i, tt.args.parts...)
			keys := viper.AllKeys()
			for _, key := range tt.want {
				found := false
				for _, k := range keys {
					if k == key {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("bindEnvs() missing key: %s", key)
				}
			}
		})
	}
}
