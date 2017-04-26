angular.module('whale.services')
.factory('Authentication', ['$q', 'Auth', 'jwtHelper', 'LocalStorage', 'StateManager', 'EndpointProvider', function AuthenticationFactory($q, Auth, jwtHelper, LocalStorage, StateManager, EndpointProvider) {
  'use strict';

  var user = {};
  return {
    init: function() {
      var jwt = LocalStorage.getJWT();
      if (jwt) {
        var tokenPayload = jwtHelper.decodeToken(jwt);
        user.username = tokenPayload.username;
        user.ID = tokenPayload.id;
        user.role = tokenPayload.role;
      }
    },
    login: function(username, password) {
      return $q(function (resolve, reject) {
        Auth.login({username: username, password: password}).$promise
        .then(function(data) {
          LocalStorage.storeJWT(data.jwt);
          var tokenPayload = jwtHelper.decodeToken(data.jwt);
          user.username = username;
          user.ID = tokenPayload.id;
          user.role = tokenPayload.role;
          resolve();
        }, function() {
          reject();
        });
      });
    },
    logout: function() {
      StateManager.clean();
      EndpointProvider.clean();
      LocalStorage.clean();
    },
    isAuthenticated: function() {
      var jwt = LocalStorage.getJWT();
      return jwt && !jwtHelper.isTokenExpired(jwt);
    },
    getUserDetails: function() {
      return user;
    },
    getSSOURI: function() {
      return $q(function (resolve, reject) {
        Auth.getSSOURI().$promise
        .then(function(data) {
          LocalStorage.storeSSOState(data.state);
          resolve(data);
        }, function() {
          reject();
        });
      });
    },
    loginSSO: function(code, state) {
      return $q(function (resolve, reject) {
        if (LocalStorage.getSSOState() !== state) {
          reject();
          return;
        }

        Auth.loginSSO({code: code}).$promise
        .then(function(data) {
          LocalStorage.storeJWT(data.jwt);
          var tokenPayload = jwtHelper.decodeToken(data.jwt);
          user.username = tokenPayload.username;
          user.ID = tokenPayload.id;
          user.role = tokenPayload.role;
          resolve();
        }, function() {
          reject();
        });
      });
    }
  };
}]);
