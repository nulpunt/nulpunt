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
	.when('/history/:sortBy', {
		templateUrl: "/html/history.html",
		controller: "HistoryCtrl"
	})
	.when('/trending', {
		templateUrl: "/html/trending.html",
		controller: "TrendingCtrl"
	})
	.when('/notifications', {
		templateUrl: "/html/notifications.html",
		controller: "NotificationsCtrl"
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
	.when('/admin/process-editmeta/:docID', {
		templateUrl: "/html/admin-process-editmeta.html",
		controller: "AdminProcessEditMetaCtrl"
	})
	.when('/admin/tags', {
		templateUrl: "/html/admin-tags.html",
		controller: "AdminTagsCtrl"
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
	$scope.documents = {
		items: [],
	};

	$scope.documents.items = [
		{title: "Title of the document", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 6, nrOfDrafts: 2, nrOfComments: 8, nrOfBookmarks: 4, tags: [{title: "Iraq"}, {title:"Conspiracy"}, {title:"Another tag"}] },
		{title: "Title of the document", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 32, nrOfDrafts: 7, nrOfComments: 18, nrOfBookmarks: 12, tags: [{title: "Random tag"}] },
		{title: "Title of the document", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 2, nrOfDrafts: 14, nrOfComments: 25, nrOfBookmarks: 4, tags: [{title: "Iraq"}] },
		{title: "Title of the document", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 10, nrOfDrafts: 55, nrOfComments: 3, nrOfBookmarks: 15, tags: [] },
	];
});

nulpunt.controller("HistoryCtrl", function($scope, $routeParams) {
	$scope.documents = {
		items: [],
	};

	if($routeParams.sortBy == "annotations") {
		//Sort by Annotations
		$scope.documents.items = [
			{title: "Title of the document 4", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 32, nrOfDrafts: 7, nrOfComments: 18, nrOfBookmarks: 12, tags: [{title: "Random tag"}] },
			{title: "Title of the document 3", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 10, nrOfDrafts: 55, nrOfComments: 3, nrOfBookmarks: 15, tags: [] },
			{title: "Title of the document 2", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 6, nrOfDrafts: 2, nrOfComments: 8, nrOfBookmarks: 4, tags: [{title: "Iraq"}, {title:"Conspiracy"}, {title:"Another tag"}] },
			{title: "Title of the document 1", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 2, nrOfDrafts: 14, nrOfComments: 25, nrOfBookmarks: 4, tags: [{title: "Iraq"}] },
		];
	}
	else if($routeParams.sortBy == "drafts") {
		$scope.documents.items = [
			{title: "Title of the document 3", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 10, nrOfDrafts: 55, nrOfComments: 3, nrOfBookmarks: 15, tags: [] },
			{title: "Title of the document 1", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 2, nrOfDrafts: 14, nrOfComments: 25, nrOfBookmarks: 4, tags: [{title: "Iraq"}] },
			{title: "Title of the document 4", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 32, nrOfDrafts: 7, nrOfComments: 18, nrOfBookmarks: 12, tags: [{title: "Random tag"}] },
			{title: "Title of the document 2", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 6, nrOfDrafts: 2, nrOfComments: 8, nrOfBookmarks: 4, tags: [{title: "Iraq"}, {title:"Conspiracy"}, {title:"Another tag"}] },
		];
	}
	else if($routeParams.sortBy == "comments") {
		$scope.documents.items = [
			{title: "Title of the document 1", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 2, nrOfDrafts: 14, nrOfComments: 25, nrOfBookmarks: 4, tags: [{title: "Iraq"}] },
			{title: "Title of the document 4", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 32, nrOfDrafts: 7, nrOfComments: 18, nrOfBookmarks: 12, tags: [{title: "Random tag"}] },
			{title: "Title of the document 2", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 6, nrOfDrafts: 2, nrOfComments: 8, nrOfBookmarks: 4, tags: [{title: "Iraq"}, {title:"Conspiracy"}, {title:"Another tag"}] },
			{title: "Title of the document 3", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 10, nrOfDrafts: 55, nrOfComments: 3, nrOfBookmarks: 15, tags: [] },
		];
	}
	else if($routeParams.sortBy == "bookmarks") {
		$scope.documents.items = [
			{title: "Title of the document 3", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 10, nrOfDrafts: 55, nrOfComments: 3, nrOfBookmarks: 15, tags: [] },
			{title: "Title of the document 4", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 32, nrOfDrafts: 7, nrOfComments: 18, nrOfBookmarks: 12, tags: [{title: "Random tag"}] },
			{title: "Title of the document 1", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 2, nrOfDrafts: 14, nrOfComments: 25, nrOfBookmarks: 4, tags: [{title: "Iraq"}] },
			{title: "Title of the document 2", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 6, nrOfDrafts: 2, nrOfComments: 8, nrOfBookmarks: 4, tags: [{title: "Iraq"}, {title:"Conspiracy"}, {title:"Another tag"}] },
		];
	}
});

nulpunt.controller("TrendingCtrl", function($scope) {
	$scope.documents = {
		items: [],
	};

	$scope.documents.items = [
		{title: "Title of the document", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 6, nrOfDrafts: 2, nrOfComments: 8, nrOfBookmarks: 4, tags: [{title: "Iraq"}, {title:"Conspiracy"}, {title:"Another tag"}] },
		{title: "Title of the document", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 32, nrOfDrafts: 7, nrOfComments: 18, nrOfBookmarks: 12, tags: [{title: "Random tag"}] },
		{title: "Title of the document", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 2, nrOfDrafts: 14, nrOfComments: 25, nrOfBookmarks: 4, tags: [{title: "Iraq"}] },
		{title: "Title of the document", description: "A short 1 or 2 sentence description of the document. Include or not?", source: "The Government", sourceDate: "01/01/2004", uploadDate: "01/11/2013", uploader: "Nulpunt", nrOfAnnotations: 10, nrOfDrafts: 55, nrOfComments: 3, nrOfBookmarks: 15, tags: [] },
	];
});

nulpunt.controller("NotificationsCtrl", function($scope) {
	$scope.notifications = {
		items: [],
	};

	$scope.notifications.items = [
		{type: "comment", time: "1 minute ago", user: "Jonas", userId: "jonas", documentTitle: "An awesome article", documentId: 1},
		{type: "annotation", time: "2 minutes ago", user: "Renée", userId: "renee", documentTitle: "Another document", documentId: 2},
		{type: "comment", time: "1 hour ago", user: "Jonas", userId: "jonas", documentTitle: "An awesome article", documentId: 1},
		{type: "annotation", time: "3 hours ago", user: "Renée", userId: "renee", documentTitle: "Another document", documentId: 2},
		{type: "comment", time: "yesterday", user: "Jonas", userId: "jonas", documentTitle: "An awesome article", documentId: 1},
		{type: "annotation", time: "28/10/ 2013 - 18:22", user: "Renée", userId: "renee", documentTitle: "Another document", documentId: 2},
	];
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

nulpunt.controller("AdminTagsCtrl", function($scope, $rootScope, $http) {
    $http.get('/service/session/admin/tags').
	success(function(data) {
	    $scope.tags = data;
	}).
	error(function(data, status, headers, config) {
	    console.log("error fetching tags");
	    console.log(data, status, headers);
	    $scope.error = data;
	});

    $scope.specify_add = function() {
	this.url = '/service/session/admin/tags';
    }
    
    $scope.specify_delete = function() {
	this.url = '/service/session/admin/tags/delete';
    }

    $scope.submit = function() {
	    $scope.done = false;
	    $scope.error = "";
	    
		$http({
		    method: 'POST', 
		    url: this.url,
		    data: { tag: $scope.tag } }).
		success(function(data, status, headers, config) {
		    // console.log(data)
		    // TODO: This doesn't update the list of available tags.
		    // We need to signal the model-viewer somehow.
		    $scope.tags = data
		}).
		error(function(data, status, headers, config) {
		    console.log("invalid response for add Tag");
			console.log(data, status, headers);
			$scope.error = data;
		});
	};
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

// THIS CONTROLLER IS AN UGLY HACK! 
// It copies uploaded-document data into new Document-record and a fake page-record. 
// Remove after the OCR-processing creates the document/pages records.
nulpunt.controller("AdminAnalyseCtrl", function($scope, $http) {
	$scope.files = [];
	$http({method: "POST", url: "/service/session/admin/getRawUploads"}).
	success(function(data) {
	    console.log(data);
		$scope.files = data.files;
	}).
	error(function(error) {
		console.log('error retrieving raw documents: ', error);
	})

    // create a new (unpublished) document to make testing document/pages possible
    $scope.analyse = function(ind) {
	$http({method: "POST", url: "/service/session/admin/insertDocument",
	       data: { 
		   //document: {
		       title:                   $scope.files[ind].filename,
		       uploader:          $scope.files[ind].uploaderUsername,
		       uploadedDate: $scope.files[ind].uploadDate,
		   //}
	       }}).
	    success(function(data) {
		console.log('success updating document.');
		alert("succes");
	    }).
	error(function(error) {
		console.log('error updating document: ', error);
	});
    }});

nulpunt.controller("AdminProcessCtrl", function($scope, $http) {
    $scope.documents = [];
    $http({method: "POST", url: "/service/getDocumentList", data: {} }).
	success(function(data) {
	    console.log(data);
	    $scope.documents = data;
	}).
	error(function(error) {
		console.log('error retrieving raw documents: ', error);
	});
});


nulpunt.controller("AdminProcessEditMetaCtrl", function($scope, $http, $routeParams) {
    $scope.done = false;
    $scope.error = "";
    console.log("DocID is " + $routeParams.docID );
    // load the requested document
    $http({method: "POST", url: "/service/getDocument", data: { DocID: $routeParams.docID } }).
	success(function(data) {
	    console.log(data);
	    $scope.document = data.document;
	    // cheat to test: $scope.document.Categories = ["irak", "test"];
	}).
	error(function(error) {
		console.log('error retrieving document: ', error);
	})

    // Helper to check the right checkboxes // take out if not needed..
    $scope.checkTag = function(tag, list) {
	console.log("List is: " + list)
	console.log("Checking tag: " +tag)
	for (i = 0; i <  list.lenght; i++) {
	    if (list[i] == tag) { 
		return true
	    }
	}
	return false;
    };

    // save the updated document
    $scope.submit = function() {
	console.log("document to submit is "+ $scope.document)
	$scope.done = false;
	$scope.error = "";
	$http({
	    method: 'POST', 
	    url: "/service/session/admin/updateDocument",
	    data: $scope.document
	}).
	    success(function(data, status, headers, config) {
		console.log(data)
		$scope.done = true
	    }).
	    error(function(data, status, headers, config) {
		console.log("error add updateDocument");
		console.log(data, status, headers);
		$scope.error = data;
	    });
    }
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