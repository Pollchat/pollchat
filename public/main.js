function hidePoll() {
	var graph = document.GetElementById("poll");

}

function sendMessage() {
	var chatMessage = document.getElementById("chat-entry");
	var chatList = document.getElementById("chat");
	
	var chatComment = document.createElement('div');
    chatComment.class = "chat-comment";
    
    var userName = document.createElement('span');
    userName.class = "chat-username";
    
    var userComment = document.createElement('span');
    userComment.class = "chat-text";
    
    userName.innerText = "Tom : ";
    userComment.innerText = chatMessage.value;
    
    chatComment.appendChild(userName);
    chatComment.appendChild(userComment);
    chatList.appendChild(chatComment);
}
