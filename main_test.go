package main

// тут тесты, я их особо не буду комментировать

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jobs/controllers"
	"jobs/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type NotesTestSuite struct {
	suite.Suite
	ts *httptest.Server
}

func TestNotesTestSuite(t *testing.T) {
	suite.Run(t, &NotesTestSuite{})
}

func (n *NotesTestSuite) SetupSuite() {
	n.ts = httptest.NewServer(controllers.GetRouter())
}

func (n *NotesTestSuite) SetupTest() {
	// создаем базу в памяти чтобы тесты быстро работали
	conn, err := gorm.Open(sqlite.Open(":memory:"))
	n.Require().Nil(err, "test db should work")

	models.Init(conn)
	models.Migrate()
}

func (n *NotesTestSuite) TearDownTest() {
	models.Truncate()
}

func (n *NotesTestSuite) TearDownSuite() {
	n.ts.Close()
}

func (n *NotesTestSuite) TestCreateJob() {
	jobName := "job1"

	r := Must(http.Post(n.ts.URL+"/jobs", "application/json", AsJSONBody(map[string]string{
		"name": jobName,
	})))

	n.Require().Equal(200, r.StatusCode)
	n.Require().Equal("application/json", r.Header.Get("Content-Type"))

	rd := &Response{}
	n.Require().Nil(json.NewDecoder(r.Body).Decode(rd))

	n.Require().True(rd.Success, rd.Message)

	n.Require().Nil(models.GetDB().First(&models.Job{}, "name = ?", jobName).Error)
}

func (n *NotesTestSuite) TestCreateTask() {
	taskName := "task1"
	jobName := "job1"

	n.Require().Nil(models.GetDB().Create(&models.Job{Name: jobName}).Error)

	r := Must(http.Post(n.ts.URL+"/tasks?parent="+jobName, "application/json", AsJSONBody(map[string]string{
		"name": taskName,
	})))

	n.Require().Equal(200, r.StatusCode)
	n.Require().Equal("application/json", r.Header.Get("Content-Type"))

	rd := &Response{}
	n.Require().Nil(json.NewDecoder(r.Body).Decode(rd))

	n.Require().True(rd.Success, rd.Message)

	n.Require().Nil(models.GetDB().First(&models.Task{}, "name = ?", taskName).Error)
}

func (n *NotesTestSuite) TestCreateWork() {
	job := &models.Job{
		Name: "job1",
		Tasks: []models.Task{
			{
				Name: "task1",
			},
		},
	}
	n.Require().Nil(models.GetDB().Create(job).Error)

	r := Must(http.Post(n.ts.URL+"/works?parent=task1", "application/json", AsJSONBody(Object{
		"name":         "work1",
		"time_elapsed": 5,
	})))

	n.Require().Equal(200, r.StatusCode)
	n.Require().Equal("application/json", r.Header.Get("Content-Type"))

	rd := &Response{}
	n.Require().Nil(json.NewDecoder(r.Body).Decode(rd))

	n.Require().True(rd.Success, rd.Message)

	n.Require().Nil(models.GetDB().First(&models.Work{}, "name = ?", "work1").Error)
}

func (n *NotesTestSuite) TestDetailJob() {
	createTestJob()

	r := Must(http.Get(n.ts.URL + "/jobs/job1"))

	n.Require().Equal(200, r.StatusCode)
	n.Require().Equal("application/json", r.Header.Get("Content-Type"))

	rd := &ResponseJob{}
	n.Require().Nil(json.NewDecoder(r.Body).Decode(rd))

	n.Require().True(rd.Success, rd.Message)

	n.Require().NotNil(rd.Data)
	n.Require().Equal(2, len(rd.Data.Tasks))
	n.Require().Equal(3, len(rd.Data.Tasks[0].Works)+len(rd.Data.Tasks[1].Works))
	n.Require().Equal(12.0, rd.Data.TimeElapsed)
	n.Require().Equal(6.0, rd.Data.TimeElapsedAverage)
}

func (n *NotesTestSuite) TestJobList() {
	n.Require().Nil(models.GetDB().Create(&models.Job{Name: "job1"}).Error)
	n.Require().Nil(models.GetDB().Create(&models.Job{Name: "job2"}).Error)

	r := Must(http.Get(n.ts.URL + "/jobs"))

	n.Require().Equal(200, r.StatusCode)
	n.Require().Equal("application/json", r.Header.Get("Content-Type"))

	rd := &ResponseJobsList{}
	n.Require().Nil(json.NewDecoder(r.Body).Decode(rd))

	n.Require().True(rd.Success, rd.Message)

	n.Require().Equal(2, len(rd.Data))
}

func createTestJob() *models.Job {
	job := &models.Job{
		Name: "job1",
		Tasks: []models.Task{
			{
				Name: "task1",
				Works: []models.Work{
					{
						Name:        "work1",
						TimeElapsed: 5,
					},
					{
						Name:        "work2",
						TimeElapsed: 3,
					},
				},
			},
			{
				Name: "task2",
				Works: []models.Work{
					{
						Name:        "work3",
						TimeElapsed: 4,
					},
				},
			},
		},
	}

	models.GetDB().Create(&job.Tasks[0])
	models.GetDB().Create(&job.Tasks[1])
	models.GetDB().Create(job)

	return job
}

func (n *NotesTestSuite) TestDeleteJob() {
	job := createTestJob()

	client := &http.Client{}

	req, _ := http.NewRequest("DELETE", fmt.Sprintf(n.ts.URL+"/jobs/"+job.Name), nil)

	r := Must(client.Do(req))

	n.Require().Equal(200, r.StatusCode)
	n.Require().Equal("application/json", r.Header.Get("Content-Type"))

	n.Require().NotNil(models.GetDB().First(&models.Job{}).Error)
}

func (n *NotesTestSuite) TestDeleteTask() {
	createTestJob()

	client := &http.Client{}

	req, _ := http.NewRequest("DELETE", fmt.Sprintf(n.ts.URL+"/tasks/task1"), nil)

	r := Must(client.Do(req))

	n.Require().Equal(200, r.StatusCode)
	n.Require().Equal("application/json", r.Header.Get("Content-Type"))

	n.Require().NotNil(models.GetDB().First(&models.Task{}, "name = task1").Error)
}

func (n *NotesTestSuite) TestDeleteWork() {
	createTestJob()

	client := &http.Client{}

	req, _ := http.NewRequest("DELETE", fmt.Sprintf(n.ts.URL+"/works/work1"), nil)

	r := Must(client.Do(req))

	n.Require().Equal(200, r.StatusCode)
	n.Require().Equal("application/json", r.Header.Get("Content-Type"))

	n.Require().NotNil(models.GetDB().First(&models.Work{}, "name = work1").Error)
}

func (n *NotesTestSuite) TestMoveTask() {
	createTestJob()

	anotherJob := &models.Job{Name: "job2"}

	models.GetDB().Create(anotherJob)
	client := &http.Client{}

	req, _ := http.NewRequest("PUT", fmt.Sprintf(n.ts.URL+"/tasks/task1?parent=job2"), nil)

	r := Must(client.Do(req))

	n.Require().Equal(200, r.StatusCode)
	n.Require().Equal("application/json", r.Header.Get("Content-Type"))

	task := &models.Task{}
	n.Require().Nil(models.GetDB().First(task, "name = ?", "task1").Error)
	n.Require().True(task.JobID == anotherJob.ID) // check we move it to another job
}

func (n *NotesTestSuite) TestMoveWork() {
	createTestJob()

	client := &http.Client{}

	req, _ := http.NewRequest("PUT", fmt.Sprintf(n.ts.URL+"/works/work1?parent=task2"), nil)

	r := Must(client.Do(req))

	n.Require().Equal(200, r.StatusCode)
	n.Require().Equal("application/json", r.Header.Get("Content-Type"))

	task := &models.Task{}
	models.GetDB().First(task, "name = ?", "task2")

	work := &models.Work{}
	n.Require().Nil(models.GetDB().First(work, "name = ?", "work1").Error)
	n.Require().True(work.TaskID == task.ID) // check we move it to another task
}

func Must(resp *http.Response, err error) *http.Response {
	if err != nil {
		panic("when working with test server: " + err.Error())
	}

	return resp
}

func AsJSONBody(obj interface{}) io.Reader {
	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(obj)
	return b
}

func (n *NotesTestSuite) TestNotFound() {
	resp := Must(http.Get(n.ts.URL + "/not-exists"))
	n.Require().Equal(404, resp.StatusCode)

	rd := &Response{}
	n.Require().Nil(json.NewDecoder(resp.Body).Decode(rd), "server should serve with valid json anyway")
	n.Require().False(rd.Success)
	n.Require().NotEmpty(rd.Message)
}

type Object map[string]interface{}

type ResponseJobsList struct {
	Response
	Data []models.Job `json:"data"`
}

type ResponseJob struct {
	Response
	Data models.Job `json:"data"`
}
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
