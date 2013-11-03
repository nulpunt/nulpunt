var nulpunt = angular.module('nulpunt', [
	// imported modules
	// please keep this list sorted
	'ngRoute',
	'ngStorage',
	'ui.bootstrap.collapse', 
	'ui.bootstrap.dropdownToggle',
	'angularFileUpload'
]);

nulpunt.config(function($routeProvider) {
	$routeProvider
	.when('/', {
		templateUrl: "/html/overview.html",
		controller: "OverviewCtrl"
	})
	.when('/inbox', {
		templateUrl: "/html/inbox.html",
		controller: "InboxCtrl"
	})
	.when('/register', {
		templateUrl: "/html/register.html",
		controller: "RegisterCtrl"
	})
	.when('/history', {
		templateUrl: "/html/history.html",
		controller: "HistoryCtrl"
	})
	.when('/sign-in', {
		templateUrl: "/html/sign-in.html",
		controller: "SignInCtrl"
	})
	.when('/sign-out', {
		templateUrl: "/html/sign-out.html",
		controller: "SignOutCtrl"
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
	.when('/admin/upload', {
		templateUrl: "/html/admin-upload.html",
		controller: "AdminUploadCtrl"
	})
	.when('/admin/analyse', {
		templateUrl: "/html/admin-analyse.html",
		controller: "AdminAnalyseCtrl"
	})
	.when('/admin/process', {
		templateUrl: "/html/admin-process.html",
		controller: "AdminProcessCtrl"
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
	$rootScope.$on("auth_changed", function() {
		$scope.account = AccountAuthService.account;
		$scope.gravatarHash = CryptoJS.MD5(AccountAuthService.account.email).toString(CryptoJS.enc.Hex);
	});

	$scope.search = function() {
		var safeSearchValue = $scope.searchValue.replace(/[\/\? ]/g, '+').replace('++', '+').trim('+');
		$location.path("search/"+safeSearchValue);
	};
});

nulpunt.controller("TabbarCtrl", function($scope, $location) {
	$scope.getClass = function(path) {
    	if ($location.path().substr(0, path.length) == path) {
	      return "active"
	    } else {
	      return ""
	    }
	}
});

nulpunt.controller("OverviewCtrl", function($scope){
	//++
});

nulpunt.controller("InboxCtrl", function($scope) {
	$scope.inbox = {
		items: [],
	};

	$scope.inbox.items = [
		{name: "bla", summary: "fdsa"},
	];
});

nulpunt.controller("HistoryCtrl", function($scope) {
	//++
});

nulpunt.controller("SearchCtrl", function($scope, $routeParams) {
	$scope.mySearch = $routeParams.searchValue.replace(/[+]/g, ' ');
});

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
	settingsPromise.then(
	function(data) { // success
		$scope.settings.testA = data.a;
		$scope.settings.testB = data.b;
	}, 
	function(error) { // error
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

nulpunt.controller("AdminUploadCtrl", function($scope, $upload) {
	$scope.uploading = false;

	$scope.onFileSelect = function($files) {
		$scope.files = [];
		_.each($files, function(file, index) {
			$scope.files.push({
				file: file,
				i: index,
				percentage: 1,
			});
		});
	};
	
	$scope.removeFile = function(index) {
		$scope.files.splice(index, 1);
	};

	$scope.uploadFiles = function() {
		$scope.uploading = true;
		_.each($scope.files, function(file, index) {
			$upload.upload({
				url: 'service/session/admin/upload',
				// headers: {'X-Nulpunt-SessionKey': 'headerValue'},
				// withCredential: true,
				data: {/*aditional data*/},
				file: file.file,
				//fileFormDataName: myFile, //(optional) sets 'Content-Desposition' formData name for file
				progress: function(evt) {
					//++ TODO: this isn't executed
					var percentage = parseInt(100.0 * evt.loaded / evt.total);
					console.log('progress: '+index+': '+percentage);
					$scope.files[index].percentage = percentage;
					$scope.$apply(); //++ is this required?
				}
			})
			.success(function(data, status, headers, config) {
				$scope.files[index].percentage = 100;
				console.log(data);
			})
			.error(function(data, status, headers, config) {
				console.log("error uploading", data);
			})
		})
	};
});

nulpunt.controller("AdminAnalyseCtrl", function($scope, $http) {
	$scope.files = [];
	$http({method: "POST", url: "/service/session/admin/getRawUploads"}).
	success(function(data) {
		console.dir(data);
		$scope.files = data.files;
	}).
	error(function(error) {
		console.log('error retrieving raw documents: ', error);
	})
});

nulpunt.controller("SignOutCtrl", function($scope, AccountAuthService, ClientSessionService) {
	$scope.username = AccountAuthService.getUsername();
	ClientSessionService.stopSession();
});

nulpunt.filter('bytes', function() {
	return function(bytes, precision) {
		if (bytes==0 || isNaN(parseFloat(bytes)) || !isFinite(bytes)) return '-';
		if (typeof precision === 'undefined') precision = 1;
		var units = ['bytes', 'kB', 'MB', 'GB', 'TB', 'PB'],
		number = Math.floor(Math.log(bytes) / Math.log(1024));
		return (bytes / Math.pow(1024, Math.floor(number))).toFixed(precision) + ' ' + units[number];
	}
});