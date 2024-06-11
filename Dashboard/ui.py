import streamlit as st
import requests
import json

# Define base URL
base_url = 'http://localhost:6942'

# Helper function to send requests
def send_request(endpoint, method='GET', data=None):
    headers = {'Content-Type': 'application/json'}
    if method == 'GET':
        response = requests.get(f"{base_url}/{endpoint}", headers=headers)
    else:
        response = requests.post(f"{base_url}/{endpoint}", headers=headers, data=json.dumps(data))
    return response.json() if response.status_code == 200 else response.text

# Sidebar options
st.sidebar.title("Database Dashboard")
option = st.sidebar.selectbox(
    "Choose an action",
    ( "View Resources", "Add Resource", "Search", "Regex Search", "Delete Resource")
)

# View Collections
if option == "View Collections":
    st.title("Collections")
    response = send_request('collections')
    if isinstance(response, str):
        st.error(response)
    else:
        st.json(response)

# View Resources
elif option == "View Resources":
    st.title("Resources")
    collection = st.text_input("Collection Name")
    if st.button("View"):
        response = send_request(f'readall?collection={collection}')
        if isinstance(response, str):
            st.error(response)
        else:
            st.json(response)

# Add Resource
elif option == "Add Resource":
    st.title("Add Resource")
    collection = st.text_input("Collection Name")
    resource = st.text_input("Resource Name")
    data = st.text_area("Resource Data (JSON format)")
    if st.button("Add"):
        try:
            data = json.loads(data)
            response = send_request(f'write?collection={collection}&resource={resource}', method='POST', data=data)
            if isinstance(response, str):
                st.error(response)
            else:
                st.success("Resource added successfully!")
        except json.JSONDecodeError:
            st.error("Invalid JSON format")

# Search
elif option == "Search":
    st.title("Search Resources")
    collection = st.text_input("Collection Name")
    query = st.text_area("Query (JSON format)")
    if st.button("Search"):
        try:
            query = json.loads(query)
            response = send_request(f'search?collection={collection}', method='POST', data=query)
            if isinstance(response, str):
                st.error(response)
            else:
                st.json(response)
        except json.JSONDecodeError:
            st.error("Invalid JSON format")

# Regex Search
elif option == "Regex Search":
    st.title("Regex Search")
    collection = st.text_input("Collection Name")
    query = st.text_area("Query (JSON format)")
    if st.button("Search"):
        try:
            query = json.loads(query)
            response = send_request(f'regexsearch?collection={collection}', method='POST', data=query)
            if isinstance(response, str):
                st.error(response)
            else:
                st.json(response)
        except json.JSONDecodeError:
            st.error("Invalid JSON format")

# Delete Resource
elif option == "Delete Resource":
    st.title("Delete Resource")
    collection = st.text_input("Collection Name")
    resource = st.text_input("Resource Name")
    if st.button("Delete"):
        response = send_request(f'delete?collection={collection}&resource={resource}', method='POST')
        if isinstance(response, str):
            st.error(response)
        else:
            st.success("Resource deleted successfully!")
