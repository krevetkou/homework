package api

import (
	"arch-demo/internal/domain"
	mock_api "arch-demo/internal/tests/api_mocks"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"io"
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

func TestListMovies(t *testing.T) {
	testCases := []struct {
		name          string
		mockInit      func(s *mock_api.MockMoviesService)
		header        http.Header
		expStatusCode int
		expErrMessage string
		expErr        bool
	}{
		{
			name: "list_movies_success",
			mockInit: func(s *mock_api.MockMoviesService) {
				s.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.Movie{
					{
						Name:        "Name",
						ReleaseDate: time.Time{},
						Country:     "Russia",
						Genre:       "Haha",
						Rating:      5,
					},
					{
						Name:        "Name2",
						ReleaseDate: time.Time{},
						Country:     "Poland",
						Genre:       "Haha",
						Rating:      3,
					},
				}).AnyTimes()
			},
			header: http.Header{
				"Content-Type": []string{
					"application/json",
				},
			},
			expStatusCode: http.StatusOK,
			expErr:        false,
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

			req := httptest.NewRequest(http.MethodGet, "/api/", io.Reader(nil))
			req.Header = tc.header
			recorder := httptest.NewRecorder()

			h.List(recorder, req)
			if recorder.Result().StatusCode != tc.expStatusCode {
				t.Errorf("expected status code: %d, got: %d", tc.expStatusCode, recorder.Result().StatusCode)
			}

			var movies []domain.Movie
			_ = json.NewDecoder(recorder.Result().Body).Decode(&movies)

		})
	}
}

func TestGetMovie(t *testing.T) {
	testCases := []struct {
		name          string
		mockInit      func(s *mock_api.MockMoviesService)
		header        http.Header
		expStatusCode int
		expErrMessage string
		expErr        bool
	}{
		{
			name: "get_actor_success",
			mockInit: func(s *mock_api.MockMoviesService) {
				s.EXPECT().Get(1).Return(domain.Movie{
					Name:        "Name",
					ReleaseDate: time.Time{},
					Country:     "Russia",
					Genre:       "Haha",
					Rating:      5,
				}, nil).AnyTimes()
			},
			header: http.Header{
				"Content-Type": []string{
					"application/json",
				},
			},
			expStatusCode: http.StatusOK,
			expErr:        false,
		},
		{
			name: "get_actor_id_doesnt_exist",
			mockInit: func(s *mock_api.MockMoviesService) {
				s.EXPECT().Get(1).Return(domain.Movie{}, domain.ErrNotFound).AnyTimes()
			},
			header: http.Header{
				"Content-Type": []string{
					"application/json",
				},
			},
			expStatusCode: http.StatusNotFound,
			expErrMessage: "actor not found\n",
			expErr:        true,
		},
		{
			name: "internal_error",
			mockInit: func(s *mock_api.MockMoviesService) {
				s.EXPECT().Get(1).Return(domain.Movie{}, errors.New("unexpected error"))
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

			req := httptest.NewRequest(http.MethodGet, "/api/", io.Reader(nil))
			req.Header = tc.header
			recorder := httptest.NewRecorder()

			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", fmt.Sprintf("%v", 1))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

			h.Get(recorder, req)
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

			var actor domain.Movie
			_ = json.NewDecoder(recorder.Result().Body).Decode(&actor)

		})
	}
}

func TestDeleteMovie(t *testing.T) {
	testCases := []struct {
		name          string
		mockInit      func(s *mock_api.MockMoviesService)
		header        http.Header
		expStatusCode int
		expErrMessage string
		expErr        bool
	}{
		{
			name: "delete_actor_success",
			mockInit: func(s *mock_api.MockMoviesService) {
				s.EXPECT().Delete(1).Return(nil)
			},
			header: http.Header{
				"Content-Type": []string{
					"application/json",
				},
			},
			expStatusCode: http.StatusAccepted,
			expErr:        false,
		},
		{
			name: "delete_actor_not_found",
			mockInit: func(s *mock_api.MockMoviesService) {
				s.EXPECT().Delete(1).Return(domain.ErrNotFound)
			},
			header: http.Header{
				"Content-Type": []string{
					"application/json",
				},
			},
			expStatusCode: http.StatusNotFound,
			expErrMessage: "movie not found\n",
			expErr:        true,
		},
		{
			name: "delete_actor_unexpected_error",
			mockInit: func(s *mock_api.MockMoviesService) {
				s.EXPECT().Delete(1).Return(errors.New("unexpected error"))
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

			req := httptest.NewRequest(http.MethodDelete, "/api/", io.Reader(nil))
			req.Header = tc.header
			recorder := httptest.NewRecorder()

			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", fmt.Sprintf("%v", 1))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

			h.Delete(recorder, req)
			if recorder.Result().StatusCode != tc.expStatusCode {
				t.Errorf("expected status code: %d, got: %d", tc.expStatusCode, recorder.Result().StatusCode)
			}

			var movies []domain.Movie
			_ = json.NewDecoder(recorder.Result().Body).Decode(&movies)

		})
	}
}

func TestGetActorsByMovie(t *testing.T) {
	testCases := []struct {
		name          string
		mockInit      func(s *mock_api.MockMoviesService)
		header        http.Header
		expStatusCode int
		expErrMessage string
		expErr        bool
	}{
		{
			name: "get_actor_success",
			mockInit: func(s *mock_api.MockMoviesService) {
				s.EXPECT().GetActorsByMovie(1).Return([]domain.Actor{{
					ID:             1,
					Name:           "Name",
					BirthYear:      1999,
					CountryOfBirth: "cob",
					Gender:         "gender"},
				}, nil).AnyTimes()
			},
			header: http.Header{
				"Content-Type": []string{
					"application/json",
				},
			},
			expStatusCode: http.StatusOK,
			expErr:        false,
		},
		{
			name: "get_actor_id_doesnt_exist",
			mockInit: func(s *mock_api.MockMoviesService) {
				s.EXPECT().GetActorsByMovie(1).Return([]domain.Actor{}, domain.ErrNotFound).AnyTimes()
			},
			header: http.Header{
				"Content-Type": []string{
					"application/json",
				},
			},
			expStatusCode: http.StatusNotFound,
			expErrMessage: "actors not found\n",
			expErr:        true,
		},
		{
			name: "internal_error",
			mockInit: func(s *mock_api.MockMoviesService) {
				s.EXPECT().GetActorsByMovie(1).Return([]domain.Actor{}, errors.New("unexpected error"))
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

			req := httptest.NewRequest(http.MethodGet, "/api/", io.Reader(nil))
			req.Header = tc.header
			recorder := httptest.NewRecorder()

			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", fmt.Sprintf("%v", 1))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

			h.GetActors(recorder, req)
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

			var actors []domain.Actor
			_ = json.NewDecoder(recorder.Result().Body).Decode(&actors)

		})
	}
}

func TestCreateActorsForMovie(t *testing.T) {
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
			name:   "actors_already_exists",
			fields: fields{},
			mockInit: func(s *mock_api.MockMoviesService) {
				s.EXPECT().CreateActorsForMovie(1, []int{1, 2, 3}).Return(1, []int{1, 2, 3}, domain.ErrNotFound).AnyTimes()
			},
			header: http.Header{
				"Content-Type": []string{
					"application/json",
				},
			},
			expStatusCode: http.StatusNotFound,
			expErrMessage: "actors not found\n",
			expErr:        true,
		},
		{
			name:   "internal_error",
			fields: fields{},
			mockInit: func(s *mock_api.MockMoviesService) {
				s.EXPECT().Create(gomock.Any()).Return(domain.Movie{}, errors.New("unexpected error")).AnyTimes()
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

			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", fmt.Sprintf("%v", 1))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

			h.CreateActorsForMovie(recorder, req)
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

			var actors []domain.Actor
			_ = json.NewDecoder(recorder.Result().Body).Decode(&actors)
		})
	}
}

func toPtr(str string) *string {
	return &str
}
