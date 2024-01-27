package api

import (
	"arch-demo/internal/domain"
	mock_api "arch-demo/internal/tests/api_mocks"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateMovie(t *testing.T) {
	type fields struct {
		name        string
		releaseDate time.Time
		country     string
		genre       string
		rating      int8
	}

	testCases := []struct {
		name          string
		fields        fields
		mockInit      func(s *mock_api.MockMoviesService)
		header        http.Header
		expStatusCode int
		expErrMessage string
		expErr        bool
	}{
		{
			name:     "fail_content_type",
			fields:   fields{},
			mockInit: func(s *mock_api.MockMoviesService) {},
			header: http.Header{
				"Content-Type": []string{
					"text/plain",
				},
			},
			expStatusCode: http.StatusUnsupportedMediaType,
			expErrMessage: "content type not allowed\n",
			expErr:        true,
		},
		{
			name:   "movie_already_exists",
			fields: fields{},
			mockInit: func(s *mock_api.MockMoviesService) {
				s.EXPECT().Create(gomock.Any()).Return(domain.Movie{}, domain.ErrExists)
			},
			header: http.Header{
				"Content-Type": []string{
					"application/json",
				},
			},
			expStatusCode: http.StatusConflict,
			expErrMessage: "movie already exists\n",
			expErr:        true,
		},
		{
			name:   "internal_error",
			fields: fields{},
			mockInit: func(s *mock_api.MockMoviesService) {
				s.EXPECT().Create(gomock.Any()).Return(domain.Movie{}, errors.New("unexpected error"))
			},
			header: http.Header{
				"Content-Type": []string{
					"application/json",
				},
			},
			expStatusCode: http.StatusInternalServerError,
			expErrMessage: "unexpected error\n",
			expErr:        true,
		},
		{
			name: "success_create",
			fields: fields{
				name:        "Name",
				releaseDate: time.Time{},
				country:     "Russia",
				genre:       "Haha",
				rating:      5,
			},
			mockInit: func(s *mock_api.MockMoviesService) {
				s.EXPECT().Create(gomock.Any()).Return(domain.Movie{
					Name:        "Name",
					ReleaseDate: time.Time{},
					Country:     "Russia",
					Genre:       "Haha",
					Rating:      5,
				}, nil)
			},
			header: http.Header{
				"Content-Type": []string{
					"application/json",
				},
			},
			expStatusCode: http.StatusCreated,
			expErr:        false,
		},
		{
			name: "not_all_fields",
			fields: fields{
				name:        "Name",
				releaseDate: time.Time{},
				country:     "Russia",
				genre:       "Haha",
			},
			mockInit: func(s *mock_api.MockMoviesService) {
				s.EXPECT().Create(gomock.Any()).Return(domain.Movie{}, domain.ErrFieldsRequired)
			},
			header: http.Header{
				"Content-Type": []string{
					"application/json",
				},
			},
			expStatusCode: http.StatusUnprocessableEntity,
			expErrMessage: "all required fields must have values\n",
			expErr:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			s := mock_api.NewMockMoviesService(ctrl)
			if tc.mockInit != nil {
				tc.mockInit(s)
			}

			h := NewMoviesHandler(s)
			payload := domain.Movie{
				Name:        tc.fields.name,
				ReleaseDate: tc.fields.releaseDate,
				Country:     tc.fields.country,
				Genre:       tc.fields.genre,
				Rating:      tc.fields.rating,
			}

			body, _ := json.Marshal(payload)
			req := httptest.NewRequest(http.MethodPost, "/api/", bytes.NewReader(body))
			req.Header = tc.header
			recorder := httptest.NewRecorder()

			h.Create(recorder, req)
			if recorder.Result().StatusCode != tc.expStatusCode {
				t.Errorf("expected status code: %d, got: %d", tc.expStatusCode, recorder.Result().StatusCode)
			}

			if tc.expErr {
				respBody := recorder.Body.Bytes()
				msg := string(respBody)

				if msg != tc.expErrMessage {
					t.Errorf("expected error message: %s, got: %s", tc.expErrMessage, msg)
				}
				return
			}

			var movie domain.Movie
			_ = json.NewDecoder(recorder.Result().Body).Decode(&movie)

			if tc.fields.name != movie.Name &&
				tc.fields.releaseDate != movie.ReleaseDate &&
				tc.fields.country != movie.Country &&
				tc.fields.genre != movie.Genre &&
				tc.fields.rating != movie.Rating {
				t.Errorf("expected movie: %v, got: %v", tc.fields, movie)
			}
		})
	}
}

func toPtr(str string) *string {
	return &str
}
