package client

import (
	"testing"
)

func TestClient_GetNodes(t *testing.T) {
	tests := []struct {
		name      string
		wantNodes []NodeInfo
		wantErr   bool
	}{
		{
			name:      "test",
			wantNodes: []NodeInfo{},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(localKubeConfig)
			_, err = c.GetNodes()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
