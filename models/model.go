package models

import (
	"fmt"
	"sync"
	"time"
)

// есть повторяющиеся части, но не использую встраивание
// тк нельзя будет написать вот так &Job{Name: "something"}

// уникальность названия на уровне обеспечивается отдельной моделью под уровень с уникальным полем name

// есть также вариант хранить все узлы в одной таблице, чтобы они были однородными
// это даст гибкость построения любых деревьев с произвольной высотой
// но немного усложнит реализацию и создаст некоторые проблемы

// Job представляет Задание
type Job struct {
	ID        int       `json:"-"`
	Name      string    `gorm:"unique" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	TimeElapsed        float64 `json:"time_elapsed" gorm:"-"`
	TimeElapsedAverage float64 `json:"time_elapsed_average" gorm:"-"`

	Tasks []Task `json:"tasks"`
}

// Load это функция для загрузки задания с просчетом суммарного времени и среднего времени
func (job *Job) Load() error {
	// обработка ошибок
	if err := db.First(job, "name = ?", job.Name).Error; err != nil {
		return fmt.Errorf("when loading job from db: %v", err)
	}
	if err := db.Find(&job.Tasks, "job_id = ?", job.ID).Error; err != nil {
		return fmt.Errorf("when loading job from db: %v", err)
	}

	// будем параллелить обработку задач для задания:
	// данные для задачи будут извлекаться и просчитываться в отдельных горутинах
	// я сделал синхронизацию с помощью примитивов синхронизации
	var wg sync.WaitGroup
	var mutex sync.Mutex

	// через эту переменную потоки будут пробрасывать ошибку наверх
	errWhenFetchTasks := error(nil)
	for i := range job.Tasks {
		i := i // создаем локальную копию чтобы на ней замкнуться
		wg.Add(1)
		go func() {
			// тут нужно пролочить тк errWhenFetchTasks могут менять разные горутины
			mutex.Lock()
			if errWhenFetchTasks != nil {
				return
			}
			if err := db.Find(&job.Tasks[i].Works, "task_id = ?", job.Tasks[i].ID).Error; err != nil {
				errWhenFetchTasks = err
				return
			}
			mutex.Unlock()

			// тут можно не лочить тк горутины работают с разными элементами массива (не разделяют память)
			task := &job.Tasks[i]
			for _, work := range task.Works {
				task.TimeElapsed += work.TimeElapsed
			}
			if len(task.Works) > 0 {
				task.TimeElapsedAverage = task.TimeElapsed / float64(len(task.Works))
			}

			// тоже критическая секция
			mutex.Lock()
			job.TimeElapsed += task.TimeElapsed
			mutex.Unlock()

			wg.Done()
		}()
	}

	wg.Wait() // ждем пока все просчитается

	if errWhenFetchTasks != nil {
		return fmt.Errorf("when fetching tasks: %v", errWhenFetchTasks)
	}
	if len(job.Tasks) > 0 {
		job.TimeElapsedAverage = job.TimeElapsed / float64(len(job.Tasks))
	}

	return nil // успех
}

// Task представляет Задачy
type Task struct {
	ID        int       `json:"-"`
	Name      string    `gorm:"unique" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	TimeElapsed        float64 `json:"time_elapsed" gorm:"-"`
	TimeElapsedAverage float64 `json:"time_elapsed_average" gorm:"-"`

	JobID int    `json:"-"`
	Works []Work `json:"works"`
}

// Work представляет Трудозатрату
type Work struct {
	ID        int       `json:"-"`
	Name      string    `gorm:"unique" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	TimeElapsed float64 `json:"time_elapsed"`
	TaskID      int     `json:"-"`
}
