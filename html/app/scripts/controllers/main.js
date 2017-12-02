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

		// Handler function for pinging the car
		$scope.requestPing = function (carNo) {			
			console.log('INFO: Handling requestPing for carno '+carNo);
			var command = {
				command: 'ping',
				carNo: carNo,
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

		// Handler function for version of the car
		$scope.requestVersion = function (carNo) {
			console.log('INFO: Handling requestVersion for carno '+carNo);
			var command = {
				command: 'ver',
				carNo: carNo,
				source: 'ui'
			};
			MainFactory.postCommand(command)
				.then(
				function (response) { // ok
					console.log('INFO: Version command submitted to server: '+response.statusText);
				},
				function (response) { // nok
					console.error('ERROR: Request failed: ' + response.statusText);
				}
				);
		};

		// Handler function for changeSpeedUp of the car (250-1000)
		$scope.changeSpeedUp = function (carNo) {
			console.log('INFO: Handling changeSpeedUp for carno '+carNo);
			var command = {
				command: 's',
				param1: ''+MainConfig.highSpeed,
				carNo: carNo,
				source: 'ui'
			};
			MainFactory.postCommand(command)
				.then(
				function (response) { // ok
					console.log('INFO: SpeedUp command submitted to server: '+response.statusText);
				},
				function (response) { // nok
					console.error('ERROR: Request failed: ' + response.statusText);
				}
				);
		};

		// Handler function for changeSpeedDown of the car (250-1000)
		$scope.changeSpeedDown = function (carNo) {
			console.log('INFO: Handling changeSpeedDown for carno '+carNo);
			var command = {
				command: 's',
				param1: ''+MainConfig.slowSpeed,
				carNo: carNo,
				source: 'ui'
			};
			MainFactory.postCommand(command)
				.then(
				function (response) { // ok
					console.log('INFO: SpeedDown command submitted to server: '+response.statusText);
				},
				function (response) { // nok
					console.error('ERROR: Request failed: ' + response.statusText);
				}
				);
		};

		// Handler function for changeSpeedStop of the car (250-1000)
		$scope.changeSpeedStop = function (carNo) {
			console.log('INFO: Handling changeSpeedStop for carno '+carNo);
			var command = {
				command: 'e',
				carNo: carNo,
				source: 'ui'
			};
			MainFactory.postCommand(command)
				.then(
				function (response) { // ok
					console.log('INFO: SpeedStop command submitted to server: '+response.statusText);
				},
				function (response) { // nok
					console.error('ERROR: Request failed: ' + response.statusText);
				}
				);
		};

		// Handler function for changeLaneLeft with the car (-68-+68)
		$scope.changeLaneLeft = function (carNo) {
			console.log('INFO: Handling changeLaneLeft for carno '+carNo);
			var command = {
				command: 'c',
				param1: '-68',
				carNo: carNo,
				source: 'ui'
			};
			MainFactory.postCommand(command)
				.then(
				function (response) { // ok
					console.log('INFO: changeLaneLeft command submitted to server: '+response.statusText);
				},
				function (response) { // nok
					console.error('ERROR: Request failed: ' + response.statusText);
				}
				);
		};

		// Handler function for changeLaneRight with the car (-68-+68)
		$scope.changeLaneRight = function (carNo) {
			console.log('INFO: Handling changeLaneRight for carno '+carNo);
			var command = {
				command: 'c',
				param1: '68',
				carNo: carNo,
				source: 'ui'
			};
			MainFactory.postCommand(command)
				.then(
				function (response) { // ok
					console.log('INFO: changeLaneRight command submitted to server: '+response.statusText);
				},
				function (response) { // nok
					console.error('ERROR: Request failed: ' + response.statusText);
				}
				);
		};

		// Handler function for showInitialise
		$scope.showInitialise = function (carNo) {
			console.log('INFO: Handling showInitialise for carno '+carNo);
		};

		// Handler function for changeLight
		$scope.changeLight = function (carNo) {
			console.log('INFO: Handling changeLight for carno '+carNo);
			var command = {
				command: 'l',
				carNo: carNo,
				source: 'ui'
			};
			MainFactory.postCommand(command)
				.then(
				function (response) { // ok
					console.log('INFO: changeLight command submitted to server: '+response.statusText);
				},
				function (response) { // nok
					console.error('ERROR: Request failed: ' + response.statusText);
				}
				);
		};

		// Handler function for changeLightPattern
		$scope.changeLightPattern = function (carNo) {
			console.log('INFO: Handling changeLightPattern for carno '+carNo);
			var command = {
				command: 'lp',
				carNo: carNo,
				source: 'ui'
			};
			MainFactory.postCommand(command)
				.then(
				function (response) { // ok
					console.log('INFO: changeLightPattern command submitted to server: '+response.statusText);
				},
				function (response) { // nok
					console.error('ERROR: Request failed: ' + response.statusText);
				}
				);
		};

		// Handler function for requestLevel
		$scope.requestLevel = function (carNo) {
			console.log('INFO: Handling requestLevel for carno '+carNo);
			var command = {
				command: 'bat',
				carNo: carNo,
				source: 'ui'
			};
			MainFactory.postCommand(command)
				.then(
				function (response) { // ok
					console.log('INFO: requestLevel command submitted to server: '+response.statusText);
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


