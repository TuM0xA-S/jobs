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

// JobsList returns list of jobs without details(not calculates total elapsed time)
func JobsList(w http.ResponseWriter, r *http.Request) {
	jobs := []models.Job{}
	// достаем список
	if err := models.GetDB().Find(&jobs).Error; err != nil {
		panic(fmt.Errorf("when fetch jobs: %v", err))
	}

	resp := util.ResponseBaseOK()
	resp["data"] = jobs

	util.RespondWithJSON(w, 200, resp)
}

// JobDetail returns job detailt view with average time etc
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

// JobCreate for creating job
func JobCreate(w http.ResponseWriter, r *http.Request) {
	// комменты в другом аналогичном контроллере
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

// TaskCreate creates task
func TaskCreate(w http.ResponseWriter, r *http.Request) {
	// комменты в другом аналогичном контроллере
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
	}

	task.JobID = job.ID
	if err := models.GetDB().Create(task).Error; err != nil {
		panic(fmt.Errorf("when creating task: %v", err))
	}

	util.RespondWithJSON(w, 200, util.ResponseBaseOK())
}

// WorkCreate creates work
func WorkCreate(w http.ResponseWriter, r *http.Request) {
	work := &models.Work{}
	defer r.Body.Close()
	// декодировали тело запроса
	if err := json.NewDecoder(r.Body).Decode(work); err != nil {
		util.RespondWithError(w, 400, "invalid request")
		return
	}
	// простая валидация
	if work.Name == "" {
		util.RespondWithError(w, 422, "empty name")
		return
	}

	if work.TimeElapsed <= 0 {
		util.RespondWithError(w, 422, "time_elapsed should be > 0")
		return
	}
	// взяли родителя и смотрим есть ли он
	task := &models.Task{Name: mux.Vars(r)["parent"]}
	if err := models.GetDB().Where(task).Take(task).Error; err == gorm.ErrRecordNotFound {
		util.RespondWithError(w, 404, "no such task")
		return
	} else if err != nil {
		panic(fmt.Errorf("when checking task: %v", err))
	}
	// подвесили
	work.TaskID = task.ID
	if err := models.GetDB().Create(work).Error; err != nil {
		panic(fmt.Errorf("when creating work: %v", err))
	}

	util.RespondWithJSON(w, 200, util.ResponseBaseOK())
}

//WorkDelete deletes work
func WorkDelete(w http.ResponseWriter, r *http.Request) {
	// сначала достаем
	// а потом удаляем
	// так делается чтобы можно было обработать ошибку когда удалять нечего
	work := &models.Work{}
	if err := models.GetDB().Take(work, "name = ?", mux.Vars(r)["work"]).Error; err == gorm.ErrRecordNotFound {
		util.RespondWithError(w, 404, "no such work")
		return
	} else if err != nil {
		panic(fmt.Errorf("when deleting work: %v", err))
	}

	if err := models.GetDB().Where(work).Delete(work).Error; err != nil {
		panic(fmt.Errorf("when deleting work: %v", err))
	}
	util.RespondWithJSON(w, 200, util.ResponseBaseOK())
}

//TaskDelete deletes task
func TaskDelete(w http.ResponseWriter, r *http.Request) {
	// достаем и пробуем удалить таск
	task := &models.Task{}
	if err := models.GetDB().Take(task, "name = ?", mux.Vars(r)["task"]).Error; err == gorm.ErrRecordNotFound {
		util.RespondWithError(w, 404, "no such task")
		return
	} else if err != nil {
		panic(fmt.Errorf("when deleting task: %v", err))
	}

	if err := models.GetDB().Where(task).Delete(task).Error; err != nil {
		panic(fmt.Errorf("when deleting task: %v", err))
	}

	util.RespondWithJSON(w, 200, util.ResponseBaseOK())
}

//JobDelete deletes job
func JobDelete(w http.ResponseWriter, r *http.Request) {
	// аналогично
	job := &models.Job{}
	if err := models.GetDB().Take(job, "name = ?", mux.Vars(r)["job"]).Error; err == gorm.ErrRecordNotFound {
		util.RespondWithError(w, 404, "no such job")
		return
	} else if err != nil {
		panic(fmt.Errorf("when deleting job: %v", err))
	}

	if err := models.GetDB().Where(job).Delete(job).Error; err != nil {
		panic(fmt.Errorf("when deleting job: %v", err))
	}

	util.RespondWithJSON(w, 200, util.ResponseBaseOK())
}

// WorkMove is a controller for move work from one task to another
func WorkMove(w http.ResponseWriter, r *http.Request) {
	// этот контроллер аналогичен другому
	workName := mux.Vars(r)["work"]
	newParentName := mux.Vars(r)["parent"]
	work := &models.Work{}
	if err := models.GetDB().Take(work, "name = ?", workName).Error; err == gorm.ErrRecordNotFound {
		util.RespondWithError(w, 404, "no such work")
		return
	} else if err != nil {
		panic(fmt.Errorf("when moving work: %v", err))
	}

	task := &models.Task{}
	if err := models.GetDB().Take(task, "name = ?", newParentName).Error; err == gorm.ErrRecordNotFound {
		util.RespondWithError(w, 404, "no such task")
		return
	} else if err != nil {
		panic(fmt.Errorf("when moving work: %v", err))
	}

	work.TaskID = task.ID
	if err := models.GetDB().Save(work).Error; err != nil {
		panic(fmt.Errorf("when saving work: %v", err))
	}

	util.RespondWithJSON(w, 200, util.ResponseBaseOK())
}

// TaskMove is a controller for move task from one job to another
func TaskMove(w http.ResponseWriter, r *http.Request) {
	// достали входные данные
	taskName := mux.Vars(r)["task"]
	newParentName := mux.Vars(r)["parent"]
	//достаем таск который будем двигать
	task := &models.Task{}
	if err := models.GetDB().Take(task, "name = ?", taskName).Error; err == gorm.ErrRecordNotFound {
		util.RespondWithError(w, 404, "no such task")
		return
	} else if err != nil {
		panic(fmt.Errorf("when moving task: %v", err))
	}
	//достаем джоб который будем двигать
	job := &models.Job{}
	if err := models.GetDB().Take(job, "name = ?", newParentName).Error; err == gorm.ErrRecordNotFound {
		util.RespondWithError(w, 404, "no such job")
		return
	} else if err != nil {
		panic(fmt.Errorf("when moving task: %v", err))
	}
	// цепляем таск на другой джоб
	task.JobID = job.ID
	if err := models.GetDB().Save(task).Error; err != nil {
		panic(fmt.Errorf("when saving task: %v", err))
	}
	// сохраняем
	models.GetDB().Save(task)

	util.RespondWithJSON(w, 200, util.ResponseBaseOK())
}
