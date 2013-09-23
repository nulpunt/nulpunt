
nulpunt.run(function() {
	//++ check if there is already authed for given ClientSession
});

nulpunt.factory('AccountAuthService', function($rootScope) {
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
		service.account = {
			username: username,
			email: "gjr19912@gmail.com"
		};
		$rootScope.$broadcast("auth_changed");
	};

	service.unAuthenticate = function() {
		service.account = emptyAuth;
		$rootScope.$broadcast('auth_changed');
	}

	return service;
});