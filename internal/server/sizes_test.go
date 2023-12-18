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

func TestSizesHandler_listSizes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name          string
		method        string
		buildStubs    func(repo *mock.MockSizer)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "200 on get",
			method: http.MethodGet,
			buildStubs: func(sizer *mock.MockSizer) {
				sizer.EXPECT().ListSizes().Times(1)
			},
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

			newSizerSrvc := packer.NewSizerService(packer.SortedSizes)
			newPackerSrvc := packer.NewPacketsService(newSizerSrvc)
			server := NewServer(newSizerSrvc, newPackerSrvc)
			recorder := httptest.NewRecorder()

			url := "/api/v1/tables"
			req, err := http.NewRequest(tc.method, url, nil)
			require.NoError(t, err)

			server.listSizesHandler(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestSizesHandler_addSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name          string
		method        string
		sizeToAdd     int
		buildStubs    func(repo *mock.MockSizer)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "200 on POST - positive num",
			method:    http.MethodGet,
			sizeToAdd: 100,
			buildStubs: func(sizer *mock.MockSizer) {
				sizer.EXPECT().AddSize(context.Background(), 0).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:      "422 on POST - 0",
			method:    http.MethodGet,
			sizeToAdd: 0,
			buildStubs: func(sizer *mock.MockSizer) {
				sizer.EXPECT().AddSize(context.Background(), 1).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name:      "422 on POST - (-1)",
			method:    http.MethodGet,
			sizeToAdd: 0,
			buildStubs: func(sizer *mock.MockSizer) {
				sizer.EXPECT().AddSize(context.Background(), 1).Times(1)
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

			url := "/api/v1/sizes"
			reqBody := map[string]int{"size": tc.sizeToAdd}
			var buf bytes.Buffer
			_ = json.NewEncoder(&buf).Encode(reqBody)
			req, err := http.NewRequest(tc.method, url, &buf)
			require.NoError(t, err)

			server.addSizeHandler(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestSizesHandler_putSizeHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name          string
		method        string
		sizesToAdd    []int
		buildStubs    func(repo *mock.MockSizer)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:       "200 on PUT - positive nums",
			method:     http.MethodPut,
			sizesToAdd: []int{100, 200, 300},
			buildStubs: func(sizer *mock.MockSizer) {
				sizer.EXPECT().PutSizes(context.Background(), []int{1, 2, 3}).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:       "422 on PUT - 0",
			method:     http.MethodGet,
			sizesToAdd: []int{100, 0, 300},
			buildStubs: func(sizer *mock.MockSizer) {
				sizer.EXPECT().AddSize(context.Background(), 1).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name:       "422 on PUT - (-1)",
			method:     http.MethodGet,
			sizesToAdd: []int{100, -1, 300},
			buildStubs: func(sizer *mock.MockSizer) {
				sizer.EXPECT().AddSize(context.Background(), 1).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
			},
		},
		{
			name:       "400 on PUT - duplicated sizes",
			method:     http.MethodGet,
			sizesToAdd: []int{100, 100, 300},
			buildStubs: func(sizer *mock.MockSizer) {
				sizer.EXPECT().AddSize(context.Background(), 1).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			url := "/api/v1/sizes"
			reqBody := map[string][]int{"sizes": tc.sizesToAdd}
			var buf bytes.Buffer
			_ = json.NewEncoder(&buf).Encode(reqBody)
			req, err := http.NewRequest(tc.method, url, &buf)
			require.NoError(t, err)

			server.putSizesHandler(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}
