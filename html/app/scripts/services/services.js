// Copyright 2018 NTT Group

// Permission is hereby granted, free of charge, to any person obtaining a copy of this 
// software and associated documentation files (the "Software"), to deal in the Software 
// without restriction, including without limitation the rights to use, copy, modify, 
// merge, publish, distribute, sublicense, and/or sell copies of the Software, and to 
// permit persons to whom the Software is furnished to do so, subject to the following 
// conditions:

// The above copyright notice and this permission notice shall be included in all copies 
// or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, 
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR 
// PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE 
// FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR 
// OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER 
// DEALINGS IN THE SOFTWARE.

'use strict';

console.log('Initialise services at: '+window.location.href);
var baseURL = '';

function startsWithCharAt(string, pattern) {
    for (var i = 0, length = pattern.length; i < length; i += 1) {
      if (pattern.charAt(i) !== string.charAt(i)) {
		  return false;
	  }
    }
    return true;
}

if (startsWithCharAt(window.location.href, 'http://localhost:9000')) {
	baseURL = 'http://localhost:8002';
	console.log('INFO: Using hard coded dev server at: ' + baseURL);
}

/**
 * @ngdoc function
 * @name htmlApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the htmlApp
 */
angular.module('htmlApp')
	.factory('MainFactory', ['$http', function ($http) {
		var MainFactory = {};
		MainFactory.getStatus = function () {
			return $http.get(baseURL + '/v1/overtake/status');
		};
		MainFactory.postCommand = function (command) {
			return $http.post(baseURL + '/v1/overtake/command', command);
		};
		return MainFactory;
	}]);

/**
 * @ngdoc function
 * @name htmlApp.controller:MainConfig
 * @description
 * # MainCtrl
 * Controller of the htmlApp
 */
angular.module('htmlApp')
.factory('MainConfig', function () {
	var MainConfig = {};
	MainConfig.lowSpeed = 250;
	MainConfig.highSpeed = 500;
	return MainConfig;
});


