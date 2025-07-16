package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

type Task struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

var fileName = "newTasks.json"
var reader = bufio.NewReader(os.Stdin)
var line = strings.Repeat("-", 40)

func main() {

	var command string
	for {
		fmt.Println("Please enter action (add,update,delete,mark,list) or cancel to stop:")
		command, _ = reader.ReadString('\n')
		command = strings.TrimSpace(command)

		switch command {
		case "add":
			// add function
			var taskDescription string
			fmt.Print("Enter task description: ")
			taskDescription, _ = reader.ReadString('\n')
			taskDescription = strings.TrimSpace(taskDescription)
			outputMessage := addTask(taskDescription, fileName)
			fmt.Println(outputMessage)
			fmt.Println(line)

		case "update":
			// update function
			var newTaskDescription string
			fmt.Print("Enter task id to be updated: ")
			taskID, errorMessage := readTaskID()
			if errorMessage != "" {
				fmt.Println(errorMessage)
				continue
			}
			fmt.Print("Enter new task description: ")
			newTaskDescription, _ = reader.ReadString('\n')
			newTaskDescription = strings.TrimSpace(newTaskDescription)
			outputMessage := updateTask(taskID-1, newTaskDescription, fileName)
			fmt.Println(outputMessage)
			fmt.Println(line)

		case "delete":
			// delete function
			fmt.Print("Enter task id to be removed: ")
			taskID, errorMessage := readTaskID()
			if errorMessage != "" {
				fmt.Println(errorMessage)
				continue
			}
			outputMessage := deleteTask(taskID-1, fileName)
			fmt.Println(outputMessage)

		case "mark":
			// mark function
			var taskProgress string
			statuses := []string{
				"1: todo",
				"2: in-progress",
				"3: done",
			}
			fmt.Print("Enter task id to be updated: ")
			taskID, errorMessage := readTaskID()
			if errorMessage != "" {
				fmt.Println(errorMessage)
				continue
			}
			fmt.Println("--------Choices--------")
			for _, status := range statuses {
				fmt.Println(status)
			}
			fmt.Print("Enter task progress:")
			taskProgress, _ = reader.ReadString('\n')
			taskProgress = strings.TrimSpace(taskProgress)
			outputMessage := markTask(taskID, taskProgress, fileName)
			fmt.Println(outputMessage)
			fmt.Println(line)

		case "list":
			// list function
			var choice string
			options := []string{
				"1: all",
				"2: not",
				"3: in-progress",
				"4: done",
			}
			fmt.Println("--------Choices--------")
			for _, option := range options {
				fmt.Println(option)
			}
			fmt.Print("Enter choice:")
			choice, _ = reader.ReadString('\n')
			choice = strings.TrimSpace(choice)
			taskList, errorMessage := listTask(choice, fileName)
			fmt.Println()
			if errorMessage != "" {
				fmt.Println(errorMessage)
			} else if len(taskList) == 0 {
				fmt.Println("No tasks found!")
			} else {
				printTasks(taskList)
			}
			fmt.Println(line)

		case "cancel":
			return
		default:
			fmt.Println("Invalid choice, please try again")
			continue
		}
	}

}

func getTask(fileName string) ([]Task, string) {
	// Read exiting tasks from file (if any)
	var tasks []Task
	var errorMessage string
	fileData, err := os.ReadFile(fileName)
	if err != nil && !os.IsNotExist(err) {
		errorMessage = "Error reading file, please try again"
		log.Fatal("Error reading file, Error code:", err)

	}

	// If there are existing tasks
	if len(fileData) != 0 {
		if err := json.Unmarshal([]byte(fileData), &tasks); err != nil {
			errorMessage = "Error unmarshalling code, please try again"
			log.Fatal("Error unmarshalling code, Error code:", err)
		}
	}
	return tasks, errorMessage
}

func writeToFile(fileName string, tasks []Task) string {
	jsonBytes, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		log.Fatal("Error marshalling task, Error code:", err)
		return "Error marshalling code, please try again"
	}

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("Error opening file, Error code:", err)
		return "Error opening file, please try again"
	}
	if err := f.Close(); err != nil {
		log.Fatal("Error closing file, Error code:", err)
		return "Error closing file, please try again"
	}
	defer f.Close()

	err = os.WriteFile(fileName, jsonBytes, 0666)
	if err != nil {
		log.Fatal("Error writing file, Error code:", err)
		return "Error writing file, please try again"
	}

	return "Success"
}

func printTasks(tasks []Task) {
	for _, task := range tasks {
		fmt.Println("---------- Task ----------")
		fmt.Printf("ID: %d\n", task.Id)
		fmt.Printf("Description: %s\n", task.Description)
		fmt.Printf("Status: %s\n", task.Status)
		fmt.Printf("Created At: %s\n", task.CreatedAt)

		// Only print UpdatedAt if it's not empty
		if task.UpdatedAt != "" {
			fmt.Printf("Updated At: %s\n", task.UpdatedAt)
		} else {
			fmt.Println("Updated At: (not updated)")
		}

		fmt.Println()
	}
}

func readTaskID() (int, string) {
	var taskID int
	var errorMessage string
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	taskID, err := strconv.Atoi(input)
	if err != nil || taskID <= 0 {
		errorMessage = "Invalid input, please try again"
	}
	return taskID, errorMessage
}

func addTask(taskDescription string, fileName string) string {
	var tasks []Task
	var errorMessage string
	tasks, errorMessage = getTask(fileName)

	if len(errorMessage) != 0 {
		return errorMessage
	}

	// Determine ID for new task
	newId := 1
	if len(tasks) > 0 {
		newId = tasks[len(tasks)-1].Id + 1
	}

	var newTask = Task{
		Id:          newId,
		Description: taskDescription,
		Status:      "Not Done",
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   "",
	}

	tasks = append(tasks, newTask)

	outputMessage := writeToFile(fileName, tasks)
	if outputMessage != "Success" {
		return outputMessage
	}
	return "Successfully Added New Task"

}

func updateTask(taskId int, newTaskDescription string, fileName string) string {
	var tasks []Task
	var errorMessage string
	tasks, errorMessage = getTask(fileName)

	// Check if any errors when opening file and taskID is valid
	if len(errorMessage) != 0 {
		return errorMessage
	}
	if len(tasks) == 0 || taskId >= len(tasks) {
		return "No tasks found, please try again"
	}

	// Update task description
	tasks[taskId].Description = newTaskDescription
	tasks[taskId].UpdatedAt = time.Now().Format(time.RFC3339)

	// Update json file
	outputMessage := writeToFile(fileName, tasks)
	if outputMessage != "Success" {
		return outputMessage
	}
	return "Successfully Updated Task"
}

func markTask(taskId int, taskProgress string, fileName string) string {
	var tasks []Task
	var errorMessage string
	tasks, errorMessage = getTask(fileName)

	// Check if any errors when opening file and taskID is valid
	if len(errorMessage) != 0 {
		return errorMessage
	}
	if len(tasks) == 0 || taskId >= len(tasks) {
		return "No tasks found, please try again"
	}

	// Update task progress
	switch taskProgress {
	case "1":
		tasks[taskId].Status = "Not Done"
	case "2":
		tasks[taskId].Status = "In-Progress"
	case "3":
		tasks[taskId].Status = "Done"
	default:
		//fmt.Println("Invalid choice, please try again")
		return "Invalid choice, please try again"
	}
	tasks[taskId].UpdatedAt = time.Now().Format(time.RFC3339)

	// Update json file
	outputMessage := writeToFile(fileName, tasks)
	if outputMessage != "Success" {
		return outputMessage
	}
	return "Successfully Marked Task"
}

func listTask(choice string, fileName string) ([]Task, string) {
	var tasks []Task
	var errorMessage string
	tasks, errorMessage = getTask(fileName)
	if errorMessage != "" {
		return nil, errorMessage
	}
	switch choice {
	case "1":
		// all
		return tasks, errorMessage
	case "2":
		// to-do
		var todoTasks []Task
		for _, task := range tasks {
			if task.Status == "Not Done" {
				todoTasks = append(todoTasks, task)
			}
		}
		return todoTasks, errorMessage
	case "3":
		// in progress
		var inProgressTasks []Task
		for _, task := range tasks {
			if task.Status == "In-Progress" {
				inProgressTasks = append(inProgressTasks, task)
			}
		}
		return inProgressTasks, errorMessage
	case "4":
		// done
		var doneTasks []Task
		for _, task := range tasks {
			if task.Status == "Done" {
				doneTasks = append(doneTasks, task)
			}
		}
		return doneTasks, errorMessage
	default:
		return nil, "Invalid choice, please try again"
	}
}

func deleteTask(taskId int, fileName string) string {
	var tasks []Task
	var errorMessage string
	tasks, errorMessage = getTask(fileName)

	// Check if any errors when opening file and taskID is valid
	if len(errorMessage) != 0 {
		return errorMessage
	}
	if len(tasks) == 0 || taskId >= len(tasks) {
		return "No tasks found, please try again"
	}
	newTaskList := append(tasks[:taskId], tasks[taskId+1:]...)
	for i := taskId; i < len(newTaskList); i++ {
		newTaskList[i].Id = i + 1
	}
	outputMessage := writeToFile(fileName, newTaskList)
	if outputMessage != "Success" {
		return outputMessage
	}
	return "Successfully Deleted Task"
}
