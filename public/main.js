function hidePoll() {
	var graph = document.getElementById("poll-results");
	graph.className = "hidden";
}

function showPoll() {
	var graph = document.getElementById("poll-results");
	graph.classList.remove("hidden");
}

function sendMessage() {
	var chatMessage = document.getElementById("chat-entry");
    if (chatMessage.value == ""){
        return
    }
	var chatList = document.getElementById("chat");
	
	var chatComment = document.createElement('div');
    chatComment.className = "chat-comment";
    
    var userName = document.createElement('span');
    userName.className = "chat-username";
    
    var userComment = document.createElement('span');
    userComment.className = "chat-text";
    
    userName.innerText = "Tom : ";
    userComment.innerText = chatMessage.value;
    
    reOrderChat();

    chatComment.appendChild(userName);
    chatComment.appendChild(userComment);
    chatList.appendChild(chatComment);

    chatMessage.value = "";

    reOrderChat();
}

function reOrderChat() {
    var chatList = document.getElementById("chat");
    var divs = chatList.children, i = divs.length - 1;
    for (; i--;){
        chatList.appendChild(divs[i]);
    }
}

function createPoll() {
    
}
