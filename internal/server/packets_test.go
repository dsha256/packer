package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dsha256/packer/internal/mock"
	"github.com/dsha256/packer/internal/packer"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestPacketsHandler_getPacks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name          string
		method        string
		items         int
		buildStubs    func(repo *mock.MockPacker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "200 on POST - positive num",
			method: http.MethodGet,
			items:  100,
			buildStubs: func(sizer *mock.MockPacker) {
				sizer.EXPECT().GetPackets(context.Background(), 0).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "422 on POST - 0",
			method: http.MethodGet,
			items:  0,
			buildStubs: func(sizer *mock.MockPacker) {
				sizer.EXPECT().GetPackets(context.Background(), 1).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name:   "422 on POST - (-1)",
			method: http.MethodGet,
			items:  0,
			buildStubs: func(sizer *mock.MockPacker) {
				sizer.EXPECT().GetPackets(context.Background(), 1).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			newSizerSrvc := packer.NewSizerService(packer.SortedSizes)
			newPackerSrvc := packer.NewPacketsService(newSizerSrvc)
			server := NewServer(newSizerSrvc, newPackerSrvc)
			recorder := httptest.NewRecorder()

			url := "/api/v1/packets"
			reqBody := map[string]int{"items": tc.items}
			var buf bytes.Buffer
			_ = json.NewEncoder(&buf).Encode(reqBody)
			req, err := http.NewRequest(tc.method, url, &buf)
			require.NoError(t, err)

			server.getPacksHandler(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}
