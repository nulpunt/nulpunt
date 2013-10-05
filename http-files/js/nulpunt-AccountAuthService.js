

nulpunt.run(function() {
	//++ check if there is already authed for given ClientSession
});

nulpunt.factory('AccountAuthService', function($rootScope, $http) {
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

	// returns "ok" on valid auth
	// returns "" on invalid auth
	// returns error on error
	service.authenticate = function(username, password) {
		//++ create a promise

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
				return "ok"; //++ accept on earlier returned promise
			} else {
				console.log(data.error);
				return data.error; //++ accept on earlier returned promise
			}
		}).
		error(function(data, status, headers, config) {
			return "Request error." //++ accept on earlier returned promise
		});

		//++ return the promise
	};

	service.unAuthenticate = function() {
		service.account = emptyAuth;
		$rootScope.$broadcast('auth_changed');
	}

	return service;
});