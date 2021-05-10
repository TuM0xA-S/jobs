package controllers

import (
	"encoding/json"
	"fmt"
	"jobs/models"
	"jobs/util"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// NotFound Handler ..
func NotFound(w http.ResponseWriter, r *http.Request) {
	util.RespondWithError(w, 404, r.URL.String()+" not found")
}

// MethodNotAllowed ...
func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	util.RespondWithError(w, 405, r.Method+" not allowed")
}

func JobsList(w http.ResponseWriter, r *http.Request) {
	jobs := []models.Job{}

	if err := models.GetDB().Find(&jobs).Error; err != nil {
		panic(fmt.Errorf("when fetch jobs: %v", err))
	}

	resp := util.ResponseBaseOK()
	resp["data"] = jobs

	util.RespondWithJSON(w, 200, resp)
}

func JobDetail(w http.ResponseWriter, r *http.Request) {
	job := models.Job{Name: mux.Vars(r)["job"]}

	if err := job.Load(); err == gorm.ErrRecordNotFound {
		util.RespondWithError(w, 404, "no such job")
	} else if err != nil {
		panic(fmt.Errorf("when get job detail: %v", err))
	} else {
		resp := util.ResponseBaseOK()
		resp["data"] = job
		util.RespondWithJSON(w, 200, resp)
	}

}

func JobCreate(w http.ResponseWriter, r *http.Request) {
	job := &models.Job{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(job); err != nil {
		util.RespondWithError(w, 400, "invalid request")
		return
	}

	if job.Name == "" {
		util.RespondWithError(w, 422, "empty name")
		return
	}

	if err := models.GetDB().Create(job).Error; err != nil {
		panic(fmt.Errorf("when creating job: %v", err))
	}

	util.RespondWithJSON(w, 200, util.ResponseBaseOK())
}

func TaskCreate(w http.ResponseWriter, r *http.Request) {
	task := &models.Task{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(task); err != nil {
		util.RespondWithError(w, 400, "invalid request")
		return
	}

	if task.Name == "" {
		util.RespondWithError(w, 422, "empty name")
	}

	job := &models.Job{Name: mux.Vars(r)["parent"]}
	if err := models.GetDB().Where(job).Take(job).Error; err == gorm.ErrRecordNotFound {
		util.RespondWithError(w, 404, "no such job")
		return
	} else if err != nil {
		panic(fmt.Errorf("when checking job: %v", err))
		return
	}

	task.JobID = job.ID
	if err := models.GetDB().Create(task).Error; err != nil {
		panic(fmt.Errorf("when creating job: %v", err))
	}

	util.RespondWithJSON(w, 200, util.ResponseBaseOK())
}

func WorkCreate(w http.ResponseWriter, r *http.Request) {
	work := &models.Work{}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(work); err != nil {
		util.RespondWithError(w, 400, "invalid request")
		return
	}

	if work.Name == "" {
		util.RespondWithError(w, 422, "empty name")
		return
	}

	if work.TimeElapsed <= 0 {
		util.RespondWithError(w, 422, "time_elapsed should be > 0")
		return
	}

	task := &models.Task{Name: mux.Vars(r)["parent"]}
	if err := models.GetDB().Where(task).Take(task).Error; err == gorm.ErrRecordNotFound {
		util.RespondWithError(w, 404, "no such task")
		return
	} else if err != nil {
		panic(fmt.Errorf("when checking task: %v", err))
	}

	work.TaskID = task.ID
	if err := models.GetDB().Create(work).Error; err != nil {
		panic(fmt.Errorf("when creating job: %v", err))
	}

	util.RespondWithJSON(w, 200, util.ResponseBaseOK())
}

func WorkDelete(w http.ResponseWriter, r *http.Request) {
	if err := models.GetDB().Delete(&models.Work{}, "name = ?", mux.Vars(r)["work"]).Error; err == gorm.ErrRecordNotFound {
		util.RespondWithError(w, 404, "no such work")
		return
	} else if err != nil {
		panic(fmt.Errorf("when deleting work: %v", err))
	}

	util.RespondWithJSON(w, 200, util.ResponseBaseOK())
}

func TaskDelete(w http.ResponseWriter, r *http.Request) {
	if err := models.GetDB().Delete(&models.Task{}, "name = ?", mux.Vars(r)["task"]).Error; err == gorm.ErrRecordNotFound {
		util.RespondWithError(w, 404, "no such task")
		return
	} else if err != nil {
		panic(fmt.Errorf("when deleting task: %v", err))
	}

	util.RespondWithJSON(w, 200, util.ResponseBaseOK())
}

func JobDelete(w http.ResponseWriter, r *http.Request) {
	if err := models.GetDB().Delete(&models.Job{}, "name = ?", mux.Vars(r)["job"]).Error; err == gorm.ErrRecordNotFound {
		util.RespondWithError(w, 404, "no such job")
		return
	} else if err != nil {
		panic(fmt.Errorf("when deleting job: %v", err))
	}

	util.RespondWithJSON(w, 200, util.ResponseBaseOK())
}
