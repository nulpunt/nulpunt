
nulpunt.factory('AccountDataService', function($http, $q, AccountAuthService) {
	var service = {};

	// loadBlob loads (and decrypts) datablob from server
	function loadBlob(name) {
		var deferred = $q.defer();

		$http({method: 'POST', url: '/service/session/dataBlobLoad', data: {name: name}}).
		success(function(data, status, headers, config) {
			// decode data
			data = angular.fromJson(data);

			// switch on ecnryption type
			switch(data.encryptionType) {
			case 'plain':
				deferred.resolve(data.blob);
				break;
			case 'encryptedABC':
				//++ TODO
				break;
			}
			deferred.resolve(data);
		}).
		error(function(data, status, headers, config) {
			deferred.reject('http error');
		});

		return deferred.promise;
	}

	// saveBlob (encrypts and) sends data to server for storage
	function saveBlob(name, blob) {
		var deferred = $q.defer();

		// different kinds of blob encryption
		var encType = 'plain';
		switch(encType) {
		case 'plain':
			// don't do anything with blob
			break;
		case 'encryptedABC':
			//++ encrypt blob
			break;
		}

		// encode data
		var blob = angular.toJson({encryptionType: encType, blob: blob});

		// make http call
		$http({method: 'POST', url: '/service/session/dataBlobSave', data: {name: name, blob: blob}}).
		success(function(data, status, headers, config) {
			deferred.resolve();
		}).
		error(function(data, status, headers, config) {
			deferred.reject('http error');
		});

		return deferred.promise;
	}

	service.getObject = function(name) {
		//++ use promise chaining
		var deferred = $q.defer();

		var blobPromise = loadBlob(name);
		blobPromise.then(function(blob) {
			var object = angular.fromJson(blob);
			deferred.resolve(object);
		},function() {
			deferred.reject('error');
		});

		return deferred.promise;
	}

	service.setObject = function(name, object) {
		//++ use promise chaining
		var deferred = $q.defer();

		var blob = angular.toJson(object);
		var promise = saveBlob(name, blob);
		promise.then(function(blob) {
			deferred.resolve();
		},function(error) {
			deferred.reject('error: '+error);
		});

		return deferred.promise;
	}

	return service;
});