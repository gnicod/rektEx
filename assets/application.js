var c=new WebSocket('ws://localhost:3000/sock/ip');
c.onopen = function(){
    c.onmessage = function(response){
        console.log(response.data);
        var newLogs = JSON.parse(response.data);
        for (var i = 0; i < newLogs.length; i++) {
            var newLog = newLogs[i];
            console.log(newLog);
            $('table > tbody ').prepend('<tr><td>'+newLog.Message+'</td><td>'+newLog.Ip+'</td><td></td></tr>');
            $('#logs').append(newLog);
            $('#logs').val('');

        }

    };
}

