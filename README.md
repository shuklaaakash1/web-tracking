# web-tracking
My Awesome Web Application
Welcome to My Awesome Web Application! This application provides various features for managing users and products, including user authentication, product creation, and recommendations.

Table of Contents
Features
Getting Started
Prerequisites
Installation
Usage
Endpoints
Contributing
User authentication using JSON Web Tokens (JWT)
Creating and managing users and products
Providing product recommendations based on interactions and categories
Getting Started
Features
postman collection 
er digrame


*************************
Prerequisites
Go programming language (https://golang.org/dl/)
MySQL database server
Installation

Clone this repository to your local machine:

bash
Copy code
git clone https://github.com/shuklaaakash1/web-tracking.git
Install the required Go packages using Go modules:

bash
Copy code
go mod tidy
Update the MySQL database connection details in main.go by modifying the ConnectDatabase function.

Run the application:

bash
Copy code:-
go run main.go auth.go authhandler.go db.go loginhandler.go modleproduct.go modlesinterc.go modlesuser.go pr.go productsorthandler.go  recommandation.go tokengen.go



Usage
Make sure the application is running by following the installation steps.


Use a tool like curl or Postman to interact with the 



**************************************
endpoints.
See the Endpoints section for available routes and how to use them.
Endpoints
POST /login: User login endpoint to get an authentication token.
GET /protected: Protected route that requires a valid token to access.
POST /create-user: Create a new user.
POST /create-product: Create a new product.
GET /product/{id}: Get details of a product by ID.
GET /products/sort: Get sorted list of products.
Contributing
Contributions are welcome! If you find any issues or want to enhance the application, feel free to create a pull request.

1 POST /login: User Login Endpoint

Description: This endpoint allows users to log in and receive an authentication token that can be used to access protected routes.
Request: Send a JSON payload with the user's username and email for authentication.
Response: Returns a JSON object containing the token for successful login.
Usage: After obtaining the token, include it in the Authorization header for protected routes.
            ex:{
"username": "john_doe;",
  "email": "aakash@gmaicom"
  
}

2   GET /protected: Protected Route

Description: This is a protected route that requires a valid authentication token to access.
Request: Requires a valid Authorization header with the Bearer token obtained from the login endpoint.
Response: Returns a response indicating successful access to the protected content.
Usage: Demonstrates how to secure routes using token-based authentication.   ex
oauth2 enter token into token aviable section


3 POST /create-user: Create User

Description: Allows the creation of a new user account.
Request: Send a JSON payload with username, email, password, and role fields.
Response: Returns a response indicating the successful creation of the user.
Usage: Used to register new users in the system.
POST /create-product: Create Product

Description: Enables the addition of a new product to the system.
Request: Send a JSON payload with name, category, and price fields.
Response: Returns a response indicating the successful creation of the product.
Usage: Allows administrators to add new products to the inventory.  ex:={    
    "username":"john_doe",
  "username": "aakash@gmail.com",
  "password": "shukla",
  "role": "user"
}


4 GET /product/{id}: Get Product Details

Description: Retrieves detailed information about a specific product by its ID.
Request: Provide the id parameter in the URL to specify the product.
Response: Returns a JSON object containing detailed information about the requested product.
Usage: Used to display product details to users.



5.GET /products/sort: Sort Products

Description: Returns a list of products sorted by the specified criteria.
Request: Provide query parameters sortBy (field to sort by) and order (ascending or descending order).
Response: Returns a JSON array of products sorted according to the specified criteria.
Usage: Allows users to retrieve a sorted list of products, e.g., by price or creation date.


6.POST /create-product: Create Product

Description: Enables the addition of a new product to the system.
Request: Send a JSON payload with name, category, and price fields.
Response: Returns a response indicating the successful creation of the product.
Usage: Allows administrators to add new products to the inventory.