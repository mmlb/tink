package grpcserver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tinkerbell/tink/db/mock"
	pb "github.com/tinkerbell/tink/protos/template"
)

const (
	template1 = `version: "0.1"
name: hello_world_workflow
global_timeout: 600
tasks:
  - name: "hello world"
    worker: "{{.device_1}}"
    actions:
    - name: "hello_world"
      image: hello-world
      timeout: 60`

	template2 = `version: "0.1"
name: hello_world_workflow
global_timeout: 600
tasks:
  - name: "hello world again"
    worker: "{{.device_1}}"
    actions:
  	- name: "hello_world_again"
      image: hello-world
      timeout: 60`
)

func TestDuplicateTemplateName(t *testing.T) {
	type (
		args struct {
			db   mock.DB
			name string
		}
		want struct {
			expectedError bool
		}
	)
	testCases := map[string]struct {
		args args
		want want
	}{
		"test_1": {
			args: args{
				db:   mock.DB{},
				name: "template_1",
			},
			want: want{
				expectedError: true,
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			s := testServer(tc.args.db)
			res, err := s.CreateTemplate(context.TODO(), &pb.WorkflowTemplate{Name: tc.args.name, Data: template1})
			assert.Nil(t, err)
			assert.NotNil(t, res)
			if err == nil {
				_, err = s.CreateTemplate(context.TODO(), &pb.WorkflowTemplate{Name: tc.args.name, Data: template2})
			}
			if err != nil {
				assert.Error(t, err)
				assert.True(t, tc.want.expectedError)
			}
		})
	}
}
