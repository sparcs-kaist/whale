angular.module('whale.services')
.factory('LocalStorage', ['localStorageService', function LocalStorageFactory(localStorageService) {
  'use strict';
  return {
    storeEndpointID: function(id) {
      localStorageService.set('ENDPOINT_ID', id);
    },
    getEndpointID: function() {
      return localStorageService.get('ENDPOINT_ID');
    },
    storeEndpointState: function(state) {
      localStorageService.set('ENDPOINT_STATE', state);
    },
    getEndpointState: function() {
      return localStorageService.get('ENDPOINT_STATE');
    },
    storeApplicationState: function(state) {
      localStorageService.set('APPLICATION_STATE', state);
    },
    getApplicationState: function() {
      return localStorageService.get('APPLICATION_STATE');
    },
    storeJWT: function(jwt) {
      localStorageService.set('JWT', jwt);
    },
    getJWT: function() {
      return localStorageService.get('JWT');
    },
    deleteJWT: function() {
      localStorageService.remove('JWT');
    },
    storePaginationCount: function(key, count) {
      localStorageService.cookie.set('pagination_' + key, count);
    },
    getPaginationCount: function(key) {
      return localStorageService.cookie.get('pagination_' + key);
    },
    clean: function() {
      localStorageService.clearAll();
    },
    storeSSOState: function(state) {
      localStorageService.set('SSO_STATE', state);
    },
    getSSOState: function() {
      return localStorageService.get('SSO_STATE');
    }
  };
}]);
