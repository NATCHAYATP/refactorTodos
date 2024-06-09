package todo

import (
	"testing"
)

func TestNewTask(t *testing.T) {
	// handler := NewTodoHandler(&gorm.DB{})
	handler := NewTodoHandler(&TestDB{})

	// there are context, now you not create because use mock context
	// w := httptest.NewRecorder()
	// payload := bytes.NewBufferString(`{"text":"sleep"}`)
	// req, _ := http.NewRequest("POST", "http://0.0.0.0:8080/todos", payload)
	// req.Header.Add("TransactionID", "testIDxxx")

	// c, _ := gin.CreateTestContext(w)
	// c.Request = req
	c := &TestContext{}

	handler.NewTask(c)

	want := "not allowed"

	if want != c.v["error"] {
		t.Errorf("want %s but get %s\n", want, c.v["error"])
	}
}

// mock store
type TestDB struct {}

func (TestDB) New(*Todo) error {
	return nil
}

// mock context
type TestContext struct{
	v map[string]interface{}
}

func (TestContext) Bind(v interface{}) error {
	*v.(*Todo) = Todo{
		Title: "sleep",
	}
	return nil
}

func (c *TestContext) JSON(code int, v interface{}) {
	c.v = v.(map[string]interface{})
}

func (TestContext) TransactionID() string {
	return "TestTransactionID"
}
func (TestContext) Audience() string {
	return "Unit Test"
}
