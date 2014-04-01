var conn;
var graphConn;
var nickname;


function hidePoll() {
	var graph = document.getElementById("poll-results");
	graph.className = "hidden";
}

function showPoll(response) {
    // send response to server based on selected option
    var pollNumber = document.getElementById("pollNumber").innerText;
    var request = new XMLHttpRequest();
    request.open("POST", "http://pollchat.co.uk:3000/poll/"+pollNumber+"/"+response, true);
    request.send();
	var graph = document.getElementById("poll-results");
	graph.classList.remove("hidden");
}

function setNickname() {
    var nicknameEntry = document.getElementById("nickname-entry");
    if (nicknameEntry.value == ""){
        return;
    }
    nickname = nicknameEntry.value;
    var nicknamediv = document.getElementById("nickname-input");
    nicknamediv.className = "hidden";

    var chat = document.getElementById("chat-input");
    chat.classList.remove("hidden");
}

function sendMessage() {
    // send message via websocket
	var chatMessage = document.getElementById("chat-entry");
    if (chatMessage.value == "" || !conn){
        return;
    }
    var comment = {
        "name": nickname,
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
    conn = new WebSocket("ws://pollchat.co.uk:3000/data/"+pollNumber);
    graphConn = new WebSocket("ws://pollchat.co.uk:3000/graph/"+pollNumber);

    conn.onmessage = function(evt){
        appendMessageToChat(evt.data);
    }

    conn.onclose = function(evt) {
        console.log("chat connection closed")
    }

    graphConn.onmessage = function(evt){
        // update graph
        var res = JSON.parse(evt.data);
        updateGraph(res.responses);
    }

    graphConn.onclose = function(evt) {
        console.log("graph connection closed")
    }

    drawGraph();
}

function drawGraph(){
    var margin = {top: 20, right: 20, bottom: 30, left: 40},
        width = 500 - margin.left - margin.right,
        height = 500 - margin.top - margin.bottom;

    var color = d3.scale.ordinal()
      .range(["#98abc5", "#8a89a6", "#7b6888", "#6b486b", "#a05d56", "#d0743c", "#ff8c00"]);

    var x = d3.scale.ordinal()
        .rangeRoundBands([0, width], .1);

    var y = d3.scale.linear()
        .range([height, 0]);

    var xAxis = d3.svg.axis()
        .scale(x)
        .orient("bottom");

    var yAxis = d3.svg.axis()
        .scale(y)
        .orient("left")
        .ticks(5);

    var svg = d3.select("#poll-results").append("svg")
        .attr("class", "d3graph")
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
      .append("g")
        .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

    var responses = [];
        responses[0] = {"response":"", "count":0};
        responses[1] = {"response":"", "count":0};
        responses[2] = {"response":"", "count":0};
        responses[3] = {"response":"", "count":0};

      x.domain(responses.map(function(d) { return d.response; }));
      y.domain([0, d3.max(responses, function(d) { return d.count; })]);

      svg.append("g")
          .attr("class", "x-axis")
          .attr("transform", "translate(0," + height + ")")
          .call(xAxis);

      svg.append("g")
          .attr("class", "y-axis")
          .call(yAxis)
        .append("text")
          .attr("transform", "rotate(-90)")
          .attr("y", 3)
          .attr("dy", ".71em")
          .style("text-anchor", "end")
          .text("Votes");

      svg.selectAll(".bar")
          .data(responses)
        .enter().append("rect")
          .attr("class", "bar")
          .attr("fill", function(d) { return color(d.response) })
          .attr("x", function(d) { return x(d.response); })
          .attr("width", x.rangeBand())
          .attr("y", function(d) { return y(d.count); })
          .attr("height", function(d) { return height - y(d.count); });
}

function updateGraph(res) {
    var margin = {top: 20, right: 20, bottom: 30, left: 40},
        width = 500 - margin.left - margin.right,
        height = 500 - margin.top - margin.bottom;

    var color = d3.scale.ordinal()
      .range(["#98abc5", "#8a89a6", "#7b6888", "#6b486b", "#a05d56", "#d0743c", "#ff8c00"]);

    var x = d3.scale.ordinal()
        .rangeRoundBands([0, width], .1);

    var y = d3.scale.linear()
        .range([height, 0]);

    var xAxis = d3.svg.axis()
        .scale(x)
        .orient("bottom");

    var yAxis = d3.svg.axis()
        .scale(y)
        .orient("left")
        .ticks(5, "votes");

    var responses = [];
        responses[0] = res["1"];
        responses[1] = res["2"];
        responses[2] = res["3"];
        responses[3] = res["4"];

    var data = res;
      x.domain(responses.map(function(d) { return d.response; }));
      y.domain([0, d3.max(responses, function(d) { return d.count; })]);

    var svg = d3.selectAll(".d3graph");

    svg.selectAll("g.y-axis")
        .call(yAxis);

    svg.selectAll("rect")
          .data(responses)
        .transition()
          .duration(1000)
          .attr("class", "bar")
          .attr("fill", function(d) { return color(d.response) })
          .attr("x", function(d) { return x(d.response); })
          .attr("width", x.rangeBand())
          .attr("y", function(d) { return y(d.count); })
          .attr("height", function(d) { return height - y(d.count); });
}

function validatePollEntry(){
  // check there are at least two responses given
  // check a question has been supplied
  // trim white space when checking values
  var question = document.getElementById("pollquestion");
  console.log("question = " + question.value);
  if (document.getElementById("pollquestion").value.trim() == ""){
    return false;
  }

  return true;
}
