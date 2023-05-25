package main

import (
	"fmt"
	"log"
	"os"

	postgr ".skillfactory/module30.8.1/pkg/storage"
)

func main() {
	dbuser := os.Getenv("dbuser")
	pwd := os.Getenv("dbpass")
	dbhost := os.Getenv("dbpath")

	db, err := postgr.New("postgres://" + dbuser + ":" + pwd + "@" + dbhost + "/module_30_8_1")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	tasks, err := db.Tasks(0, 0)
	if err != nil {
		log.Fatalf("Unable to get tasks: %v\n", err)
	}
	fmt.Printf("tasks: %v\n", tasks)

	// add some tasks
	makeTestData(db, "test task", "testing", 1)
	makeTestData(db, "task for user", "", 2)
	tasks, err = db.Tasks(0, 0)
	if err != nil {
		log.Fatalf("Unable to get tasks: %v\n", err)
	}
	fmt.Printf("tasks: %v\n", tasks)

	fmt.Println("assign task to user")
	task := tasks[len(tasks)-2]
	task.AssignedID = 3
	db.UpdateTask(task.ID, task)
	tasks, err = db.Tasks(0, 0)
	if err != nil {
		log.Fatalf("Unable to get tasks: %v\n", err)
	}
	fmt.Printf("tasks: %v\n", tasks)

	fmt.Println("deleting test tasks")
	db.DeletTask(tasks[len(tasks)-2].ID)
	db.DeletTask(tasks[len(tasks)-1].ID)
	tasks, err = db.Tasks(0, 0)
	if err != nil {
		log.Fatalf("Unable to get tasks: %v\n", err)
	}
	fmt.Printf("tasks: %v\n", tasks)
}

func makeTestData(db *postgr.Storage, title, content string, author int64) {
	var task postgr.Task
	task.Title = title
	task.Content = content
	task.AuthorID = author

	db.NewTask(task)
}
