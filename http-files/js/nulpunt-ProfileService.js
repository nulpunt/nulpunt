// Profile, Tags and SearchDocument services.

nulpunt.factory('TagService', function($http, $q) {
    var service = {};
    
    service.getTags = function() {
     	var deferred = $q.defer();
     	$http({
     	    method: "GET", 
     	    url: "/service/session/get-tags", 
     	}).success(function(data) {
     	    deferred.resolve(data);
     	}).error(function(error) {
     	    console.log('error retrieving tags ', error);
     	    deferred.reject(error);
     	});
     	return deferred.promise;
     };

    service.addTag = function(tag) {
 	console.log('adding tag: '+ tag);
    	var deferred = $q.defer();
 	$http({
 	    method: 'POST', 
 	    //url: '/service/session/admin/add-tag',
 	    url: '/service/add-tag',
 	    data: { tag: tag }
      	}).success(function(data) {
     	    deferred.resolve(data);
     	}).error(function(error) {
     	    console.log('error adding tag ', error);
     	    deferred.reject(error);
     	});
     	return deferred.promise;
    }	

   service.deleteTag = function(tag) {
 	console.log('deleting tag: '+ tag);
    	var deferred = $q.defer();
 	$http({
 	    method: 'POST', 
 	    url: '/service/session/admin/delete-tag',
 	    data: { tag: tag }
      	}).success(function(data) {
     	    deferred.resolve(data);
     	}).error(function(error) {
     	    console.log('error adding tag ', error);
     	    deferred.reject(error);
     	});
     	return deferred.promise;
    }	
    
    return service;
});

nulpunt.factory('ProfileService', function($http, $q) {
    var service = {};
    
    service.getProfile = function() {
     	var deferred = $q.defer();
     	$http({
     	    method: "GET", 
     	    url: "/service/session/get-profile", 
     	    // no parameters, the server uses the session.account.username value.
     	}).success(function(data) {
	    console.log("ProfileService got: ", data)
     	    deferred.resolve(data.profile);
     	}).error(function(error) {
     	    console.log('ProfileService got error retrieving profile ', error);
     	    deferred.reject({error: error});
     	});
     	return deferred.promise;
     };

    return service;
});


nulpunt.factory('SearchDocumentService', function($http, $q) {
    var service = {};
    
    service.searchDocuments = function(tags) {
     	var deferred = $q.defer();
      	$http({method: "POST", url: "/service/session/get-documents-by-tags", data: {tags: tags} }).
 	    success(function(data) {
      		deferred.resolve(data);
      	    }).error(function(error) {
      		console.log('error retrieving profile ', error);
      		deferred.reject({error: error});
      	    });
      	return deferred.promise;
    };
    
    return service;
});
