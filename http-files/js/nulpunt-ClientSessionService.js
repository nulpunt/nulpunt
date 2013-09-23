

nulpunt.run(function(ClientSessionService) {
	ClientSessionService.startSession();
});

nulpunt.factory('ClientSessionService', function($rootScope, $http, $sessionStorage) {
	var service = {
		sessionKey: ""
	};

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
				if(data.result) {
					console.log("got sessionKey from browser sessionStorage");
					setKey(sessionKey);
				} else {
					initSession();
				}
			}).
			error(function(data, status, headers, config) {
				console.error(status, data);
				initSession();
			});
		} else {
			initSession();
		}
	}

	// stop the session: destroy on server, destroy locally
	service.stopSession = function() {
		console.log($http.defaults.headers);
		$http({method: 'GET', url: '/service/sessionDestroy'}).
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