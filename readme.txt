refactor todo
1. move gorm leave handler 
-> in func NewTask change db to use store
-> rename func Create to New
-> create interface store, in interface have func New()
-> func (t *TodoHandler) NewTask(c *gin.Context), you see *TodoHandler, present it not have store
you must add it, *TodoHandler must have not use gorm
-> then you have error at error because New() return error, you can fix variable
-> fix NewTodoHandler, it not need to receive gorm

2. create Store for use when call t.store.New(&todo)
-> create file gorm.go in directory todo
-> create func New() return ....Create(todo)....
-> use receiver parameter
func (s *GormStore) New(todo *Todo) error 
-> create struct GormStore for use db
-> create func newGormStore for receive db form main.go

3. move gin leave handler
-> change c. to my method (todo.go)
-> create interface for my method
-> create gin.go for create func all my method
-> create NewMyContext for receive gin form main.go
-> r.POST("/todos", handler.NewTask) error because r want gin, so that we want to convert my Context to gin Context
can create func for convert in todo.go
-> edit call handler in main.go

4. cleaning up 
-> remove gorm.model in todo.go 
-> ใช้ค่าใน model ของ gorm มาใส่ไปเลย
type Todo struct {
	Title     string `json:"text" binding:"required"`
	ID        uint   `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

-> remove gin.H with map[string]interface{}
-> when error .Model, remove because I not use gorm.Model but use struct in gorm model can call .ID .Update

5. test
in todo.go have 2 dependency are store and Context
-> create mock struct in todo_test.go
-> add code in func bind for receive data in todo_test.go
-> add code in func JSON for send to frontend in todo_test.go but you need create struct TestContext
for get data
type TestContext struct{
	v map[string]interface{}
}
-> call handler 
handler := NewTodoHandler(&TestDB{})
-> call context
c := &TestContext{}
-> use them 
handler.NewTask(c)
want := "not allowed"
if want != c.v["error"] {
	t.Errorf("want %s but get %s\n", want, c.v["error"])
}

6. Moving Packages, make dependency have Package
-> store, move gorm
create directory store, move gorm, change package name, and
func (s *GormStore) New(todo *todo.Todo) error {...}
in main.go change todo. to store.
-> router, move gin
create directory router, move gin, change package name, same move gorm

7. Cleaning up Routers
-> move router in main.go to router/gin.go func NewMyRouter (-> 3)
-> create struct MyRouter in gin.go
-> create func NewMyRouter in gin.go
-> call router in main.go
r := router.NewMyRouter()
-> r.POST("/todos", router.NewGinHandler(handler.NewTask)) is hard too look fix it 
create func POST in gin.go, add code, call new code in main.go
r.POST("/todos", handler.NewTask)