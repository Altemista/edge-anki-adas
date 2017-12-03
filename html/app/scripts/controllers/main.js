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

var timer;

function refreshLoop($scope, $timeout, MainFactory) {
	timer = $timeout(function () {
		//console.log('INFO: Timer set');
	}, 300);
	timer.then(function () {
		//console.log('INFO: Timer triggered');
		MainFactory.getStatus()
			.then(
			function (response) { // ok
				$scope.status = response.data;
			},
			function (response) { // nok
				console.error('ERROR: Request failed: ' + response.statusText);
			}
			);
		if ($scope.poll) {
			refreshLoop($scope, $timeout, MainFactory);
		}
	}, function () {
		console.error('ERROR: Timer rejected!');
	});
}

/**
 * @ngdoc function
 * @name htmlApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the htmlApp
 */
angular.module('htmlApp')
	.controller('MainCtrl', ['$scope', '$timeout', 'MainFactory', 'MainConfig',
		function ($scope, $timeout, MainFactory, MainConfig) {
			this.awesomeThings = [
				'HTML5 Boilerplate',
				'AngularJS',
				'Karma'
			];

		// Handler function for togging of status polling
		$scope.togglePoll = function () {
			console.log('INFO: Handling togglePoll');
			if ($scope.poll) {
				$scope.text = 'Start polling';
				$scope.poll = false;
			} else {
				$scope.text = 'Stop polling';
				$scope.poll = true;
				refreshLoop($scope, $timeout, MainFactory);
			}
		};

		// Handler function starting car1
		$scope.startCar = function (carNo) {			
			console.log('INFO: Handling requestPing for carno '+carNo);
			var param = '';
			if (carNo===1) {
				param = ''+MainConfig.lowSpeed;
			} else if (carNo===2) {
				param = ''+MainConfig.highSpeed;
			}
			var command = {
				command: 's',
				param1: param,
				carNo: ''+carNo,
				source: 'ui'
			};
			MainFactory.postCommand(command)
				.then(
				function (response) { // ok
					console.log('INFO: Ping command submitted to server: '+response.statusText);
				},
				function (response) { // nok
					console.error('ERROR: Request failed: ' + response.statusText);
				}
				);
		};

		// Handler function starting car2
		$scope.stopCar = function (carNo) {			
			console.log('INFO: Handling requestPing for carno '+carNo);
			var command = {
				command: 'e',
				carNo: ''+carNo,
				source: 'ui'
			};
			MainFactory.postCommand(command)
				.then(
				function (response) { // ok
					console.log('INFO: Ping command submitted to server: '+response.statusText);
				},
				function (response) { // nok
					console.error('ERROR: Request failed: ' + response.statusText);
				}
				);
		};

		// Initialise
		$scope.lastUpdate = 'N/A';
		$scope.poll = false;
		$scope.text = 'Start polling';
		$scope.status = [{},{},{}];
	}]);


