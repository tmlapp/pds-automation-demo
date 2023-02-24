package lib

import (
	"database/sql"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	mysqlHost string = os.Getenv("MYSQL_HOST")
	mysqlUser string = os.Getenv("MYSQL_USER")
	mysqlPass string = os.Getenv("MYSQL_PASS")
	mysqlPort string = os.Getenv("MYSQL_PORT")
)

type PastOrders struct {
	OrderId    int
	Main       string
	Side1      string
	Side2      string
	Drink      string
	Restaurant string
	Date       string
}

type myOrderData struct {
	PageTitle string
	Orders    []PastOrders
}

func generateOrderID() int {
	rand.Seed(time.Now().UnixNano())
	orderId := rand.Intn(100000)
	return (orderId)
}

func orderStatus(w http.ResponseWriter, r *http.Request, messageData PageData) {
	fmt.Println("method:", r.Method, "on URL:", r.URL)
	session, _ := store.Get(r, "cookie-name")
	t, _ := template.ParseFiles("./static/order_status.html")
	if r.Method == "POST" {
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			t, _ = template.ParseFiles("./static/external_order_status.html")
		}
		t.Execute(w, messageData)
	}
}

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method, "on URL:", r.URL)
	session, _ := store.Get(r, "cookie-name")
	fmt.Println(session.Values["authenticated"])
	fmt.Println(session.Values["email"])
	fmt.Println(r.Method)
	//generate Order ID
	//orderNum := generateOrderID()

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		fmt.Println("Not Authenticated")
	} else {
		fmt.Println("Authenticated")
	}

	if r.Method == "GET" {
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			fmt.Println("Order form requested, but unauthenticated; redirecting to login page.")
			t, _ := template.ParseFiles("./static/login.html")
			t.Execute(w, nil)
		} else {
			fmt.Printf("should allow order")
			t, _ := template.ParseFiles("./static/order.html")
			t.Execute(w, nil)
		}
	}
}

func MyOrderHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method, "on URL:", r.URL)
	session, _ := store.Get(r, "cookie-name")
	fmt.Println(session.Values["email"])
	email := session.Values["email"].(string)

	//myOrders := GetMyOrders(email)
	myOrders := MyOrderHistory(email)

	data := myOrderData{
		PageTitle: "My Order History",
		Orders:    myOrders,
	}

	t, _ := template.ParseFiles("./static/order_history.html")
	t.Execute(w, data)

}

func PxbbqOrderHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method, "on URL:", r.URL)
	session, _ := store.Get(r, "cookie-name")
	fmt.Println(session.Values["authenticated"])
	fmt.Println(session.Values["email"])
	fmt.Println(r.Method)

	//generate Order ID
	orderNum := generateOrderID()

	//verify the user is authenticated and if not send them to the login page instead
	if r.Method == "GET" {
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			fmt.Println("Order form requested, but unauthenticated; redirecting to login page.")
			t, _ := template.ParseFiles("./static/login.html")
			t.Execute(w, nil)
		} else {
			fmt.Printf("should allow order")
			t, _ := template.ParseFiles("./static/pxbbq_order.html")
			t.Execute(w, nil)
		}
	} else {
		r.ParseForm()
		statusData := PageData{
			PageTitle: "Order Status",
			Message:   fmt.Sprintf("Your order has been received. Order number %v", orderNum),
		}

		//Organize form submission data for a write to storage
		currentTime := time.Now() //used in the order submission
		email := session.Values["email"].(string)
		main := r.FormValue("main")
		side1 := r.FormValue("side1")
		side2 := r.FormValue("side2")
		drink := r.FormValue("drink")
		restaurant := r.FormValue("restaurant")

		//retrieve address from customer's account
		myAddress := GetAddress(email)
		street1 := myAddress.Street1
		street2 := myAddress.Street2
		city := myAddress.City
		state := myAddress.State
		zipcode := myAddress.Zipcode
		orderDate := currentTime.Format("2 January 2006")

		//log order to std out - Can be used for troubleshooting
		fmt.Printf("Order submitted by: ")
		fmt.Println(session.Values["email"].(string))
		fmt.Println("Order Taken by :" + restaurant)
		fmt.Println("Ordered : " + main)
		fmt.Println("Ordered : " + side1)
		fmt.Println("Ordered : " + side2)
		fmt.Println("Ordered : " + drink)
		fmt.Println("Street1 : " + street1)
		fmt.Println("Street2 : " + street2)
		fmt.Println("City  : " + city)
		fmt.Println("State  : " + state)
		fmt.Println("Zip  : " + zipcode)
		fmt.Print("########")

		//submit order to storage
		SubmitOrder(orderNum, orderDate, email, restaurant, main, side1, side2, drink, street1, street2, city, state, zipcode)

		//Display Operation Status Page to User
		orderStatus(w, r, statusData)
	}
}

func McdOrderHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method, "on URL:", r.URL)
	session, _ := store.Get(r, "cookie-name")
	fmt.Println(session.Values["authenticated"])
	fmt.Println(session.Values["email"])
	fmt.Println(r.Method)

	//generate Order ID
	orderNum := generateOrderID()

	//verify the user is authenticated and if not send them to the login page instead
	if r.Method == "GET" {
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			fmt.Println("Order form requested, but unauthenticated; redirecting to login page.")
			t, _ := template.ParseFiles("./static/login.html")
			t.Execute(w, nil)
		} else {
			fmt.Printf("should allow order")
			t, _ := template.ParseFiles("./static/mcd_order.html")
			t.Execute(w, nil)
		}
	} else {
		r.ParseForm()
		statusData := PageData{
			PageTitle: "Order Status",
			Message:   fmt.Sprintf("Your order has been received. Order number %v", orderNum),
		}

		//Organize form submission data for a write to storage
		currentTime := time.Now() //used in the order submission
		email := session.Values["email"].(string)
		main := r.FormValue("main")
		side1 := r.FormValue("side1")
		side2 := r.FormValue("side2")
		drink := r.FormValue("drink")
		restaurant := r.FormValue("restaurant")

		//retrieve address from customer's account
		myAddress := GetAddress(email)
		street1 := myAddress.Street1
		street2 := myAddress.Street2
		city := myAddress.City
		state := myAddress.State
		zipcode := myAddress.Zipcode
		orderDate := currentTime.Format("2 January 2006")

		//log order to std out - Can be used for troubleshooting
		fmt.Printf("Order submitted by: ")
		fmt.Println(session.Values["email"].(string))
		fmt.Println("Order Taken by :" + restaurant)
		fmt.Println("Ordered : " + main)
		fmt.Println("Ordered : " + side1)
		fmt.Println("Ordered : " + side2)
		fmt.Println("Ordered : " + drink)
		fmt.Println("Street1 : " + street1)
		fmt.Println("Street2 : " + street2)
		fmt.Println("City  : " + city)
		fmt.Println("State  : " + state)
		fmt.Println("Zip  : " + zipcode)
		fmt.Print("########")

		//submit order to storage
		SubmitOrder(orderNum, orderDate, email, restaurant, main, side1, side2, drink, street1, street2, city, state, zipcode)

		//Display Operation Status Page to User
		orderStatus(w, r, statusData)
	}
}

func CentralperkOrderHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method, "on URL:", r.URL)
	session, _ := store.Get(r, "cookie-name")
	fmt.Println(session.Values["authenticated"])
	fmt.Println(session.Values["email"])
	fmt.Println(r.Method)

	//generate Order ID
	orderNum := generateOrderID()

	//verify the user is authenticated and if not send them to the login page instead
	if r.Method == "GET" {
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			fmt.Println("Order form requested, but unauthenticated; redirecting to login page.")
			t, _ := template.ParseFiles("./static/login.html")
			t.Execute(w, nil)
		} else {
			fmt.Printf("should allow order")
			t, _ := template.ParseFiles("./static/centralperk_order.html")
			t.Execute(w, nil)
		}
	} else {
		r.ParseForm()
		statusData := PageData{
			PageTitle: "Order Status",
			Message:   fmt.Sprintf("Your order has been received. Order number %v", orderNum),
		}

		//Organize form submission data for a write to storage
		currentTime := time.Now() //used in the order submission
		email := session.Values["email"].(string)
		main := r.FormValue("main")
		side1 := r.FormValue("side1")
		side2 := r.FormValue("side2")
		drink := r.FormValue("drink")
		restaurant := r.FormValue("restaurant")

		//retrieve address from customer's account
		myAddress := GetAddress(email)
		street1 := myAddress.Street1
		street2 := myAddress.Street2
		city := myAddress.City
		state := myAddress.State
		zipcode := myAddress.Zipcode
		orderDate := currentTime.Format("2 January 2006")

		//log order to std out - Can be used for troubleshooting
		fmt.Printf("Order submitted by: ")
		fmt.Println(session.Values["email"].(string))
		fmt.Println("Order Taken by : " + restaurant)
		fmt.Println("Ordered : " + main)
		fmt.Println("Ordered : " + side1)
		fmt.Println("Ordered : " + side2)
		fmt.Println("Ordered : " + drink)
		fmt.Println("Street1 : " + street1)
		fmt.Println("Street2 : " + street2)
		fmt.Println("City  : " + city)
		fmt.Println("State  : " + state)
		fmt.Println("Zip  : " + zipcode)
		fmt.Print("########")

		//submit order to storage
		SubmitOrder(orderNum, orderDate, email, restaurant, main, side1, side2, drink, street1, street2, city, state, zipcode)

		//Display Operation Status Page to User
		orderStatus(w, r, statusData)
	}
}

func MyOrderHistory(email string) []PastOrders {
	dsn := mysqlUser + ":" + mysqlPass + "@tcp(" + mysqlHost + ":" + mysqlPort + ")/delivery"
	//fmt.Println("DSN is : " + dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		//log.Fatal(err)
		fmt.Println("Could not connect to Order History Database")
		return (nil)
	}

	defer db.Close()

	rows, err := db.Query("select orderid, main, side1, side2, drink, restaurant, date from orders where email = " + "'" + email + "'")
	if err != nil {
		//log.Fatal(err)
		fmt.Println("No Orders History exists for user: " + email)
		return (nil)
	}

	defer rows.Close()

	var myOrders []PastOrders

	for rows.Next() {
		var order PastOrders
		if err != nil {
			fmt.Println("Cannot push row into array")
		}
		//myOrders = append(myOrders, order)
		err := rows.Scan(&order.OrderId, &order.Main, &order.Side1, &order.Side2, &order.Drink, &order.Restaurant, &order.Date)
		if err != nil {
			fmt.Println(err)
		} else {
			myOrders = append(myOrders, order)

			fmt.Println("ATTEMPTING TO PRINT FROM STRUCT AS A TEST")
			fmt.Println(order)
		}

	}

	return (myOrders)
}
