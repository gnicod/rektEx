var c=new WebSocket('ws://localhost:3000/sock');
c.onopen = function(){
	c.onmessage = function(response){
		console.log(response.data);
		var newLog = JSON.parse(response.data)
		$('table > tbody > tr:first').before('<tr><td>'+newLog.Message+'</td><td>'+newLog.Ip+'</td><td></td></tr>');
		$('#logs').append(newLog);
		$('#logs').val('');
	};
}

