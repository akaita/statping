package handlers

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/statping/statping/types"
	"github.com/statping/statping/types/services"
	"github.com/statping/statping/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestApiServiceRoutes(t *testing.T) {

	since := utils.Now().Add(-30 * types.Day)
	startEndQuery := fmt.Sprintf("?start=%d&end=%d", since.Unix(), utils.Now().Unix())

	tests := []HTTPTest{
		{
			Name:             "Statping All Public and Private Services",
			URL:              "/api/services",
			Method:           "GET",
			ExpectedContains: []string{`"name":"Google"`},
			ExpectedStatus:   200,
			ResponseLen:      5,
			BeforeTest:       SetTestENV,
			FuncTest: func(t *testing.T) error {
				count := len(services.Services())
				if count != 5 {
					return errors.Errorf("incorrect services count: %d", count)
				}
				return nil
			},
		},
		{
			Name:             "Statping All Public Services",
			URL:              "/api/services",
			Method:           "GET",
			ExpectedContains: []string{`"name":"Google"`},
			ExpectedStatus:   200,
			ResponseLen:      4,
			BeforeTest:       UnsetTestENV,
			FuncTest: func(t *testing.T) error {
				count := len(services.Services())
				if count != 5 {
					return errors.Errorf("incorrect services count: %d", count)
				}
				return nil
			},
		},
		{
			Name:             "Statping Public Service 1",
			URL:              "/api/services/1",
			Method:           "GET",
			ExpectedContains: []string{`"name":"Google"`},
			ExpectedStatus:   200,
			BeforeTest:       UnsetTestENV,
		},
		{
			Name:             "Statping Private Service 1",
			URL:              "/api/services/2",
			Method:           "GET",
			ExpectedContains: []string{`"error":"not authenticated"`},
			ExpectedStatus:   200,
			BeforeTest:       UnsetTestENV,
		},
		{
			Name:             "Statping Service 1 with Private responses",
			URL:              "/api/services/1",
			Method:           "GET",
			ExpectedContains: []string{`"name":"Google"`},
			ExpectedStatus:   200,
			BeforeTest:       SetTestENV,
		},
		{
			Name:           "Statping Service Failures",
			URL:            "/api/services/1/failures",
			Method:         "GET",
			ResponseLen:    125,
			ExpectedStatus: 200,
		},
		{
			Name:           "Statping Service Failures Limited",
			URL:            "/api/services/1/failures?limit=1",
			Method:         "GET",
			ResponseLen:    1,
			ExpectedStatus: 200,
		},
		{
			Name:           "Statping Service 1 Data",
			URL:            "/api/services/1/hits_data" + startEndQuery,
			Method:         "GET",
			ResponseLen:    73,
			ExpectedStatus: 200,
		},
		{
			Name:           "Statping Service 1 Ping Data",
			URL:            "/api/services/1/ping_data" + startEndQuery,
			Method:         "GET",
			ResponseLen:    73,
			ExpectedStatus: 200,
		},
		{
			Name:           "Statping Service 1 Failure Data - 24 Hour",
			URL:            "/api/services/1/failure_data" + startEndQuery + "&group=24h",
			Method:         "GET",
			ResponseLen:    4,
			ExpectedStatus: 200,
		},
		{
			Name:           "Statping Service 1 Failure Data - 12 Hour",
			URL:            "/api/services/1/failure_data" + startEndQuery + "&group=12h",
			Method:         "GET",
			ResponseLen:    7,
			ExpectedStatus: 200,
		},
		{
			Name:           "Statping Service 1 Failure Data - 1 Hour",
			URL:            "/api/services/1/failure_data" + startEndQuery + "&group=1h",
			Method:         "GET",
			ResponseLen:    73,
			ExpectedStatus: 200,
		},
		{
			Name:           "Statping Service 1 Failure Data - 15 Minute",
			URL:            "/api/services/1/failure_data" + startEndQuery + "&group=15m",
			Method:         "GET",
			ResponseLen:    124,
			ExpectedStatus: 200,
		},
		{
			Name:           "Statping Service 1 Hits",
			URL:            "/api/services/1/hits_data" + startEndQuery,
			Method:         "GET",
			ResponseLen:    73,
			ExpectedStatus: 200,
		},
		{
			Name:           "Statping Service 1 Failure Data",
			URL:            "/api/services/1/failure_data" + startEndQuery,
			Method:         "GET",
			ResponseLen:    73,
			ExpectedStatus: 200,
		},
		{
			Name:           "Statping Reorder Services",
			URL:            "/api/reorder/services",
			Method:         "POST",
			Body:           `[{"service":1,"order":1},{"service":4,"order":2},{"service":2,"order":3},{"service":3,"order":4}]`,
			ExpectedStatus: 200,
			HttpHeaders:    []string{"Content-Type=application/json"},
			SecureRoute:    true,
		},
		{
			Name:        "Statping Create Service",
			URL:         "/api/services",
			HttpHeaders: []string{"Content-Type=application/json"},
			Method:      "POST",
			Body: `{
					"name": "New Private Service",
					"domain": "https://statping.com",
					"expected": "",
					"expected_status": 200,
					"check_interval": 30,
					"type": "http",
					"public": false,
					"group_id": 1,
					"method": "GET",
					"post_data": "",
					"port": 0,
					"timeout": 30,
					"order_id": 0
				}`,
			ExpectedStatus:   200,
			ExpectedContains: []string{`"status":"success","type":"service","method":"create"`, `"public":false`, `"group_id":1`},
			FuncTest: func(t *testing.T) error {
				count := len(services.Services())
				if count != 6 {
					return errors.Errorf("incorrect services count: %d", count)
				}
				return nil
			},
			SecureRoute: true,
		},
		{
			Name:        "Statping Update Service",
			URL:         "/api/services/1",
			HttpHeaders: []string{"Content-Type=application/json"},
			Method:      "POST",
			Body: `{
					"name": "Updated New Service",
					"domain": "https://google.com",
					"expected": "",
					"expected_status": 200,
					"check_interval": 60,
					"type": "http",
					"method": "GET",
					"post_data": "",
					"port": 0,
					"timeout": 10,
					"order_id": 0
				}`,
			ExpectedStatus:   200,
			ExpectedContains: []string{`"status":"success"`, `"name":"Updated New Service"`, `"method":"update"`},
			FuncTest: func(t *testing.T) error {
				item, err := services.Find(1)
				require.Nil(t, err)
				if item.Interval != 60 {
					return errors.Errorf("incorrect service check interval: %d", item.Interval)
				}
				return nil
			},
			SecureRoute: true,
		},
		{
			Name:             "Statping Delete Failures",
			URL:              "/api/services/1/failures",
			Method:           "DELETE",
			ExpectedStatus:   200,
			ExpectedContains: []string{`"status":"success"`, `"method":"delete_failures"`},
			FuncTest: func(t *testing.T) error {
				item, err := services.Find(1)
				require.Nil(t, err)
				fails := item.AllFailures().Count()
				if fails != 0 {
					return errors.Errorf("incorrect service failures count: %d", fails)
				}
				return nil
			},
			SecureRoute: true,
		},
		{
			Name:             "Statping Delete Service",
			URL:              "/api/services/1",
			Method:           "DELETE",
			ExpectedStatus:   200,
			ExpectedContains: []string{`"status":"success"`, `"method":"delete"`},
			FuncTest: func(t *testing.T) error {
				count := len(services.Services())
				if count != 5 {
					return errors.Errorf("incorrect services count: %d", count)
				}
				return nil
			},
			SecureRoute: true,
		}}

	for _, v := range tests {
		t.Run(v.Name, func(t *testing.T) {
			_, t, err := RunHTTPTest(v, t)
			assert.Nil(t, err)
		})
	}
}
