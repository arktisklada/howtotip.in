angular.module('HowToTip', ['ngMaterial'])

.controller("AppController", function($scope, $window) {
  $scope.$watch("allCountries", function() {
    $scope.allCountries.map(function(country) {
      country.lowercase = angular.lowercase(country.name);
    });
  });

  $scope.searchCountries = function(query) {
    return query ? $scope.allCountries.filter( createLowercaseFilterFor(query) ) : $scope.allCountries;
  };
  function createLowercaseFilterFor(query) {
    var lowercaseQuery = angular.lowercase(query);
    return function(country) {
      return (country.lowercase.indexOf(lowercaseQuery) === 0);
    };
  }

  $scope.searchItemSelected = function(item) {
    $window.location.href = "/" + item.slug;
  };
})

.directive("backgroundImage", function() {
  return {
    restrict: "A",
    link: function($scope, $el, $attrs) {
      imagePath = "/assets/img/" + $attrs.backgroundImage;
      $el.css({"background-image": "url(" + imagePath + ")"});
    }
  };
})

.directive("linkTo", function($window) {
  return {
    restrict: "A",
    link: function($scope, $el, $attrs) {
      $el.on("click", function() {
        $window.location.href = $attrs.linkTo;
      });
    }
  };
})
