// not sure how I'm going to organize this yet. none of this code has even been run, let alone tested

// actual object:
socko = {
	start: function(endPoint, callback) {
		if(!window["WebSocket"]) {
			console.log("This browser does not support WebSockets.")
			return
		}

		// creates new socket
		var sock = new socket(endPoint)
		callback(sock) // let user register event handlers first

	}
}

function socket(endPoint) {
	var conn = new WebSocket(endpoint)

	conn.onopen = function(evt) {
		this.onOpen(evt)
	}

	conn.onclose = function(evt) {
		this.onClose(evt)
	}

	conn.onmessage = function(evt) {
		msg = JSON.parse(evt.data)
		console.log(msg)

		if(msg.type in this.events) {
			this.events[msg.type](msg.content)
		} else  {
			console.log("No handler is defined for this message type.")
		}
	}
}

socket.prototype.onOpen = function(fn) {
	this.onOpen = fn
}

socket.prototype.onClose = function(fn) {
	this.onClose = fn
}

socket.prototype.onEvent = function(eventName, fn) {
	this.events[eventName] = fn()
}