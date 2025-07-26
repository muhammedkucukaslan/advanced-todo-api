package domain

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestNewTodo(t *testing.T) {
	type args struct {
		userId uuid.UUID
		title  string
	}
	tests := []struct {
		name    string
		args    args
		want    *Todo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTodo(tt.args.userId, tt.args.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTodo() = %v, want %v", got, tt.want)
			}
		})
	}
}
