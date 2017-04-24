angular.module('whale.rest')
.factory('Config', ['$resource', 'CONFIG_ENDPOINT', function ConfigFactory($resource, CONFIG_ENDPOINT) {
  return $resource(CONFIG_ENDPOINT).get();
}]);
