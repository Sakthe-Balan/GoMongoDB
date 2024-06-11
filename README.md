

GoMongoDB: File-Based Database System in Go
===========================================

Introduction
------------

GoMongoDB is a lightweight, file-based database system written in Go. It provides an easy-to-use interface for performing CRUD (Create, Read, Update, Delete) operations on data collections. Inspired by MongoDB, it supports a range of query operators for searching and filtering data.
![mmo](https://github.com/Sakthe-Balan/GoMongoDB/assets/103580234/c76c86c2-aba5-46a5-ba58-967bc0a4d25f)

Features
--------

* **Simple CRUD Operations:** Perform basic Create, Read, Update, and Delete operations on data collections.
* **Flexible Querying:** Utilize MongoDB-like query operators for advanced search and filtering capabilities.
* **Concurrency Safety:** Ensures data integrity and prevents race conditions using mutexes.
* **File-Based Storage:** Stores data persistently in files, ensuring durability and portability.

Data Integrity and Concurrency Management
-----------------------------------------

GoMongoDB employs mutexes to ensure data integrity and manage concurrent access safely. Each collection has its own mutex, allowing for concurrent operations without compromising data consistency.

Getting Started
---------------

To get started with GoMongoDB, follow these steps:

1. **Clone the repository:**
    ```bash
    git clone https://github.com/Sakthe-Balan/GoMongoDB.git
    ```

2. **Run the main.go file:**
    ```bash
    go run main.go
    ```


Alternatively, you can use the Dockerfile to build and run the application:

    # Build the Docker image
    docker build -t go-mongodb-app .
    
    # Run the Docker container
    docker run -p 6942:6942 go-mongodb-app

### Collections and Resources

In GoMongoDB, data is organized into collections, which act as containers for storing related documents. Each document within a collection is a JSON-like object representing a single record or entity. 

- **Collections:** Collections are analogous to tables in relational databases or folders in a file system. They provide a way to logically group similar data together.

- **Resources:** Resources refer to individual documents within a collection. Each resource contains fields with corresponding values, similar to attributes in a table row or properties in an object.

When performing CRUD (Create, Read, Update, Delete) operations, you specify both the collection and the resource. This allows for efficient management and retrieval of data within GoMongoDB.



# Instructions to Use Dashboard

To interact with the database through the dashboard, follow these steps:

1. **Navigate to the dashboard directory in your terminal.**

2. **Put your Streamlit UI file (let's call it `ui.py`) in this directory.**

3. **Run the Streamlit app using the following command:**

        streamlit run ui.py


This command will start the Streamlit server and launch your dashboard UI in a web browser. You can then use the dashboard to interact with the database, perform CRUD operations, and visualize data as needed.

    

Endpoints
---------

### 1\. Create Resource

**Endpoint:** /write

**Method:** POST

**Description:** Adds a new resource to a specified collection.

**Parameters:**

* `collection`: The name of the collection.
* `resource`: The name of the resource.

**Example Usage:**

    curl -X POST "http://localhost:6942/write?collection=<Collection>&resource=<resource>" \
    -H "Content-Type: application/json" \
    -d '{"name":"John Doe","age":35,"city":"New York"}'

### Example Usage (Python)

    import requests
    
    url = 'http://localhost:6942/write?collection=&resource='
    headers = {'Content-Type': 'application/json'}
    data = {
        "name": "John Doe",
        "age": 35,
        "city": "New York",
       
    }
    
    response = requests.post(url, json=data, headers=headers)
    print(response)
    

### 2\. Read Resource

**Endpoint:** /read

**Method:** GET

**Description:** Retrieves a specific resource from a collection.

**Parameters:**

* `collection`: The name of the collection.
* `resource`: The name of the resource.

**Example Usage:**

    curl -X GET "http://localhost:6942/read?collection=<Collection>&resource=<resource>"

### Example Usage (Python)

    import requests
    
    url = 'http://localhost:6942/read?collection=&resource='
    headers = {'Content-Type': 'application/json'}
    
    
    response = requests.get(url, headers=headers)
    print(response.json())
    

### 3\. Read All Resources

**Endpoint:** /readAll

**Method:** GET

**Description:** Retrieves all resources from a collection.

**Parameters:**

* `collection`: The name of the collection.

**Example Usage:**

    curl -X GET "http://localhost:6942/readall?collection=<Collection>"

### Example Usage (Python)
    
    
    import requests
    
    url = 'http://localhost:6942/readall?collection='
    headers = {'Content-Type': 'application/json'}
    
    
    response = requests.get(url, headers=headers)
    print(response.json())

### 4\. Delete Resource

**Endpoint:** /delete

**Method:** DELETE

**Description:** Removes a specific resource from a collection.

**Parameters:**

* `collection`: The name of the collection.
* `resource`: The name of the resource.

**Example Usage:**

    curl -X DELETE "http://localhost:6942/delete?collection=<Collection>&resource=<resource>"
    
### Example Usage (Python)

    import requests
    
    url = 'http://localhost:6942/delete?collection=&resource='
    headers = {'Content-Type': 'application/json'}
  
    
    response = requests.delete(url, headers=headers)
    print(response)

### 5\. Delete All Resources

**Endpoint:** /deleteAll

**Method:** DELETE

**Description:** Removes all resources from a collection.

**Parameters:**

* `collection`: The name of the collection.

**Example Usage:**

    curl -X DELETE "http://localhost:6942/deleteall?collection=<Collection>"

### Example Usage (Python)

    
    import requests
    
    url = 'http://localhost:6942/deleteall?collection='
    headers = {'Content-Type': 'application/json'}
    
    
    response = requests.delete(url, headers=headers)
    print(response)
    

### 6\. Search

**Endpoint:** /search

**Method:** POST

**Description:** Search for resources using MongoDB-like query operators.

**Parameters:**

* `collection`: The name of the collection.
* `query`: MongoDB-like query to filter resources.

**Example Usage:**

    curl -X POST "http://localhost:6942/search?collection=<Collection>" \
    -H "Content-Type: application/json" \
    -d '{"age":{"$lte":30}}'

    
### Example Usage (Python)

    import requests

    
    url = 'http://localhost:6942/search?collection='
    headers = {'Content-Type': 'application/json'}

    data = {"age": {"$lte": 35}}

    response = requests.post(url, json=data, headers=headers)
    print(response.text)

### 7\. RegexSearch

**Endpoint:** /regexsearch

**Method:** POST

**Description:** Search for resources using regular expressions to match fields.

**Parameters:**

* `collection`: The name of the collection.
* `query`: A map of field names to regex patterns.

**Example Usage:**

    curl -X POST "http://localhost:6942/regexsearch?collection=<Collection>" \
    -H "Content-Type: application/json" \
    -d '{"name": "^John", "email": ".*@example\\.com$"}'

    
### Example Usage (Python)

    import requests
    
    url = 'http://localhost:6942/regexsearch?collection=<Collection>'
    headers = {'Content-Type': 'application/json'}
    
    data = {
        "name": "^John",
        "email": ".*@example\\.com$"
    }
    
    response = requests.post(url, json=data, headers=headers)
    print(response.text)

**MongoDB-like Query Operators:**

The search endpoint supports various MongoDB-like query operators such as `$eq`, `$ne`, `$gt`, `$gte`, `$lt`, `$lte`, `$in`, etc. Refer to the MongoDB documentation for more details
 on query operators.

Features to be Added
--------------------

* [X]  **Regex and Other Useful Functions:** Enhance querying capabilities by adding support for regular expressions and other useful functions.
* [ ]  **Distributed Framework for Handling Large Data:** Implement a distributed framework to handle large datasets, ensuring scalability and performance.
* [ ]  **Framework for Database Access:** Develop a framework to simplify and streamline access to the database, making integration and development more efficient.
* [ ]  **Replication and Sharding Logic:** Integrate replication and sharding logic to improve data availability, reliability, and distribution across multiple nodes.
* [X]  **Dashboard for Handling Data:**


Contributing
------------

Contributions are welcome! Feel free to open issues or pull requests on the GitHub repository.
