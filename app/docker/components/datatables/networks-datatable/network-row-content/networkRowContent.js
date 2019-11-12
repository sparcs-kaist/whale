import { ResourceControlOwnership as RCO } from 'Portainer/models/resourceControl/resourceControlOwnership';

angular.module('portainer.docker')
.directive('networkRowContent', [function networkRowContent() {
  var directive = {
    templateUrl: './networkRowContent.html',
    restrict: 'A',
    transclude: true,
    scope: {
      item: '<',
      parentCtrl: '<',
      allowCheckbox: '<',
      allowExpand: '<'
    },
    controller: ($scope) => {
      $scope.RCO = RCO;
    }
  };
  return directive;
}]);
