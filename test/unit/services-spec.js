'use strict';

describe('services', function() {
  var socket;

  beforeEach(module('uchiwa'));
  beforeEach(inject(function(_socket_) {
    socket = _socket_;
  }));

  describe('Page', function() {

    it('should have a title method', inject(function(Page) {
      expect(Page.title).toBeDefined();
    }));

    describe('title()', function() {
      it('should suffix the application title', inject(function(Page) {
        var title = 'Test';
        Page.setTitle(title);
        expect(Page.title()).toBe(title + ' | Uchiwa');
      }));
    });

  });

  describe('stashesService', function() {

    it('should have a stash method', inject(function(stashesService) {
      expect(stashesService.stash).toBeDefined();
    }));

  });

  describe('routingService', function() {

    it('should have a go method', inject(function(routingService) {
      expect(routingService.go).toBeDefined();
    }));
    it('should have a deleteEmptyParameter method', inject(function(routingService) {
      expect(routingService.deleteEmptyParameter).toBeDefined();
    }));
    it('should have a initFilters method', inject(function(routingService) {
      expect(routingService.initFilters).toBeDefined();
    }));
    it('should have a permalink method', inject(function(routingService) {
      expect(routingService.permalink).toBeDefined();
    }));
    it('should have a updateFilters method', inject(function(routingService) {
      expect(routingService.updateFilters).toBeDefined();
    }));
    it('should have a updateValue method', inject(function(routingService) {
      expect(routingService.updateValue).toBeDefined();
    }));
  });

});