package main

import (
	"github.com/golang/glog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseFlags(t *testing.T) {
	tcs := []struct{
		name string
		filePath string
		configPath string
		expErr string
	}{{
		name: "no flag set",
		expErr: "no file passed, pass file using -filePath flag",
	}}

	for _, tc := range tcs{
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseFlags()
			assert.Equal(t, err, tc.expErr)
		})
	}
	glog.Info(tcs)
	*configPath = "testConfig"
	main()
}
