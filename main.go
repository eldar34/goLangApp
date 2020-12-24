package main

import (
    "fmt"
    "net/http"
    "html/template"

    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"
    "os"
    "github.com/gorilla/mux"
)

type Task struct{
    Id uint16
    Title, Content string
}

var tasks = []Task{}
var singleTask = Task{}

// init is invoked before main()
func init() {
    // loads values from .env into the system
    if err := godotenv.Load(); err != nil {
        fmt.Println("No .env file found")
    }
}

func index(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")

    if err != nil {
        fmt.Fprintf(w, err.Error())
    }

    db, err := sql.Open("mysql", fmt.Sprintf(
        "%s:%s@tcp(%s:%s)/%s", 
        os.Getenv("DB_USERNAME"), 
        os.Getenv("DB_PASSWORD"), 
        os.Getenv("DB_HOST"), 
        os.Getenv("DB_PORT"), 
        os.Getenv("DB_NAME")))
    
    if err != nil {
        panic(err)
    }

    defer db.Close()

    // Create Table If Not Exists
    stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS `article` ("+        
            "id INT(11) NOT NULL AUTO_INCREMENT,"+
            "title VARCHAR(250) NULL DEFAULT NULL COLLATE 'utf8mb4_unicode_ci',"+
            "content TEXT(65535) NULL DEFAULT NULL COLLATE 'utf8mb4_unicode_ci',"+
            "PRIMARY KEY (`id`) USING BTREE)"+
        "COLLATE='utf8mb4_unicode_ci'"+
        "ENGINE=InnoDB")
            
    if err != nil {
        panic(err)
    }

    _, err = stmt.Exec()

    if err != nil {
        panic(err)
    }

    // Get data
    res, err := db.Query("SELECT * FROM `article`")
    if err != nil {
        panic(err)
    }

    tasks = []Task{}
    
    for res.Next() {
        var task Task
        err = res.Scan(&task.Id, &task.Title, &task.Content)
        if err != nil {
            panic(err)
        }

        // Slice content
        task.Content = task.Content[:7]

        tasks = append(tasks, task)

    }

    t.ExecuteTemplate(w, "index", tasks)
}

func add_task(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")

    if err != nil {
        fmt.Fprintf(w, err.Error())
    }

    t.ExecuteTemplate(w, "create", nil)
}

func save_task(w http.ResponseWriter, r *http.Request) {
    title := r.FormValue("title")
    content := r.FormValue("content")

    if title == "" || content == "" {
        fmt.Fprintf(w, "Fields can't be empty")
    }else{

        db, err := sql.Open("mysql", fmt.Sprintf(
            "%s:%s@tcp(%s:%s)/%s", 
            os.Getenv("DB_USERNAME"), 
            os.Getenv("DB_PASSWORD"), 
            os.Getenv("DB_HOST"), 
            os.Getenv("DB_PORT"), 
            os.Getenv("DB_NAME")))
    
        if err != nil {
            panic(err)
        }

        defer db.Close()

        // Insert data
        insert, err := db.Query(fmt.Sprintf("INSERT INTO `article` (`title`, `content`) VALUES('%s', '%s')", title, content))
        
        if err != nil {
            panic(err)
        }

        defer insert.Close()

        http.Redirect(w, r, "/", http.StatusSeeOther)
    }

    
}

func get_task(w http.ResponseWriter, r *http.Request) {
    
    vars := mux.Vars(r)

    t, err := template.ParseFiles("templates/single_task.html", "templates/header.html", "templates/footer.html")

    if err != nil {
        fmt.Fprintf(w, err.Error())
    }
    
    db, err := sql.Open("mysql", fmt.Sprintf(
        "%s:%s@tcp(%s:%s)/%s", 
        os.Getenv("DB_USERNAME"), 
        os.Getenv("DB_PASSWORD"), 
        os.Getenv("DB_HOST"), 
        os.Getenv("DB_PORT"), 
        os.Getenv("DB_NAME")))

    if err != nil {
        panic(err)
    }

    defer db.Close()

    // Get data
    res, err := db.Query(fmt.Sprintf("SELECT * FROM `article` WHERE `id` = '%s'", vars["id"]))

    if err != nil {
        panic(err)
    }

    singleTask = Task{}
    
    for res.Next() {
        var task Task
        err = res.Scan(&task.Id, &task.Title, &task.Content)
        if err != nil {
            panic(err)
        }

        singleTask = task

    }

    t.ExecuteTemplate(w, "single_task", singleTask)

}

func handleFunc() {
    rtr := mux.NewRouter()
    rtr.HandleFunc("/", index).Methods("GET")
    rtr.HandleFunc("/add_task/", add_task).Methods("GET")
    rtr.HandleFunc("/save_task/", save_task).Methods("POST")
    rtr.HandleFunc("/task/{id:[0-9]+}", get_task).Methods("GET")

    http.Handle("/", rtr)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
    http.ListenAndServe(":8080", nil)
}

func main() {

    handleFunc()

}