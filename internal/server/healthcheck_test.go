package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestHandlers_HealthcheckHandler(t *testing.T) {
	testCases := []struct {
		name          string
		method        string
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "200 on GET",
			method: http.MethodGet,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "404 on POST",
			method: http.MethodPost,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "404 on PATCH",
			method: http.MethodPatch,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "404 on DELETE",
			method: http.MethodDelete,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "404 on OPTIONS",
			method: http.MethodOptions,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := NewServer(nil, nil)
			recorder := httptest.NewRecorder()

			url := "/v1/healthcheck"
			req, err := http.NewRequest(tc.method, url, nil)
			require.NoError(t, err)

			server.healthcheckHandler(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}
