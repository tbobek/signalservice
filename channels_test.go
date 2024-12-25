package main

import (
	"testing"
)

// to run this test successfully, the file channels.json must exist in the config folder
// and the signal service has to run with this config file loaded
func TestReadChannels(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want Channels
	}{
		{
			name: "simple test",
			args: args{
				filename: "./testconfig/channels.json",
			},
			want: Channels{
				ModelName:    "Test",
				ModelId:      "",
				ModelVersion: "",
				Variables:    []Variable{},
				Locations:    []Location{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadChannels(tt.args.filename)
			if err != nil {
				t.Errorf("ReadChannels() error = %v", err)
			}
			if got.ModelName != tt.want.ModelName {
				t.Errorf("ReadChannels() = %v, want %v", got, tt.want)
			}
		})
	}
}
