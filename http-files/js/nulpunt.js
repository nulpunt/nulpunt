var nulpunt = angular.module('nulpunt', [
	// imported modules
	// please keep this list sorted
	'ngRoute',
	'ngStorage',
	'ui.bootstrap.collapse', 
	'ui.bootstrap.dropdownToggle'
]);

nulpunt.config(function($routeProvider) {
	$routeProvider
	.when('/', {
		template: 'Welcome home',
		controller: "HomeCtrl" //++ rename to Overview?
	})
	.when('/register', {
		templateUrl: "/html/register.html",
		controller: "RegisterCtrl"
	})
	.when('/sign-in', {
		templateUrl: "/html/sign-in.html",
		controller: "SignInCtrl"
	})
	.when('/sign-out', {
		templateUrl: "/html/sign-out.html",
		controller: "SignOutCtrl"
	})
	.when('/topics', {
		templateUrl: "/html/topics.html",
		controller: "TopicsCtrl"
	})
	.when('/search/:searchValue', {
		templateUrl: '/html/search.html',
		controller: "SearchCtrl"
	})
	.when('/profile/:userID', {
		templateUrl: '/html/profile.html',
		controller: "ProfileCtrl"
	})
	.when('/settings', {
		templateUrl: '/html/settings.html',
		controller: "SettingsCtrl"
	})
	.otherwise({
		templateUrl: "/html/not-found.html",
		controller: "NotFoundCtrl",
	});
});

nulpunt.controller("MainCtrl", function($scope) {
	//++
});

nulpunt.controller("NavbarCtrl", function($scope, $rootScope, $location, AccountAuthService) {
	$scope.loc = 'home';

	$rootScope.$on("auth_changed", function() {
		$scope.account = AccountAuthService.account;
		$scope.gravatarHash = CryptoJS.MD5(AccountAuthService.account.email).toString(CryptoJS.enc.Hex);
	});

	$scope.search = function() {
		var safeSearchValue = $scope.searchValue.replace(/[\/\? ]/g, '+').replace('++', '+').trim('+');
		$location.path("search/"+safeSearchValue);
	};
});

nulpunt.controller("HomeCtrl", function($scope){
	//++
});

nulpunt.controller("TopicsCtrl", function($scope) {
	//++
});

nulpunt.controller("SearchCtrl", function($scope, $routeParams) {
	$scope.mySearch = $routeParams.searchValue.replace(/[+]/g, ' ');
})

nulpunt.controller("ProfileCtrl", function() {
	//++
});

nulpunt.controller('NotFoundCtrl', function($scope, $location) {
	$scope.path = $location.url()
});

nulpunt.controller("RegisterCtrl", function($scope, $rootScope, $http) {
	$scope.submit = function() {
		$http({method: 'POST', url: '/service/session/registerAccount', data: {
			username: $scope.username,
			email: $scope.email,
			password: $scope.password
		}}).
		success(function(data, status, headers, config) {
			if(data.success) {
				$scope.done = true
			} else {
				$scope.error = data.error;
			}
		}).
		error(function(data, status, headers, config) {
			console.log("invalid response for registerAccount")
			console.log(data, status, headers);
			$scope.error = "Could not make request to server";
		});
	};
});

nulpunt.controller("SettingsCtrl", function($scope, AccountDataService) {
	// defaults
	$scope.settings = {
		testA: "emptyA",
		testB: "emptyB",
	};

	// get settings from server
	var settingsPromise = AccountDataService.getObject("settings");
	settingsPromise.then(function(data) {
		$scope.settings.testA = data.a;
		$scope.settings.testB = data.b;
	}, function(error) {
		console.error(error);
	});

	// saveSettings function
	$scope.saveSettings = function() {
		var data = {
			a: $scope.settings.testA,
			b: $scope.settings.testB,
		};
		var donePromise = AccountDataService.setObject("settings", data);
		donePromise.then(function() {
			console.log('saved');
		}, function(error) {
			console.error(error);
		})
	}
});

nulpunt.controller("SignInCtrl", function($scope, $rootScope, AccountAuthService) {
	$scope.submit = function() {
		$scope.success = false;
		$scope.wrong = false;
		$scope.error = "";
		var prom = AccountAuthService.authenticate($scope.username, $scope.password);

		prom.then(function() {
				$scope.success = true;
			}, function(error) {
				if(error == "") {
					// no success, but also no error: credentials are wrong.
					$scope.wrong = true;
				} else {
					$scope.error = result;
				}
				//++ need to do some "digest" on $scope ?? or $scope.$apply()?
				//++ find out what good convention is
			}
		);
	};
	
	$rootScope.$on("auth_changed", function() {
		$scope.account = AccountAuthService.account;
	});
});

nulpunt.controller("SignOutCtrl", function($scope, AccountAuthService, ClientSessionService) {
	$scope.username = AccountAuthService.getUsername();
	ClientSessionService.stopSession();
});