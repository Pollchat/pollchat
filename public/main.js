var conn;
var graphConn;


function hidePoll() {
	var graph = document.getElementById("poll-results");
	graph.className = "hidden";
}

function showPoll(response) {
    // send response to server based on selected option
    var pollNumber = document.getElementById("pollNumber").innerText;
    var request = new XMLHttpRequest();
    request.open("POST", "http://localhost:3000/poll/"+pollNumber+"/"+response, true);
    request.send();
	var graph = document.getElementById("poll-results");
	graph.classList.remove("hidden");
}

function sendMessage() {
    // send message via websocket
	var chatMessage = document.getElementById("chat-entry");
    if (chatMessage.value == "" || !conn){
        return;
    }
    var comment = {
        "name": "tom",
        "message": chatMessage.value
    }
    conn.send(JSON.stringify(comment));

	chatMessage.value = "";
}

function appendMessageToChat(message) {
    var chatList = document.getElementById("chat");
    
    var chatComment = document.createElement('div');
    chatComment.className = "chat-comment";
    
    var userName = document.createElement('span');
    userName.className = "chat-username";
    
    var userComment = document.createElement('span');
    userComment.className = "chat-text";

    var comment = JSON.parse(message);
    
    userName.innerText = comment.name;
    userComment.innerText = comment.message;
    
    reOrderChat();

    chatComment.appendChild(userName);
    chatComment.appendChild(userComment);
    chatList.appendChild(chatComment);

    reOrderChat();
}

function reOrderChat() {
    var chatList = document.getElementById("chat");
    if (chatList.children.length == 0){
        return;
    }
    var divs = chatList.children, i = divs.length - 1;
    for (; i--;){
        chatList.appendChild(divs[i]);
    }
}

function onPollLoad() {
    var pollNumber = document.getElementById("pollNumber").innerText;
    conn = new WebSocket("ws://localhost:3000/data/"+pollNumber);
    graphConn = new WebSocket("ws://localhost:3000/graph/"+pollNumber);

    conn.onmessage = function(evt){
        appendMessageToChat(evt.data);
    }

    conn.onclose = function(evt) {
        console.log("chat connection closed")
    }

    graphConn.onmessage = function(evt){
        // update graph
        console.log("graph update");
        console.log(JSON.parse(evt.data));
    }

    graphConn.onclose = function(evt) {
        console.log("graph connection closed")
    }
}

function createPoll() {
    
}
