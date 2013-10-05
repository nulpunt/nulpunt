

nulpunt.run(function() {
	//++ check if there is already authed for given ClientSession
});

nulpunt.factory('AccountAuthService', function($rootScope, $http, $q) {
	var emptyAuth = {
		username: "",
		email: ""
	};

	var service = {
		account: emptyAuth
	};

	$rootScope.account = service.account;

	service.getUsername = function() {
		if(service.account.username == undefined) {
			return "";
		}
		return service.account.username;
	}

	service.authenticate = function(username, password) {
		var defered = $q.defer();

		$http({method: 'POST', url: '/service/session/authenticateAccount', data: {username: username, password: password}}).
		success(function(data, status, headers, config) {
			if(data.success) {
				//++ retrieve account details from server
				service.account = {
					username: username,
					email: "gjr19912@gmail.com"
				};
				$rootScope.$broadcast("auth_changed");

				// all done
				defered.resolve();
			} else {
				console.log(data.error);
				defered.reject(data.error);
			}
		}).
		error(function(data, status, headers, config) {
			defered.reject("Request error.");
		});

		return defered.promise;
	};

	service.unAuthenticate = function() {
		service.account = emptyAuth;
		$rootScope.$broadcast('auth_changed');
	}

	return service;
});