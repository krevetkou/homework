package api

import (
	"arch-demo/internal/domain"
	mock_api "arch-demo/internal/tests/api_mocks"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreate(t *testing.T) {
	type fields struct {
		name           string
		birthYear      int
		countryOfBirth string
		gender         string
	}

	testCases := []struct {
		name          string
		fields        fields
		mockInit      func(s *mock_api.MockActorsService)
		header        http.Header
		expStatusCode int
		expErrMessage string
		expErr        bool
	}{
		{
			name:     "fail_content_type",
			fields:   fields{},
			mockInit: func(s *mock_api.MockActorsService) {},
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
			name:   "user_already_exists",
			fields: fields{},
			mockInit: func(s *mock_api.MockActorsService) {
				s.EXPECT().Create(gomock.Any()).Return(domain.Actor{}, domain.ErrExists)
			},
			header: http.Header{
				"Content-Type": []string{
					"application/json",
				},
			},
			expStatusCode: http.StatusConflict,
			expErrMessage: "actor already exists\n",
			expErr:        true,
		},
		{
			name:   "internal_error",
			fields: fields{},
			mockInit: func(s *mock_api.MockActorsService) {
				s.EXPECT().Create(gomock.Any()).Return(domain.Actor{}, errors.New("unexpected error"))
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
				name:           "name",
				birthYear:      1900,
				countryOfBirth: "cob",
				gender:         "gender",
			},
			mockInit: func(s *mock_api.MockActorsService) {
				s.EXPECT().Create(gomock.Any()).Return(domain.Actor{
					Name:           "name",
					BirthYear:      1900,
					CountryOfBirth: "cob",
					Gender:         "gender",
				}, nil)
			},
			header: http.Header{
				"Content-Type": []string{
					"application/json",
				},
			},
			expStatusCode: http.StatusCreated,
			expErrMessage: "unexpected error\n",
			expErr:        false,
		},
		{
			name: "not_all_fields",
			fields: fields{
				birthYear:      1900,
				countryOfBirth: "cob",
				gender:         "gender",
			},
			mockInit: func(s *mock_api.MockActorsService) {
				s.EXPECT().Create(gomock.Any()).Return(domain.Actor{
					Name:           "name",
					BirthYear:      1900,
					CountryOfBirth: "cob",
					Gender:         "gender",
				}, domain.ErrFieldsRequired)
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

			s := mock_api.NewMockActorsService(ctrl)
			if tc.mockInit != nil {
				tc.mockInit(s)
			}

			h := NewActorsHandler(s)
			payload := domain.Actor{
				Name:           tc.fields.name,
				BirthYear:      tc.fields.birthYear,
				CountryOfBirth: tc.fields.countryOfBirth,
				Gender:         tc.fields.gender,
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

			var actor domain.Actor
			_ = json.NewDecoder(recorder.Result().Body).Decode(&actor)

			if tc.fields.name != actor.Name &&
				tc.fields.birthYear != actor.BirthYear &&
				tc.fields.countryOfBirth != actor.CountryOfBirth &&
				tc.fields.gender != actor.Gender {
				t.Errorf("expected actor: %v, got: %v", tc.fields, actor)
			}
		})
	}
}

func TestList(t *testing.T) {

	testCases := []struct {
		name          string
		mockInit      func(s *mock_api.MockActorsService)
		header        http.Header
		expStatusCode int
		expErrMessage string
		expErr        bool
	}{
		{
			name: "get_actors_success",
			mockInit: func(s *mock_api.MockActorsService) {
				s.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return([]domain.Actor{
					{
						Name:           "Cheburek",
						BirthYear:      900,
						CountryOfBirth: "Klkd",
						Gender:         "elephant",
					},
					{
						Name:           "Kek",
						BirthYear:      2012,
						CountryOfBirth: "Lock",
						Gender:         "male",
					},
					{
						Name:           "Lol",
						BirthYear:      1909,
						CountryOfBirth: "Sos",
						Gender:         "female",
					},
				}, nil)
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

			s := mock_api.NewMockActorsService(ctrl)
			if tc.mockInit != nil {
				tc.mockInit(s)
			}
			h := NewActorsHandler(s)

			actors := []domain.Actor{
				{
					Name:           "Cheburek",
					BirthYear:      900,
					CountryOfBirth: "Klkd",
					Gender:         "elephant",
				},
				{
					Name:           "Kek",
					BirthYear:      2012,
					CountryOfBirth: "Lock",
					Gender:         "male",
				},
				{
					Name:           "Lol",
					BirthYear:      1909,
					CountryOfBirth: "Sos",
					Gender:         "female",
				},
			}

			req := httptest.NewRequest(http.MethodGet, "/api/", io.Reader(nil))
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

			var actorsTest string
			_ = json.NewDecoder(recorder.Result().Body).Decode(&actorsTest)

			fmt.Println(actorsTest)
			fmt.Println(actors)

		})
	}
}
