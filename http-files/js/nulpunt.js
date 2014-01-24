var nulpunt = angular.module('nulpunt', [
	// imported modules
	// please keep this list sorted
	'ngRoute',
	'ngStorage',
	'ui.bootstrap.collapse',
	'ui.bootstrap.dropdownToggle',
	'angularFileUpload',
	'checklist-model'
]);

nulpunt.config(function($routeProvider) {
	$routeProvider
	.when('/', {
		templateUrl: "/html/trending.html",
		controller: "TrendingCtrl"
	})
	.when('/documents', {
		templateUrl: "/html/dashboard.html",
		controller: "DashboardCtrl"
	})

	.when('/document/:docID', {
		templateUrl: "/html/document.html",
		controller: "DocumentCtrl"
	})
	.when('/register', {
		templateUrl: "/html/register.html",
		controller: "RegisterCtrl"
	})
	.when('/history/:sortBy', {
		templateUrl: "/html/history.html",
		controller: "HistoryCtrl"
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
	.when('/sign-out-successful', {
		templateUrl: "/html/sign-out.html",
	})
	.when('/search/:searchValue', {
		templateUrl: '/html/search.html',
		controller: "SearchCtrl"
	})
	.when('/profile', {
		templateUrl: '/html/profile.html',
		controller: "EmptyCtrl"
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
	.when('/about', {
		templateUrl: "/html/about.html",
		controller: "AboutCtrl"
	})
	.when('/contact', {
		templateUrl: "/html/contact.html",
		controller: "ContactCtrl"
	})
	.when('/colophon', {
		templateUrl: "/html/colophon.html",
		controller: "ColophonCtrl"
	})
	.otherwise({
		templateUrl: "/html/not-found.html",
		controller: "NotFoundCtrl",
	});
});

nulpunt.controller("EmptyCtrl", function() {
	// Empty controller can be used to when a template specifies the controllers in-line.
});

nulpunt.controller("MainCtrl", function($scope, $rootScope, AccountAuthService) {

	$rootScope.$on("auth_changed", function() {
		$scope.account = AccountAuthService.account;
	});
});

nulpunt.controller("NavbarCtrl", function($scope, $rootScope, $location, AccountAuthService) {

	$rootScope.$on("auth_changed", function() {
		$scope.account = AccountAuthService.account;
	});

	$scope.search = function() {
		var safeSearchValue = $scope.searchValue.replace(/[\/\? ]/g, '+').replace('++', '+').trim('+');
		$location.path("search/"+safeSearchValue);
	};

	$scope.getClass = function(path) {
		if ($location.path().substr(0, path.length) == path) {
		  return "active";
		} else {
		  return "";
		}
	}
});

nulpunt.controller("OverviewCtrl", function($scope){
	//++
});


nulpunt.controller("DashboardCtrl", function($scope, $http, SearchDocumentService) {
	$scope.documents = [];
	$scope.searchTags = [];
	   // Don't show anything at page load
	   // Leave it in for when we will default on the users' profile.
	// $http({method: "POST", url: "/service/getDocuments", data: {} }).
	// success(function(data) {
	// 	console.log(data);
	// 	$scope.documents = data.documents;
	// }).
	// error(function(error) {
	//     console.log('error retrieving raw documents: ', error);
	// });

	// Tagsearch gets the tag to add or remove.
	$scope.TagSearch = function(tags, tag) {
		console.log("TagSearch has: ", tags, tag)
		//var tags = profile_tags.filter(function(x) {return true}); // copy into new array to make it idempotent.
		var index = tags.indexOf(tag)
		if (index > -1) { // found it, remove from tags list
			tags.splice(index, 1);
		} else { // not in there, add it
			tags.push(tag);
		};
		console.log("TagSearch has: ", tag, " -> ", tags)
		SearchDocumentService.searchDocuments(tags).then(
			function(data) {
				//console.log("TagSearch got from SearchDoc promise: ", data);
				$scope.documents = data.documents;
			},
				function(error) {
					console.log('error retrieving raw documents: ', error);
					deferred.reject('error');
			});
	};
	
	// To assist in show/hide
	$scope.isElement = function(tags, tag) {
		if(tags == undefined) {
			return false;
		}
		var index = tags.indexOf(tag);
		if (index == -1) {
		return "np-notselected"
		} else {
		return "np-selected"
		}
	};
});


nulpunt.controller("InboxCtrl", function() {
});

// This controllers is used on inbox-page to query on users' selected Tags
nulpunt.controller("DocumentsByTagsCtrl", function ($scope, $http, ProfileService, SearchDocumentService) {
	console.log("DocumentsByTagsCtrl has found profile: ", $scope.profiles);
	ProfileService.getProfile().then(
	function(profile) {
		
		console.log("DocumentsByTagCtrl got from Profile promise: ", profile)
			SearchDocumentService.searchDocuments(profile.Tags).then(
		function(data) {
			console.log("DocumentsByTagCtrl got from SearchDoc promise: ", data);
			$scope.documents = data.documents;
		},
			function(error) {
				console.log('error retrieving raw documents: ', error);
				deferred.reject('error');
		}),
		function(error) {
			console.log('error retrieving raw documents: ', error);
			deferred.reject('error');
		}});
	
	// Tagsearch gets the tag to add or remove.
	$scope.TagSearch = function(tags, tag) {
		//console.log("TagSearch has: ", profile_tags, tag)
		//var tags = profile_tags.filter(function(x) {return true}); // copy into new array to make it idempotent.
		var index = tags.indexOf(tag)
		if (index > -1) { // found it, remove from tags list
			tags.splice(index, 1);
		} else { // not in there, add it
			tags.push(tag);
		};
		//console.log("TagSearch has: ", profile_tags, tag, " -> ", tags)
		SearchDocumentService.searchDocuments(tags).then(
			function(data) {
				console.log("TagSearch got from SearchDoc promise: ", data);
				$scope.documents = data.documents;
			},
				function(error) {
					console.log('error retrieving raw documents: ', error);
					deferred.reject('error');
			});
	};
	
	// To assist in ng-show/hide
	$scope.isElement = function(tags, tag) {
		if(tags == undefined) {
			return false;
		}
		var index = tags.indexOf(tag);
		return index != -1;
	};
});

nulpunt.controller("DocumentCtrl", function($scope, $http, $routeParams) {
	$scope.currentPage = {
		number: 1,
		data: {},
	};

	$scope.nextPage = function() {
		$scope.currentPage.number++;
	}
	$scope.prevPage = function() {
		$scope.currentPage.number--;
	}

	$scope.$watch('currentPage.number', function() {
		if($scope.document != undefined && $scope.currentPage.number > $scope.document.PageCount) {
			//++ TODO WARNING: this check is skipped when document wasn't loaded yet..
			$scope.currentPage.number = $scope.document.PageCount;
			return;
		}
		if($scope.currentPage.number < 1) {
			$scope.currentPage.number = 1;
			return;
		}

		loadPage();
	});

	function loadPage() {
		$http({method: 'POST', url: "/service/getPage", data: {documentID: $routeParams.docID, pageNumber: $scope.currentPage.number}}).
			success(function(data) {
					console.log(data);
					$scope.currentPage.data = data;
			}).
			error(function(error) {
					console.error('error retrieving page information: ', error);
			});
	}

	$http({method: "POST", url: "/service/getDocument", data: { docID: $routeParams.docID } }).
		success(function(data) {
			console.log(data);
			$scope.document = data.document;
			$scope.annotations = data.annotations;
			$scope.twitter = {
				url: "https://nulpunt.nu/#/document/"+data.document.ID,
				text: data.document.Title,
			};
		}).error(function(error) {
			console.log('error retrieving raw documents: ', error);
		});

	var canvas = document.getElementById("cvPage");
	var ctx = canvas.getContext("2d");

	// page location and size
	var pageOffsetX;
	var pageOffsetY;
	var pageWidth;
	var pageHeight;

	// box locations
	var boxStartX;
	var boxStartY;
	var boxStopX;
	var boxStopY;

	var highlight = {x1: 0, x2: 0, y1: 0, y2: 0};

	var isDown = false;

	function handleMouseDown(e) {
		// udpate page and canvas width/height
		if(document.defaultView) {
			pageWidth = parseInt(document.defaultView.getComputedStyle(document.getElementById('pageBox'), "").getPropertyValue("width"));
			pageHeight = parseInt(document.defaultView.getComputedStyle(document.getElementById('pageBox'), "").getPropertyValue("height"));
		} else if(document.getElementById('pageBox').currentStyle) {
			pageWidth = parseInt(document.getElementById('pageBox').currentStyle["width"]);
			pageHeight = parseInt(document.getElementById('pageBox').currentStyle["height"]);
		} else {
			console.error('Could not update width/height on canvas element.');
			//++ TODO: automated error reporting (after user consent).
			alert('There seems to be a problem. Please report this problem to the nulpunt development team.');
			// fake and wrong values
			pageWidth = 42;
			pageHeight = 42;
		}

		// set canvas width/height
		document.getElementById("cvPage").width = pageWidth;
		document.getElementById("cvPage").height = pageHeight;

		// update offset details
		var canvasOffset = $("#pageBox").offset();
		pageOffsetX = canvasOffset.left;
		pageOffsetY = canvasOffset.top;
		console.log('pageOffsetX: '+pageOffsetX+' pageOffsetY: '+pageOffsetY);

		// save mouse location
		boxStartX = parseInt(e.clientX - pageOffsetX);
		boxStartY = parseInt(e.clientY - pageOffsetY + $(window).scrollTop());
		console.log('mouseX: '+boxStartX+' mouseY: '+boxStartY);

		// all done
		isDown = true;
	}

	function handleMouseUp(e) {
		// save mouse location
		boxStopX = parseInt(e.clientX - pageOffsetX);
		boxStopY = parseInt(e.clientY - pageOffsetY + $(window).scrollTop());
		console.log('mouseX: '+boxStopX+' mouseY: '+boxStopY);

		highlight.x1 = boxStartX/pageWidth*100;
		highlight.x2 = boxStopX/pageWidth*100;
		highlight.y1 = boxStartY/pageHeight*100;
		highlight.y2 = boxStopY/pageHeight*100;
		console.log(highlight);

		// all done
		isDown = false;
	}

	function handleMouseMove(e) {
		// only do stuff if mouse is down (update rect)
		if (!isDown) {
			return;
		}

		// save mouse positions
		boxStopX = parseInt(e.clientX - pageOffsetX);
		boxStopY = parseInt(e.clientY - pageOffsetY + $(window).scrollTop());

		// update highlight
		updateHighlight();
	}

	function updateHighlight() {
		// clear canvas
		ctx.clearRect(0, 0, canvas.width, canvas.height);

		// draw new rectangle
		var width = boxStopX - boxStartX;
		var height = boxStopY - boxStartY;
		ctx.beginPath();
		ctx.rect(boxStartX, boxStartY, width, height);
		ctx.globalAlpha = "0.5";
		ctx.fillStyle = "#FFFF00";
		ctx.fill();
	}

	$("#cvPage").mousedown(handleMouseDown);
	$("#cvPage").mousemove(handleMouseMove);
	$("#cvPage").mouseup(handleMouseUp);
});

nulpunt.controller("AnnotationSubmitCtrl", function($scope, $http) {
    $scope.showForm = false;
	$scope.submit = function() {
		$http({method: 'POST', url: '/service/session/add-annotation', data: {
			documentId: $scope.document.ID,
			locations: $scope.locations,
			annotationText: $scope.annotationText,
		}}).
		success(function(data, status, headers, config) {
		    $scope.annotations.push(data)
		    $scope.showForm = false;
		}).
		error(function(data, status, headers, config) {
		    console.log("invalid response for AnnotationSubmit:");
			console.log(data, status, headers);
			$scope.error = "Could not make request to server";
		});
	};
});

nulpunt.controller("CommentSubmitCtrl", function($scope, $http) {
    $scope.showForm = false;
	$scope.submit = function() {
		$http({method: 'POST', url: '/service/session/add-comment', data: {
			annotationId: $scope.annotation.ID,
			commentText: $scope.commentText,
			// parentId: $scope.parentID, // is for threaded comments
		}}).
		success(function(data, status, headers, config) {
		    $scope.annotation.Comments.push(data);
		    $scope.showForm = false;
		}).
		error(function(data, status, headers, config) {
			console.log("invalid response for CommentSubmit:")
			console.log(data, status, headers);
			$scope.error = "Could not make request to server";
		});
	};
});

nulpunt.controller("HistoryCtrl", function($scope, $routeParams) {
	$scope.documents = [];
});

nulpunt.controller("TrendingCtrl", function($scope) {
	$scope.documents = {
		items: [],
	};

	$scope.documents.items = [
		{
			title: "Verzoek militaire bijstand inzet Raven UAV Culemborg", 
			description: "Verzoek militaire bijstand in het kader van de handhaving van de openbare orde op grond van de Politiewet, art. 59, met de inzet van een Raven mini UAV op 18, 21 en 31 december 2012 in de gemeente Culemborg tijdens de opname van televisie programma 'Pauw in Culemborg' door de Burgemeester van de gemeente Culemborg.", 
			source: "Ministerie van Binnenlandse Zaken en Koninkrijksrelaties", 
			country: "NL",
			sourceDate: "01/12/2011", 
			uploadDate: "05/12/2013",
			docID: "52bee9390b4aec49b4a50be2",
			uploader: "0.",
			uploaderColor: "#4effa4",
			requester: "Rejo Zenger",
			type: "Verzoek Externe Bijstand",
			nrOfPages: 2,
			nrOfAnnotations: 6,
			nrOfDrafts: 2,
			nrOfComments: 8,
			nrOfBookmarks: 4,
			tags: [
				{title: "Drones"}, 
				{title:"Police"}, 
				{title:"Surveillance"}
			],
			annotations: [
				{annotationDate: "2013-11-24", annotator: "rick", annotation: "Quas illaboritati ius de plit prae vid maxim que dendae re ne plaborio. Facideb itatur ressiment apiendae. Itatemo luptaestius am essimi, te rem volorum sed maximintiis si remporp oremperatia dit incitati dolorposse provitas ad ut fuga. Hillore nobitemquis et ma si con commol"}
			]
		},
		{
			title: "Verdict US Foreign Intelligence Surveillance Court NSA", 
			description: "The Foreign Intelligence Surveillance Court verdict on the bulk acquisition of metadata by the National Security Agency (NSA). The matter discussed is the gathering of large amounts of data, which for years have exceeded the previously authorized acquisition.", 
			source: "United States Foreign Intelligence Surveillance Court",
			country: "USA",
			sourceDate: "unknown", 
			uploadDate: "12/11/2013",
			docID: "52c1ac7d0b4aec49b4a50eaf",			
			uploader: "0.",
			uploaderColor: "#ffb060",
/* 			requester: "xnone", */
/* 			type: "congressional report", */
			nrOfPages: 11,
			nrOfAnnotations: 32,
			nrOfDrafts: 7,
			nrOfComments: 18,
			nrOfBookmarks: 12, 
			tags: [
				{title: "Data Retention"},
				{title: "Espionage"},
				{title: "Surveillance"},
				{title: "Privacy"}													
			],
			annotations: [
				{annotationDate: "2013-08-20", annotator: "rick", annotation: "Quas illaboritati ius de plit prae vid maxim que dendae re ne plaborio. Facideb itatur ressiment apiendae. Itatemo luptaestius am essimi, te rem volorum sed maximintiis si remporp oremperatia dit incitati dolorposse provitas ad ut fuga. Hillore nobitemquis et ma si con commol"}
			] },
		{
			title: "Contract Koninklijke Landmacht en Blackwater", 
/* 			description: "A short 1 or 2 sentence description of the document. Include or not?",  */
			source: "Ministerie van Defensie",
			country: "NL",
			sourceDate: "12/01/2005", 
			uploadDate: "03/12/2013",
			docID: "52bf00290b4aec49b4a50bfc",
			uploader: "0.",
			type: "contract report",
			nrOfPages: 12,
			uploaderColor: "#00b7ff",
			requester : "R.G.C. Bik, Bureau Jansen & Janssen",
			nrOfAnnotations: 2,
			nrOfDrafts: 14,
			nrOfComments: 25,
			nrOfBookmarks: 4,
			tags: [
				{title: "Corporate State"},
				{title: "Military"},
				{title: "State Terrorism"}								
				],
			annotations: [
				{annotationDate: "some day", annotator: "rick", annotation: "Quas illaboritati ius de plit prae vid maxim que dendae re ne plaborio. Facideb itatur ressiment apiendae. Itatemo luptaestius am essimi, te rem volorum sed maximintiis si remporp oremperatia dit incitati dolorposse provitas ad ut fuga. Hillore nobitemquis et ma si con commol"}
			] },
	];
});

nulpunt.controller("NotificationsCtrl", function($scope) {
	$scope.notifications = [];
});

nulpunt.controller("SearchCtrl", function($scope, $routeParams) {
	$scope.mySearch = $routeParams.searchValue.replace(/[+]/g, ' ');
});

nulpunt.controller("ProfileCtrl", function($scope, $http) {
	$scope.done = false;
	$scope.error = "";
	// load the users' profile
	$http({
		method: "GET", 
		url: "/service/session/get-profile", 
		// no parameters, the server uses the session.account.username value.
	}).
	success(function(data) {
		console.log(data);
			 // UGLY HACK: 
			// Each user has only one profile, yet  we create an array.
			// This is so that the inbox.html template can use a ng-repeat
			// That makes the dependencies between that and this controller clear to Angular.
		$scope.profile = data.profile;
		$scope.profiles = [ data.profile ];
	}).
	error(function(error) {	
		console.log('error retrieving profile ', error);
		$scope.error = error;
	})

	// save the updated document
	$scope.submit = function() {
		$scope.done = false;
		$scope.error = "";
		$http({
			method: 'POST', 
			url: "/service/session/update-profile",
			data: $scope.profile
		}).
		success(function(data, status, headers, config) {
			console.log(data)
			$scope.done = true
		}).
		error(function(data, status, headers, config) {
			console.log("error updateProfile");
			console.log(data, status, headers);
			$scope.error = data;
		});
	}
});

nulpunt.controller('NotFoundCtrl', function($scope, $location) {
	$scope.path = $location.url();
});

nulpunt.controller('AboutCtrl', function($scope, $location) {
	$scope.path = $location.url();
});
nulpunt.controller('ContactCtrl', function($scope, $location) {
	$scope.path = $location.url();
});
nulpunt.controller('ColophonCtrl', function($scope, $location) {
	$scope.path = $location.url();
});

nulpunt.controller("RegisterCtrl", function($scope, $rootScope, $http) {
	$scope.submit = function() {
		$http({method: 'POST', url: '/service/session/registerAccount', data: {
			username: $scope.username,
			email: $scope.email,
		    password: $scope.password,
		    color: $scope.color
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
				window.location.href = "/#/documents";
				window.location.reload();
			}, function(error) {
				if(error == "") {
					// no success, but also no error: credentials are wrong.
					$scope.wrong = true;
				} else {
					$scope.error = error;
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

nulpunt.controller("AdminTagsCtrl", function($scope, $rootScope, $http, TagService) {
	TagService.getTags().then(
	function(data) {
		console.log("AdminTagsCtrl received data: ", data);
		$scope.tags = data.tags;
	},
	function(error) {
		console.log(error);
	}
	);
	
	$scope.add_tag = function() {
	console.log('adding tag: ', $scope.tag);
	TagService.addTag($scope.tag).then(
		function(data) {
		console.log(data);
		$scope.tags = data.tags;
		$scope.done = true;
		},
		function(error) {
		console.log(error);
		}
	)};
	
	$scope.delete_tag = function(tagname) {
	console.log('deleting tag: '+tagname);
	TagService.deleteTag(tagname).then(
		function(data) {
			console.log(data);
			//var index = $scope.tags.indexOf($scope.tag)
			//$scope.tags.splice(index, 1);
			$scope.tags = data.tags;
			$scope.done = true;
		},
		function(error) {
			console.log(error);
		}
	)};
});
	

nulpunt.controller("AdminUploadCtrl", function($scope, $upload) {
	$scope.uploading = false;
	$scope.language = "nl_NL"; // default

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
				data:  { language: $scope.language },
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
// nulpunt.controller("AdminAnalyseCtrl", function($scope, $http) {
// 	$scope.files = [];
// 	$http({method: "POST", url: "/service/session/admin/getRawUploads"}).
// 	success(function(data) {
// 	    console.log(data);
// 		$scope.files = data.files;
// 	}).
// 	error(function(error) {
// 		console.log('error retrieving raw documents: ', error);
// 	})

//     // create a new (unpublished) document to make testing document/pages possible
//     $scope.analyse = function(ind) {
// 	$http({method: "POST", url: "/service/session/admin/insertDocument",
// 	       data: { 
// 		   //document: {
// 		       title:                   $scope.files[ind].filename,
// 		       uploaderUsername:          $scope.files[ind].uploader,
// 		       uploadDate: $scope.files[ind].uploadDate,
// 		       language:         $scope.files[ind].language,
// 		   //}
// 	       }}).
// 	    success(function(data) {
// 		console.log('success updating document.');
// 		alert("succes");
// 	    }).
// 	error(function(error) {
// 		console.log('error updating document: ', error);
// 	});
//     }});

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


nulpunt.controller("AdminProcessEditMetaCtrl", function($scope, $http, $routeParams, $filter, $window) {
	$scope.done = false;
	$scope.error = "";
	console.log("DocID is " + $routeParams.docID );
	// load the requested document
	$http({
		method: "POST", 
		url: "/service/getDocument", 
		data: { 
			DocID: $routeParams.docID,
		}
	}).
	success(function(data) {
		console.log(data);

		$scope.OriginalDateString = String($filter('date')(data.document.OriginalDate, 'dd-MM-yyyy'));
		$scope.document = data.document;
		if($scope.document.Tags == undefined) { $scope.document.Tags = []; };
		console.log(data);
	}).
	error(function(error) {
		console.log('error retrieving document: ', error);
	});

	// save the updated document
	$scope.submit = function() {
		console.log("originalDateString: "+$scope.OriginalDateString);
		
		var dateInfo = $scope.OriginalDateString.split('-');
		var day = dateInfo[0] || "01";
		var month = dateInfo[1] || "01";
		var year = dateInfo[2] || "0001";

		var newStr = year + '-' + month + '-' + day + 'T00:00:00Z';
		console.log('Saving in doc: '+newStr);
		console.log('Reverse: '+String($filter('date')(newStr, 'dd-MM-yyyy')))
		
		var doc = $scope.document;
		doc.OriginalDate = newStr;
	
		$scope.done = false;
		$scope.error = "";
		$http({
			method: 'POST', 
			url: "/service/session/admin/updateDocument",
			data: doc
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
	
	$("#originalDate").datepicker({format: 'dd-mm-yyyy'});

	// disable the publish checkbox when OCR is not completed yet
	$scope.isDisabled = function(state) {
	return (state != "completed");
	};

	$scope.deleteDocument = function(docID) {
		doit = confirm("Delete this document,\nall annotations and comments.\n\nDeleting is permanent.\n");
		if (doit == true) {
			$http({
			method: 'POST', 
			url: "/service/session/admin/deleteDocument",
			data: { DocID: docID }
			}).
			success(function(data, status, headers, config) {
				console.log(data);
				alert("Your document is gone (forever).");
				$scope.done = true;
			}).
			error(function(data, status, headers, config) {
				console.log("error add updateDocument");
				console.log(data, status, headers);
				alert("Deletion gave an error. Your document might or might not be there.");
				$scope.error = data;
			});
			$window.location.href = "/#/admin/process";
		} else {
			$scope.deleteflag = false;
		};
	};
});

nulpunt.controller("SignOutCtrl", function($scope, $location, AccountAuthService, ClientSessionService) {
	$scope.username = AccountAuthService.getUsername();
	ClientSessionService.stopSession().then(function() {
		window.location.href = '/#/sign-out-successful';
		window.location.reload();
	}, function() {
		console.error('Could not destroy session. Internet connection lost?');
		alert('Could not destroy session. Internet connection lost?');
	});
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

nulpunt.filter('urlencode', function() {
	return window.escape;
});