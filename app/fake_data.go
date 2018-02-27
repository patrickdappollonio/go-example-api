package app

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/icrowley/fake"
)

type user struct {
	ID   int     `json:"id"`
	User details `json:"user_details"`
	Addr address `json:"address"`
}

type details struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	JobTitle  string `json:"job_title"`
}

type address struct {
	Country string `json:"country"`
	State   string `json:"state"`
	City    string `json:"city"`
	Address string `json:"address"`
	Zip     string `json:"zip_code"`
}

type post struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Paragraphs string    `json:"content"`
	Date       time.Time `json:"date"`
	User       int       `json:"user_id"`
}

type domain struct {
	ID     int    `json:"id"`
	Domain string `json:"domain"`
	User   int    `json:"user_id"`
}

type product struct {
	ID        int     `json:"id"`
	Product   string  `json:"product_name"`
	Brand     string  `json:"brand"`
	StoreAddr address `json:"store_addr"`
	User      int     `json:"user_id"`
}

type slicer interface {
	length() int
	slice(x, y int) slicer
	uniq(pos int) interface{}
}

type users []user

func (u users) length() int                 { return len(u) }
func (u users) slice(begin, end int) slicer { return users(u[begin:end]) }
func (u users) uniq(pos int) interface{}    { return u[pos] }

type posts []post

func (u posts) length() int                 { return len(u) }
func (u posts) slice(begin, end int) slicer { return posts(u[begin:end]) }
func (u posts) uniq(pos int) interface{}    { return u[pos] }

type domains []domain

func (u domains) length() int                 { return len(u) }
func (u domains) slice(begin, end int) slicer { return domains(u[begin:end]) }
func (u domains) uniq(pos int) interface{}    { return u[pos] }

type products []product

func (u products) length() int                 { return len(u) }
func (u products) slice(begin, end int) slicer { return products(u[begin:end]) }
func (u products) uniq(pos int) interface{}    { return u[pos] }

var (
	systemUsers    = make(users, 126)
	systemPosts    = make(posts, 293)
	systemDomains  = make(domains, 408)
	systemProducts = make(products, 321)
)

func rnddate() time.Time {
	min := time.Date(2017, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2018, 2, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	return time.Unix(rand.Int63n(delta)+min, 0)
}

func generateRandomData() {
	for i := 0; i < len(systemUsers); i++ {
		systemUsers[i] = user{
			ID: i + 1,
			User: details{
				FirstName: fake.FirstName(),
				LastName:  fake.LastName(),
				Username:  fake.UserName(),
				Password:  fake.SimplePassword(),
				Email:     strings.ToLower(fake.EmailAddress()),
				JobTitle:  fake.JobTitle(),
			},
			Addr: address{
				Country: "USA",
				State:   fake.State(),
				City:    fake.City(),
				Address: fmt.Sprintf("%d %s", fake.Year(1100, 2000), fake.Street()),
				Zip:     fake.Zip(),
			},
		}
	}

	for i := 0; i < len(systemPosts); i++ {
		systemPosts[i] = post{
			ID:         i + 1,
			Title:      fake.Sentence(),
			Paragraphs: fake.ParagraphsN(6),
			Date:       rnddate(),
			User:       rand.Intn(len(systemUsers)-1) + 1,
		}
	}

	for i := 0; i < len(systemDomains); i++ {
		systemDomains[i] = domain{
			ID:     i + 1,
			Domain: fmt.Sprintf("%s.%s", strings.ToLower(fake.UserName()), fake.TopLevelDomain()),
			User:   rand.Intn(len(systemUsers)-1) + 1,
		}
	}

	for i := 0; i < len(systemProducts); i++ {
		systemProducts[i] = product{
			ID:      i + 1,
			Product: fake.ProductName(),
			Brand:   fake.Brand(),
			StoreAddr: address{
				Country: "USA",
				State:   fake.State(),
				City:    fake.City(),
				Address: fmt.Sprintf("%d %s", fake.Year(1100, 2000), fake.Street()),
				Zip:     fake.Zip(),
			},
			User: rand.Intn(len(systemUsers)-1) + 1,
		}
	}
}
