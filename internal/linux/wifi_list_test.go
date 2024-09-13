package linux

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanWiFiPoints(t *testing.T) {
	type args struct {
		sc syscalWiFiList
	}
	tests := []struct {
		name       string
		args       args
		wantPoints int
		wantErr    bool
	}{
		{
			name:       "should return full wifi list",
			args:       args{sc: newTestSyscallWifiList("fulllist")},
			wantPoints: 15,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ScanWiFiPoints(tt.args.sc)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanWiFiPoints() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantPoints {
				t.Errorf("ScanWiFiPoints() = %v, want %v points", got, tt.wantPoints)
			}
		})
	}
}

type testSyscallWifiList struct {
	filename string
	cache    []byte
}

func newTestSyscallWifiList(filename string) *testSyscallWifiList {
	return &testSyscallWifiList{
		filename: filename,
	}
}

func (t *testSyscallWifiList) call() ([]byte, error) {
	if t.cache != nil {
		return t.cache, nil
	}
	fn := filepath.Join("./testdata", t.filename)
	data, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	t.cache = data
	return data, nil
}
