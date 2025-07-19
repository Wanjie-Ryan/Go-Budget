package common

import (
	"math"
	"net/http"
	"strconv"

	// "github.com/Wanjie-Ryan/Go-Budget/internal/models"
	"gorm.io/gorm"
)

// define a struct that will hold the data that the FE needs to pass when requesting for data

type Pagination struct {
	Limit     int         `query:"limit" json:"limit"` // how much data does the user retrieve per request? 10, 20, 30 etc etc. How many items per page
	Page      int         `query:"page" json:"page"`   // defines the page the user wants to get. Which page the user wants.
	Sort      string      `query:"sort"`               //sorts results like ASC or DESC
	TotalRows int64       `json:"total_rows"`          // how many records exist in the DB
	TotalPage int         `json:"total_page"`          // total number of pages based on total records and limit
	Items     interface{} `json:"items"`
}

// SCENARIO
// we have 100 categories
// want to retrieve a limit of 10 --> limit of 10
// for the first page, page 1 --> page of 1
// total rows are 100
// Total page = Total rows / limit

// define functions to set the pagination params
// this function will return to us the page that the user actually passed
func (p *Pagination) GetPage() int {

	if p.Page <= 0 {
		p.Page = 1
	}

	return p.Page
}

// function to get the maximum number of items user can get per page
func (p *Pagination) GetLimit() int {
	if p.Limit > 100 {
		p.Limit = 100
	} else if p.Limit <= 0 {
		p.Limit = 10
	}
	return p.Limit
}

// offset method
// calculates how many rows to skip in the DB to get to the desired page
// for page 1, skip 0 rows (1-1) * 10 = 0
// for page 2, skip 10 rows (2-1) * 10 = 10
func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

// method to instantiate the pagination  struct
// the function will accept the model that we want to paginate, which will be an interface, meaning it may hold any type of data, also pass in the request, and the DB
// it will return a pagination struct
// initializes pagination object
func NewPagination(models interface{}, r *http.Request, db *gorm.DB) *Pagination {

	// declare a pah=gination variable that will hold the value
	var pagination Pagination
	// define a  request to retrieve the query urls
	q := r.URL.Query()
	// get the page from the url query
	page, _ := strconv.Atoi(q.Get("page")) // the page number the FE wants to retrieve data from

	// get the limit from the url query
	limit, _ := strconv.Atoi(q.Get("limit")) // how much data per page the FE wants to retrieve

	// define the total number of rows from the db using the count method
	// variable to hold the total number of rows
	var totalRows int64
	db.Model(models).Count(&totalRows)

	pagination.Page = page
	pagination.Limit = limit
	pagination.TotalRows = totalRows

	// after storing the variables, the next thing is to calculate the total pages

	// use the math.Ceil to ensure we always round up
	// if you have 91 rows and a limit of 10 -> 91/10 = 9.1 which ceiling is 10 pages
	totalPage := int(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))

	pagination.TotalPage = totalPage

	return &pagination
}

// the function below returns a **closure** which is a function that returns another function
// the returned function accepts a GORM object and applies both the offset and limit to it.
func (p *Pagination) Paginate() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.GetOffset()).Limit(p.GetLimit())
	}
}

// offset --> determines the starting point of the records to be retreived from a data source. it helps fetch specific set or records based on page number and number of records per page (limit).
// offset refers to how many records to skip before starting to return results
// if you want to retrieve results for a particular page say page 3, the offset tells the system where to begin retrieving the records from, based on the page number and limit
// lets say you have 10 items (rows) you want to show 10 items per page.
// page 1 will retrieve items 1-10 --> strating from offset = 0
// page 2 will retrieve items 11-20 --> starting from offset = 10
// page 3 will retrieve items 21-30 --> starting from offset = 20

// limit = 10 items per page
// page = 1,2,3 (requested page)
// offset formula = (page-1) * limit

// offset for page 1: (1-1) * 10 = 0
// offset for page 2: (2-1) * 10 = 10
// offset for page 3: (3-1) * 10 = 20
