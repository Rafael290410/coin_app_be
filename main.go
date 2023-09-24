package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := connectDB()
	r := gin.Default()
	// mysql_pass: abc123!!!
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:19006"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/coin-rate", func(c *gin.Context) {
		// call coingecko api to get exchange rate
		resp, err := http.Get("https://api.coingecko.com/api/v3/exchange_rates")
		if err != nil {
			log.Fatalln(err)
		}

		// try to read body from coingecko response
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		// declare interface variable to store response as a JSON
		var jsonResponse map[string]interface{}

		//decoding the json data and storing in the variable
		err2 := json.Unmarshal([]byte(body), &jsonResponse)

		jsonResponse["updated_at"] = time.Now()
		jsonResponse["name"] = "rafael nadal"

		coins := jsonResponse["rates"]
		log.Println(coins)
		//iterateStructFields(coins)
		// for i, s := range coins {
		// 	fmt.Println(i, s)
		// }

		//Checks whether the error is nil or not
		if err2 != nil {
			//Prints the error if not nil
			fmt.Println("Error while decoding the data", err.Error())
		}

		// insert coin data to database\
		query := "INSERT INTO coins (alias,name,type,unit,status) VALUES (\"test\",\"test\", \"fiat\" , 99, \"inactive\")"
		insertResult, err := db.ExecContext(context.Background(), query)
		if err != nil {
			log.Fatalf("impossible insert coins: %s", err)
		}
		id, err := insertResult.LastInsertId()
		if err != nil {
			log.Fatalf("impossible to retrieve last inserted id: %s", err)
		}
		log.Printf("inserted id: %d", id)
		defer db.Close()

		// response data to client (web browser)
		c.JSON(200, jsonResponse)
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

func rafael_nadal() {
	panic("unimplemented")
}

func connectDB() *sql.DB {
	// build the DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "root", "abc123!!!", "127.0.0.1", 3306, "coin_watch")
	// Open the connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("impossible to create the connection: %s", err)
	} else {
		log.Printf("Mysql Connected")
	}

	return db
}

func iterateStructFields(input interface{}) {
	value1 := reflect.ValueOf(input)
	value := reflect.ValueOf(value1)
	numFields := value.NumField()
	log.Println(numFields)
	fmt.Printf("Number of fields: %d\n", numFields)

	structType := value.Type()

	for i := 0; i < numFields; i++ {
		field := structType.Field(i)
		fieldValue := value.Field(i)

		fmt.Printf("Field %d: %s (%s) = %v\n", i+1, field.Name, field.Type, fieldValue)
	}
}
