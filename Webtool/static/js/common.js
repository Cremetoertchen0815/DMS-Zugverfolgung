//HTTP Status Codes
var STATUS_OK = 200; //The request was successful
var STATUS_CREATED = 201; //The request was successful and a new resource was created
var STATUS_BAD_REQUEST = 400; //The request parameters were invalid
var STATUS_UNAUTHORIZED = 401; //The action requires authentication/session expired
var STATUS_FORBIDDEN = 403; //The action is not allowed
var STATUS_NOT_FOUND = 404; //The requested resource was not found
var STATUS_CONFLICT = 409; //The to be created resource already exists or the resource to be deleted is still in use
var STATUS_INTERNAL_SERVER_ERROR = 500; //An internal server error occurred

function createAPICall(type, url, callback) {
    let xhr = new XMLHttpRequest();
    xhr.open(type, url);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onload = () => {
        //Check if session expired
        if (xhr.status === STATUS_UNAUTHORIZED) {
            alert("Session expired!");
            window.location.href = "./login.html";
            return;
        }

        //If session is still valid, call callback
        callback(xhr.status, xhr.responseText);
    };
    xhr.send();
}

function createAPICallWithBody(type, url, body, callback) {
    let xhr = new XMLHttpRequest();
    xhr.open(type, url);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onload = () => {
        //Check if session expired
        if (xhr.status === STATUS_UNAUTHORIZED) {
            alert("Session expired!");
            window.location.href = "./login.html";
            return;
        }

        //If session is still valid, call callback
        callback(xhr.status, xhr.responseText);
    };
    xhr.send(body);
}