

nulpunt.run(function(ClientSessionService) {
	ClientSessionService.startSession();
});

nulpunt.factory('ClientSessionService', function($rootScope, $http, $sessionStorage, $timeout) {
	var service = {
		sessionKey: ""
	};

	function sendPing() {
		$http({method: 'GET', url:"/service/session/ping"}).
		success(function(data, status, headers, config) {
			//++ ok !
		}).
		error(function(data, status, headers, config) {
			//++ error !
		});
	}
	
	function timeoutPing() {
		sendPing();
		$timeout(timeoutPing, 3*60*1000);
	}
	timeoutPing();

	// setKey updates the key on multiple locations (either a valid key, or empty string for unset)
	function setKey(newKey) {
		service.sessionKey = newKey;
		$http.defaults.headers.common['X-Nulpunt-SessionKey'] = newKey;
		$sessionStorage.sessionKey = newKey;
	}

	// init a new session
	function initSession() {
		$http({method: 'GET', url: '/service/sessionInit'}).
		success(function(data, status, headers, config) {
			setKey(data.sessionKey);
		}).
		error(function(data, status, headers, config) {
			console.error(status, data);
		});
	}

	// start session (try to continiue existing session)
	service.startSession = function() {
		// get key from session storage, check if it is valid
		sessionKey = $sessionStorage.sessionKey;
		if(sessionKey!=undefined && sessionKey.length>0) {
			$http({method: 'POST', url: '/service/sessionCheck', data: {sessionKey: sessionKey}}).
			success(function(data, status, headers, config) {
				console.log("got sessionKey from browser sessionStorage ...");
				if(data.valid) {
					console.log("... and guess what!? It's still valid! :D");
					setKey(sessionKey);
				} else {
					console.log("... but that sessionKey was invalid");
					initSession();
				}
			}).
			error(function(data, status, headers, config) {
				console.log("invalid response")
				console.log(data, status, headers);
				initSession();
			});
		} else {
			initSession();
		}
	}

	// stop the session: destroy on server, destroy locally
	service.stopSession = function() {
		console.log($http.defaults.headers);
		$http({method: 'GET', url: '/service/session/destroy'}).
		success(function(data, status, headers, config) {
			console.log(status, data);
			setKey("");
		}).
		error(function(data, status, headers, config) {
			console.error(status, data);
			setKey("");
		});
	};

	//++ add websocket (?)
	//++ add timed pings

	return service;
});