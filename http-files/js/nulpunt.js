var nulpunt = angular.module('nulpunt', [
	// imported modules
	// please keep this list sorted
	"ngRoute",
	'ui.bootstrap.collapse', 
	'ui.bootstrap.dropdownToggle'
]);

nulpunt.config(function($routeProvider) {
	$routeProvider
	.when('/', {
		template: 'Welcome home',
		controller: "HomeCtrl" //++ rename to Overview?
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
	
});

nulpunt.controller('NotFoundCtrl', function($scope, $location) {
	$scope.path = $location.url()
});

nulpunt.controller("SignInCtrl", function($scope, $rootScope, AccountAuthService) {
	$scope.submit = function() {
		AccountAuthService.authenticate($scope.username, $scope.password);
	};
	
	$rootScope.$on("auth_changed", function() {
		$scope.account = AccountAuthService.account;
	});
});

nulpunt.controller("SignOutCtrl", function($scope, AccountAuthService) {
	$scope.username = AccountAuthService.getUsername();
	AccountAuthService.unAuthenticate();
});