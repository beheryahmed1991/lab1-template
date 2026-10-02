package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/bmstu-rsoi/lab1-template.git/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type mockPersonService struct {
	create  func(context.Context, model.PersonRequest) (int64, error)
	getByID func(context.Context, int64) (*model.Person, error)
	getAll  func(context.Context) ([]model.Person, error)
	update  func(context.Context, int64, model.PersonRequest) (*model.Person, error)
	delete  func(context.Context, int64) error
}

var _ personService = (*mockPersonService)(nil)

func (m *mockPersonService) Create(ctx context.Context, req model.PersonRequest) (int64, error) {
	return m.create(ctx, req)
}
func (m *mockPersonService) GetByID(ctx context.Context, id int64) (*model.Person, error) {
	return m.getByID(ctx, id)
}
func (m *mockPersonService) GetAll(ctx context.Context) ([]model.Person, error) {
	return m.getAll(ctx)
}
func (m *mockPersonService) Update(ctx context.Context, id int64, req model.PersonRequest) (*model.Person, error) {
	return m.update(ctx, id, req)
}
func (m *mockPersonService) Delete(ctx context.Context, id int64) error {
	return m.delete(ctx, id)
}

func request(t *testing.T, svc *mockPersonService, method, path, body string, status int) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	h := NewPersonHandler(svc)
	router := gin.New()
	router.POST("/api/v1/persons", h.Create)
	router.GET("/api/v1/persons", h.GetAll)
	router.GET("/api/v1/persons/:id", h.GetByID)
	router.PATCH("/api/v1/persons/:id", h.Update)
	router.DELETE("/api/v1/persons/:id", h.Delete)
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != status {
		t.Fatalf("status = %d, want %d; body: %s", w.Code, status, w.Body.String())
	}
	return w
}

func assertJSON(t *testing.T, w *httptest.ResponseRecorder, want any) {
	t.Helper()
	if !strings.Contains(w.Header().Get("Content-Type"), "application/json") {
		t.Fatal("response must have JSON content type")
	}
	var got, expected any
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, &expected); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("JSON = %v, want %v", got, expected)
	}
}

func TestCreatePerson(t *testing.T) {
	t.Run("created", func(t *testing.T) {
		called := false
		svc := &mockPersonService{create: func(ctx context.Context, req model.PersonRequest) (int64, error) {
			called = true
			if ctx == nil || req.Name != "Ahmed" || req.Age == nil || *req.Age != 35 || req.Address == nil || *req.Address != "Moscow" || req.Work == nil || *req.Work != "Developer" {
				t.Fatalf("unexpected create arguments: %+v", req)
			}
			return 1, nil
		}}
		w := request(t, svc, http.MethodPost, "/api/v1/persons", `{"name":"Ahmed","age":35,"address":"Moscow","work":"Developer"}`, 201)
		if !called || w.Header().Get("Location") != "/api/v1/persons/1" || w.Body.Len() != 0 {
			t.Fatal("expected service call, Location header and empty response body")
		}
	})
	for _, body := range []string{`{`, `{}`, `{"name":""}`, `{"name":"Ahmed","age":"invalid"}`} {
		t.Run("invalid "+body, func(t *testing.T) {
			request(t, &mockPersonService{}, http.MethodPost, "/api/v1/persons", body, 400)
		})
	}
	t.Run("service failure", func(t *testing.T) {
		svc := &mockPersonService{create: func(context.Context, model.PersonRequest) (int64, error) { return 0, errors.New("failure") }}
		request(t, svc, http.MethodPost, "/api/v1/persons", `{"name":"Ahmed"}`, 500)
	})
}

func TestGetPersonByID(t *testing.T) {
	for _, tc := range []struct {
		name   string
		err    error
		status int
	}{
		{"found", nil, 200}, {"missing", fmt.Errorf("lookup: %w", pgx.ErrNoRows), 404}, {"service failure", errors.New("failure"), 500},
	} {
		t.Run(tc.name, func(t *testing.T) {
			person := &model.Person{ID: 7, Name: "Ahmed"}
			called := false
			svc := &mockPersonService{getByID: func(_ context.Context, id int64) (*model.Person, error) {
				called = true
				if id != 7 {
					t.Fatalf("id = %d, want 7", id)
				}
				return person, tc.err
			}}
			w := request(t, svc, http.MethodGet, "/api/v1/persons/7", "", tc.status)
			if !called {
				t.Fatal("service was not called")
			}
			if tc.err == nil {
				assertJSON(t, w, person)
			}
		})
	}
	t.Run("invalid id", func(t *testing.T) { request(t, &mockPersonService{}, http.MethodGet, "/api/v1/persons/abc", "", 400) })
}

func TestGetAllPersons(t *testing.T) {
	for _, persons := range [][]model.Person{{{ID: 1, Name: "Ahmed"}, {ID: 2, Name: "Ali"}}, {}} {
		t.Run(fmt.Sprintf("%d persons", len(persons)), func(t *testing.T) {
			svc := &mockPersonService{getAll: func(context.Context) ([]model.Person, error) { return persons, nil }}
			assertJSON(t, request(t, svc, http.MethodGet, "/api/v1/persons", "", 200), persons)
		})
	}
	t.Run("service failure", func(t *testing.T) {
		svc := &mockPersonService{getAll: func(context.Context) ([]model.Person, error) { return nil, errors.New("failure") }}
		request(t, svc, http.MethodGet, "/api/v1/persons", "", 500)
	})
}

func TestUpdatePerson(t *testing.T) {
	for _, tc := range []struct {
		name   string
		err    error
		status int
	}{
		{"updated", nil, 200}, {"missing", fmt.Errorf("update: %w", pgx.ErrNoRows), 404}, {"service failure", errors.New("failure"), 500},
	} {
		t.Run(tc.name, func(t *testing.T) {
			age := 35
			person := &model.Person{ID: 7, Name: "Updated", Age: &age}
			called := false
			svc := &mockPersonService{update: func(_ context.Context, id int64, req model.PersonRequest) (*model.Person, error) {
				called = true
				if id != 7 || req.Name != "Updated" || req.Age != nil || req.Work != nil {
					t.Fatalf("unexpected PATCH arguments: id=%d request=%+v", id, req)
				}
				return person, tc.err
			}}
			w := request(t, svc, http.MethodPatch, "/api/v1/persons/7", `{"name":"Updated"}`, tc.status)
			if !called {
				t.Fatal("service was not called")
			}
			if tc.err == nil {
				assertJSON(t, w, person)
			}
		})
	}
	t.Run("invalid id", func(t *testing.T) {
		request(t, &mockPersonService{}, http.MethodPatch, "/api/v1/persons/abc", `{"name":"Updated"}`, 400)
	})
	for _, body := range []string{`{`, `{}`, `{"name":""}`, `{"name":"Updated","age":"invalid"}`} {
		t.Run("invalid "+body, func(t *testing.T) { request(t, &mockPersonService{}, http.MethodPatch, "/api/v1/persons/7", body, 400) })
	}
}

func TestDeletePerson(t *testing.T) {
	for _, tc := range []struct {
		name   string
		err    error
		status int
	}{
		{"deleted", nil, 204}, {"missing", fmt.Errorf("delete: %w", pgx.ErrNoRows), 404}, {"service failure", errors.New("failure"), 500},
	} {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			svc := &mockPersonService{delete: func(_ context.Context, id int64) error {
				called = true
				if id != 7 {
					t.Fatalf("id = %d, want 7", id)
				}
				return tc.err
			}}
			w := request(t, svc, http.MethodDelete, "/api/v1/persons/7", "", tc.status)
			if !called {
				t.Fatal("service was not called")
			}
			if tc.err == nil && w.Body.Len() != 0 {
				t.Fatal("204 response must be empty")
			}
		})
	}
	t.Run("invalid id", func(t *testing.T) {
		request(t, &mockPersonService{}, http.MethodDelete, "/api/v1/persons/abc", "", 400)
	})
}
