package visit

import (
	"be/configs"
	"be/delivery/controllers/auth"
	"be/entities"
	"be/repository/visit"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/calendar/v3"
)

type mockSuccess struct{}

func (m *mockSuccess) CreateVal(doctor_uid, patient_uid string, req entities.Visit) (entities.Visit, error) {
	return entities.Visit{}, nil
}

func (m *mockSuccess) Update(visit_uid string, req entities.Visit) (entities.Visit, error) {
	return entities.Visit{}, nil
}

func (m *mockSuccess) Delete(visit_uid string) (entities.Visit, error) {
	return entities.Visit{}, nil
}

func (m *mockSuccess) GetVisitsVer1(kind, uid, status, date, grouped string) (visit.Visits, error) {
	return visit.Visits{}, nil
}

func (m *mockSuccess) GetVisitList(visit_uid string) (visit.VisitCalendar, error) {
	return visit.VisitCalendar{}, nil
}

type mockFail struct{}

func (m *mockFail) CreateVal(doctor_uid, patient_uid string, req entities.Visit) (entities.Visit, error) {
	return entities.Visit{}, errors.New("")
}

func (m *mockFail) Update(visit_uid string, req entities.Visit) (entities.Visit, error) {
	return entities.Visit{}, errors.New("")
}

func (m *mockFail) Delete(visit_uid string) (entities.Visit, error) {
	return entities.Visit{}, errors.New("")
}

func (m *mockFail) GetVisitsVer1(kind, uid, status, date, grouped string) (visit.Visits, error) {
	return visit.Visits{}, errors.New("")
}

func (m *mockFail) GetVisitList(visit_uid string) (visit.VisitCalendar, error) {
	return visit.VisitCalendar{}, errors.New("")
}

type MockAuthLib struct{}

func (m *MockAuthLib) Login(userName string, password string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"data": "abc",
		"type": "clinic",
	}, nil
}

type MockCal struct{}

func (m *MockCal) CreateEvent(visit_uid string) (*calendar.Event, error) {
	return &calendar.Event{}, nil
}

func (m *MockCal) InsertEvent(event *calendar.Event) error {
	return nil
}

func TestCreate(t *testing.T) {

	var jwt string
	t.Run("success login", func(t *testing.T) {
		e := echo.New()

		reqBody, _ := json.Marshal(map[string]string{
			"userName": "anonim@123",
			"password": "anonim123",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/login")

		authCont := auth.New(&MockAuthLib{})
		authCont.Login()(context)

		resp := auth.LoginRespFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)

		jwt = resp.Data["token"].(string)
	})

	t.Run("success", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"doctor_uid":  "doctor",
			"patient_uid": "patient",
			"date":        "05-05-2022",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.SetPath("/doctor")

		var controller = New(&mockSuccess{}, &MockCal{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Create())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 201, response.Code)
	})

	t.Run("binding", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"doctor_uid":  123,
			"patient_uid": "patient",
			"date":        "05-05-2022",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &MockCal{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Create())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
	})

	t.Run("validator", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"doctor_uid":  "",
			"patient_uid": "patient",
			"date":        "05-05-2022",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &MockCal{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Create())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
		// log.Info(response.Message)
	})

	t.Run("internal server", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"doctor_uid":  "doctor",
			"patient_uid": "patient",
			"date":        "05-05-2022",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)

		var controller = New(&mockFail{}, &MockCal{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.Create())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 500, response.Code)
		// log.Info(response.Message)
	})

}

func TestUpdate(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"complaint": "sick",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/visit")
		context.Param("visit_uid")
		context.SetParamNames("visit_uid")
		context.SetParamValues("visit 123")
		// log.Info(context.ParamNames())

		var controller = New(&mockSuccess{}, &MockCal{})
		controller.Update()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)
		// log.Info(response)
		assert.Equal(t, 202, response.Code)
	})

	t.Run("binding", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"complaint": 123,
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockSuccess{}, &MockCal{})
		controller.Update()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, 400, response.Code)
	})

	t.Run("internal server", func(t *testing.T) {
		var e = echo.New()

		var reqBody, _ = json.Marshal(map[string]interface{}{
			"complaint": "sick",
		})

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockFail{}, &MockCal{})
		controller.Update()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 500, response.Code)
		// log.Info(response.Message)
	})

}

func TestDelete(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/visit")
		context.Param("visit_uid")
		context.SetParamNames("visit_uid")
		context.SetParamValues("visit 123")
		// log.Info(context.ParamNames())

		var controller = New(&mockSuccess{}, &MockCal{})
		controller.Delete()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)
		// log.Info(response)
		assert.Equal(t, 202, response.Code)
	})

	t.Run("internal server", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)

		var controller = New(&mockFail{}, &MockCal{})
		controller.Delete()(context)

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 500, response.Code)
		// log.Info(response.Message)
	})

}

func TestGetVisits(t *testing.T) {

	var jwt string
	t.Run("success login", func(t *testing.T) {
		e := echo.New()

		reqBody, _ := json.Marshal(map[string]string{
			"userName": "anonim@123",
			"password": "anonim123",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/login")

		authCont := auth.New(&MockAuthLib{})
		authCont.Login()(context)

		resp := auth.LoginRespFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &resp)

		jwt = resp.Data["token"].(string)
	})

	t.Run("success", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)
		context.QueryParams().Add("status", "pending")

		var controller = New(&mockSuccess{}, &MockCal{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetVisits())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)
		// log.Info(response)
		assert.Equal(t, 200, response.Code)
	})

	t.Run("internal server", func(t *testing.T) {
		var e = echo.New()

		var req = httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(nil))
		var res = httptest.NewRecorder()
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwt))

		context := e.NewContext(req, res)

		var controller = New(&mockFail{}, &MockCal{})
		if err := middleware.JWT([]byte(configs.JWT_SECRET))(controller.GetVisits())(context); err != nil {
			log.Fatal(err)
			return
		}

		var response = ResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		// log.Info(response)
		assert.Equal(t, 500, response.Code)
		// log.Info(response.Message)
	})

}
