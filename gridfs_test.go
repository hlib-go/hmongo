package hmongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"testing"
)

func TestFSDelete(t *testing.T) {
	type args struct {
		db     *mongo.Database
		fileId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := FSDelete(tt.args.db, tt.args.fileId); (err != nil) != tt.wantErr {
				t.Errorf("FSDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFSDownload(t *testing.T) {
	type args struct {
		db     *mongo.Database
		fileId string
	}
	tests := []struct {
		name      string
		args      args
		wantBytes []byte
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBytes, err := FSDownload(tt.args.db, tt.args.fileId)
			if (err != nil) != tt.wantErr {
				t.Errorf("FSDownload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBytes, tt.wantBytes) {
				t.Errorf("FSDownload() gotBytes = %v, want %v", gotBytes, tt.wantBytes)
			}
		})
	}
}

func TestFSUpload(t *testing.T) {
	type args struct {
		db        *mongo.Database
		filename  string
		fileBytes []byte
	}
	tests := []struct {
		name       string
		args       args
		wantFileId string
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFileId, err := FSUpload(tt.args.db, tt.args.filename, tt.args.fileBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("FSUpload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFileId != tt.wantFileId {
				t.Errorf("FSUpload() gotFileId = %v, want %v", gotFileId, tt.wantFileId)
			}
		})
	}
}
