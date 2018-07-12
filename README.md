### 使用Golang和MongoDB构建 RESTful API

记录一下创建 **RESTful API**使用 `Go`开发语言和 `MongoDB`后台数据库  

[完整代码 Github](https://github.com/coderminer/restful)  

![](http://7xplrz.com1.z0.glb.clouddn.com//corderminer/website/golang_restful_api.png)

#### 安装依赖  

```
go get github.com/globalsign/mgo  // MongoDB的Go语言驱动
go get github.com/gorilla/mux   // Go语言的http路由库
```

#### API 结构预览

`app.go`  

```
package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func AllMovies(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented !")
}

func FindMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented !")
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented !")
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented !")
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented !")
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/movies", AllMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", FindMovie).Methods("GET")
	r.HandleFunc("/movies", CreateMovie).Methods("POST")
	r.HandleFunc("/movies", UpdateMovie).Methods("PUT")
	r.HandleFunc("/movies", DeleteMovie).Methods("DELETE")

	http.ListenAndServe(":8080", r)
}

```

随着后续路由的增加，需要重构一个这个文件，仿照 `beego`的工程目录，我们也分别创建对应的 `controllers`
和 `routes`,改造一下上面的代码

* 控制器

创建 `controllers` 文件夹和对应的文件 `movies.go`  

`movies.go`

```
func AllMovies(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented !")
}

func FindMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented !")
}

func CreateMovie(w http.ResponseWriter, t *http.Request) {
	fmt.Fprintln(w, "not implemented !")
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented !")
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented !")
}
```

* 路由

创建一个 `routes`文件夹，并创建对应的文件 `routes.go`  

`routes.go`  

```
package routes

import (
	"net/http"

	"github.com/coderminer/restful/controllers"
	"github.com/gorilla/mux"
)

type Route struct {
	Method     string
	Pattern    string
	Handler    http.HandlerFunc
	Middleware mux.MiddlewareFunc
}

var routes []Route

func init() {
	register("GET", "/movies", controllers.AllMovies, nil)
	register("GET", "/movies/{id}", controllers.FindMovie, nil)
	register("POST", "/movies", controllers.CreateMovie, nil)
	register("PUT", "/movies", controllers.UpdateMovie, nil)
	register("DELETE", "/movies", controllers.DeleteMovie, nil)
}

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	for _, route := range routes {
		r.Methods(route.Method).
			Path(route.Pattern).
			Handler(route.Handler)
		if route.Middleware != nil {
			r.Use(route.Middleware)
		}
	}
	return r
}

func register(method, pattern string, handler http.HandlerFunc, middleware mux.MiddlewareFunc) {
	routes = append(routes, Route{method, pattern, handler, middleware})
}

```

重构后的 `app.go`  

```
package main

import (
	"net/http"

	"github.com/coderminer/restful/routes"
)

func main() {
	r := routes.NewRouter()

	http.ListenAndServe(":8080", r)
}
```

#### Models

* 创建 `models` 文件夹和对应的文件 `db.go`(数据层)，封装对`MongoDB`的封装  

```
package models

import (
	"log"

	"github.com/globalsign/mgo"
)

const (
	host   = "127.0.0.1:27017"
	source = "admin"
	user   = "user"
	pass   = "123456"
)

var globalS *mgo.Session

func init() {
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{host},
		Source:   source,
		Username: user,
		Password: pass,
	}
	s, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatalln("create session error ", err)
	}
	globalS = s
}

func connect(db, collection string) (*mgo.Session, *mgo.Collection) {
	s := globalS.Copy()
	c := s.DB(db).C(collection)
	return s, c
}

func Insert(db, collection string, docs ...interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Insert(docs...)
}

func FindOne(db, collection string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).One(result)
}

func FindAll(db, collection string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Select(selector).All(result)
}

func Update(db, collection string, query, update interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Update(query, update)
}

func Remove(db, collection string, query interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Remove(query)
}

```

* 业务逻辑层  `models/movies.go`

```
package models

import "github.com/globalsign/mgo/bson"

type Movies struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	CoverImage  string        `bson:"cover_image" json:"cover_image"`
	Description string        `bson:"description" json:"description"`
}

const (
	db         = "Movies"
	collection = "MovieModel"
)

func (m *Movies) InsertMovie(movie Movies) error {
	return Insert(db, collection, movie)
}

func (m *Movies) FindAllMovies() ([]Movies, error) {
	var result []Movies
	err := FindAll(db, collection, nil, nil, &result)
	return result, err
}

func (m *Movies) FindMovieById(id string) (Movies, error) {
	var result Movies
	err := FindOne(db, collection, bson.M{"_id": bson.ObjectIdHex(id)}, nil, &result)
	return result, err
}

func (m *Movies) UpdateMovie(movie Movies) error {
	return Update(db, collection, bson.M{"_id": movie.Id}, movie)
}

func (m *Movies) RemoveMovie(id string) error {
	return Remove(db, collection, bson.M{"_id": bson.ObjectIdHex(id)})
}

```

#### 控制器逻辑  

* CreateMovie

```
func CreateMovie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movie models.Movies
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		responseWithJson(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	movie.Id = bson.NewObjectId()
	if err := dao.InsertMovie(movie); err != nil {
		responseWithJson(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseWithJson(w, http.StatusCreated, movie)
}
```

使用 `Postman` 或 `curl`进行测试

* AllMovies

```
func AllMovies(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var movies []models.Movies
	movies, err := dao.FindAllMovies()
	if err != nil {
		responseWithJson(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseWithJson(w, http.StatusOK, movies)

}
```

* FindMovie

```
func FindMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	result, err := dao.FindMovieById(id)
	if err != nil {
		responseWithJson(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseWithJson(w, http.StatusOK, result)
}
```

> http://127.0.0.1:8080/movies/5b45ef3a9d5e3e197c15e774  

* UpdateMovie

```
func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var params models.Movies
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		responseWithJson(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.UpdateMovie(params); err != nil {
		responseWithJson(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}
```

* DeleteMovie

```
func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if err := dao.RemoveMovie(id); err != nil {
		responseWithJson(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	responseWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}
```
