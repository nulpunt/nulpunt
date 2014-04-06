var nulpunt = angular.module('nulpunt', [
	// imported modules
	// please keep this list sorted
	'ngRoute',
	'ngStorage',
	'ui.bootstrap',
	// 'ui.bootstrap.collapse',
	// 'ui.bootstrap.dropdownToggle',
	'angularFileUpload',
	'checklist-model'
]);

nulpunt.factory('LoginFactory', function($modal, ClientSessionService) {
	this.showLogin = function() {
		var loginModalInstance = $modal.open({
			templateUrl: 'html/sign-in.html',
			controller: "SignInCtrl",
		});
	}

	this.signOut = function() {
		//$scope.username = AccountAuthService.getUsername();
		ClientSessionService.stopSession().then(function() {
			window.location.href = '/#/';
			window.location.reload();
		}, function() {
			console.error('Could not destroy session. Internet connection lost?');
			alert('Could not destroy session. Internet connection lost?');
		});
	}

	return this;
});

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
		controller: "ProfileCtrl"
	})
	.when('/bookmarks', {
		templateUrl: '/html/bookmarks.html',
		controller: "BookmarksCtrl"
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

// EmptyCtrl for when no controller is required
nulpunt.controller("EmptyCtrl", function() {
	// Empty controller can be used to when a template specifies the controllers in-line.
});

nulpunt.controller("MainCtrl", function($scope, $rootScope, AccountAuthService) {
	// change account in scope on auth_changed event
	$rootScope.$on("auth_changed", function() {
		$scope.account = AccountAuthService.account;
	});
});

// NavbarCtrl manages the top navigation bar
nulpunt.controller("NavbarCtrl", function($scope, $rootScope, $location, LoginFactory, AccountAuthService) {
	// change account in scope on auth_changed event
	$rootScope.$on("auth_changed", function() {
		$scope.account = AccountAuthService.account;
	});

	// search handler
	$scope.search = function() {
		var safeSearchValue = $scope.searchValue.replace(/[\/\? ]/g, '+').replace('++', '+').trim('+');
		$location.path("search/"+safeSearchValue);
	};

	// returns wether given page path is currently active
	$scope.isActivePage = function(path) {
		if ($location.path().substr(0, path.length) == path) {
			return "active";
		} else {
			return "";
		}
	}

	$scope.loginFactory = LoginFactory;
});

nulpunt.controller("OverviewCtrl", function($scope){
	//++
});


nulpunt.controller("DashboardCtrl", function($scope, $http, SearchDocumentService) {
	$scope.documents = [];
	$scope.searchTags = [];
	// TODO MARKED FOR REMOVAL
	// // Don't show anything at page load
	// // Leave it in for when we will default on the users' profile.
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
		// TODO MARKED FOR REMOVAL
		// var tags = profile_tags.filter(function(x) {return true}); // copy into new array to make it idempotent.
		var index = tags.indexOf(tag)
		if (index > -1) {
			// found it, remove from tags list
			tags.splice(index, 1);
		} else {
			// not in there, add it
			tags.push(tag);
		}
		console.log("TagSearch has: ", tag, " -> ", tags)
		SearchDocumentService.searchDocuments(tags).then(
			function(data) {
				//console.log("TagSearch got from SearchDoc promise: ", data);
				$scope.documents = data.documents;
			},
			function(error) {
					console.log('error retrieving raw documents: ', error);
					deferred.reject('error');
			}
		);
	};
	
	// To assist in show/hide
	$scope.isElement = function(tags, tag) {
		if(tags == undefined) {
			return false;
		}
		var index = tags.indexOf(tag);
		if (index == -1) {
			return "np-notselected";
		} else {
			return "np-selected";
		}
	};
});


nulpunt.controller("InboxCtrl", function() {
	//++
});

// This controllers is used on inbox-page to query on users' selected Tags
nulpunt.controller("DocumentsByTagsCtrl", function ($scope, $http, ProfileService, SearchDocumentService) {
	console.log("DocumentsByTagsCtrl has found profile: ", $scope.profiles);
	ProfileService.getProfile().then(
		function(profile) {
			console.log("DocumentsByTagCtrl got from Profile promise: ", profile);
			SearchDocumentService.searchDocuments(profile.Tags).then(
				function(data) {
					console.log("DocumentsByTagCtrl got from SearchDoc promise: ", data);
					$scope.documents = data.documents;
				},
				function(error) {
					console.log('error retrieving raw documents: ', error);
					deferred.reject('error');
				}
			);
		},
		function(error) {
			console.log('error retrieving profile: ', error);
			deferred.reject('error');
		}
	);
	
	// Tagsearch gets the tag to add or remove.
	$scope.TagSearch = function(tags, tag) {
		// TODO MARKED FOR REMOVAL
		//console.log("TagSearch has: ", profile_tags, tag)
		// TODO MARKED FOR REMOVAL
		//var tags = profile_tags.filter(function(x) {return true}); // copy into new array to make it idempotent.
		var index = tags.indexOf(tag)
		if (index > -1) {
			// found it, remove from tags list
			tags.splice(index, 1);
		} else {
			// not in there, add it
			tags.push(tag);
		};
		// TODO MARKED FOR REMOVAL
		//console.log("TagSearch has: ", profile_tags, tag, " -> ", tags)
		SearchDocumentService.searchDocuments(tags).then(
			function(data) {
				console.log("TagSearch got from SearchDoc promise: ", data);
				$scope.documents = data.documents;
			},
			function(error) {
				console.log('error retrieving raw documents: ', error);
				deferred.reject('error');
			}
		);
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

nulpunt.controller("DocumentCtrl", function($scope, $http, $routeParams, $modal, $window, LoginFactory) {
	$scope.currentPage = {
		number: 1,
		data: {},
	};

	$scope.loginFactory = LoginFactory;

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

    // get any annotation that has coordinates at the given pageNr.
    $scope.annotationsOnPage = function(pageNr) {
		//console.log("filter annotations on page: ", pageNr);
		//console.log("annotatations in scope: ", $scope.annotations);
		annotations = _.filter($scope.annotations, function(ann) {
		    return _.some(ann.Locations, function(loc) { 
				//console.log("found: ", loc);
				return loc.PageNumber == pageNr;
		    })
		})
		//console.log("returning: ", annotations);
		return annotations
    }

    function clearHighlights() {
		document.getElementById("cvPage").width = 0;
		document.getElementById("cvPage").height = 0;
    }

	function loadPage() {
		console.log("loadpage");
		clearHighlight();	// Clear the highlight you were creating
		clearHighlights();	// Clear all previously created highlights
		$http({method: 'POST', url: "/service/getPage", data: {documentID: $routeParams.docID, pageNumber: $scope.currentPage.number}}).
			success(function(data) {
				console.log("loadpage data");
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

	// get canvas
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
	$scope.highlight = highlight;

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
		// console.log('pageOffsetX: '+pageOffsetX+' pageOffsetY: '+pageOffsetY);

		// save mouse location
		boxStartX = parseInt(e.clientX - pageOffsetX);
		boxStartY = parseInt(e.clientY - pageOffsetY + $(window).scrollTop());
		// console.log('mouseX: '+boxStartX+' mouseY: '+boxStartY);

		// Hide + button
		$('#annotation-add-btn').hide();

		// all done
		isDown = true;
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

	function handleMouseUp(e) {
		// save mouse location
		boxStopX = parseInt(e.clientX - pageOffsetX);
		boxStopY = parseInt(e.clientY - pageOffsetY + $(window).scrollTop());

		// Calculate box size
		boxWidth = Math.abs(boxStopX - boxStartX);
		boxHeight = Math.abs(boxStopY - boxStartY);

		if(boxWidth < 10 || boxHeight < 10) {
			// clear selection if it's too small
			clearHighlight();
		} else {
			// set hightlight
			highlight.pagenumber = $scope.currentPage.number;
			highlight.x1 = boxStartX/pageWidth*100;
			highlight.x2 = boxStopX/pageWidth*100;
			highlight.y1 = boxStartY/pageHeight*100;
			highlight.y2 = boxStopY/pageHeight*100;

			// set bottom and right to show "add annotation" box
			highlight.xMax = Math.max(boxStartX, boxStopX);
			highlight.yMax = Math.max(boxStartY, boxStopY);
		}

		$scope.$apply();

		console.log(highlight);
		// TODO: check if highlight is not too large.
		// In case of large highlight, give notification and color red..

		// Show + button
		$('#annotation-add-btn').show();

		// all done
		isDown = false;
	}

	function clearHighlight() {
		highlight.pagenumber = $scope.currentPage.number;
		highlight.x1 = 0;
		highlight.x2 = 0;
		highlight.y1 = 0;
		highlight.y2 = 0;

		highlight.xMax = 0;
		highlight.yMax = 0;
	}


	// update highlight on screen with latest info
	function updateHighlight() {
		// clear canvas
		ctx.clearRect(0, 0, canvas.width, canvas.height);

		// draw new rectangle
		var width = boxStopX - boxStartX;
		var height = boxStopY - boxStartY;
		ctx.beginPath();
		ctx.rect(boxStartX, boxStartY, width, height);
		ctx.globalAlpha = "0.6";
		ctx.fillStyle = $scope.account.color;
		ctx.fill();
	}

	// attach mouse handlers
	$("#cvPage").mousedown(handleMouseDown);
	$("#cvPage").mousemove(handleMouseMove);
	$("#cvPage").mouseup(handleMouseUp);


	$scope.addAnnotation = function () {

		var modalInstance = $modal.open({
			templateUrl: 'html/new-annotation-modal.html',
				controller: "NewAnnotationModal",
				resolve: {
				highlight: function () {
					return highlight;
				},
				documentId: function() {
					return $scope.document.ID;
				},
				pageNr: function() {
					return $scope.currentPage.number;
				},
			}
		});

		modalInstance.result.then(function (annotationText) {
			console.log('annotation result: '+annotationText);
			$http({
				method: 'POST',
				url: '/service/session/add-annotation',
				data: {
					documentId: $scope.document.ID,
					annotationText: annotationText,
				        locations: [  $scope.highlight ],
				}
			}).
			success(function(data, status, headers, config) {
				$scope.annotations.push(data);
				$scope.showForm = false;
				loadPage();
			}).
			error(function(data, status, headers, config) {
				console.log("invalid response for AnnotationSubmit:");
				console.log(data, status, headers);
				$scope.error = "Could not make request to server";
			});
		}, function (info) {
			console.log('modal dismissed because: '+info);
		});
	};

	$scope.activateHighlight = function(annotationId) {
		$('#' + annotationId).addClass('active-highlight');
		$('#annotation_' + annotationId).addClass('active-highlight');
	};

	$scope.deactivateHighlight = function() {
		$('.active-highlight').removeClass('active-highlight');
	};

    $scope.shareTwitter = function () {
        $window.open('https://twitter.com/share?url='+encodeURIComponent($scope.twitter.url)+'&text='+encodeURIComponent($scope.twitter.text)+'&hashtags=nulpunt','das','location=no,links=no,scrollbars=no,toolbar=no,width=750,height=300');
    }
	$scope.shareDiaspora = function () {
		$window.open('http://sharetodiaspora.github.io/?url='+encodeURIComponent($scope.twitter.url)+'&title='+encodeURIComponent($scope.twitter.text),'das','location=no,links=no,scrollbars=no,toolbar=no,width=620,height=550');
	}
    $scope.shareGooglePlus = function () {
        $window.open('https://plus.google.com/share?url='+encodeURIComponent($scope.twitter.url)+'&title=title&text=text&message=message','das','location=no,links=no,scrollbars=no,toolbar=no,width=520,height=500');
    }
    $scope.shareReddit = function () {
        $window.open('http://www.reddit.com/submit?url='+encodeURIComponent($scope.twitter.url)+'&title='+encodeURIComponent($scope.twitter.text),'das','location=no,links=no,scrollbars=no,toolbar=no,width=850,height=550');
    }
    $scope.shareVK = function () {
        $window.open('https://vk.com/share.php?url='+encodeURIComponent($scope.twitter.url)+'&title='+encodeURIComponent($scope.twitter.text),'das','location=no,links=no,scrollbars=no,toolbar=no,width=550,height=375');
    }
    $scope.shareLinkedIn = function () {
        $window.open('http://www.linkedin.com/shareArticle?mini=true&url='+encodeURIComponent($scope.twitter.url)+'&title='+encodeURIComponent($scope.twitter.text),'das','location=no,links=no,scrollbars=no,toolbar=no,width=600,height=500');
    }

    // add a bookmark
	$scope.bookmark = function(documentId) {
	    console.log("adding documentId to bookmarks");
	    $http({
		method: 'POST', 
		url: "/service/session/add-bookmark",
		data: { "DocumentID": documentId },
		}).
		success(function(data, status, headers, config) {
		    console.log(data)
		    $scope.done = true
		}).
		error(function(data, status, headers, config) {
		    console.log("error addBookmark");
 		    console.log(data, status, headers);
 		    $scope.error = data;
 		});
	}
});

nulpunt.controller("NewAnnotationModal", function($scope, $modalInstance, highlight, documentId, pageNr) {
	// highlight location for crop
	// TODO: use documentId, pageNr and highlight area for page crop
	$scope.highlight = highlight;
	$scope.documentId = documentId;
	$scope.pageNr = pageNr;

	$scope.annotation = {
		text: "",
	};

	// save new annotation
	$scope.save = function () {
		$modalInstance.close($scope.annotation.text);
	};

	// cancel new annocation
	$scope.cancel = function () {
		$modalInstance.dismiss('cancel');
	};
});

// MARKED FOR DELETION
// MARKED FOR DELETION
// MARKED FOR DELETION
// MARKED FOR DELETION
// MARKED FOR DELETION
// nulpunt.controller("AnnotationSubmitCtrl", function($scope, $http) {
// 	$scope.showForm = false;
// 	$scope.submit = function() {
// 		$http({
// 			method: 'POST',
// 			url: '/service/session/add-annotation',
// 			data: {
// 				documentId: $scope.document.ID,
// 				locations: $scope.locations,
// 				annotationText: $scope.annotationText,
// 			}
// 		}).
// 		success(function(data, status, headers, config) {
// 			$scope.annotations.push(data)
// 			$scope.showForm = false;
// 		}).
// 		error(function(data, status, headers, config) {
// 			console.log("invalid response for AnnotationSubmit:");
// 			console.log(data, status, headers);
// 			$scope.error = "Could not make request to server";
// 		});
// 	};
// });

nulpunt.controller("CommentSubmitCtrl", function($scope, $http) {
	$scope.showForm = false;
	$scope.submit = function() {
		$http({method: 'POST', url: '/service/session/add-comment', data: {
			annotationId: $scope.annotation.ID,
			commentText: $scope.commentText,
			// TODO MARKED FOR REMOVAL
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

// TrendingCtrl retrieves latest trending data from server
nulpunt.controller("TrendingCtrl", function($scope, $http, $routeParams) {
	$http({
		method: "GET", 
		url: "/service/get-trending",
		//data: {
		//    offset: $routeParams.offset,
		//    limit: $routeParams.limit,
		//},
	}).
	success(function(data) {
	    console.log("Trending:", data);
		$scope.trending = data.trending;
	}).
	error(function(error) {	
		console.log('error retrieving trending ', error);
		$scope.error = error;
	});
});

// NotificationsCtrl (TODO)
nulpunt.controller("NotificationsCtrl", function($scope) {
	$scope.notifications = [];
});

// SearchCtrl (TODO) makes a search request at the server and displays the data through search.html
nulpunt.controller("SearchCtrl", function($scope, $routeParams) {
	$scope.mySearch = $routeParams.searchValue.replace(/[+]/g, ' ');
});

// ProfileCtrl
nulpunt.controller("ProfileCtrl", function($scope, $http, $routeParams) {
	$scope.done = false;
	$scope.error = "";
    $scope.selectedTags = [];
	// load the users' profile
	$http({
		method: "GET", 
		url: "/service/session/get-profile",
	}).
	success(function(data) {
		console.log(data);
		// UGLY HACK: 
		// Each user has only one profile, yet  we create an array.
		// This is so that the profile.html template can use a ng-repeat
		// That makes the dependencies between that and this controller clear to Angular.
	        // It solves the race condition between get-profile and get-tags
		$scope.profile = data.profile;
		$scope.profiles = [ data.profile ];
        $.each(data.profile.Tags, function(index, tag) {
            $scope.selectedTags.push(tag);
        });
	}).
	error(function(error) {	
		console.log('error retrieving profile ', error);
		$scope.error = error;
	});

    $scope.isSelectedTag = function(tags, tag) {
        if(tags == undefined) {
            return false;
        }
        var index = tags.indexOf(tag);
        if (index == -1) {
            return "np-notselected";
        }
        else {
            return "np-selected";
        }
    }

    $scope.toggleTag = function(tags, tag) {
        var index = tags.indexOf(tag);
        console.log(index);
        console.log(tags);
        if (index > -1) {
            // Found it, remove from selectedTags
            tags.splice(index, 1);
            $(this).closest('input[type=checkbox]').prop('checked', false);
        }
        else {
            // Add it
            tags.push(tag);
            $(this).closest('input[type=checkbox]').prop('checked', true);
        }
    }

	// save the updated document
    //console.log("submitting new profile to service");
	$scope.submit = function() {
		$scope.done = false;
		$scope.error = "";
        $scope.profile.Tags = $scope.selectedTags;
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

// BookmarksCtrl
nulpunt.controller("BookmarksCtrl", function($scope, $http) {
	$scope.done = false;
	$scope.error = "";
	// load the users' bookmarks
	$http({
		method: "GET", 
		url: "/service/session/get-bookmarks",
	}).
	success(function(data) {
	    console.log(data);
	    $scope.bookmarks = data.bookmarks;
	}).
	error(function(error) {	
		console.log('error retrieving bookmark ', error);
		$scope.error = error;
	});
});


// NotFoundCtrl prepares information for the not-found.html page
nulpunt.controller('NotFoundCtrl', function($scope, $location) {
	$scope.path = $location.url();
});

// AboutCtrl controls the about page
nulpunt.controller('AboutCtrl', function($scope, $location) {
	// empty controller
});

// ContactCtrl controls the contact page
nulpunt.controller('ContactCtrl', function($scope, $location) {
	// empty controller
});

// ColophonCtrl controls the colophon page
nulpunt.controller('ColophonCtrl', function($scope, $location) {
	// empty controller
});

// RegisterCtrl (TODO) checks registration input and sends the registration request
nulpunt.controller("RegisterCtrl", function($scope, $rootScope, $http) {
	// TODO: check input on-change (passwords match etc. etc.)

	// submit sends registration request to the server
	$scope.submit = function() {
		$http({method: 'POST', url: '/service/session/registerAccount', data: {
			username: $scope.username,
			email: $scope.email,
			password: $scope.password,
			color: $scope.color
		}}).
		success(function(data, status, headers, config) {
			if(data.success) {
				// set error to null in case of previous error
				$scope.error = null;
				// registration is done
				$scope.done = true;
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

// SettingsCtrl fetches and stores settings
nulpunt.controller("SettingsCtrl", function($scope, AccountDataService) {
	// TODO: fix
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

// SignInCtrl manages user sign-in
nulpunt.controller("SignInCtrl", function($scope, $rootScope, $modalInstance, AccountAuthService) {
	$scope.signin = {
		username: "",
		password: "",
	};

	$scope.register = function () {
		$modalInstance.close();
		window.location.href = "#/register";
	};

	$scope.submit = function() {
		// reset state on scope
		$scope.success = false;
		$scope.wrong = false;
		$scope.error = "";

		// authenticate to server
		var prom = AccountAuthService.authenticate($scope.signin.username, $scope.signin.password);
		prom.then(
			function() {
				$scope.success = true;
				$modalInstance.close();
				//++ TODO: let user choose to go to dashboard or to go to the page he/she came from?
				// TODO MARKED FOR REMOVAL
				// window.location.href = "/#/dashboard";
				// window.location.reload();
			},
			function(error) {
				if(error == "") {
					// no success, but also no error: credentials are wrong.
					$scope.wrong = true;
				} else {
					$scope.error = error;
				}
				//++ TODO: need to do some "digest" on $scope ?? or $scope.$apply()? find out what good convention is
			}
		);
	};

	$scope.cancel = function () {
		$modalInstance.dismiss('cancel');
	};
	
	// watch auth_changed event and set scope if required
	// TODO: how is this cleaned up when controller is destroyed????
	$rootScope.$on("auth_changed", function() {
		$scope.account = AccountAuthService.account;
	});
});

// AdminTagsCtrl does stuff TODO
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
		);
	};
	
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
		);
	};
});
	
// AdminUploadCtrl manages document uploads
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

// TODO MARKED FOR REMOVAL
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

// AdminProcessCtrl to process the files
nulpunt.controller("AdminProcessCtrl", function($scope, $http) {
	$scope.documents = [];
	//$http({method: "POST", url: "/service/getDocumentList", data: {} }).
	$http({method: "GET", url: "/service/getDocumentList", data: {} }).
	success(function(data) {
		console.log(data);
		$scope.documents = data;
	}).
	error(function(error) {
		console.log('error retrieving raw documents: ', error);
	});
});

// AdminProcessEditMetaCtrl to edit the meta data
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

// SignOutCtrl kills the complete user session (effectively logging out)
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

// bytes filter converts number of bytes to human readable value
nulpunt.filter('bytes', function() {
	return function(bytes, precision) {
		if (bytes==0 || isNaN(parseFloat(bytes)) || !isFinite(bytes)){
			return '-';
		}
		if (typeof precision === 'undefined') {
			precision = 1;
		}
		var units = ['bytes', 'kB', 'MB', 'GB', 'TB', 'PB'],
		number = Math.floor(Math.log(bytes) / Math.log(1024));
		return (bytes / Math.pow(1024, Math.floor(number))).toFixed(precision) + ' ' + units[number];
	}
});

// urlencode filter url escapes the given string
nulpunt.filter('urlencode', function() {
	return window.escape;
});