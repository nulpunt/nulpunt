var nulpunt = angular.module('nulpunt', [
	// imported modules
	// please keep this list sorted
	"ngRoute",
	'ui.bootstrap.collapse', 
	'ui.bootstrap.dropdownToggle'
]);

nulpunt.config(function($routeProvider) {
	$routeProvider.when('/', {
		template: 'Welcome home',
		controller: "HomeCtrl"
	})
	.when('/account', {
		template: 'Account info',
		controller: "AccountCtrl"
	})
	.when('/not-found', {
		template: "Page not found",
	})
	.otherwise({
		redirectTo: 'not-found' //++ write not-found template+controller
	});
});

nulpunt.controller("HomeCtrl", function(){
	//++
});

nulpunt.controller("AccountCtrl", function(){
	//++
});

nulpunt.controller("NavCtrl", function($scope) {
	$scope.loc = 'home';

	$scope.gravatarHash = CryptoJS.MD5("gjr19912@gmail.com").toString(CryptoJS.enc.Hex);
});

nulpunt.controller("MainCtrl", function($scope) {
	//++
});