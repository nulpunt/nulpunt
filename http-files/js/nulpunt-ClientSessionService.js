nulpunt.run(function($http, ClientSessionService) {
	$http({method: 'GET', url: '/service/init'}).
	success(function(data, status, headers, config) {
		ClientSessionService.setSessionKey(data.sessionKey);
	}).
	error(function(data, status, headers, config) {
		console.error(status, data);
	});
})

nulpunt.factory('ClientSessionService', function($rootScope, $http) {
	var service = {
		sessionKey: ""
	};

	service.setSessionKey = function(newKey) {
		service.sessionKey = newKey;
		$http.defaults.headers.common['X-nulpunt-sessionKey'] = newKey;
	};

	//++ add websocket (?)
	//++ add timed pings

	return service;
})