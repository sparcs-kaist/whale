angular.module('whale.rest')
.factory('Auth', ['$resource', 'AUTH_ENDPOINT', function AuthFactory($resource, AUTH_ENDPOINT) {
  'use strict';
  return $resource(AUTH_ENDPOINT + '/:type', {}, {
    login: {method: 'POST', params: {type: 'local'}},
    getSSOURI: {method: 'GET', params: {type: 'sso'}},
    loginSSO: {method: 'POST', params: {type: 'sso'}}
  });
}]);
