var c=new WebSocket('ws://localhost:3000/sock');
c.onopen = function(){
	c.onmessage = function(response){
		console.log(response.data);
		var newLog = $('<li>').text(response.data);
		$('#logs').append(newLog);
		$('#logs').val('');
	};
}

